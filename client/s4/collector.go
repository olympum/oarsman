package s4

type EventCollector struct {
	channel             <-chan AggregateEvent
	totalTimeSeconds    int64
	distanceMeters      uint64
	maximumSpeed        uint64
	calories            uint64
	averageHeartRateBpm uint64
	maximumHeartRateBpm uint64
	nSamples            uint64
	start               int64
	events              []AggregateEvent
}

func NewEventCollector(aggregateEventChannel <-chan AggregateEvent) *EventCollector {
	return &EventCollector{channel: aggregateEventChannel, events: []AggregateEvent{}}
}

func (collector *EventCollector) Run() {
	collector.start = 0
	collector.events = []AggregateEvent{}
	for {
		event := <-collector.channel
		if event == EndAggregateEvent {
		} else {
			if collector.start == 0 {
				collector.start = event.Time
			}
			collector.totalTimeSeconds = event.Time - collector.start
			if event.Total_distance_meters > collector.distanceMeters {
				collector.distanceMeters = event.Total_distance_meters
			}
			if event.Speed_cm_s > collector.maximumSpeed {
				collector.maximumSpeed = event.Speed_cm_s
			}
			if event.Calories > collector.calories {
				collector.calories = event.Calories
			}
			if event.Heart_rate > collector.maximumHeartRateBpm {
				collector.maximumHeartRateBpm = event.Heart_rate
			}
			if collector.averageHeartRateBpm == 0 {
				collector.averageHeartRateBpm = event.Heart_rate
			} else {
				collector.averageHeartRateBpm += event.Heart_rate
			}
			collector.events = append(collector.events, event)
			collector.nSamples++
		}
	}
}
