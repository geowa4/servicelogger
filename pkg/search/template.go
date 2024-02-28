package search

import (
	"github.com/geowa4/servicelogger/pkg/templates"
)

type ListableTemplate struct {
	templates.Template
}

func NewListableTemplate(template *templates.Template) *ListableTemplate {
	lt := &ListableTemplate{}
	lt.Severity = template.Severity
	lt.ServiceName = template.ServiceName
	lt.Summary = template.Summary
	lt.Desc = template.Desc
	lt.LogType = template.LogType
	lt.InternalOnly = template.InternalOnly
	lt.EventStreamId = template.EventStreamId
	lt.DocReferences = template.DocReferences
	lt.Tags = template.Tags
	lt.SourcePath = template.SourcePath
	return lt
}

func (t *ListableTemplate) ToTemplate() *templates.Template {
	return &templates.Template{
		Severity:      t.Severity,
		ServiceName:   t.ServiceName,
		Summary:       t.Summary,
		Desc:          t.Desc,
		LogType:       t.LogType,
		InternalOnly:  t.InternalOnly,
		EventStreamId: t.EventStreamId,
		DocReferences: t.DocReferences,
		Tags:          t.Tags,
		SourcePath:    t.SourcePath,
	}
}

func (t *ListableTemplate) FilterValue() string {
	return t.String()
}

func (t *ListableTemplate) Title() string {
	return t.Summary
}

func (t *ListableTemplate) Description() string {
	return t.Desc
}
