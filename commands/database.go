package commands

import (
	"github.com/olympum/oarsman/db"
	"github.com/spf13/viper"
)

func workoutDatabase() (*db.OarsmanDB, error) {
	workingFolder := viper.GetString("DbFolder")
	database, e := db.OpenDatabase(workingFolder)
	database.InitializeDatabase()
	return database, e
}
