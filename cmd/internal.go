package cmd

import (
	"github.com/geowa4/servicelogger/pkg/config"
	"github.com/geowa4/servicelogger/pkg/internalservicelog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var internalServiceLogCmd = &cobra.Command{
	Use:   "internal",
	Short: "Send an internal service log",
	Long: `Prompt for comments and send internal service log.

` + commonSendArgLongHelpStanza + `
Since there can be a lot of text to send in an internal service log, the form provides hot keys to open the default "EDITOR".

Example:
    servicelogger internal -u 'https://api.openshift.com' -t \"$(ocm token)\" -c $CLUSTER_ID
`,
	Args: cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindSendArgsToViper(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet(config.ClusterIdsKey) {
			viper.Set(config.ClusterIdsKey, []string{viper.GetString(config.ClusterIdKey)})
		}
		cobra.CheckErr(checkRequiredArgsExist(config.OcmUrlKey, config.OcmTokenKey, config.ClusterIdsKey))

		desc, confirmation, err := internalservicelog.Program()
		cobra.CheckErr(err)

		if confirmation {
			sendServiceLogsToManyClusters(viper.GetStringSlice(config.ClusterIdsKey), func(cId string) error {
				return sendInternalServiceLog(
					viper.GetString(config.OcmUrlKey),
					viper.GetString(config.OcmTokenKey),
					cId,
					desc,
				)
			})
		}
	},
}

func init() {
	setSendArgsOnCmd(internalServiceLogCmd)
	rootCmd.AddCommand(internalServiceLogCmd)
}
