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

var workoutCmd = &cobra.Command{
	Use:   "workout",
	Short: "Start a rowing workout",
	Long: `
Send workout instructions to rowing monitor and start collecting
rowing event data till workout is completed. Data is not saved in
the database (use the import command to save it in the database).`,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		eventChannel := make(chan s4.AtomicEvent)
		aggregateEventChannel := make(chan s4.AggregateEvent)
		collector := s4.NewEventCollector(aggregateEventChannel)
		go collector.Run()

		stamp := util.MillisToZulu(time.Now().UTC().Unix())
		tempFile := viper.GetString("TempFolder") + string(os.PathSeparator) + stamp + ".log"
		go s4.Logger(eventChannel, tempFile)
		workout := s4.NewS4Workout()
		workout.AddSingleWorkout(duration, distance)
		s := s4.NewS4(eventChannel, aggregateEventChannel, debug)

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

		s.Run(&workout)

	},
}

func init() {
	workoutCmd.Flags().BoolVar(&debug, "debug", false, "debug communication data packets")
	workoutCmd.Flags().Uint64Var(&distance, "distance", 2000, "distance of workout (in meters)")
	workoutCmd.Flags().DurationVar(&duration, "duration", 0, "duration of workout (e.g. 1800s or 45m)")
}
