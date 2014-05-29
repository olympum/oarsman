package s4

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	t "time"
)

type ReplayS4 struct {
	scanner      *bufio.Scanner
	eventChannel chan Event
	replay       bool
	debug        bool
}

func NewReplayS4(eventChannel chan Event, debug bool, replayfile string, replay bool) S4Interface {
	f, err := os.Open(replayfile)
	if err != nil {
		log.Fatalf("Could not read %s", replayfile)
	}
	log.Printf("Reading from %s", f.Name())
	s := bufio.NewScanner(f)
	return &ReplayS4{scanner: s, eventChannel: eventChannel, replay: replay, debug: debug}
}

func (s4 *ReplayS4) Run(workout Workout) {
	for s4.scanner.Scan() {
		line := s4.scanner.Text()
		tokens := strings.Split(line, " ")
		time, _ := strconv.ParseInt(tokens[0], 10, 64)
		values := strings.Split(tokens[1], ":")
		label := values[0]
		value, _ := strconv.ParseUint(values[1], 10, 64)
		event := Event{Time: time, Label: label, Value: value}
		if s4.debug {
			log.Print(event)
		}
		s4.eventChannel <- event
		if s4.replay {
			t.Sleep(t.Millisecond * 25)
		}
	}
}

func (s4 *ReplayS4) Exit() {
}
