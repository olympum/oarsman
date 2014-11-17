package db

import (
	"database/sql"
	"github.com/olympum/oarsman/s4"
	jww "github.com/spf13/jwalterweatherman"
)

var fields = `
start_time_milliseconds,
start_time_seconds,
start_time_zulu,
parent_start_time_milliseconds,
total_time_seconds,
distance_meters,
maximum_speed_m_s,
average_speed_m_s,
kcalories,
average_heart_rate_bpm,
maximum_heart_rate_bpm,
average_cadence_rpm,
maximum_cadence_rpm,
average_power_watts,
maximum_power_watts
`

var insertString = `

INSERT INTO activity
(` + fields +
	`)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)


`

var deleteString = `

DELETE FROM activity
WHERE start_time_milliseconds = ?

`

var selectAllActivitiesString = `

SELECT` + fields + `
FROM activity
WHERE parent_start_time_milliseconds = -1

`

var selectActivityString = `

SELECT` + fields + `
FROM activity
WHERE parent_start_time_milliseconds = -1
AND start_time_milliseconds = ?

`
var selectAllLapsString = `

SELECT` + fields + `
FROM activity
WHERE parent_start_time_milliseconds = ?

`

var createTableString = `

CREATE TABLE activity (
start_time_milliseconds INTEGER,
start_time_seconds INTEGER,
start_time_zulu VARCHAR,
parent_start_time_milliseconds INTEGER,
total_time_seconds INTEGER,
distance_meters INTEGER,
maximum_speed_m_s REAL,
average_speed_m_s REAL,
kcalories INTEGER,
average_heart_rate_bpm INTEGER,
maximum_heart_rate_bpm INTEGER,
average_cadence_rpm INTEGER,
maximum_cadence_rpm INTEGER,
average_power_watts INTEGER,
maximum_power_watts INTEGER
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

func (db *OarsmanDB) ListActivities() []*s4.Activity {
	jww.DEBUG.Println(selectAllActivitiesString)
	rows, err := db.odb.Query(selectAllActivitiesString)
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}
	laps := parseLaps(rows)

	activities := []*s4.Activity{}
	if len(laps) == 0 {
		jww.DEBUG.Println("No activities found")
		return nil
	}
	for _, lap := range laps {
		activity := s4.NewActivity(lap, nil)
		activities = append(activities, activity)
		jww.DEBUG.Println("Converted lap into activity", activity)
	}
	return activities
}

func (db *OarsmanDB) FindActivityById(id int64) *s4.Activity {
	jww.DEBUG.Printf("Looking for activity %d", id)
	rows, err := db.odb.Query(selectActivityString, id)
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}
	laps := parseLaps(rows)
	if len(laps) > 0 {
		jww.DEBUG.Printf("Activity %d found", id)
		return s4.NewActivity(laps[0], nil)
	} else {
		jww.DEBUG.Printf("Activity %d not found", id)
		return nil
	}
}

func (db *OarsmanDB) FindLapsByParentId(id int64) []*s4.Lap {
	jww.DEBUG.Printf("Looking for laps for activity %d", id)
	rows, err := db.odb.Query(selectAllLapsString, id)
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}
	laps := parseLaps(rows)
	if len(laps) > 0 {
		jww.DEBUG.Printf("Laps for activity %d found", id)
		return laps
	} else {
		jww.DEBUG.Printf("Laps for activity %d not found", id)
		return nil
	}
}

func (db *OarsmanDB) RemoveActivityById(id int64) *s4.Activity {
	jww.DEBUG.Printf("Removing activity %d", id)
	activity := db.FindActivityById(id)
	if activity != nil {
		_, error := db.odb.Exec(deleteString, id)
		if error != nil {
			jww.ERROR.Println(error)
		} else {
			jww.INFO.Printf("Activity %d deleted", activity.StartTimeMilliseconds)
		}
	}
	return activity
}

func parseLaps(rows *sql.Rows) []*s4.Lap {
	laps := []*s4.Lap{}
	for rows.Next() {

		lap := s4.NewLap()
		var id int64

		rows.Scan(&lap.StartTimeMilliseconds,
			&lap.StartTimeSeconds,
			&lap.StartTimeZulu,
			&id,
			&lap.TotalTimeSeconds,
			&lap.DistanceMeters,
			&lap.MaximumSpeedMs,
			&lap.AverageSpeedMs,
			&lap.KCalories,
			&lap.AverageHeartRateBpm,
			&lap.MaximumHeartRateBpm,
			&lap.AverageCadenceRpm,
			&lap.MaximumCadenceRpm,
			&lap.AveragePowerWatts,
			&lap.MaximumPowerWatts,
		)

		jww.DEBUG.Printf("Parsed lap with %v start time, parent id %v: %v", lap.StartTimeMilliseconds, id, lap)

		laps = append(laps, &lap)
	}
	jww.DEBUG.Println("Laps parsed", len(laps))
	return laps
}

func (db *OarsmanDB) InsertActivity(activity *s4.Activity) *s4.Activity {

	if db.FindActivityById(activity.StartTimeMilliseconds) != nil {
		jww.ERROR.Printf("Activity already exists in database, ignoring %d\n", activity.StartTimeMilliseconds)
		return nil
	}

	result, err := db.odb.Exec(insertString,
		activity.StartTimeMilliseconds,
		activity.StartTimeSeconds,
		activity.StartTimeZulu,
		-1,
		activity.TotalTimeSeconds,
		activity.DistanceMeters,
		activity.MaximumSpeedMs,
		activity.AverageSpeedMs,
		activity.KCalories,
		activity.AverageHeartRateBpm,
		activity.MaximumHeartRateBpm,
		activity.AverageCadenceRpm,
		activity.MaximumCadenceRpm,
		activity.AveragePowerWatts,
		activity.MaximumPowerWatts,
	)
	if err != nil {
		jww.ERROR.Printf("Could not insert activity with id %v into database: %v", activity.StartTimeMilliseconds, err)
		return nil
	} else {
		for _, lap := range activity.Laps() {
			result, err := db.odb.Exec(insertString,
				lap.StartTimeMilliseconds,
				lap.StartTimeSeconds,
				lap.StartTimeZulu,
				activity.StartTimeMilliseconds,
				lap.TotalTimeSeconds,
				lap.DistanceMeters,
				lap.MaximumSpeedMs,
				lap.AverageSpeedMs,
				lap.KCalories,
				lap.AverageHeartRateBpm,
				lap.MaximumHeartRateBpm,
				lap.AverageCadenceRpm,
				lap.MaximumCadenceRpm,
				lap.AveragePowerWatts,
				lap.MaximumPowerWatts,
			)
			if err != nil {
				jww.ERROR.Println("Could not insert lap in the database", err)
			}
			jww.DEBUG.Println("Inserted lap", lap, result)
		}
	}

	jww.DEBUG.Println("Inserted activity", activity, result)

	return activity
}
