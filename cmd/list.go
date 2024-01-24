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
		Long: `Display a filterable list of service logs

` + "Example: `osdctl servicelog list $CLUSTER_ID | servicelogger list`",
		Run: func(cmd *cobra.Command, args []string) {
			cobra.CheckErr(checkRequiredStringArgs("ocm_url", "ocm_token", "cluster_id"))
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
				serviceLogList, err = client.ListServiceLogs(viper.GetString("cluster_id"), "")
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
	_ = viper.BindPFlag("ocm_url", listCmd.Flags().Lookup("ocm-url"))
	listCmd.Flags().StringP("ocm-token", "t", "", "OCM token (falls back to $OCM_TOKEN)")
	_ = viper.BindPFlag("ocm_token", listCmd.Flags().Lookup("ocm-token"))
	listCmd.Flags().StringP("cluster-id", "c", "", "internal cluster ID (defaults to $CLUSTER_ID)")
	_ = viper.BindPFlag("cluster_id", listCmd.Flags().Lookup("cluster-id"))

	rootCmd.AddCommand(listCmd)
}
