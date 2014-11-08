package commands

import (
	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"os"
)

var format string

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
		exportActivity(activityId)
	},
}

func exportActivity(activityId int64) {
	database, error := workoutDatabase()
	if error != nil {
		// TODO
		return
	}
	defer database.Close()

	if activityId == 0 {
		return
	}

	activity := database.FindActivityById(activityId)
	if activity == nil {
		jww.ERROR.Printf("Activity %d not found\n", activityId)
		return
	}

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

	prefix := viper.GetString("TempFolder") + string(os.PathSeparator) + fileName
	if format == "TCX" {
		s4.ExportCollectorEvents(collector, prefix+".tcx", s4.TCXWriter)
	} else if format == "CSV" {
		s4.ExportCollectorEvents(collector, prefix+".csv", s4.CSVWriter)
	} else {
		jww.ERROR.Printf("Unknow export file format %s\n", format)
	}
}

func init() {
	exportCmd.Flags().Int64Var(&activityId, "id", 0, "id of activity to export")
	exportCmd.Flags().StringVar(&format, "format", "TCX", "format to export activity as, TCX or CSV")
}
