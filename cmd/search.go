package cmd

import (
	"github.com/geowa4/servicelogger/pkg/labels"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for a service log",
	Long:  `Run an interactive TUI to search and discover service log templates`,
	Run: func(cmd *cobra.Command, args []string) {
		labels.SearchProgram()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
