package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/geowa4/servicelogger/pkg/teaspoon"
	"github.com/geowa4/servicelogger/pkg/templates"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
)

var sendServiceLogCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a service log",
	Long: `Send service log to the customer from JSON template passed via stdin

Example: servicelogger search | servicelogger send -u 'https://api.openshift.com' -t \"$(ocm token)\" -c $CLUSTER_ID"`,
	Args: cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		bindSendArgsToViper(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("cluster_ids") {
			viper.Set("cluster_ids", []string{viper.GetString("cluster_id")})
		}
		cobra.CheckErr(checkRequiredArgsExist("ocm_url", "ocm_token", "cluster_ids"))

		var template templates.Template
		input, err := io.ReadAll(os.Stdin)
		cobra.CheckErr(err)
		err = json.Unmarshal(input, &template)
		cobra.CheckErr(err)
		fmt.Println(teaspoon.RenderMarkdown(template.String()))

		clusterIds := viper.GetStringSlice("cluster_ids")
		confirmation := false
		err = huh.NewForm(huh.NewGroup(huh.NewConfirm().Value(&confirmation).Title(fmt.Sprintf("Send this service log to %v cluster(s)?", len(clusterIds))).Affirmative("Send").Negative("Cancel"))).Run()
		cobra.CheckErr(err)

		if confirmation {
			sendServiceLogsToManyClusters(viper.GetStringSlice("cluster_ids"), func(cId string) error {
				return sendServiceLog(
					viper.GetString("ocm_url"),
					viper.GetString("ocm_token"),
					cId,
					template,
				)
			})
		} else {
			_, _ = fmt.Fprint(os.Stderr, "Service log canceled")
		}
	},
}

func init() {
	setSendArgsOnCmd(sendServiceLogCmd)
	rootCmd.AddCommand(sendServiceLogCmd)
}
