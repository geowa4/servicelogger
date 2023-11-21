package templates

import (
	"github.com/geowa4/servicelogger/pkg/config"
	"os"
	"regexp"
	"strings"
)

const (
	managedNotificationsDirName = "managed-notifications"
	opsSOPDirName               = "ops-sop"
)

type Template struct {
	Severity      string   `json:"severity"`
	ServiceName   string   `json:"service_name"`
	Summary       string   `json:"summary"`
	Description   string   `json:"description"`
	InternalOnly  bool     `json:"internal_only"`
	EventStreamId string   `json:"event_stream_id,omitempty"`
	Tags          []string `json:"_tags,omitempty"`
	SourcePath    string   `json:"-"`
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
	return cloneTarget + string(os.PathSeparator)
}

// GetOsdServiceLogTemplatesDir returns the directory to use to find templates for OSD
// and returns empty string if there was an unlikely error
func GetOsdServiceLogTemplatesDir() string {
	return GetServiceLogTemplatesDir() + string(os.PathSeparator) + "osd"
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
