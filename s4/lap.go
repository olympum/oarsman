package s4

import (
	"github.com/olympum/oarsman/util"
)

type Lap struct {
	events          []AggregateEvent
	sumHeartRateBpm uint64
	sumCadenceRpm   uint64
	sumPowerWatts   uint64

	StartTimeMilliseconds int64
	StartTimeSeconds      int64
	StartTimeZulu         string
	TotalTimeSeconds      int64
	DistanceMeters        uint64
	MaximumSpeedMs        float64
	AverageSpeedMs        float64
	KCalories             uint64
	AverageHeartRateBpm   uint64
	MaximumHeartRateBpm   uint64
	AverageCadenceRpm     uint64
	MaximumCadenceRpm     uint64
	AveragePowerWatts     uint64
	MaximumPowerWatts     uint64
}

func NewLap() Lap {
	return Lap{events: []AggregateEvent{}}
}

func (lap *Lap) AddEvent(event AggregateEvent) {
	if event.Time == 0 {
		return
	}
	if lap.StartTimeMilliseconds == 0 {
		lap.StartTimeMilliseconds = event.Time
	}
	lap.events = append(lap.events, event)

	if event.Speed_m_s > lap.MaximumSpeedMs {
		lap.MaximumSpeedMs = event.Speed_m_s
	}
	if event.Heart_rate > lap.MaximumHeartRateBpm {
		lap.MaximumHeartRateBpm = event.Heart_rate
	}
	if event.Watts > lap.MaximumPowerWatts {
		lap.MaximumPowerWatts = event.Watts
	}
	if event.Stroke_rate > lap.MaximumCadenceRpm {
		lap.MaximumCadenceRpm = event.Stroke_rate
	}
	lap.sumHeartRateBpm += event.Heart_rate
	lap.sumPowerWatts += event.Watts
	lap.sumCadenceRpm += event.Stroke_rate

	lap.calculate()
}

func (lap *Lap) calculate() {
	numSamples := len(lap.events)

	first := lap.events[0]
	last := lap.events[numSamples-1]

	lap.StartTimeSeconds = lap.StartTimeMilliseconds / 1000
	lap.StartTimeZulu = util.MillisToZulu(lap.StartTimeMilliseconds)
	lap.TotalTimeSeconds = (last.Time - first.Time) / 1000
	lap.DistanceMeters = last.Total_distance_meters - first.Total_distance_meters
	lap.AverageSpeedMs = float64(lap.DistanceMeters) / float64(lap.TotalTimeSeconds)
	lap.KCalories = (last.Calories - first.Calories) / 1000
	lap.AverageHeartRateBpm = uint64(float64(lap.sumHeartRateBpm) / float64(numSamples))
	lap.AverageCadenceRpm = uint64(float64(lap.sumCadenceRpm) / float64(numSamples))
	lap.AveragePowerWatts = uint64(float64(lap.sumPowerWatts) / float64(numSamples))

}
