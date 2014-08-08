package main

import (
	"flag"
	"fmt"
	"github.com/olympum/gorower/client/s4"
	"os"
	"os/signal"
	"time"
)

type S4Options struct {
	Total_distance_meters uint64
	Duration              time.Duration
	CSV                   string
	TCX                   string
	In                    string
	Out                   string
	Replay                bool
	Debug                 bool
}

func ParseCLI(options S4Options) {
	flag.Uint64Var(&options.Total_distance_meters, "distance", 0, "distance to row in meters")
	flag.DurationVar(&options.Duration, "duration", 0, "duration to row (e.g. 1800s or 45m)")
	flag.BoolVar(&options.Debug, "debug", false, "debug communication data packets")
	flag.StringVar(&options.CSV, "csv", "", "filename to output 10Hz aggregate data as CSV")
	flag.StringVar(&options.TCX, "tcx", "", "filename to output 10Hz aggregate data as TCX")
	flag.StringVar(&options.In, "in", "", "filename with raw S4 data to read")
	flag.StringVar(&options.Out, "out", "", "filename where to output raw S4 data, default stdout")
	flag.BoolVar(&options.Replay, "replay", false, "replay input data in accurate time")

	flag.Parse()
	if !flag.Parsed() || len(os.Args) == 1 {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	options := S4Options{}
	ParseCLI(options)

	// client mode
	fmt.Println("CLI mode - Press CTRL+C to interrupt")

	eventChannel := make(chan s4.AtomicEvent)
	aggregateEventChannel := make(chan s4.AggregateEvent)
	collector := s4.NewEventCollector(aggregateEventChannel)
	go collector.Run()

	go s4.Logger(eventChannel, options.Out)

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
			fmt.Printf("Terminating process (received %s signal)", sig.String())
			s.Exit()
			os.Exit(0)
		}
	}()

	s.Run(workout)

	if options.CSV != "" {
		s4.ExportCollectorEvents(collector, options.CSV, s4.CSVWriter)
	}

	if options.TCX != "" {
		s4.ExportCollectorEvents(collector, options.TCX, s4.TCXWriter)
	}

}
