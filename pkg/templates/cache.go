package templates

import (
	"github.com/charmbracelet/log"
	"github.com/geowa4/servicelogger/pkg/config"
	"os"
	"os/exec"
)

func cacheGitRepo(dirName, gitURL string) {
	cloneTarget, err := config.GetCacheDir(dirName)
	if err != nil {
		log.Error("failed to load cache directory", "error", err)
		return
	}
	log.Info("syncing", "repo", gitURL, "directory", cloneTarget)
	_, statErr := os.Stat(cloneTarget)
	var cmdErr error
	if statErr != nil {
		cmd := exec.Command("git", "clone", gitURL, cloneTarget)
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
	log.Info("sync complete", "repo", gitURL, "directory", cloneTarget)
}

func CacheManagedNotifications() {
	cacheGitRepo(managedNotificationsDirName, "git@github.com:openshift/managed-notifications.git")
}

func CacheOpsSOP() {
	cacheGitRepo(opsSOPDirName, "git@github.com:openshift/ops-sop.git")
}
