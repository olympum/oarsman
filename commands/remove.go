package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an activity from the database",
	Long: `
Deletes an activity from the data. The log file(s) remain
intact in the workouts folder and temp folder.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		removeActivity(activityId)
	},
}

func removeActivity(activityId int64) {
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	if activityId == -1 {
		return
	}

	activity := database.RemoveActivityById(activityId)

	if activity != nil {
		fmt.Printf("Row %d deleted: %s,%d,%f,%f\n",
			activity.StartTimeMilliseconds,
			activity.StartTimeZulu,
			activity.DistanceMeters,
			activity.AverageSpeedMs,
			activity.MaximumSpeedMs)
	} else {
		jww.ERROR.Printf("Activity %d not found", activityId)
	}
}

func init() {
	removeCmd.Flags().Int64Var(&activityId, "id", -1, "id of activity to remove")
}
