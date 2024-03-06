package client

import (
	"fmt"

	"github.com/Trendyol/es-alert-cli/pkg/model"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/gommon/log"
)

type ElasticsearchAPIClient struct {
	Client *BaseClient
}

func NewElasticsearchAPI(client string, auth *BasicAuth) (*ElasticsearchAPIClient, error) {
	return &ElasticsearchAPIClient{Client: NewBaseClient(client, auth)}, nil
}

type ElasticsearchAPIClientInterface interface {
	FetchDestinations() (map[string]model.Destination, error)
	FetchMonitors() (map[string]model.Monitor, mapset.Set, error)
	UpdateMonitors(preparedMonitors map[string]model.Monitor) []model.UpdateMonitorResponse
	CreateMonitors(monitors map[string]model.Monitor) []model.UpdateMonitorResponse
}

type ElasticsearchQuery map[string]interface{}

func (es *ElasticsearchAPIClient) FetchMonitors() (map[string]model.Monitor, mapset.Set, error) {
	// Since this is very simple call to match all maximum monitors which is 1000 for now
	alertQuery := ElasticsearchQuery{
		"size": 1000,
		"query": ElasticsearchQuery{
			"match_all": ElasticsearchQuery{},
		},
	}

	var response model.ElasticFetchResponse

	res, err := es.Client.Post("/_opendistro/_alerting/monitors/_search", alertQuery)
	if err != nil {
		return nil, nil, fmt.Errorf("err while getting monitor response: %s", err)
	}

	err = es.Client.Bind(res.Body(), &response)
	if err != nil {
		return nil, nil, fmt.Errorf("err while binding monitor response: %s", err)
	}

	monitors := make(map[string]model.Monitor)
	remoteMonitorsSet := mapset.NewSet()
	for _, hit := range response.Hits.Hits {
		hit.Source.Monitor.ID = hit.ID
		monitors[hit.Source.Monitor.Name] = hit.Source.Monitor
		remoteMonitorsSet.Add(hit.Source.Monitor.Name)
	}

	return monitors, remoteMonitorsSet, nil
}

func (es *ElasticsearchAPIClient) FetchDestinations() (map[string]model.Destination, error) {
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
	res, err := es.Client.Post("/_opendistro/_alerting/monitors/_search", query)
	if err != nil {
		log.Fatal(fmt.Errorf("err while getting destination response: %s", err))
	}

	err = es.Client.Bind(res.Body(), &response)
	if err != nil {
		log.Fatal(fmt.Errorf("err while binding destination response: %s", err))
	}

	destinations := make(map[string]model.Destination)
	for _, hit := range response.Hits.Hits {
		if hit.Source.Destination.Name != "" {
			destination := hit.Source.Destination
			if destination.ID == "" {
				destination.ID = hit.ID
			}
			// note: if the destinationID does not come from remote,
			// we should be able to operate with destinationName
			// so that we can operate from destinationName when changing the destination and creating a new monitor occur.
			destinations[hit.Source.Destination.Name] = destination
		}
	}

	return destinations, nil
}

func (es *ElasticsearchAPIClient) UpdateMonitors(preparedMonitors map[string]model.Monitor) []model.UpdateMonitorResponse {
	response := make([]model.UpdateMonitorResponse, len(preparedMonitors))
	for monitorName, currentMonitor := range preparedMonitors {
		// Select monitor
		log.Debug("Running monitor: ", monitorName)

		// Send the request to the Elasticsearch cluster
		path := fmt.Sprintf("/_opendistro/_alerting/monitors/%s", currentMonitor.ID)
		res, err := es.Client.Put(path, currentMonitor)
		if err != nil {
			log.Fatal(fmt.Errorf("err while updating monitor: %s", err))
		}

		// Bind the response
		var monitorResponse model.UpdateMonitorResponse
		err = es.Client.Bind(res.Body(), &monitorResponse)
		if err != nil {
			log.Fatal(fmt.Errorf("err while binding monitor update response, response: %s", err))
		}
	}
	log.Info("Monitors updated.")
	return response
}

func (es *ElasticsearchAPIClient) CreateMonitors(preparedMonitors map[string]model.Monitor) []model.UpdateMonitorResponse {
	response := make([]model.UpdateMonitorResponse, len(preparedMonitors))
	for monitorName, currentMonitor := range preparedMonitors {
		// Select monitor
		log.Debug("Running monitor: ", monitorName)

		// Send the request to the Elasticsearch cluster
		path := fmt.Sprintf("/_opendistro/_alerting/monitors/%s", currentMonitor.ID)
		res, err := es.Client.Post(path, currentMonitor)
		if err != nil {
			log.Fatal(fmt.Errorf("err while posting to create monitor: %s", err))
		}

		// Bind the response
		var monitorResponse model.UpdateMonitorResponse
		err = es.Client.Bind(res.Body(), &monitorResponse)
		if err != nil {
			log.Fatal(fmt.Errorf("err while binding create monitor response: %s", err))
		}
		response = append(response, monitorResponse)
	}
	log.Info("Monitors created.")
	return response
}
