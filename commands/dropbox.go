package commands

import (
	"github.com/spf13/cobra"
)

var token string

var dropboxCmd = &cobra.Command{
	Use:   "dropbox",
	Short: "Share workout via Dropbox",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		InitializeConfig()
		dropbox(fileName)
	},
}

func dropbox(fileName string) {

}

func init() {
	dropboxCmd.Flags().StringVar(&fileName, "file", "", "tcx file to upload")
	dropboxCmd.Flags().StringVar(&token, "token", "", "Dropbox auth token")
}
