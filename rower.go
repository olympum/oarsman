package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("Press CTRL+C to interrupt")

	var distanceFlag = flag.Int64("distance", 10000, "distance to row in meters")
	var durationFlag = flag.Duration("duration", 0, "duration to row (e.g. 1800s or 45m")
	var debug = flag.Bool("debug", false, "debug communication data packets")
	flag.Parse()
	if !flag.Parsed() {
		log.Fatal(flag.ErrHelp)
	}

	logCallback := func(event Event) {
		fmt.Printf("%d %s:%d\n", event.time, event.label, event.value)
	}
	workout := NewWorkout(*durationFlag, *distanceFlag)

	s4 := NewS4(logCallback, *debug)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			log.Print("Received signal " + sig.String())
			s4.Exit()
			os.Exit(0)
		}
	}()

	// TODO allow goroutine channel to interrupt workout
	s4.Run(workout)
}
