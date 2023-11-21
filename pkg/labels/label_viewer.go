package labels

import (
	"github.com/geowa4/servicelogger/pkg/templates"
)

func FindFilesWithTags() map[string][]*templates.Template {
	tagMap := map[string][]*templates.Template{}
	templates.WalkTemplates(func(template *templates.Template) {
		for _, tag := range template.Tags {
			tagMap[tag] = append(tagMap[tag], template)
		}
		if template.Tags == nil {
			tagMap["untagged"] = append(tagMap["untagged"], template)
		}
	})
	return tagMap
}
