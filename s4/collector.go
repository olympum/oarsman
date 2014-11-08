package s4

type EventCollector struct {
	channel  <-chan AggregateEvent
	events   []AggregateEvent
	Activity *Activity
}

func NewEventCollector(aggregateEventChannel <-chan AggregateEvent) *EventCollector {
	return &EventCollector{channel: aggregateEventChannel, events: []AggregateEvent{}, Activity: &Activity{}}
}

func (collector *EventCollector) Run() {
	collector.Activity.StartTimeMilliseconds = 0
	collector.events = []AggregateEvent{}
	for {
		event := <-collector.channel
		if event.Time != 0 {
			n := collector.Activity.numSamples

			if collector.Activity.StartTimeMilliseconds == 0 {
				collector.Activity.StartTimeMilliseconds = event.Time
			}
			collector.Activity.TotalTimeMilliseconds = event.Time - collector.Activity.StartTimeMilliseconds
			if event.Total_distance_meters > collector.Activity.DistanceMeters {
				collector.Activity.DistanceMeters = event.Total_distance_meters
			}
			if event.Speed_cm_s > collector.Activity.MaximumSpeed_cm_s {
				collector.Activity.MaximumSpeed_cm_s = event.Speed_cm_s
			}
			if event.Calories > collector.Activity.Calories {
				collector.Activity.Calories = event.Calories
			}
			if event.Heart_rate > collector.Activity.MaximumHeartRateBpm {
				collector.Activity.MaximumHeartRateBpm = event.Heart_rate
			}
			if collector.Activity.AverageHeartRateBpm == 0 {
				collector.Activity.AverageHeartRateBpm = float64(event.Heart_rate)
			} else {
				collector.Activity.AverageHeartRateBpm = (collector.Activity.AverageHeartRateBpm*float64(n) + float64(event.Heart_rate)) / (float64(n) + 1.0)
			}
			if collector.Activity.AveragePower == 0 {
				collector.Activity.AveragePower = float64(event.Watts)
			} else {
				collector.Activity.AveragePower = (collector.Activity.AveragePower*float64(n) + float64(event.Watts)) / (float64(n) + 1.0)
			}
			if event.Watts > collector.Activity.MaximumPower {
				collector.Activity.MaximumPower = event.Watts
			}
			if collector.Activity.AverageCadence == 0 {
				collector.Activity.AverageCadence = float64(event.Stroke_rate)
			} else {
				collector.Activity.AverageCadence = (collector.Activity.AverageCadence*float64(n) + float64(event.Stroke_rate)) / (float64(n) + 1.0)
			}
			if event.Stroke_rate > collector.Activity.MaximumCadence {
				collector.Activity.MaximumCadence = event.Stroke_rate
			}
			collector.events = append(collector.events, event)
			collector.Activity.numSamples++
		}
	}
}
