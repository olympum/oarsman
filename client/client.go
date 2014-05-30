package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/olympum/gorower/client/s4"
	"io"
	"log"
	"os"
	"os/signal"
	"time"
)

type S4Options struct {
	Total_distance_meters uint64
	Duration              time.Duration
	CSV                   string
	In                    string
	Out                   string
	Replay                bool
	Debug                 bool
}

func aggregateEventWriter(aggregateEventChannel chan s4.AggregateEvent, csv string) {
	var w io.Writer
	f, err := os.Create(csv)
	if err != nil {
		log.Fatalf("Could not create %s", csv)
	}
	defer f.Close()
	w = bufio.NewWriter(f)
	log.Printf("Writing aggregate data CSV to %s", f.Name())
	fmt.Fprint(w, "time,total_distance_meters,stroke_rate,watts,calories,speed_cm_s,heart_rate\n")
	for {
		event := <-aggregateEventChannel
		fmt.Fprintf(w, "%d,%d,%d,%d,%d,%d,%d\n",
			event.Time,
			event.Total_distance_meters,
			event.Stroke_rate,
			event.Watts,
			event.Calories,
			event.Speed_cm_s,
			event.Heart_rate)
	}
}

func logger(ch chan s4.AtomicEvent, out string) {
	var writer *os.File
	if out != "" {
		f, err := os.Create(out)
		if err != nil {
			log.Fatal(err)
		}
		writer = f
	} else {
		writer = os.Stdout
	}

	log.Printf("Writing to %s", writer.Name())

	for {
		event := <-ch
		fmt.Fprintf(writer, "%d %s:%d\n", event.Time, event.Label, event.Value)
	}
}

func main() {
	options := S4Options{}
	flag.Uint64Var(&options.Total_distance_meters, "distance", 0, "distance to row in meters")
	flag.DurationVar(&options.Duration, "duration", 0, "duration to row (e.g. 1800s or 45m)")
	flag.BoolVar(&options.Debug, "debug", false, "debug communication data packets")
	flag.StringVar(&options.CSV, "csv", "", "filename to output 10Hz aggregate data as CSV")
	flag.StringVar(&options.In, "in", "", "filename with raw S4 data to read")
	flag.StringVar(&options.Out, "out", "", "filename where to output raw S4 data, default stdout")
	flag.BoolVar(&options.Replay, "replay", false, "replay input data in accurate time")

	flag.Parse()
	if !flag.Parsed() {
		log.Fatal(flag.ErrHelp)
	}

	// client mode
	log.Println("CLI mode - Press CTRL+C to interrupt")

	var aggregateEventChannel chan s4.AggregateEvent

	if options.CSV != "" {
		aggregateEventChannel = make(chan s4.AggregateEvent)
		go aggregateEventWriter(aggregateEventChannel, options.CSV)
	}

	eventChannel := make(chan s4.AtomicEvent)
	go logger(eventChannel, options.Out)

	var workout s4.S4Workout

	var s s4.S4Interface
	if options.In != "" {
		s = s4.NewReplayS4(eventChannel, aggregateEventChannel, options.Debug, options.In, options.Replay)
	} else {
		workout = s4.NewS4Workout()
		workout.AddSingleWorkout(options.Duration, options.Total_distance_meters)
		s = s4.NewS4(eventChannel, aggregateEventChannel, options.Debug)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Terminating process (received %s signal)", sig.String())
			s.Exit()
			os.Exit(0)
		}
	}()

	s.Run(workout)
}
