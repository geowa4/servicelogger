package config

import (
	"github.com/spf13/viper"
	"os"
)

const commonDirectoryName = "servicelogger"

func getDir(cacheOrConfig string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home + string(os.PathSeparator) + "." + cacheOrConfig + string(os.PathSeparator) + commonDirectoryName, nil
}

func GetConfigDir() (string, error) {
	return getDir("config")
}

func GetDefaultCacheDir() (string, error) {
	cacheDir, err := getDir("cache")
	if err != nil {
		return "", err
	}
	return cacheDir, nil
}

func GetCacheDir(subDir string) (string, error) {
	var err error
	cacheDir := viper.GetString("cache_directory")
	if cacheDir == "" {
		cacheDir, err = GetDefaultCacheDir()
		if err != nil {
			return "", err
		}
	}
	return cacheDir + string(os.PathSeparator) + subDir, nil
}
