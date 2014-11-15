package s4

import (
	"bufio"
	jww "github.com/spf13/jwalterweatherman"
	"os"
	"strconv"
	"strings"
	t "time"
)

type ReplayS4 struct {
	scanner    *bufio.Scanner
	aggregator Aggregator
	replay     bool
	debug      bool
}

func NewReplayS4(eventChannel chan<- AtomicEvent, aggregateEventChannel chan<- AggregateEvent, debug bool, replayfile string, replay bool) (S4Interface, error) {
	f, err := os.Open(replayfile)
	if err != nil {
		jww.FATAL.Printf("Could not read %s\n", replayfile)
		return nil, err
	}
	jww.DEBUG.Printf("Reading from %s\n", f.Name())
	s := bufio.NewScanner(f)
	aggregator := newAggregator(eventChannel, aggregateEventChannel)
	return &ReplayS4{scanner: s, aggregator: aggregator, replay: replay, debug: debug}, nil
}

func (s4 *ReplayS4) Run(workout *S4Workout) {
	for s4.scanner.Scan() {
		line := s4.scanner.Text()
		tokens := strings.Split(line, " ")
		if len(tokens) < 2 {
			continue
		}
		time, _ := strconv.ParseInt(tokens[0], 10, 64)
		if time == 0 {
			// skip incorrect rows
			continue
		}
		values := strings.Split(tokens[1], ":")
		if len(values) < 2 {
			continue
		}
		label := values[0]
		value, _ := strconv.ParseUint(values[1], 10, 64)
		event := AtomicEvent{Time: time, Label: label, Value: value}
		if s4.debug {
			jww.DEBUG.Println(event)
		}
		s4.aggregator.consume(event)
		if s4.replay {
			t.Sleep(t.Millisecond * 25)
		}
	}
	s4.aggregator.complete()
}

func (s4 *ReplayS4) Exit() {
}
