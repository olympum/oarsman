package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var VERSION = "v0.1"
var DEV = "201411-SNAPSHOT"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long: `
The version number for this SmarterRower build`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SmarterRower for WaterRower", VERSION, "-- Revision ", DEV)
	},
}
