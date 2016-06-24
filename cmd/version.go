package cmd

import (
	"github.com/spf13/cobra"

	"github.com/antham/doc-hunt/ui"
	"github.com/antham/doc-hunt/util"
)

var version = "0.0.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Run: func(cmd *cobra.Command, args []string) {
		ui.Info("v" + version)

		util.SuccessExit()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
