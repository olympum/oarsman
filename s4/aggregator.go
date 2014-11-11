package s4

type AggregateEvent struct {
	Time                  int64
	Total_distance_meters uint64
	Stroke_rate           uint64
	Watts                 uint64
	Calories              uint64
	Speed_cm_s            uint64
	Heart_rate            uint64
}

const RESOLUTION = 1000

type Aggregator struct {
	reftime               int64
	event                 AggregateEvent
	atomicEventChannel    chan<- AtomicEvent
	aggregateEventChannel chan<- AggregateEvent
}

func newAggregator(atomicEventChannel chan<- AtomicEvent, aggregateEventChannel chan<- AggregateEvent) Aggregator {
	return Aggregator{
		atomicEventChannel:    atomicEventChannel,
		aggregateEventChannel: aggregateEventChannel,
		reftime:               0,
		event:                 AggregateEvent{}}
}

func (aggregator *Aggregator) complete() {
	if aggregator.reftime != 0 {
		aggregator.aggregateEventChannel <- aggregator.event
		aggregator.reftime = 0
	}
}

func (aggregator *Aggregator) consume(event AtomicEvent) {
	if aggregator.atomicEventChannel != nil {
		aggregator.atomicEventChannel <- event
	}

	if aggregator.aggregateEventChannel == nil {
		return
	}

	if aggregator.reftime == 0 {
		aggregator.reftime = event.Time - event.Time%RESOLUTION
		aggregator.event.Time = aggregator.reftime + RESOLUTION
	}

	v := event.Value
	switch event.Label {
	case "total_distance_meters":
		aggregator.event.Total_distance_meters = v
	case "stroke_rate":
		aggregator.event.Stroke_rate = v
	case "watts":
		if v > 0 {
			aggregator.event.Watts = v
		}
	case "calories":
		aggregator.event.Calories = v
	case "speed_cm_s":
		aggregator.event.Speed_cm_s = v
	case "heart_rate":
		if v > 0 {
			aggregator.event.Heart_rate = v
		}
	}

	if event.Time-aggregator.reftime >= RESOLUTION {
		aggregator.complete()
	}

}
