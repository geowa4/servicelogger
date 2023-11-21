package cmd

import (
	"encoding/csv"
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/labels"
	"github.com/spf13/cobra"
	"os"
)

var tagCmd = &cobra.Command{
	Use:   "find-by-tag",
	Short: "Find service log by tag",
	Long:  `Inspects each service log template for _tags and provides navigation`,
	Run: func(cmd *cobra.Command, args []string) {
		csvWriter := csv.NewWriter(os.Stdout)
		_ = csvWriter.Write([]string{"Subject", "Tag", "Description", "Path"})
		for k, v := range labels.FindFilesWithTags() {
			for _, template := range v {
				_ = csvWriter.Write([]string{
					template.Summary,
					k,
					template.Description,
					template.SourcePath,
				})
			}
			log.Info("tagMap", "k", k, "v", len(v))
		}
		csvWriter.Flush()
	},
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
