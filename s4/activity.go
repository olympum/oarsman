package s4

type Activity struct {
	Lap
	laps []*Lap
}

func NewActivity(lap *Lap, laps []*Lap) *Activity {
	var _laps []*Lap
	if laps != nil {
		_laps = laps
	} else {
		_laps = []*Lap{}
	}

	var activity *Activity
	if lap == nil {
		activity = &Activity{laps: _laps}
	} else {
		activity = &Activity{Lap: *lap, laps: _laps}
	}
	return activity.update()
}

func (activity *Activity) Laps() []*Lap {
	return activity.laps
}

func (activity *Activity) addLap() *Lap {
	lap := NewLap()
	activity.laps = append(activity.laps, &lap)
	return &lap
}

func (activity *Activity) firstLap() *Lap {
	if len(activity.laps) == 0 {
		return nil
	}

	return activity.laps[0]
}

func (activity *Activity) lastLap() *Lap {
	if len(activity.laps) == 0 {
		return nil
	}

	return activity.laps[len(activity.laps)-1]
}

func (activity *Activity) update() *Activity {
	first := activity.firstLap()
	last := activity.lastLap()

	if first == nil || last == nil {
		return activity
	}

	activity.StartTimeMilliseconds = first.StartTimeMilliseconds
	activity.StartTimeSeconds = first.StartTimeSeconds
	activity.StartTimeZulu = first.StartTimeZulu
	activity.KCalories = last.KCalories

	for _, l := range activity.laps {
		activity.TotalTimeSeconds += l.TotalTimeSeconds
		activity.DistanceMeters += l.DistanceMeters

		if l.MaximumCadenceRpm > activity.MaximumCadenceRpm {
			activity.MaximumCadenceRpm = l.MaximumCadenceRpm
		}
		if l.MaximumHeartRateBpm > activity.MaximumHeartRateBpm {
			activity.MaximumHeartRateBpm = l.MaximumHeartRateBpm
		}
		if l.MaximumPowerWatts > activity.MaximumPowerWatts {
			activity.MaximumPowerWatts = l.MaximumPowerWatts
		}
		if l.MaximumSpeedMs > activity.MaximumSpeedMs {
			activity.MaximumSpeedMs = l.MaximumSpeedMs
		}

		weight := uint64(l.TotalTimeSeconds)
		activity.AverageCadenceRpm += l.AverageCadenceRpm * weight
		activity.AverageHeartRateBpm += l.AverageHeartRateBpm * weight
		activity.AveragePowerWatts += l.AveragePowerWatts * weight
	}

	activity.AverageCadenceRpm = uint64(float64(activity.AverageCadenceRpm) / float64(activity.TotalTimeSeconds))
	activity.AverageHeartRateBpm = uint64(float64(activity.AverageHeartRateBpm) / float64(activity.TotalTimeSeconds))
	activity.AveragePowerWatts = uint64(float64(activity.AveragePowerWatts) / float64(activity.TotalTimeSeconds))
	activity.AverageSpeedMs = float64(activity.DistanceMeters) / float64(activity.TotalTimeSeconds)

	return activity

}
