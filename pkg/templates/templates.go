package templates

import (
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/config"
	"os"
	"os/exec"
)

const ManagedNotificationsGitURL = "git@github.com:openshift/managed-notifications.git"

// GetOsdServiceLogTemplatesDir returns the directory to use to find templates for OSD
// and returns empty string if there was an unlikely error
func GetOsdServiceLogTemplatesDir() string {
	cloneTarget, err := config.GetCacheDir("managed-notifications")
	if err != nil {
		return ""
	}
	return cloneTarget + string(os.PathSeparator) + "osd"
}

func CacheManagedNotifications() {
	cloneTarget, err := config.GetCacheDir("managed-notifications")
	if err != nil {
		log.Error("failed to load cache directory", "error", err)
		return
	}
	log.Info("syncing", "repo", ManagedNotificationsGitURL, "directory", cloneTarget)
	_, statErr := os.Stat(cloneTarget)
	var cmdErr error
	if statErr != nil {
		cmd := exec.Command("git", "clone", ManagedNotificationsGitURL, cloneTarget)
		cmdErr = cmd.Run()
		if cmdErr != nil {
			log.Error("failed to clone", "error", cmdErr)
			return
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = cloneTarget
		cmdErr = cmd.Run()
		if cmdErr != nil {
			log.Error("failed to pull updates", "error", cmdErr)
			return
		}
	}
	log.Info("sync complete", "repo", ManagedNotificationsGitURL, "directory", cloneTarget)
}
