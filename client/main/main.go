package main

import (
	"flag"
	"github.com/olympum/gorower/client"
	"log"
	"os"
	"os/signal"
)

func main() {
	s4Options := client.S4Options{}
	flag.Uint64Var(&s4Options.Total_distance_meters, "distance", 0, "distance to row in meters")
	flag.DurationVar(&s4Options.Duration, "duration", 0, "duration to row (e.g. 1800s or 45m)")
	flag.BoolVar(&s4Options.Debug, "debug", false, "debug communication data packets")
	flag.StringVar(&s4Options.CSV, "csv", "", "filename to output 10Hz aggregate data as CSV")
	flag.StringVar(&s4Options.In, "in", "", "filename with raw S4 data to read")
	flag.StringVar(&s4Options.Out, "out", "", "filename where to output raw S4 data, default stdout")
	flag.BoolVar(&s4Options.Replay, "replay", false, "replay input data in accurate time")

	flag.Parse()
	if !flag.Parsed() {
		log.Fatal(flag.ErrHelp)
	}

	// client mode
	log.Println("CLI mode - Press CTRL+C to interrupt")
	client := client.NewS4Client(s4Options)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Printf("Terminating process (received %s signal)", sig.String())
			client.Exit()
			os.Exit(0)
		}
	}()

	client.Run()
}
