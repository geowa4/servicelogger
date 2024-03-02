package config

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type formValueWrapper struct {
	OcmUrl   string `yaml:"ocm_url,omitempty"`
	CacheDir string `yaml:"cache_directory,omitempty"`
}

func Program() error {
	var (
		cacheDir = viper.GetString(CacheDirectoryKey)
		ocmUrl   = viper.GetString(OcmUrlKey)
	)
	defaultCacheDir, err := GetDefaultCacheDir()
	if err != nil {
		return err
	}
	form := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("OCM URL").
			Description("The base OCM URL for making API calls").
			Options(
				huh.NewOption("Production", "https://api.openshift.com"),
				huh.NewOption("Stage", "https://api.stage.openshift.com"),
				huh.NewOption("Integration", "https://api.integration.openshift.com"),
				huh.NewOption("Unset this property", ""),
			).
			Value(&ocmUrl),
		huh.NewInput().
			Title("Repo Cache Directory").
			Description("Optionally set this to a directory where you already have managed-notifications and ops-sop cloned").
			Placeholder("Leave blank for default: "+defaultCacheDir).
			Value(&cacheDir),
	))
	if err = form.Run(); err != nil {
		return err
	}
	viper.Set(OcmUrlKey, ocmUrl)
	if strings.HasPrefix(cacheDir, "~/") {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		cacheDir = strings.Replace(cacheDir, "~", home, 1)
	}
	viper.Set(CacheDirectoryKey, cacheDir)
	log.Info("Saving config", "file", viper.ConfigFileUsed())
	yamlBytes, err := yaml.Marshal(formValueWrapper{
		OcmUrl:   viper.GetString(OcmUrlKey),
		CacheDir: viper.GetString(CacheDirectoryKey),
	})
	if err != nil {
		return err
	}
	return os.WriteFile(viper.ConfigFileUsed(), yamlBytes, 0x644)
}
