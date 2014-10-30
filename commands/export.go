package commands

import (
	"fmt"
	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var activityId int64

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export workout data from database",
	Long: `
Exports one or multiple workouts from the database
as RAW (40Hz JSON formatted feed), CSV or TCX (Garmin
Training Center). CSV and TCX files are aggregated at 10Hz
before export.`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		database, error := WorkoutDatabase()
		if error != nil {
			// TODO
			return
		}
		defer database.Close()

		if activityId == 0 {
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

		activity := database.FindActivityById(activityId)
		eventChannel := make(chan s4.AtomicEvent)
		aggregateEventChannel := make(chan s4.AggregateEvent)
		collector := s4.NewEventCollector(aggregateEventChannel)
		go collector.Run()

		fileName := util.MillisToZulu(activity.StartTimeMilliseconds)
		inputFile := viper.GetString("WorkoutFolder") + string(os.PathSeparator) + fileName + ".log"
		s, err := s4.NewReplayS4(eventChannel, aggregateEventChannel, false, inputFile, false)
		if err != nil {
			// TODO
			return
		}
		fqOfn := viper.GetString("TempFolder") + string(os.PathSeparator) + randomId() + ".log"
		go s4.Logger(eventChannel, fqOfn)

		s.Run(nil)

		tcx := viper.GetString("TempFolder") + string(os.PathSeparator) + fileName + ".tcx"
		s4.ExportCollectorEvents(collector, tcx, s4.TCXWriter)
	},
}

func init() {
	exportCmd.Flags().Int64Var(&activityId, "id", 0, "id of activity to export")
}
