package templates

import (
	"fmt"
	"github.com/geowa4/servicelogger/pkg/config"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

const (
	managedNotificationsDirName = "managed-notifications"
	opsSOPDirName               = "ops-sop"
)

var (
	templateVarRegexp = regexp.MustCompile("\\$\\{[A-Z0-9_]+}")
)

type Template struct {
	Severity      string   `json:"severity"`
	ServiceName   string   `json:"service_name"`
	LogType       string   `json:"log_type,omitempty"`
	Summary       string   `json:"summary"`
	Desc          string   `json:"description"`
	InternalOnly  bool     `json:"internal_only"`
	EventStreamId string   `json:"event_stream_id,omitempty"`
	DocReferences []string `json:"doc_references,omitempty"`
	Tags          []string `json:"_tags,omitempty"`
	SourcePath    string   `json:"-"`
}

func (t *Template) String() string {
	// Main body
	mainBody := fmt.Sprintf(
		"# %s\n\n%s",
		t.Summary,
		t.Desc,
	)

	// Extended, optional fields
	var extendedFields strings.Builder
	if len(t.Tags) > 0 {
		extendedFields.WriteString(fmt.Sprintf("\n\n_Tags_: %s", strings.Join(t.Tags, ", ")))
	}
	if t.LogType != "" {
		extendedFields.WriteString(fmt.Sprintf("\n\n_Log Type_: %s", t.LogType))
	}
	if t.EventStreamId != "" {
		extendedFields.WriteString(fmt.Sprintf("\n\n_Event Stream ID_: %s", t.EventStreamId))
	}
	if len(t.DocReferences) > 0 {
		extendedFields.WriteString(fmt.Sprintf("\n\n_Doc Refs_: %s", strings.Join(t.DocReferences, ", ")))
	}
	md := mainBody
	if extendedFields.Len() > 0 {
		md += extendedFields.String()
	}
	return md
}

func (t *Template) GetVariables() []string {
	summaryMatches := templateVarRegexp.FindAllString(t.Summary, -1)
	descriptionMatches := templateVarRegexp.FindAllString(t.Desc, -1)
	eventStreamMatches := templateVarRegexp.FindAllString(t.EventStreamId, -1)
	allMatches := make([]string, 0)
	for _, match := range append(summaryMatches, append(descriptionMatches, eventStreamMatches...)...) {
		if !slices.Contains(allMatches, match) {
			allMatches = append(allMatches, match)
		}
	}
	return allMatches
}

func GetRelativePathForManagedNotifications(filePath string) string {
	split := strings.Split(filePath, string(os.PathSeparator)+managedNotificationsDirName+string(os.PathSeparator))
	re := regexp.MustCompile("^.*master" + string(os.PathSeparator))
	return re.ReplaceAllString(split[1], "")
}

// GetServiceLogTemplatesDir returns the directory to use to find all templates
// and returns empty string if there was an unlikely error
func GetServiceLogTemplatesDir() string {
	cloneTarget, err := config.GetCacheDir(managedNotificationsDirName)
	if err != nil {
		return ""
	}
	return cloneTarget
}

// GetOsdServiceLogTemplatesDir returns the directory to use to find templates for OSD
// and returns empty string if there was an unlikely error
func GetOsdServiceLogTemplatesDir() string {
	return filepath.Join(GetServiceLogTemplatesDir(), "osd")
}

// GetOpsSOPDir returns the directory to use to find Ops SOPs
// and returns empty string if there was an unlikely error
func GetOpsSOPDir() string {
	cloneTarget, err := config.GetCacheDir(opsSOPDirName)
	if err != nil {
		return ""
	}
	return cloneTarget
}
