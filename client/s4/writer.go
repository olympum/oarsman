package s4

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

type WriterFunc func(collector *EventCollector, writer *bufio.Writer)

func CSVWriter(collector *EventCollector, writer *bufio.Writer) {
	fmt.Fprint(writer, "time,total_distance_meters,stroke_rate,watts,calories,speed_cm_s,heart_rate\n")
	for _, event := range collector.events {
		fmt.Fprintf(writer, "%d,%d,%d,%d,%d,%d,%d\n",
			event.Time,
			event.Total_distance_meters,
			event.Stroke_rate,
			event.Watts,
			event.Calories,
			event.Speed_cm_s,
			event.Heart_rate)
	}
}

func millisToZulu(millis int64) string {
	return time.Unix(millis/1000, millis%1000*1000).UTC().Format(time.RFC3339)
}

func TCXWriter(collector *EventCollector, writer *bufio.Writer) {
	// header
	w := writer
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
	fmt.Fprintf(w, "<Value>%d</Value>\n", collector.averageHeartRateBpm/collector.nSamples)
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

func ExportCollectorEvents(collector *EventCollector, filename string, writerFunc WriterFunc) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Could not create %s", filename)
	}
	defer f.Close()

	var w *bufio.Writer
	w = bufio.NewWriter(f)
	log.Printf("Writing aggregate data to %s", f.Name())
	writerFunc(collector, w)
}
