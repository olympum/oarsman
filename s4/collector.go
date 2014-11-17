package s4

import (
	jww "github.com/spf13/jwalterweatherman"
)

type EventCollector struct {
	channel  <-chan AggregateEvent
	activity *Activity
}

func NewEventCollector(aggregateEventChannel <-chan AggregateEvent) *EventCollector {
	return &EventCollector{channel: aggregateEventChannel, activity: NewActivity(nil, nil)}
}

func (collector *EventCollector) Run() {
	activity := collector.activity
	activity.addLap()

	for {
		event := <-collector.channel
		jww.DEBUG.Printf("Received event to collect: %v", event)
		activity.lastLap().AddEvent(event)
		if event.Total_distance_meters > 0 && event.Total_distance_meters%2000 == 0 {
			lap := activity.addLap()
			jww.DEBUG.Printf("Added auto-lap at %d meters", event.Total_distance_meters)
			lap.AddEvent(event)
		}
	}
}

func (collector *EventCollector) Activity() *Activity {
	activity := collector.activity
	if activity.firstLap() == nil || len(activity.firstLap().events) == 0 {
		return nil
	}
	return activity.update()
}
