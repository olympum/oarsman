package commands

import (
	"github.com/olympum/oarsman/util"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"
	"os"
	"os/user"
)

var CfgFile string
var Verbose bool

var RootCmd = &cobra.Command{
	Use:   "oarsman",
	Short: "Oarsman is a workout management tool for the WaterRower S4",
	Long: `
A log capturing tool connecting to the S4 via USB, will store and
process workout data into a database, and allows exporting as CSV
and TCX (Garmin Training Center).`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			os.Exit(0)
		}
		InitializeConfig()
	},
}

func InitializeConfig() {
	if Verbose {
		jww.SetStdoutThreshold(jww.LevelDebug)
	} else {
		jww.SetStdoutThreshold(jww.LevelInfo)
	}

	if len(CfgFile) > 0 {
		viper.SetConfigFile(CfgFile)
		err := viper.ReadInConfig()
		if err != nil {
			jww.ERROR.Println("Using configuration defaults, config file not found:", CfgFile)
		} else {
			jww.INFO.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		jww.INFO.Println("Using configuration defaults")
	}

	user, error := user.Current()
	if error != nil {
		jww.ERROR.Printf("Unable to fetch current user", error)
	}

	workingFolder := user.HomeDir + string(os.PathSeparator) + ".oarsman"
	SetupFolder(workingFolder, "WorkingFolder", "Working folder:")

	dbFolder := workingFolder + string(os.PathSeparator) + "db"
	SetupFolder(dbFolder, "DbFolder", "Db folder:")

	workoutFolder := workingFolder + string(os.PathSeparator) + "workouts"
	SetupFolder(workoutFolder, "WorkoutFolder", "Workout folder:")

	tempFolder := os.TempDir() + "com.olympum.Oarsman"
	SetupFolder(tempFolder, "TempFolder", "Temp folder:")
}

func SetupFolder(folder string, configName string, logMessage string) {
	viper.SetDefault(configName, folder)
	util.EnsureFolderExists(viper.GetString(configName))
	jww.INFO.Println(logMessage, folder)
}

func Execute() {
	AddCommands()
	RootCmd.Execute()
}

func AddCommands() {
	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(initCmd)
	RootCmd.AddCommand(workoutCmd)
	RootCmd.AddCommand(exportCmd)
	RootCmd.AddCommand(importCmd)
}

func init() {
	RootCmd.PersistentFlags().StringVar(&CfgFile, "config", "", "config file (overrides default config params)")
	RootCmd.PersistentFlags().BoolVar(&Verbose, "verbose", false, "verbose logging")
}
