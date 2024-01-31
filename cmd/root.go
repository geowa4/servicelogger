package cmd

import (
	"errors"
	"fmt"
	"github.com/geowa4/servicelogger/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "servicelogger",
	Short: "Find and use service logs to send to troublesome clusters",

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		configDir, err := config.GetConfigDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config.yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			cobra.CheckErr(fmt.Errorf("bad config file (%s): %q", viper.ConfigFileUsed(), err))
		}
	}
}

func checkRequiredArgsExist(args ...string) error {
	for _, arg := range args {
		if !viper.IsSet(arg) {
			return fmt.Errorf(
				"argument --%s or environment variable %s not set",
				strings.ReplaceAll(arg, "_", "-"),
				strings.ToUpper(arg),
			)
		}
	}
	return nil
}
