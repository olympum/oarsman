package commands

import (
	"github.com/olympum/oarsman/s4"
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"time"
)

var distance uint64
var duration time.Duration
var debug bool

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "Start a rowing workout activity",
	Long: `
Send workout instructions to rowing monitor and start collecting
rowing event data till workout is completed. Data is not saved in
the database (use the import command to save it in the database).`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		eventChannel := make(chan s4.AtomicEvent)

		stamp := util.MillisToZulu(time.Now().UnixNano() / 1000000)
		tempFile := viper.GetString("TempFolder") + string(os.PathSeparator) + stamp + ".log"
		go s4.Logger(eventChannel, tempFile)
		workout := s4.NewS4Workout()
		workout.AddSingleWorkout(duration, distance)
		s := s4.NewS4(eventChannel, nil, debug)

		// TODO we should detect a workout completition, not use OS signals
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, os.Kill)
		go func() {
			for sig := range ch {
				jww.INFO.Printf("Terminating workout (received %s signal)\n", sig.String())
				s.Exit()
				os.Exit(0)
			}
		}()

		go s.Run(&workout)

		jww.INFO.Println(">>> Press RETURN to end workout ... <<<")
		var buffer [1]byte
		os.Stdin.Read(buffer[:])

		s.Exit()

		jww.INFO.Println("Workout completed successfully")

		activity := importActivity(tempFile, false)

		if activity != nil {
			exportActivity(activity.StartTimeMilliseconds)
		}
	},
}

func init() {
	trainCmd.Flags().BoolVar(&debug, "debug", false, "debug communication data packets")
	trainCmd.Flags().Uint64Var(&distance, "distance", 2000, "distance of workout (in meters)")
	trainCmd.Flags().DurationVar(&duration, "duration", 0, "duration of workout (e.g. 1800s or 45m)")
}
