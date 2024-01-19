package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/charmbracelet/huh/spinner"
	"github.com/geowa4/servicelogger/pkg/list"
	"github.com/geowa4/servicelogger/pkg/ocm"
	"github.com/spf13/cobra"
)

var (
	accessToken  string
	refreshToken string
	clusterID    string
	listCmd      = &cobra.Command{
		Use:   "list",
		Short: "Display service logs",
		Long: `Display a filterable list of service logs

` + "Example: `osdctl servicelog list $CLUSTER_ID | servicelogger list`",
		Run: func(cmd *cobra.Command, args []string) {
			serviceLogList := []ocm.ServiceLog{}
			var routineErr error
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				defer cancel()
				conn, err := ocm.NewConnection(accessToken, refreshToken)
				if err != nil {
					routineErr = fmt.Errorf("could not parse input: %v\n", err)
					return
				}
				defer conn.Close()
				client := ocm.NewClient(conn)
				serviceLogList, err = client.ListServiceLogs(clusterID, "")
				if err != nil {
					routineErr = fmt.Errorf("could not get serviceLogs: %v\n", err)
					return
				}
				if len(serviceLogList) == 0 {
					routineErr = fmt.Errorf("no service logs to view")
					return
				}
			}()

			err := spinner.New().Title("loading service logs from ocm...").Context(ctx).Run()
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not get servicelogs: %v", err)
				os.Exit(0)
			}
			if routineErr != nil {
				fmt.Fprintf(os.Stderr, "could not get servicelogs: %v", err)
				os.Exit(0)
			}

			slices.Reverse(serviceLogList)
			md, err := list.Program(serviceLogList)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not run bubble program: %v", err)
				os.Exit(0)
			}
			fmt.Println(md)
			os.Exit(0)
		},
	}
)

func init() {
	listCmd.Flags().StringVarP(&clusterID, "cluster-id", "c", "", "Internal Cluster ID")
	listCmd.Flags().StringVarP(&accessToken, "access-token", "a", "", "Your OCM acccess token")
	listCmd.Flags().StringVarP(&refreshToken, "refresh-token", "r", "", "Your OCM refresh token")
	rootCmd.AddCommand(listCmd)
}
