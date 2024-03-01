package cmd

import (
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
		if !viper.IsSet("cluster_ids") {
			viper.Set("cluster_ids", []string{viper.GetString("cluster_id")})
		}
		cobra.CheckErr(checkRequiredArgsExist("ocm_url", "ocm_token", "cluster_ids"))

		desc, confirmation, err := internalservicelog.Program()
		cobra.CheckErr(err)

		if confirmation {
			sendServiceLogsToManyClusters(viper.GetStringSlice("cluster_ids"), func(cId string) error {
				return sendInternalServiceLog(
					viper.GetString("ocm_url"),
					viper.GetString("ocm_token"),
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
