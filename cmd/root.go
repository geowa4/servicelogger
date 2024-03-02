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
	Use: "servicelogger",
	Long: `Find and use service logs to send to clusters.

When this is first installed, the cache-update command must be executed. Please see its help for more.

Then, assuming the proper environment variables are set, usage can be as simple as the following.

    servicelogger search | servicelogger send

    servicelogger list

    servicelogger internal

Please see each command's long form help for more details.
`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "default location is ~/.config/servicelogger/config.yaml")
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
		viper.SetConfigType(config.FileType)
		viper.SetConfigName(config.FileName)
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
