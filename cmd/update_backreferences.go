package cmd

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const sopPrefix = "sop_"

var updateBackRefsCmd = &cobra.Command{
	Use:   "update-backreferences",
	Short: "Update managed notifications to include referencing SOPs as tags",
	Long:  `Update managed notifications to include referencing SOPs as tags`,
	Run: func(cmd *cobra.Command, args []string) {
		slToReferencingSOP := templates.FindReferencingV4SOPs()
		templates.WalkTemplates(func(template *templates.Template) {
			referencingSOPs := slToReferencingSOP[template.SourcePath]
			newTags := make([]string, 0)
			for _, tag := range template.Tags {
				// Seed with non-sop tags
				if !strings.HasPrefix(tag, sopPrefix) {
					newTags = append(newTags, tag)
				}
			}
			for _, sop := range referencingSOPs {
				// Ensure current references are tagged
				if !slices.Contains(newTags, sopPrefix+sop) {
					newTags = append(newTags, sopPrefix+sop)
				}
			}
			template.Tags = newTags
			templateJson, _ := json.MarshalIndent(template, "", "  ")
			_ = os.WriteFile(filepath.Join(templates.GetServiceLogTemplatesDir(), template.SourcePath), templateJson, 0644)
			log.Info("updated service log template", "file", filepath.Join(templates.GetServiceLogTemplatesDir(), template.SourcePath))
		})
	},
}

func init() {
	rootCmd.AddCommand(updateBackRefsCmd)
}
