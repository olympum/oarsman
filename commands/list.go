package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workout activities in the database",
	Long: `
Lists all the activities stored in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		if activityId > 0 {
			listLaps(activityId)
		} else {
			listActivities()
		}
	},
}

func listLaps(activityId int64) {
	jww.DEBUG.Println("Looking for laps for activity", activityId)
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	laps := database.FindLapsByParentId(activityId)
	if laps == nil {
		return
	}

	fmt.Println("id,start_time,distance,duration,ave_speed,max_speed,ave_cadence,max_cadence,ave_power,max_power,calories,ave_hr,max_hr")
	for _, lap := range laps {
		fmt.Printf("%d,%s,%d,%d,%.2f,%.2f,%v,%v,%v,%v,%v,%v,%v\n",
			lap.StartTimeMilliseconds,
			lap.StartTimeZulu,
			lap.DistanceMeters,
			lap.TotalTimeSeconds,
			lap.AverageSpeedMs,
			lap.MaximumSpeedMs,
			lap.AverageCadenceRpm,
			lap.MaximumCadenceRpm,
			lap.AveragePowerWatts,
			lap.MaximumPowerWatts,
			lap.KCalories,
			lap.AverageHeartRateBpm,
			lap.MaximumHeartRateBpm)
	}
	return
}

func listActivities() {
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	activities := database.ListActivities()
	if len(activities) == 0 {
		jww.INFO.Println("No activities found")
		return
	}
	fmt.Println("id,start_time,distance,duration,ave_speed,max_speed,ave_cadence,max_cadence,ave_power,max_power,calories,ave_hr,max_hr")
	for _, activity := range activities {
		fmt.Printf("%d,%s,%d,%d,%.2f,%.2f,%v,%v,%v,%v,%v,%v,%v\n",
			activity.StartTimeMilliseconds,
			activity.StartTimeZulu,
			activity.DistanceMeters,
			activity.TotalTimeSeconds,
			activity.AverageSpeedMs,
			activity.MaximumSpeedMs,
			activity.AverageCadenceRpm,
			activity.MaximumCadenceRpm,
			activity.AveragePowerWatts,
			activity.MaximumPowerWatts,
			activity.KCalories,
			activity.AverageHeartRateBpm,
			activity.MaximumHeartRateBpm)
	}
	return

}

func init() {
	listCmd.Flags().Int64Var(&activityId, "id", -1, "id of activity to export")
}
