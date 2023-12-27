package client

import (
	"errors"
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/gommon/log"
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

	var response model.ElasticFetchResponse

	res, err := es.client.POST("/_opendistro/_alerting/monitors/_search", alertQuery)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	err = es.client.Bind(res.Body(), &response)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	//TODO: this one is setting wrong values for action destination id and action id so I just commented out them. If it is not needed lets move it completely.
	//destinations, err := es.FetchDestinations()
	//if err != nil {
	//	return nil, nil, errors.New(fmt.Sprintf("Error getting destination response: %s", err))
	//}

	monitors := make(map[string]model.Monitor)
	remoteMonitorsSet := mapset.NewSet()
	for _, hit := range response.Hits.Hits {
		hit.Source.Monitor.Id = hit.Id
		monitors[hit.Source.Monitor.Name] = hit.Source.Monitor
		remoteMonitorsSet.Add(hit.Source.Monitor.Name)
	}

	return monitors, remoteMonitorsSet, nil
}

func (es ElasticsearchAPIClient) FetchDestinations() (map[string]model.Destination, error) {

	// Since this is very simple call to match all maximum monitors which is 1000 for now
	query := ElasticsearchQuery{
		"query": ElasticsearchQuery{
			"bool": ElasticsearchQuery{
				"must": ElasticsearchQuery{
					"exists": ElasticsearchQuery{
						"field": "destination",
					},
				},
			},
		},
	}

	var response model.ElasticFetchResponse

	// Send the request to the Elasticsearch cluster
	res, err := es.client.POST("/_opendistro/_alerting/monitors/_search", query)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	err = es.client.Bind(res.Body(), &response)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	destinations := make(map[string]model.Destination)
	for _, hit := range response.Hits.Hits {
		if hit.Source.Destination.Name != "" {
			destination := hit.Source.Destination
			//note: if the destinationId does not come from remote, we should be able to operate with destinationName so that we can operate from destinationName when changing the destination and creating a new monitor occur.
			destinations[hit.Source.Destination.Name] = destination
		}
	}

	return destinations, nil
}

func (es ElasticsearchAPIClient) UpdateMonitors(preparedMonitors map[string]model.Monitor) {
	for monitorName, currentMonitor := range preparedMonitors {
		// Select monitor
		log.Debug("Running monitor: ", monitorName)

		// Send the request to the Elasticsearch cluster
		path := fmt.Sprintf("/_opendistro/_alerting/monitors/%s", currentMonitor.Id)
		res, err := es.client.PUT(path, currentMonitor)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error posting monitor: %s", err)))
		}

		// Bind the response
		var monitorResponse model.UpdateMonitorResponse
		err = es.client.Bind(res.Body(), &monitorResponse)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error getting response: %s", err)))
		}
	}
}

func (es ElasticsearchAPIClient) CreateMonitors(preparedMonitors map[string]model.Monitor) {
	for monitorName, currentMonitor := range preparedMonitors {
		// Select monitor
		log.Debug("Running monitor: ", monitorName)

		// Send the request to the Elasticsearch cluster
		path := fmt.Sprintf("/_opendistro/_alerting/monitors/%s", currentMonitor.Id)
		res, err := es.client.POST(path, currentMonitor)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error posting monitor: %s", err)))
		}

		// Bind the response
		var monitorResponse model.UpdateMonitorResponse
		err = es.client.Bind(res.Body(), &monitorResponse)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error getting response: %s", err)))
		}
	}
}
