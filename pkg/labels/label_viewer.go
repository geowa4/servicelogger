package labels

import (
	"encoding/json"
	"github.com/geowa4/servicelogger/pkg/templates"
	"io/fs"
	"os"
	"path/filepath"
)

type Template struct {
	Severity     string   `json:"severity"`
	ServiceName  string   `json:"service_name"`
	Summary      string   `json:"summary"`
	Description  string   `json:"description"`
	InternalOnly bool     `json:"internal_only"`
	Tags         []string `json:"_tags,omitempty"`
	SourcePath   string
}

func findFilesWithTags(dir string) map[string][]*Template {
	tagMap := map[string][]*Template{}
	_ = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		file, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		template := &Template{SourcePath: path}
		err = json.Unmarshal(file, template)
		if err != nil {
			return err
		}
		if template.Tags != nil {
			for _, tag := range template.Tags {
				tagMap[tag] = append(tagMap[tag], template)
			}
		}
		return nil
	})
	return tagMap
}

func FindFilesWithTags() map[string][]*Template {
	return findFilesWithTags(templates.GetOsdServiceLogTemplatesDir())
}
