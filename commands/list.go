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
	fmt.Println("id,start_time,distance,ave_speed,max_speed")
	for _, activity := range activities {
		fmt.Printf("%d,%s,%d,%f,%f\n",
			activity.StartTimeMilliseconds,
			activity.StartTimeZulu(),
			activity.DistanceMeters,
			activity.AverageSpeed(),
			activity.MaximumSpeed())
	}
	return

}
