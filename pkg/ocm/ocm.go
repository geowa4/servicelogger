package ocm

import (
	"errors"
	"fmt"
	"time"

	sdk "github.com/openshift-online/ocm-sdk-go"
	clv1 "github.com/openshift-online/ocm-sdk-go/servicelogs/v1"
)

type ServiceLog struct {
	ClusterId     string    `json:"cluster_id"`
	ClusterUuid   string    `json:"cluster_uuid"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `json:"created_by"`
	Desc          string    `json:"description"`
	EventStreamId string    `json:"event_stream_id"`
	Href          string    `json:"href"`
	Id            string    `json:"id"`
	InternalOnly  bool      `json:"internal_only"`
	Kind          string    `json:"kind"`
	LogType       string    `json:"log_type"`
	ServiceName   string    `json:"service_name"`
	Severity      string    `json:"severity"`
	Summary       string    `json:"summary"`
	Timestamp     time.Time `json:"timestamp"`
	Username      string    `json:"username"`
}

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

func (c ocmClient) ListServiceLogs(clusterID string, query ...string) ([]ServiceLog, error) {
	queryString := ""
	for i, s := range query {
		if i != 0 {
			queryString += fmt.Sprintf(" and %s", s)
		}
	}

	list := []ServiceLog{}
	page := 1
	size := 50
	for {
		resp, err := c.clusterLogsClient.List().
			ClusterID(clusterID).
			Search(queryString).
			Size(size).
			Page(page).
			Send()
		if err != nil {
			return []ServiceLog{}, err
		}

		resp.Items().Each(func(logEntry *clv1.LogEntry) bool {
			list = append(list, ServiceLog{
				logEntry.ClusterID(),
				logEntry.ClusterUUID(),
				logEntry.CreatedAt(),
				logEntry.CreatedBy(),
				logEntry.Description(),
				logEntry.EventStreamID(),
				logEntry.HREF(),
				logEntry.ID(),
				logEntry.InternalOnly(),
				logEntry.Kind(),
				string(logEntry.LogType()),
				logEntry.ServiceName(),
				string(logEntry.Severity()),
				logEntry.Summary(),
				logEntry.Timestamp(),
				logEntry.Username(),
			})
			return true
		})

		if resp.Size() < size {
			break
		}
		page++
	}

	return list, nil
}
