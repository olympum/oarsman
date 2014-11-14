package s4

import (
	"bufio"
	"fmt"
	"github.com/olympum/oarsman/util"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

type WriterFunc func(collector *EventCollector, writer *bufio.Writer)

func CSVWriter(collector *EventCollector, writer *bufio.Writer) {
	laps := collector.laps
	if len(laps) == 0 {
		jww.INFO.Println("Empty activity")
		return
	} else {
		jww.INFO.Printf("Writing %d laps in CSV", len(laps))
	}
	fmt.Fprint(writer, "time,total_distance_meters,stroke_rate,watts,calories,speed_m_s,heart_rate\n")
	for n, lap := range laps {
		jww.INFO.Printf("Writing lap %d (%v meters)", n, lap.DistanceMeters)
		for _, event := range lap.events {
			fmt.Fprintf(writer, "%d,%d,%d,%d,%d,%.2f,%d\n",
				event.Time,
				event.Total_distance_meters,
				event.Stroke_rate,
				event.Watts,
				event.Calories,
				event.Speed_m_s,
				event.Heart_rate)
		}
	}
}

func TCXWriter(collector *EventCollector, writer *bufio.Writer) {
	laps := collector.laps
	if len(laps) == 0 {
		jww.INFO.Println("Empty activity")
		return
	} else {
		jww.INFO.Printf("Writing %d laps in TCX", len(laps))
	}

	// header
	w := writer
	fmt.Fprintln(w, "<?xml version=\"1.0\"?>")
	fmt.Fprintln(w, "<TrainingCenterDatabase xmlns=\"http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xsi:schemaLocation=\"http://www.garmin.com/xmlschemas/ActivityExtension/v2 http://www.garmin.com/xmlschemas/ActivityExtensionv2.xsd http://www.garmin.com/xmlschemas/TrainingCenterDatabase/v2 http://www.garmin.com/xmlschemas/TrainingCenterDatabasev2.xsd\">")
	fmt.Fprintln(w, "<Activities>")
	fmt.Fprintln(w, "<Activity Sport=\"Other\">")
	fmt.Fprintf(w, "<Id>%s</Id>\n", laps[0].StartTimeZulu)
	fmt.Fprint(w, "<Creator><Name>Oarsman (WaterRower S4)</Name></Creator>")
	for n, lap := range laps {
		jww.INFO.Printf("Writing lap %d (%v meters)", n, lap.DistanceMeters)
		fmt.Fprintf(w, "<Lap StartTime=\"%s\">\n", lap.StartTimeZulu)
		fmt.Fprintf(w, "<TotalTimeSeconds>%d</TotalTimeSeconds>\n", lap.TotalTimeSeconds)
		fmt.Fprintf(w, "<DistanceMeters>%d</DistanceMeters>\n", lap.DistanceMeters)
		fmt.Fprintf(w, "<MaximumSpeed>%f</MaximumSpeed>\n", lap.MaximumSpeedMs)
		fmt.Fprintf(w, "<Calories>%d</Calories>\n", lap.KCalories)
		fmt.Fprintln(w, "<AverageHeartRateBpm>")
		fmt.Fprintf(w, "<Value>%d</Value>\n", lap.AverageHeartRateBpm)
		fmt.Fprintln(w, "</AverageHeartRateBpm>")
		fmt.Fprintln(w, "<MaximumHeartRateBpm>")
		fmt.Fprintf(w, "<Value>%d</Value>\n", lap.MaximumHeartRateBpm)
		fmt.Fprintln(w, "</MaximumHeartRateBpm>")
		fmt.Fprintln(w, "<Intensity>Active</Intensity>")
		fmt.Fprintln(w, "<TriggerMethod>Manual</TriggerMethod>")
		fmt.Fprintln(w, "<Track>")

		for _, e := range lap.events {
			fmt.Fprintln(w, "<Trackpoint>")
			fmt.Fprintf(w, "<Time>%s</Time>\n", util.MillisToZulu(e.Time))
			fmt.Fprintf(w, "<DistanceMeters>%d</DistanceMeters>\n", e.Total_distance_meters)
			fmt.Fprintln(w, "<HeartRateBpm xsi:type=\"HeartRateInBeatsPerMinute_t\">")
			fmt.Fprintf(w, "<Value>%d</Value>\n", e.Heart_rate)
			fmt.Fprintln(w, "</HeartRateBpm>")
			fmt.Fprintf(w, "<Cadence>%d</Cadence>\n", e.Stroke_rate)
			fmt.Fprintln(w, "<Extensions>")
			fmt.Fprintln(w, "<TPX xmlns=\"http://www.garmin.com/xmlschemas/ActivityExtension/v2\">")
			fmt.Fprintf(w, "<Speed>%.2f</Speed>\n", e.Speed_m_s)
			fmt.Fprintf(w, "<Watts>%d</Watts>\n", e.Watts)
			fmt.Fprintln(w, "</TPX>")
			fmt.Fprintln(w, "</Extensions>")
			fmt.Fprintln(w, "</Trackpoint>")
		}

		fmt.Fprintln(w, "</Track>")
		fmt.Fprintln(w, "</Lap>")
	}
	fmt.Fprintln(w, "</Activity>")
	fmt.Fprintln(w, "</Activities>")
	fmt.Fprintln(w, "</TrainingCenterDatabase>")

	w.Flush()
}

func ExportCollectorEvents(collector *EventCollector, filename string, writerFunc WriterFunc) {
	f, err := os.Create(filename)
	if err != nil {
		jww.FATAL.Printf("Could not create %s\n", filename)
	}
	defer f.Close()

	var w *bufio.Writer
	w = bufio.NewWriter(f)
	jww.INFO.Printf("Writing aggregate data to %s\n", f.Name())
	writerFunc(collector, w)
}
