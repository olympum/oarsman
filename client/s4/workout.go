package s4

import (
	"container/list"
	"fmt"
	"log"
	"time"
)

type S4Workout struct {
	workoutPackets *list.List
	state          int
}

func NewS4Workout() S4Workout {
	workout := S4Workout{workoutPackets: list.New(), state: Unset}
	return workout
}

func (workout S4Workout) AddSingleWorkout(duration time.Duration, distanceMeters uint64) {
	// prepare workout instructions
	durationSeconds := uint64(duration.Seconds())
	var workoutPacket Packet

	if durationSeconds > 0 {
		log.Printf("Starting single duration workout: %d seconds", durationSeconds)
		if durationSeconds >= 18000 {
			log.Fatalf("Workout time must be less than 18,000 seconds (was %d)", durationSeconds)
		}
		payload := fmt.Sprintf("%04X", durationSeconds)
		workoutPacket = Packet{cmd: WorkoutSetDurationRequest, data: []byte(payload)}
	} else if distanceMeters > 0 {
		log.Printf("Starting single distance workout: %d meters", distanceMeters)
		if distanceMeters >= 64000 {
			log.Fatalf("Workout distance must be less than 64,000 meters (was %d)", distanceMeters)
		}
		payload := Meters + fmt.Sprintf("%04X", distanceMeters)
		workoutPacket = Packet{cmd: WorkoutSetDistanceRequest, data: []byte(payload)}
	} else {
		log.Fatal("Undefined workout")
	}
	workout.workoutPackets.PushFront(workoutPacket)
}
