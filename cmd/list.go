package cmd

import (
	"context"
	"errors"
	"fmt"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/spf13/viper"
	"os"
	"slices"

	"github.com/charmbracelet/huh/spinner"
	"github.com/geowa4/servicelogger/pkg/list"
	"github.com/geowa4/servicelogger/pkg/ocm"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "Display service logs",
		Long: `Display a filterable list of service logs.
` + commonOcmFlagLongHelpStanza + `
The cluster ID must also be specified via the "--cluster-id" flag or the "CLUSTER_ID" environment variable.

Example:
    servicelogger list -u 'https://api.openshift.com' -t "$(ocm token)" -c "$CLUSTER_ID"
`,
		Args: cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlag("ocm_url", cmd.Flags().Lookup("ocm-url"))
			_ = viper.BindPFlag("ocm_token", cmd.Flags().Lookup("ocm-token"))
			_ = viper.BindPFlag("cluster_id", cmd.Flags().Lookup("cluster-id"))
		},
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(checkRequiredArgsExist("ocm_url", "ocm_token", "cluster_id"))
			serviceLogList := make([]ocm.ServiceLog, 0)
			var routineErr error
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				defer cancel()
				conn, err := ocm.NewConnectionWithTemporaryToken(
					viper.GetString("ocm_url"),
					viper.GetString("ocm_token"),
				)
				if err != nil {
					routineErr = fmt.Errorf("could not parse input: %v\n", err)
					return
				}
				defer func(conn *sdk.Connection) {
					_ = conn.Close()
				}(conn)
				client := ocm.NewClient(conn)

				var queryStrings []string
				internalFlag, _ := cmd.Flags().GetBool("internal-only")
				if internalFlag {
					queryStrings = append(queryStrings, "internal_only='true'")
				}
				allFlag, _ := cmd.Flags().GetBool("all")
				if !allFlag {
					queryStrings = append(queryStrings, "service_name='SREManualAction'")
				}

				serviceLogList, err = client.ListServiceLogs(viper.GetString("cluster_id"), queryStrings...)

				if err != nil {
					routineErr = fmt.Errorf("could not get serviceLogs: %v\n", err)
					return
				}
				if len(serviceLogList) == 0 {
					routineErr = errors.New("no service logs to view")
					return
				}
			}()

			err := spinner.New().Title("loading service logs from ocm...").Context(ctx).Run()
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not get servicelogs: %v", err)
				os.Exit(1)
			}
			if routineErr != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not get servicelogs: %v", err)
				os.Exit(1)
			}

			slices.Reverse(serviceLogList)
			md, err := list.Program(serviceLogList)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "could not run bubble program: %v", err)
				os.Exit(1)
			}
			fmt.Println(md)
			os.Exit(0)
		},
	}
)

func init() {
	listCmd.Flags().StringP("ocm-url", "u", "https://api.openshift.com", "OCM URL (falls back to $OCM_URL and then 'https://api.openshift.com')")
	listCmd.Flags().StringP("ocm-token", "t", "", "OCM token (falls back to $OCM_TOKEN)")
	listCmd.Flags().StringP("cluster-id", "c", "", "internal cluster ID (defaults to $CLUSTER_ID)")
	listCmd.Flags().BoolP("all", "a", false, "Whether to return all service logs. By default, results are filtered on service_name='SREManualAction'")
	listCmd.Flags().BoolP("internal-only", "i", false, "Whether to only return internal service logs.")

	listCmd.MarkFlagsMutuallyExclusive("all", "internal-only")

	rootCmd.AddCommand(listCmd)
}
