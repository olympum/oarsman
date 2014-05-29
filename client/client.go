package client

import (
	"bufio"
	"fmt"
	"github.com/olympum/gorower/client/s4"
	"io"
	"log"
	"os"
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

type S4Client struct {
	s4      s4.S4Interface
	workout s4.Workout
}

func NewS4Client(options S4Options) S4Client {
	var collector s4.Collector
	if options.CSV != "" {
		aggregateEventWriter := func(aggregateEventChannel chan s4.AggregateEvent) {
			var w io.Writer
			f, err := os.Create(options.CSV)
			if err != nil {
				log.Fatalf("Could not create %s", options.CSV)
			}
			defer f.Close()
			w = bufio.NewWriter(f)
			log.Printf("Writing aggregate data CSV to %s", f.Name())
			fmt.Fprint(w, "time,total_distance_meters,stroke_rate,watts,calories,speed_cm_s,heart_rate\n")
			collector = s4.NewCollector(aggregateEventChannel)
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
		go aggregateEventWriter(make(chan s4.AggregateEvent))
	}

	logger := func(ch chan s4.Event) {
		var writer *os.File
		if options.Out != "" {
			f, err := os.Create(options.Out)
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
			if &collector != nil {
				collector.Consume(event)
			}
		}
	}
	eventChannel := make(chan s4.Event)
	go logger(eventChannel)

	var workout s4.Workout

	var s s4.S4Interface
	if options.In != "" {
		s = s4.NewReplayS4(eventChannel, options.Debug, options.In, options.Replay)
	} else {
		workout = s4.NewWorkout(options.Duration, options.Total_distance_meters)
		s = s4.NewS4(eventChannel, options.Debug)
	}

	return S4Client{s4: s, workout: workout}
}

func (client S4Client) Run() {
	client.s4.Run(client.workout)
}

func (client S4Client) Exit() {
	client.s4.Exit()
}
