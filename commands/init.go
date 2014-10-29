package commands

import (
	"database/sql"
	"github.com/olympum/oarsman/db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the database",
	Long: `
Creates or upgrades a database; if the database already exists, it
will be upgraded to the latest schema, no data will be destroyed`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		database, e := WorkoutDatabase()
		if e != nil {
			// TODO
			return
		}

		defer database.Close()

		db.InitializeDatabase(database)
	},
}

func WorkoutDatabase() (*sql.DB, error) {
	workingFolder := viper.GetString("DbFolder")
	database, e := db.OpenDatabase(workingFolder)
	return database, e
}
