package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/labels"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "find-by-tag",
	Short: "Find service log by tag",
	Long:  `Inspects each service log template for _tags and provides navigation`,
	//PreRun: func(cmd *cobra.Command, args []string) {
	//	bindViper(cmd)
	//},
	Run: func(cmd *cobra.Command, args []string) {
		for k, v := range labels.FindFilesWithTags() {
			log.Info("tagMap", "k", k, "v", len(v))
		}
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
