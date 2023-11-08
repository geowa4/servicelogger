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
	Long:  `Downloads the latest templates from openshift/managed-notifications into a local cache directory`,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindViper(cmd)
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

func bindViper(cmd *cobra.Command) {
	_ = viper.BindPFlag("cache_directory", cmd.Flags().Lookup("directory"))
}
