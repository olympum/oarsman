package commands

import (
	"os"
	"time"

	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen and log the workout until Keypad OK long press is detected.",
	Long: `
	Listen and log the workout until Keypad OK long press is detected.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		eventChannel := make(chan s4.AtomicEvent)
		// todo this belongs to logger
		stamp := util.MillisToZulu(time.Now().UnixNano() / 1000000)
		tempFile := viper.GetString("TempFolder") + string(os.PathSeparator) + stamp + ".log"

		go s4.Logger(eventChannel, tempFile)
		workout := s4.NewS4Workout()
		s := s4.ListenS4(eventChannel, nil, debug)

		go util.Ready("S4 Found. Let's row.. ")

		s.Run(&workout)

		activity := importActivity(tempFile, false)

		if activity != nil {
			exportActivity(activity.StartTimeMilliseconds)
		}
		os.Exit(0)
	},
}
