package client

import (
	"errors"
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	mapset "github.com/deckarep/golang-set"
)

type ElasticsearchAPIClient struct {
	client *BaseClient
}

func NewElasticsearchAPI(client string) (*ElasticsearchAPIClient, error) {
	return &ElasticsearchAPIClient{client: NewBaseClient(client)}, nil
}

type ElasticsearchQuery map[string]interface{}

func (es ElasticsearchAPIClient) FetchMonitors() (map[string]model.Monitor, mapset.Set, error) {

	// Since this is very simple call to match all maximum monitors which is 1000 for now
	alertQuery := ElasticsearchQuery{
		"size": 1000,
		"query": ElasticsearchQuery{
			"match_all": ElasticsearchQuery{},
		},
	}

	var monitorResponse model.MonitorResponse

	// Send the request to the Elasticsearch cluster
	res, err := es.client.POST("/_opendistro/_alerting/monitors/_search", alertQuery)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	err = es.client.Bind(res.Body(), &monitorResponse)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	var monitors map[string]model.Monitor
	remoteMonitorsSet := mapset.NewSet()
	for _, hit := range monitorResponse.Hits.Hits {
		monitors[hit.Source.Monitor.Name] = hit.Source.Monitor
		remoteMonitorsSet.Add(hit.Source.Monitor.Name)
	}

	return monitors, remoteMonitorsSet, nil
}
