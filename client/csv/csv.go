package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AggregateEvent struct {
	time                  int64
	total_distance_meters int64
	stroke_rate           int64
	watts                 int64
	calories              int64
	speed_cm_s            int64
	heart_rate            int64
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("time,total_distance_meters,stroke_rate,watts,calories,speed_cm_s,heart_rate")
	var reftime int64
	var event = AggregateEvent{}
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Split(line, " ")
		time, _ := strconv.ParseInt(tokens[0], 10, 64)
		values := strings.Split(tokens[1], ":")
		label := values[0]
		value, _ := strconv.ParseInt(values[1], 10, 64)

		if reftime == 0 {
			reftime = time - time%100
			event.time = reftime + 100
		}
		switch label {
		case "total_distance_meters":
			event.total_distance_meters = value
		case "stroke_rate":
			event.stroke_rate = value
		case "watts":
			if value > 0 {
				event.watts = value
			}
		case "calories":
			event.calories = value
		case "speed_cm_s":
			event.speed_cm_s = value
		case "heart_rate":
			event.heart_rate = value
		}

		if time-reftime >= 100 {
			fmt.Printf("%d,%d,%d,%d,%d,%d,%d\n",
				event.time,
				event.total_distance_meters,
				event.stroke_rate,
				event.watts,
				event.calories,
				event.speed_cm_s,
				event.heart_rate)
			reftime = 0
		}
	}
}
