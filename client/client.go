package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/olympum/gorower/client/s4"
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

type S4Writer struct {
	writer *bufio.Writer
	file   *os.File
}

func NewS4Writer(filename string) *S4Writer {
	var w *bufio.Writer
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create %s", filename)
	}
	w = bufio.NewWriter(f)
	log.Printf("Writing aggregate data to %s", f.Name())
	return &S4Writer{writer: w, file: f}
}

func (writer *S4Writer) WriteCSV(collector *EventCollector) {
	fmt.Fprint(writer.writer, "time,total_distance_meters,stroke_rate,watts,calories,speed_cm_s,heart_rate\n")
	for _, event := range collector.events {
		fmt.Fprintf(writer.writer, "%d,%d,%d,%d,%d,%d,%d\n",
			event.Time,
			event.Total_distance_meters,
			event.Stroke_rate,
			event.Watts,
			event.Calories,
			event.Speed_cm_s,
			event.Heart_rate)
	}
}

func (writer *S4Writer) WriteTCX(collector *EventCollector) {
	// header
	w := writer.writer
	fmt.Fprintln(w, "<?xml version=\"1.0\"?>")
	fmt.Fprintln(w, "<TrainingCenterDatabase xmlns=\"http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.garmin.com/xmlschemas/ActivityExtension/v2 http://www.garmin.com/xmlschemas/ActivityExtensionv2.xsd http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd\">")
	fmt.Fprintln(w, "<Activities>")
	fmt.Fprintln(w, "<Activity Sport=\"Rowing\">")
	fmt.Fprintf(w, "<Id>%s</Id>\n", millisToZulu(collector.start))
	fmt.Fprintf(w, "<Lap StartTime=\"%s\">\n", millisToZulu(collector.start))
	fmt.Fprintf(w, "<TotalTimeSeconds>%d</TotalTimeSeconds>\n", collector.totalTimeSeconds/1000)
	fmt.Fprintf(w, "<DistanceMeters>%d</DistanceMeters>\n", collector.distanceMeters)
	fmt.Fprintf(w, "<MaximumSpeed>%f</MaximumSpeed>\n", float64(collector.maximumSpeed)/100.0)
	fmt.Fprintf(w, "<Calories>%d</Calories>\n", collector.calories/1000)
	fmt.Fprintln(w, "<AverageHeartRateBpm>")
	fmt.Fprintf(w, "<Value>%d</Value>\n", collector.averageHeartRateBpm/collector.n)
	fmt.Fprintln(w, "</AverageHeartRateBpm>")
	fmt.Fprintln(w, "<MaximumHeartRateBpm>")
	fmt.Fprintf(w, "<Value>%d</Value>\n", collector.maximumHeartRateBpm)
	fmt.Fprintln(w, "</MaximumHeartRateBpm>")
	fmt.Fprintln(w, "<Intensity>Active</Intensity>")
	fmt.Fprintln(w, "<TriggerMethod>Manual</TriggerMethod>")
	fmt.Fprintln(w, "<Track>")

	for _, e := range collector.events {
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
}

func (writer *S4Writer) Close() {
	writer.file.Close()
}

type EventCollector struct {
	channel             <-chan s4.AggregateEvent
	totalTimeSeconds    int64
	distanceMeters      uint64
	maximumSpeed        uint64
	calories            uint64
	averageHeartRateBpm uint64
	maximumHeartRateBpm uint64
	n                   uint64
	start               int64
	events              []s4.AggregateEvent
}

func NewEventCollector(aggregateEventChannel <-chan s4.AggregateEvent) *EventCollector {
	return &EventCollector{channel: aggregateEventChannel, events: []s4.AggregateEvent{}}
}

func (collector *EventCollector) Run() {
	collector.start = 0
	collector.events = []s4.AggregateEvent{}
	for {
		event := <-collector.channel
		if event == s4.EndAggregateEvent {
		} else {
			if collector.start == 0 {
				collector.start = event.Time
			}
			collector.totalTimeSeconds = event.Time - collector.start
			if event.Total_distance_meters > collector.distanceMeters {
				collector.distanceMeters = event.Total_distance_meters
			}
			if event.Speed_cm_s > collector.maximumSpeed {
				collector.maximumSpeed = event.Speed_cm_s
			}
			if event.Calories > collector.calories {
				collector.calories = event.Calories
			}
			if event.Heart_rate > collector.maximumHeartRateBpm {
				collector.maximumHeartRateBpm = event.Heart_rate
			}
			if collector.averageHeartRateBpm == 0 {
				collector.averageHeartRateBpm = event.Heart_rate
			} else if event.Heart_rate != 0 {
				collector.averageHeartRateBpm += event.Heart_rate
			}
			collector.events = append(collector.events, event)
			collector.n++
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

	eventChannel := make(chan s4.AtomicEvent)
	aggregateEventChannel := make(chan s4.AggregateEvent)
	collector := NewEventCollector(aggregateEventChannel)
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
			log.Printf("Terminating process (received %s signal)", sig.String())
			s.Exit()
			os.Exit(0)
		}
	}()

	s.Run(workout)

	if options.CSV != "" {
		file := options.CSV
		writer := NewS4Writer(file)
		writer.WriteCSV(collector)
		writer.Close()
	}

	if options.TCX != "" {
		file := options.TCX
		writer := NewS4Writer(file)
		writer.WriteTCX(collector)
		writer.Close()
	}

}
