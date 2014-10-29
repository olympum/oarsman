package db

import (
	"database/sql"
	"github.com/olympum/oarsman/s4"
	jww "github.com/spf13/jwalterweatherman"
)

var insertString = `

INSERT INTO activity (
start_time_milliseconds,
total_time_milliseconds,
distance_meters,
maximum_speed_cm_s,
calories,
average_heart_rate,
maximum_heart_rate,
average_power,
maximum_power,
average_cadence,
maximum_cadence)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)


`

var queryString = `

SELECT
start_time_milliseconds,
total_time_milliseconds,
distance_meters,
maximum_speed_cm_s,
calories,
average_heart_rate,
maximum_heart_rate,
average_power,
maximum_power,
average_cadence,
maximum_cadence
FROM activity

`

var selectString = `

SELECT
start_time_milliseconds,
total_time_milliseconds,
distance_meters,
maximum_speed_cm_s,
calories,
average_heart_rate,
maximum_heart_rate,
average_power,
maximum_power,
average_cadence,
maximum_cadence
FROM activity
WHERE start_time_milliseconds = ?

`

var createTableString = `

CREATE TABLE activity (
       start_time_milliseconds INTEGER PRIMARY KEY,
       total_time_milliseconds INTEGER,
       distance_meters REAL,
       maximum_speed_cm_s REAL,
       calories INTEGER,
       average_heart_rate INTEGER,
       maximum_heart_rate INTEGER,
       average_power INTEGER,
       maximum_power INTEGER,
       average_cadence INTEGER,
       maximum_cadence INTEGER
);

`

func CreateTables(db *sql.DB) error {
	_, err := db.Exec(createTableString)
	if err != nil {
		jww.ERROR.Printf("%q: %s\n", err, createTableString)
		return nil
	}

	jww.INFO.Println("Created table schema")

	return nil
}

func InitializeDatabase(db *sql.DB) error {
	defer db.Close()

	e := CreateTables(db)
	if e != nil {
		return e
	}
	return nil
}

func ListActivities(db *sql.DB) []s4.Activity {
	rows, err := db.Query(queryString)
	if err != nil {
		jww.ERROR.Println(err)
		var empty []s4.Activity
		return empty
	}
	return parseActivities(rows)
}

func FindActivityById(db *sql.DB, id int64) *s4.Activity {
	rows, err := db.Query(selectString, id)
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}
	activities := parseActivities(rows)
	if len(activities) > 0 {
		return &activities[0]
	} else {
		return nil
	}
}

func parseActivities(rows *sql.Rows) []s4.Activity {
	var activities []s4.Activity
	for rows.Next() {
		var start_time_milliseconds int64
		var total_time_milliseconds int64
		var distance_meters uint64
		var maximum_speed_cm_s uint64
		var calories uint64
		var average_heart_rate_bpm float64
		var maximum_heart_rate_bpm uint64
		var average_power float64
		var maximum_power uint64
		var average_cadence float64
		var maxmimum_cadence uint64

		rows.Scan(&start_time_milliseconds,
			&total_time_milliseconds,
			&distance_meters,
			&maximum_speed_cm_s,
			&calories,
			&average_heart_rate_bpm,
			&maximum_heart_rate_bpm,
			&average_power,
			&maximum_power,
			&average_cadence,
			&maxmimum_cadence)
		activity := s4.Activity{
			StartTimeMilliseconds: start_time_milliseconds,
			TotalTimeMilliseconds: total_time_milliseconds,
			DistanceMeters:        distance_meters,
			MaximumSpeed_cm_s:     maximum_speed_cm_s,
			Calories:              calories,
			AverageHeartRateBpm:   average_heart_rate_bpm,
			MaximumHeartRateBpm:   maximum_heart_rate_bpm,
			AveragePower:          average_power,
			MaximumPower:          maximum_power,
			AverageCadence:        average_cadence,
			MaximumCadence:        maxmimum_cadence,
		}
		activities = append(activities, activity)
	}
	return activities
}

func InsertActivity(db *sql.DB, activity *s4.Activity) {

	if FindActivityById(db, activity.StartTimeMilliseconds) != nil {
		jww.ERROR.Printf("Activity already exists in database, ignoring %d\n", activity.StartTimeMilliseconds)
		return
	}

	jww.INFO.Println(activity)
	result, err := db.Exec(insertString,
		activity.StartTimeMilliseconds,
		activity.TotalTimeMilliseconds,
		activity.DistanceMeters,
		activity.MaximumSpeed_cm_s,
		activity.Calories,
		activity.AverageHeartRateBpm,
		activity.MaximumHeartRateBpm,
		activity.AveragePower,
		activity.MaximumPower,
		activity.AverageCadence,
		activity.MaximumCadence,
	)
	if err != nil {
		jww.ERROR.Print(err)
	}
	jww.INFO.Print(result)

	rows, _ := db.Query(queryString)
	for rows.Next() {
		var start_time int
		rows.Scan(&start_time)
		jww.INFO.Print(start_time)
	}
}