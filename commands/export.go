package commands

import (
	"os"
	"bufio"
	"github.com/olympum/oarsman/s4"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
)

var inputFile string

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export workout data from log",
	Long: `
Exports one or multiple workouts from the log file to TCX (Garmin
Training Center) format`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		exportActivity(inputFile)
	},
}

func exportActivity(fileName string) {
	if inputFile == "" {
		jww.ERROR.Println("Nothing to export")
		return
	}
	jww.INFO.Printf("Exporting %s\n", fileName)

	eventChannel := make(chan s4.AtomicEvent)
	aggregateEventChannel := make(chan s4.AggregateEvent)
	collector := s4.NewEventCollector(aggregateEventChannel)
	go collector.Run()


	s, err := s4.NewReplayS4(eventChannel, aggregateEventChannel, false, fileName, false)
	if err != nil {
		// TODO
		return
	}
	fqOfn := viper.GetString("TempFolder") + string(os.PathSeparator) + "export.log" //TODO
	go s4.Logger(eventChannel, fqOfn)

	s.Run(nil)
	
	exportFileName := viper.GetString("ExportFolder") + string(os.PathSeparator) + "export-yeah.tcx" //TODO

	f, err := os.Create(exportFileName)
	if err != nil {
		jww.FATAL.Printf("Could not create %s\n", exportFileName)
	}
	defer f.Close()

	var w *bufio.Writer
	w = bufio.NewWriter(f)
	jww.INFO.Printf("Writing aggregate data to %s\n", f.Name())

	s4.TCXWriter(collector.Activity(), w)
}

func init() {
	exportCmd.Flags().StringVar(&inputFile, "file", "", "activity log")
}
