package ocm

import (
	"errors"
	"fmt"

	sdk "github.com/openshift-online/ocm-sdk-go"
	clv1 "github.com/openshift-online/ocm-sdk-go/servicelogs/v1"
)

type ocmClient struct {
	// conn sdk.Connection undecided if necessary here or not
	// Ocm Cluster Client
	// clusterClient *cm v1.ClustersClient
	clusterLogsClient *clv1.ClustersClusterLogsClient
}

// Not sure if I want this to be part of the ocmClient Struct yet.
// Any ways it needs to be exposed to the user for them to close the connection
func NewClient(conn *sdk.Connection) ocmClient {
	return ocmClient{conn.ServiceLogs().V1().Clusters().ClusterLogs()}
}

func NewConnection(accessToken, refreshToken string) (*sdk.Connection, error) {
	connection, err := sdk.NewConnectionBuilder().Tokens(accessToken, refreshToken).Build()
	if err != nil {
		return &sdk.Connection{}, errors.New(fmt.Sprintf("Error building ocm sdk connection :: %s \n", err.Error()))
	}

	return connection, nil
}

func (c ocmClient) ListServiceLogs(clusterID string, query ...string) (*clv1.ClustersClusterLogsListResponse, error) {
	queryString := ""
	for i, s := range query {
		if i != 0 {
			queryString += fmt.Sprintf(" and %s", s)
		}
	}
	return c.clusterLogsClient.List().ClusterID(clusterID).Search(queryString).Send()
}
