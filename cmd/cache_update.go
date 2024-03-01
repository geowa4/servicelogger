package cmd

import (
	"github.com/geowa4/servicelogger/pkg/config"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cacheUpdateCmd = &cobra.Command{
	Use:   "cache-update",
	Short: "Initialize or update the cached service log templates",
	Long: `Downloads the latest templates from openshift/managed-notifications into a local cache directory.

This will execute "git clone git@github.com:openshift/managed-notifications.git" if the directory does not exist or it will "git pull" if it does. The same will occur for openshift/ops-sop though if this fails, it will only impact the "update-backreferences" command.

The default cache directory may be changed by specifying "--cache-directory", setting the environment variable "CACHE_DIRECTORY", or setting "cache_directory" in the global config file. 

*It is HIGHLY recommend this command be executed often to keep the templates up to date.* 
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("cache_directory", cmd.Flags().Lookup("directory"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		templates.CacheManagedNotifications()
		templates.CacheOpsSOP()
	},
}

func init() {
	rootCmd.AddCommand(cacheUpdateCmd)
	cacheDir, err := config.GetDefaultCacheDir()
	cobra.CheckErr(err)
	cacheUpdateCmd.Flags().StringP("directory", "d", cacheDir, "Cache directory root")
}
