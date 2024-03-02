package cmd

import (
	"github.com/geowa4/servicelogger/pkg/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Read and set the global config file",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		// TODO
	},
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(config.Program())
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
