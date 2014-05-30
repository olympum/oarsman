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

type Collector struct {
	reftime               int64
	event                 AggregateEvent
	atomicEventChannel    chan AtomicEvent
	aggregateEventChannel chan AggregateEvent
}

func NewCollector(atomicEventChannel chan AtomicEvent, aggregateEventChannel chan AggregateEvent) Collector {
	return Collector{
		atomicEventChannel:    atomicEventChannel,
		aggregateEventChannel: aggregateEventChannel,
		reftime:               0,
		event:                 AggregateEvent{}}
}

func (collector *Collector) Consume(event AtomicEvent) {
	collector.atomicEventChannel <- event

	if collector.aggregateEventChannel == nil {
		return
	}

	if collector.reftime == 0 {
		collector.reftime = event.Time - event.Time%100
		collector.event.Time = collector.reftime + 100
	}

	v := event.Value
	switch event.Label {
	case "total_distance_meters":
		collector.event.Total_distance_meters = v
	case "stroke_rate":
		collector.event.Stroke_rate = v
	case "watts":
		if v > 0 {
			collector.event.Watts = v
		}
	case "calories":
		collector.event.Calories = v
	case "speed_cm_s":
		collector.event.Speed_cm_s = v
	case "heart_rate":
		collector.event.Heart_rate = v
	}

	if event.Time-collector.reftime >= 100 {
		collector.aggregateEventChannel <- collector.event
		collector.reftime = 0
	}

}
