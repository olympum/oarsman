package s4

import (
	"time"
)

type Activity struct {
	numSamples            uint64
	StartTimeMilliseconds int64
	TotalTimeMilliseconds int64
	DistanceMeters        uint64
	MaximumSpeed_cm_s     uint64
	Calories              uint64
	AverageHeartRateBpm   float64
	MaximumHeartRateBpm   uint64
	AverageCadence        float64
	MaximumCadence        uint64
	AveragePower          float64
	MaximumPower          uint64
}

func millisToZulu(millis int64) string {
	return time.Unix(millis/1000, millis%1000*1000).UTC().Format(time.RFC3339)
}

func (activity Activity) StartTimeSeconds() int64 {
	return activity.StartTimeMilliseconds / 1000
}

func (activity Activity) StartTimeZulu() string {
	return millisToZulu(activity.StartTimeMilliseconds)
}

func (Activity Activity) TotalTimeSeconds() int64 {
	return Activity.TotalTimeMilliseconds / 1000
}

func (activity Activity) AverageSpeed() float64 {
	return float64(activity.DistanceMeters) / float64(activity.TotalTimeMilliseconds/1000)
}

func (activity Activity) MaximumSpeed() float64 {
	return float64(activity.MaximumSpeed_cm_s) / 100
}

func (activity Activity) KCalories() uint64 {
	return activity.Calories / 1000
}
