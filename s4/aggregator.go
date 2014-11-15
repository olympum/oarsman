package s4

import (
	jww "github.com/spf13/jwalterweatherman"
)

type AggregateEvent struct {
	Time_start            int64
	Time                  int64
	Start_distance_meters uint64
	Total_distance_meters uint64
	Stroke_rate           uint64
	Watts                 uint64
	Calories              uint64
	Speed_m_s             float64
	Heart_rate            uint64
}

const MAX_RESOLUTION_MILLIS = 10000

type Aggregator struct {
	event                 *AggregateEvent
	atomicEventChannel    chan<- AtomicEvent
	aggregateEventChannel chan<- AggregateEvent
}

func newAggregator(atomicEventChannel chan<- AtomicEvent, aggregateEventChannel chan<- AggregateEvent) Aggregator {
	return Aggregator{
		atomicEventChannel:    atomicEventChannel,
		aggregateEventChannel: aggregateEventChannel,
		event: &AggregateEvent{}}
}

func (aggregator Aggregator) send(event *AggregateEvent) bool {
	if aggregator.aggregateEventChannel == nil {
		return false
	}

	aggregator.aggregateEventChannel <- *event
	aggregator.event = &AggregateEvent{}
	aggregator.event.Start_distance_meters = event.Total_distance_meters
	aggregator.event.Total_distance_meters = event.Total_distance_meters
	jww.DEBUG.Print("Sent aggregate event", event)
	return true
}

func (aggregator *Aggregator) complete() {
	e := aggregator.event
	delta_time := float64(e.Time - e.Time_start)
	delta_distance := float64(e.Total_distance_meters - e.Start_distance_meters)
	if delta_time > 0 && delta_distance > 0 {
		e.Speed_m_s = delta_distance * 1000.0 / delta_time
		aggregator.send(e)
	}

}

func (aggregator *Aggregator) consume(event AtomicEvent) {
	if aggregator.atomicEventChannel != nil {
		aggregator.atomicEventChannel <- event
		jww.DEBUG.Print("Sent atomic event", event)
	}

	if aggregator.aggregateEventChannel == nil {
		return
	}

	e := aggregator.event
	e.Time = event.Time

	if e.Time_start == 0 {
		e.Time_start = e.Time
	}

	v := event.Value
	switch event.Label {
	case "total_distance_meters":
		e.Total_distance_meters = v
		if v == 0 {
			aggregator.send(e)
		}
	case "stroke_rate":
		e.Stroke_rate = v
	case "watts":
		if v > 0 {
			e.Watts = v
		}
	case "calories":
		e.Calories = v
	case "heart_rate":
		if v > 0 {
			e.Heart_rate = v
		}
	}

	if e.Time-e.Time_start >= MAX_RESOLUTION_MILLIS {
		aggregator.complete()
	}

	// auto-laps every 2000 meters
	if e.Total_distance_meters%2000 == 0 {
		aggregator.complete()
	}

}
