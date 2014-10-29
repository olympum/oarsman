package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	jww "github.com/spf13/jwalterweatherman"
	"os"
)

var dbName = "oarsman.db"

func OpenDatabase(workingFolder string) (*sql.DB, error) {
	err := os.Chdir(workingFolder)
	if err != nil {
		jww.ERROR.Println("Error accessing working folder", err)
		return nil, err
	}
	// note that the sqlite drier ensures that the database file exists
	db, e := sql.Open("sqlite3", dbName)
	if e != nil {
		jww.ERROR.Println("Could not open database file", e)
		return nil, e
	}

	return db, nil
}
