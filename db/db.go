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

var deleteString = `

DELETE FROM activity
WHERE start_time_milliseconds = ?

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

type OarsmanDB struct {
	odb *sql.DB
}

func (db *OarsmanDB) Close() error {
	return db.odb.Close()
}

func (db *OarsmanDB) CreateTables() error {
	_, err := db.odb.Exec(createTableString)
	if err != nil {
		jww.ERROR.Printf("%q: %s\n", err, createTableString)
		return nil
	}

	jww.INFO.Println("Created table schema")

	return nil
}

func (db *OarsmanDB) InitializeDatabase() {
	// TODO database migrations
	q := `SELECT name FROM sqlite_master WHERE type='table' AND name='activity'`
	var name string
	err := db.odb.QueryRow(q).Scan(&name)
	switch {
	case err == sql.ErrNoRows:
		jww.INFO.Println("Initializing database for the first time ...")
		e := db.CreateTables()
		if e != nil {
			jww.ERROR.Println(e)
		}
	case err != nil:
		jww.ERROR.Println(err)
	default:
		jww.DEBUG.Println("Activity table alreay exists in database")
	}
}

func (db *OarsmanDB) ListActivities() []s4.Activity {
	rows, err := db.odb.Query(queryString)
	if err != nil {
		jww.ERROR.Println(err)
		var empty []s4.Activity
		return empty
	}
	return parseActivities(rows)
}

func (db *OarsmanDB) FindActivityById(id int64) *s4.Activity {
	rows, err := db.odb.Query(selectString, id)
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

func (db *OarsmanDB) RemoveActivityById(id int64) *s4.Activity {
	activity := db.FindActivityById(id)
	if activity != nil {
		_, error := db.odb.Exec(deleteString, id)
		if error != nil {
			jww.ERROR.Println(error)
		} else {
			jww.INFO.Printf("Rows deleted %d", activity.StartTimeMilliseconds)
		}
	}
	return activity
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

func (db *OarsmanDB) InsertActivity(activity *s4.Activity) {

	if db.FindActivityById(activity.StartTimeMilliseconds) != nil {
		jww.ERROR.Printf("Activity already exists in database, ignoring %d\n", activity.StartTimeMilliseconds)
		return
	}

	jww.DEBUG.Println(activity)
	result, err := db.odb.Exec(insertString,
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
	jww.DEBUG.Print(result)

}
