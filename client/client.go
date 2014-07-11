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
	TCX                   string
	In                    string
	Out                   string
	Replay                bool
	Debug                 bool
}

func csvWriter(aggregateEventChannel <-chan s4.AggregateEvent, csv string) {
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

var totalTimeSeconds int64
var distanceMeters uint64
var maximumSpeed uint64
var calories uint64
var averageHeartRateBpm uint64
var maximumHeartRateBpm uint64
var n uint64
var start int64
var events []s4.AggregateEvent
var w *bufio.Writer
var f *os.File

func tcxWriter(aggregateEventChannel <-chan s4.AggregateEvent, tcx string) {
	start = 0
	var err error
	f, err = os.Create(tcx)
	if err != nil {
		log.Fatalf("Could not create %s", tcx)
	}
	w = bufio.NewWriter(f)
	events = []s4.AggregateEvent{}
	for {
		event := <-aggregateEventChannel
		if event == s4.EndAggregateEvent {
		} else {
			if start == 0 {
				start = event.Time
			}
			totalTimeSeconds = event.Time - start
			if event.Total_distance_meters > distanceMeters {
				distanceMeters = event.Total_distance_meters
			}
			if event.Speed_cm_s > maximumSpeed {
				maximumSpeed = event.Speed_cm_s
			}
			if event.Calories > calories {
				calories = event.Calories
			}
			if event.Heart_rate > maximumHeartRateBpm {
				maximumHeartRateBpm = event.Heart_rate
			}
			if averageHeartRateBpm == 0 {
				averageHeartRateBpm = event.Heart_rate
			} else if event.Heart_rate != 0 {
				averageHeartRateBpm += event.Heart_rate
			}
			events = append(events, event)
			n++
		}
	}
}

func millisToZulu(millis int64) string {
	return time.Unix(millis/1000, millis%1000*1000).UTC().Format(time.RFC3339)
}

func main() {
	options := S4Options{}
	flag.Uint64Var(&options.Total_distance_meters, "distance", 0, "distance to row in meters")
	flag.DurationVar(&options.Duration, "duration", 0, "duration to row (e.g. 1800s or 45m)")
	flag.BoolVar(&options.Debug, "debug", false, "debug communication data packets")
	flag.StringVar(&options.CSV, "csv", "", "filename to output 10Hz aggregate data as CSV")
	flag.StringVar(&options.TCX, "tcx", "", "filename to output 10Hz aggregate data as TCX")
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
		go csvWriter(aggregateEventChannel, options.CSV)
	}

	if options.TCX != "" {
		aggregateEventChannel = make(chan s4.AggregateEvent)
		go tcxWriter(aggregateEventChannel, options.TCX)
	}

	eventChannel := make(chan s4.AtomicEvent)
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
			log.Printf("Terminating process (received %s signal)", sig.String())
			s.Exit()
			os.Exit(0)
		}
	}()

	s.Run(workout)

	if options.TCX != "" {
		// header
		fmt.Fprintln(w, "<?xml version=\"1.0\"?>")
		fmt.Fprintln(w, "<TrainingCenterDatabase xmlns=\"http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.garmin.com/xmlschemas/ActivityExtension/v2 http://www.garmin.com/xmlschemas/ActivityExtensionv2.xsd http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd\">")
		fmt.Fprintln(w, "<Activities>")
		fmt.Fprintln(w, "<Activity Sport=\"Rowing\">")
		fmt.Fprintf(w, "<Id>%s</Id>\n", millisToZulu(start))
		fmt.Fprintf(w, "<Lap StartTime=\"%s\">\n", millisToZulu(start))
		fmt.Fprintf(w, "<TotalTimeSeconds>%d</TotalTimeSeconds>\n", totalTimeSeconds/1000)
		fmt.Fprintf(w, "<DistanceMeters>%d</DistanceMeters>\n", distanceMeters)
		fmt.Fprintf(w, "<MaximumSpeed>%f</MaximumSpeed>\n", float64(maximumSpeed)/100.0)
		fmt.Fprintf(w, "<Calories>%d</Calories>\n", calories/1000)
		fmt.Fprintln(w, "<AverageHeartRateBpm>")
		fmt.Fprintf(w, "<Value>%d</Value>\n", averageHeartRateBpm/n)
		fmt.Fprintln(w, "</AverageHeartRateBpm>")
		fmt.Fprintln(w, "<MaximumHeartRateBpm>")
		fmt.Fprintf(w, "<Value>%d</Value>\n", maximumHeartRateBpm)
		fmt.Fprintln(w, "</MaximumHeartRateBpm>")
		fmt.Fprintln(w, "<Intensity>Active</Intensity>")
		fmt.Fprintln(w, "<TriggerMethod>Manual</TriggerMethod>")
		fmt.Fprintln(w, "<Track>")

		for _, e := range events {
			fmt.Fprintln(w, "<Trackpoint>")
			fmt.Fprintf(w, "<Time>%s</Time>\n", millisToZulu(e.Time))
			fmt.Fprintf(w, "<DistanceMeters>%d</DistanceMeters>\n", e.Total_distance_meters)
			fmt.Fprintln(w, "<HeartRateBpm xsi:type=\"HeartRateInBeatsPerMinute_t\">")
			fmt.Fprintf(w, "<Value>%d</Value>\n", e.Heart_rate)
			fmt.Fprintln(w, "</HeartRateBpm>")
			fmt.Fprintf(w, "<Cadence>%d</Cadence>\n", e.Stroke_rate)
			fmt.Fprintln(w, "<Extensions>")
			fmt.Fprintln(w, "<TPX xmlns=\"http://www.garmin.com/xmlschemas/ActivityExtension/v2\">")
			fmt.Fprintf(w, "<Speed>%f</Speed>\n", float64(e.Speed_cm_s)/100.0)
			fmt.Fprintf(w, "<Watts>%d</Watts>\n", e.Watts)
			fmt.Fprintln(w, "</TPX>")
			fmt.Fprintln(w, "</Extensions>")
			fmt.Fprintln(w, "</Trackpoint>")
		}

		fmt.Fprintln(w, "</Track>")
		fmt.Fprintln(w, "</Lap>")
		fmt.Fprintln(w, "</Activity>")
		fmt.Fprintln(w, "</Activities>")
		fmt.Fprintln(w, "</TrainingCenterDatabase>")

		w.Flush()

		defer f.Close()
	}
}
