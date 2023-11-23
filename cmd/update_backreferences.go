package cmd

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/spf13/cobra"
	"os"
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
			amendedTags := false
			newTags := make([]string, 0)
			for _, tag := range template.Tags {
				// Filter out old references
				if strings.HasPrefix(tag, sopPrefix) && !slices.Contains(referencingSOPs, sopPrefix+tag) {
					amendedTags = true
					continue
				}
				newTags = append(newTags, sopPrefix+tag)
			}
			for _, sop := range referencingSOPs {
				// Ensure current references are tagged
				if !slices.Contains(template.Tags, sopPrefix+sop) {
					template.Tags = append(template.Tags, sopPrefix+sop)
					amendedTags = true
				}
			}
			if amendedTags {
				templateJson, _ := json.MarshalIndent(template, "", "  ")
				_ = os.WriteFile(templates.GetServiceLogTemplatesDir()+string(os.PathSeparator)+template.SourcePath, templateJson, 0644)
				log.Info("updated service log template", "file", templates.GetOsdServiceLogTemplatesDir()+string(os.PathSeparator)+template.SourcePath)
				templateJson, _ = json.Marshal(template)
				log.Info(string(templateJson))
			} else {
				log.Info("skipping service log", "template", template.SourcePath)
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(updateBackRefsCmd)
}
