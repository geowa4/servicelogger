package templates

import (
	"encoding/json"
	"io/fs"
	"mvdan.cc/xurls/v2"
	"os"
	"path/filepath"
	"strings"
)

func readFileContents(path string, info fs.FileInfo, err error) ([]byte, error) {
	if err != nil || info.IsDir() {
		return []byte{}, err
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return fileBytes, nil
}

func WalkTemplates(processTemplate func(template *Template)) {
	_ = filepath.Walk(GetOsdServiceLogTemplatesDir(), func(path string, info fs.FileInfo, err error) error {
		fileBytes, err := readFileContents(path, info, err)
		if err != nil || len(fileBytes) == 0 {
			return err
		}
		template := &Template{SourcePath: GetRelativePathForManagedNotifications(path)}
		err = json.Unmarshal(fileBytes, template)
		if err != nil {
			return err
		}
		processTemplate(template)
		return nil
	})
}

func FindReferencingV4SOPs() map[string][]string {
	urls := map[string][]string{}
	_ = filepath.Walk(GetOpsSOPDir()+string(os.PathSeparator)+"v4"+string(os.PathSeparator)+"alerts", func(path string, info fs.FileInfo, err error) error {
		fileContents, err := readFileContents(path, info, err)
		if err != nil {
			return err
		}
		rxStrict := xurls.Strict()
		for _, link := range rxStrict.FindAllString(string(fileContents), -1) {
			if strings.Contains(link, managedNotificationsDirName) {
				normalizedManagedNotificationPath := GetRelativePathForManagedNotifications(link)
				if _, ok := urls[normalizedManagedNotificationPath]; !ok {
					urls[normalizedManagedNotificationPath] = make([]string, 0, 1)
				}
				urls[normalizedManagedNotificationPath] = append(
					urls[normalizedManagedNotificationPath],
					strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
				)
			}
		}
		return nil
	})
	return urls
}
