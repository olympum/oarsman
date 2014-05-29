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
	Reftime int64
	Event   AggregateEvent
	Channel chan AggregateEvent
}

func NewCollector(channel chan AggregateEvent) Collector {
	return Collector{Channel: channel, Reftime: 0, Event: AggregateEvent{}}
}

func (collector *Collector) Consume(event Event) {
	if collector.Reftime == 0 {
		collector.Reftime = event.Time - event.Time%100
		collector.Event.Time = collector.Reftime + 100
	}

	v := event.Value
	switch event.Label {
	case "total_distance_meters":
		collector.Event.Total_distance_meters = v
	case "stroke_rate":
		collector.Event.Stroke_rate = v
	case "watts":
		if v > 0 {
			collector.Event.Watts = v
		}
	case "calories":
		collector.Event.Calories = v
	case "speed_cm_s":
		collector.Event.Speed_cm_s = v
	case "heart_rate":
		collector.Event.Heart_rate = v
	}

	if event.Time-collector.Reftime >= 100 {
		collector.Channel <- collector.Event
		collector.Reftime = 0
	}
}
