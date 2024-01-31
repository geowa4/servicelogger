package cmd

import (
	"context"
	"fmt"
	"github.com/charmbracelet/huh/spinner"
	"github.com/geowa4/servicelogger/pkg/internalservicelog"
	"github.com/geowa4/servicelogger/pkg/ocm"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

var internalServiceLogCmd = &cobra.Command{
	Use:   "internal",
	Short: "Send an internal service log",
	Long: `Prompt for and send internal service log 

` + "Example: `servicelogger internal -u 'https://api.openshift.com' -t \"$(ocm token)\" -c $CLUSTER_ID`",
	Args: cobra.NoArgs,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("ocm_url", cmd.Flags().Lookup("ocm-url"))
		_ = viper.BindPFlag("ocm_token", cmd.Flags().Lookup("ocm-token"))
		_ = viper.BindPFlag("cluster_id", cmd.Flags().Lookup("cluster-id"))
	},
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(checkRequiredArgsExist("ocm_url", "ocm_token", "cluster_id"))
		desc, confirmation, err := internalservicelog.Program()
		cobra.CheckErr(err)
		if confirmation {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			var errSendSL error
			go func() {
				defer cancel()
				errSendSL = sendInternalServiceLog(
					viper.GetString("ocm_url"),
					viper.GetString("ocm_token"),
					viper.GetString("cluster_id"),
					desc)
			}()
			err = spinner.New().Title("Sending service log").Context(ctx).Run()
			cobra.CheckErr(errSendSL)
			cobra.CheckErr(err)
			_, _ = fmt.Fprint(os.Stderr, "Service log sent")
		}
	},
}

func init() {
	internalServiceLogCmd.Flags().StringP("ocm-url", "u", "https://api.openshift.com", "OCM URL (falls back to $OCM_URL and then 'https://api.openshift.com')")
	internalServiceLogCmd.Flags().StringP("ocm-token", "t", "", "OCM token (falls back to $OCM_TOKEN)")
	internalServiceLogCmd.Flags().StringP("cluster-id", "c", "", "internal cluster ID (defaults to $CLUSTER_ID)")

	rootCmd.AddCommand(internalServiceLogCmd)
}

func sendInternalServiceLog(url, token, clusterId, description string) error {
	conn, err := ocm.NewConnectionWithTemporaryToken(url, token)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error creating ocm connection: %q", err)
	}
	defer func(conn *sdk.Connection) {
		_ = conn.Close()
	}(conn)
	client := ocm.NewClient(conn)
	err = client.PostInternalServiceLog(clusterId, description)
	return err
}
