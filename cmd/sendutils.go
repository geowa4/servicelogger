package cmd

import (
	"fmt"
	"github.com/geowa4/servicelogger/pkg/config"
	"github.com/geowa4/servicelogger/pkg/ocm"
	"github.com/geowa4/servicelogger/pkg/templates"
	sdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"sync"
)

func setSendArgsOnCmd(cmd *cobra.Command) {
	cmd.Flags().StringP("ocm-url", "u", "https://api.openshift.com", "OCM URL (falls back to $OCM_URL and then 'https://api.openshift.com')")
	cmd.Flags().StringP("ocm-token", "t", "", "OCM token (falls back to $OCM_TOKEN)")
	cmd.Flags().StringP("cluster-id", "c", "", "internal cluster ID (defaults to $CLUSTER_ID)")
	cmd.Flags().StringSlice("cluster-ids", nil, "internal cluster IDs (defaults to $CLUSTER_IDS, space separated)")

	cmd.MarkFlagsMutuallyExclusive("cluster-id", "cluster-ids")
}

func bindSendArgsToViper(cmd *cobra.Command) {
	_ = viper.BindPFlag(config.OcmUrlKey, cmd.Flags().Lookup("ocm-url"))
	_ = viper.BindPFlag(config.OcmTokenKey, cmd.Flags().Lookup("ocm-token"))
	_ = viper.BindPFlag(config.ClusterIdKey, cmd.Flags().Lookup("cluster-id"))
	_ = viper.BindPFlag(config.ClusterIdsKey, cmd.Flags().Lookup("cluster-ids"))
}

func sendServiceLogsToManyClusters(clusterIds []string, sendFunc func(cId string) error) {
	var serviceLogWaitGroup sync.WaitGroup
	var printMutex sync.Mutex

	for _, clusterId := range clusterIds {
		serviceLogWaitGroup.Add(1)
		go func(cId string) {
			defer serviceLogWaitGroup.Done()

			errSendSl := sendFunc(cId)

			result := fmt.Sprintf("%s\t", cId)
			if errSendSl == nil {
				result += "success"
			} else {
				result += fmt.Sprintf("failure\t%v", errSendSl.Error())
			}

			printMutex.Lock()
			defer printMutex.Unlock()
			fmt.Println(result)
		}(clusterId)
	}

	serviceLogWaitGroup.Wait()
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
