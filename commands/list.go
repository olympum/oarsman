package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all workout activities in the database",
	Long: `
Lists all the activities stored in the database`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		listActivities()
	},
}

func listActivities() {
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	activities := database.ListActivities()
	fmt.Println("id,start_time,distance,duration,ave_speed,max_speed,ave_cadence,max_cadence,ave_power,max_power,calories,ave_hr,max_hr")
	for _, activity := range activities {
		fmt.Printf("%d,%s,%d,%d,%v,%v,%v,%v,%v,%v,%v,%v,%v\n",
			activity.StartTimeMilliseconds,
			activity.StartTimeZulu(),
			activity.DistanceMeters,
			activity.TotalTimeSeconds(),
			activity.AverageSpeed(),
			activity.MaximumSpeed(),
			activity.AverageCadence,
			activity.MaximumCadence,
			activity.AveragePower,
			activity.MaximumPower,
			activity.KCalories(),
			activity.AverageHeartRateBpm,
			activity.MaximumHeartRateBpm)
	}
	return

}
