package s4

import (
	jww "github.com/spf13/jwalterweatherman"
)

type EventCollector struct {
	channel <-chan AggregateEvent
	laps    []*Lap
}

func NewEventCollector(aggregateEventChannel <-chan AggregateEvent) *EventCollector {
	return &EventCollector{channel: aggregateEventChannel}
}

func (collector *EventCollector) Run() {
	l := NewLap()
	collector.laps = append(collector.laps, &l)

	for {
		event := <-collector.channel
		lap := collector.laps[len(collector.laps)-1]
		lap.AddEvent(event)
		if event.Total_distance_meters > 0 && event.Total_distance_meters%2000 == 0 {
			l2 := NewLap()
			collector.laps = append(collector.laps, &l2)
			jww.DEBUG.Printf("Added %d auto-lap at %d meters", len(collector.laps), event.Total_distance_meters)
			l2.AddEvent(event)
		}
	}
}

func (collector *EventCollector) Laps() []*Lap {
	return collector.laps
}

func (collector *EventCollector) Activity() *Lap {
	if len(collector.laps) == 0 {
		return nil
	}
	lap := NewLap()

	// StartTimeMilliseconds int64
	// StartTimeSeconds      int64
	// StartTimeZulu         string
	// TotalTimeSeconds      int64
	// DistanceMeters        uint64
	// MaximumSpeedMs        float64
	// AverageSpeedMs        float64
	// KCalories             uint64
	// AverageHeartRateBpm   uint64
	// MaximumHeartRateBpm   uint64
	// AverageCadenceRpm     uint64
	// MaximumCadenceRpm     uint64
	// AveragePowerWatts     uint64
	// MaximumPowerWatts     uint64

	first := collector.laps[0]
	last := collector.laps[len(collector.laps)-1]

	lap.StartTimeMilliseconds = first.StartTimeMilliseconds
	lap.StartTimeSeconds = first.StartTimeSeconds
	lap.StartTimeZulu = first.StartTimeZulu
	lap.KCalories = last.KCalories

	for _, l := range collector.laps {
		lap.TotalTimeSeconds += l.TotalTimeSeconds
		lap.DistanceMeters += l.DistanceMeters

		if l.MaximumCadenceRpm > lap.MaximumCadenceRpm {
			lap.MaximumCadenceRpm = l.MaximumCadenceRpm
		}
		if l.MaximumHeartRateBpm > lap.MaximumHeartRateBpm {
			lap.MaximumHeartRateBpm = l.MaximumHeartRateBpm
		}
		if l.MaximumPowerWatts > lap.MaximumPowerWatts {
			lap.MaximumPowerWatts = l.MaximumPowerWatts
		}
		if l.MaximumSpeedMs > lap.MaximumSpeedMs {
			lap.MaximumSpeedMs = l.MaximumSpeedMs
		}

		weight := uint64(l.TotalTimeSeconds)
		lap.AverageCadenceRpm += l.AverageCadenceRpm * weight
		lap.AverageHeartRateBpm += l.AverageHeartRateBpm * weight
		lap.AveragePowerWatts += l.AveragePowerWatts * weight
	}

	lap.AverageCadenceRpm = uint64(float64(lap.AverageCadenceRpm) / float64(lap.TotalTimeSeconds))
	lap.AverageHeartRateBpm = uint64(float64(lap.AverageHeartRateBpm) / float64(lap.TotalTimeSeconds))
	lap.AveragePowerWatts = uint64(float64(lap.AveragePowerWatts) / float64(lap.TotalTimeSeconds))
	lap.AverageSpeedMs = float64(lap.DistanceMeters) / float64(lap.TotalTimeSeconds)
	return &lap
}
