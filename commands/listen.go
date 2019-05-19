package commands

import (
	"os"

	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen and log the workout until Keypad three OK press is detected.",
	Long: `
	Listen and log the workout until Keypad three OK press is detected.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		eventChannel := make(chan s4.AtomicEvent)
		tempFile := util.NewTempFilename()

		go s4.Logger(eventChannel, tempFile)
		workout := s4.NewS4Workout()
		s := s4.ListenS4(eventChannel, nil, false)

		go util.Ready("S4 Found. Let's row.. ")

		s.Run(&workout)

		exportActivity(tempFile)
		
		os.Exit(0)
	},
}
