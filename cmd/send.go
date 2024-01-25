package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/geowa4/servicelogger/pkg/ocm"
	"github.com/geowa4/servicelogger/pkg/teaspoon"
	"github.com/geowa4/servicelogger/pkg/templates"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"time"
)

var sendServiceLogCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a service log",
	Long: `Send service log to the customer from JSON template passed via stdin

` + "Example: `servicelogger search | servicelogger send -u 'https://api.openshift.com' -t \"$(ocm token)\" -c $CLUSTER_ID`",
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(checkRequiredStringArgs("ocm_url", "ocm_token", "cluster_id"))

		var template templates.Template
		input, err := io.ReadAll(os.Stdin)
		cobra.CheckErr(err)
		err = json.Unmarshal(input, &template)
		cobra.CheckErr(err)
		fmt.Println(teaspoon.RenderMarkdown(template.String()))

		confirmation := false
		err = huh.NewForm(huh.NewGroup(huh.NewConfirm().Value(&confirmation).Title("Send this service log?").Affirmative("Send").Negative("Cancel"))).Run()
		if err != nil {
			return
		}

		if confirmation {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			var errSendSL error
			go func() {
				defer cancel()
				errSendSL = sendServiceLog(
					viper.GetString("ocm_url"),
					viper.GetString("ocm_token"),
					viper.GetString("cluster_id"),
					template)
			}()
			err = spinner.New().Title("Sending service log").Context(ctx).Run()
			cobra.CheckErr(errSendSL)
			cobra.CheckErr(err)
			_, _ = fmt.Fprint(os.Stderr, "Service log sent")
		} else {
			_, _ = fmt.Fprint(os.Stderr, "Service log canceled")
		}
	},
}

func init() {
	sendServiceLogCmd.Flags().StringP("ocm-url", "u", "https://api.openshift.com", "OCM URL (falls back to $OCM_URL and then 'https://api.openshift.com')")
	_ = viper.BindPFlag("ocm_url", sendServiceLogCmd.Flags().Lookup("ocm-url"))
	sendServiceLogCmd.Flags().StringP("ocm-token", "t", "", "OCM token (falls back to $OCM_TOKEN)")
	_ = viper.BindPFlag("ocm_token", sendServiceLogCmd.Flags().Lookup("ocm-token"))
	sendServiceLogCmd.Flags().StringP("cluster-id", "c", "", "internal cluster ID (defaults to $CLUSTER_ID)")
	_ = viper.BindPFlag("cluster_id", sendServiceLogCmd.Flags().Lookup("cluster-id"))

	rootCmd.AddCommand(sendServiceLogCmd)
}

func sendServiceLog(url, token, clusterId string, t templates.Template) error {
	conn, err := ocm.NewConnectionWithTemporaryToken(url, token)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error creating ocm connection: %q", err)
	}
	defer func(conn *sdk.Connection) {
		_ = conn.Close()
	}(conn)
	client := ocm.NewClient(conn)
	err = client.PostServiceLog(clusterId, t)
	return err
}
