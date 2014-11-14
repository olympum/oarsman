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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)


`

var deleteString = `

DELETE FROM activity
WHERE start_time_milliseconds = ?

`

var queryString = `

SELECT` + fields + `FROM activity

`

var selectString = queryString + `

WHERE start_time_milliseconds = ?

`

var createTableString = `

CREATE TABLE activity (
start_time_milliseconds INTEGER PRIMARY KEY,
start_time_seconds INTEGER,
start_time_zulu VARCHAR,
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

func (db *OarsmanDB) ListActivities() []s4.Lap {
	jww.DEBUG.Println(queryString)
	rows, err := db.odb.Query(queryString)
	if err != nil {
		jww.ERROR.Println(err)
		var empty []s4.Lap
		return empty
	}
	return parseActivities(rows)
}

func (db *OarsmanDB) FindActivityById(id int64) *s4.Lap {
	jww.DEBUG.Printf("Looking for activity %d", id)
	rows, err := db.odb.Query(selectString, id)
	if err != nil {
		jww.ERROR.Println(err)
		return nil
	}
	activities := parseActivities(rows)
	if len(activities) > 0 {
		jww.DEBUG.Printf("Activity %d found", id)
		return &activities[0]
	} else {
		jww.DEBUG.Printf("Activity %d not found", id)
		return nil
	}
}

func (db *OarsmanDB) RemoveActivityById(id int64) *s4.Lap {
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

func parseActivities(rows *sql.Rows) []s4.Lap {
	var activities []s4.Lap
	for rows.Next() {

		lap := s4.NewLap()

		rows.Scan(&lap.StartTimeMilliseconds,
			&lap.StartTimeSeconds,
			&lap.StartTimeZulu,
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

		jww.DEBUG.Printf("Parsed activity with %v start time", lap.StartTimeMilliseconds)

		activities = append(activities, lap)
	}
	return activities
}

func (db *OarsmanDB) InsertActivity(activity *s4.Lap) *s4.Lap {

	if db.FindActivityById(activity.StartTimeMilliseconds) != nil {
		jww.ERROR.Printf("Activity already exists in database, ignoring %d\n", activity.StartTimeMilliseconds)
		return nil
	}

	result, err := db.odb.Exec(insertString,
		activity.StartTimeMilliseconds,
		activity.StartTimeSeconds,
		activity.StartTimeZulu,
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
	}
	jww.DEBUG.Print(result)

	return activity
}
