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

	// Send the request to the Elasticsearch cluster
	res, err := es.client.POST("/_opendistro/_alerting/monitors/_search", alertQuery)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	err = es.client.Bind(res.Body(), &response)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Error getting response: %s", err))
	}

	//destinations, err := es.FetchDestinations()
	//if err != nil {
	//	return nil, nil, errors.New(fmt.Sprintf("Error getting destination response: %s", err))
	//}

	monitors := make(map[string]model.Monitor)
	remoteMonitorsSet := mapset.NewSet()
	for _, hit := range response.Hits.Hits {
		//for i, trigger := range hit.Source.Monitor.Triggers {
		//	hit.Source.Monitor.Triggers[i].Actions[0].DestinationId = destinations[trigger.Actions[0].DestinationId].Name
		//	hit.Source.Monitor.Triggers[i].Actions[0].Id = hit.Id
		//	hit.Source.Monitor.Triggers[i].Actions[0].Id = destinations[trigger.Actions[0].DestinationId].Name
		//}
		hit.Source.Monitor.Id = hit.Id
		monitors[hit.Source.Monitor.Name] = hit.Source.Monitor
		remoteMonitorsSet.Add(hit.Source.Monitor.Name)
	}

	return monitors, remoteMonitorsSet, nil
}

func (es ElasticsearchAPIClient) FetchDestinations() (map[string]model.Destination, error) {

	// Since this is very simple call to match all maximum monitors which is 1000 for now
	alertQuery := ElasticsearchQuery{
		"size": 1000,
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
	res, err := es.client.POST("/_opendistro/_alerting/monitors/_search", alertQuery)
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
			destinations[hit.Source.Destination.Id] = hit.Source.Destination
		}
	}

	return destinations, nil
}

func (es ElasticsearchAPIClient) PushMonitors(monitorsToBeUpdated mapset.Set, preparedMonitors map[string]model.Monitor) {
	for currentMonitor := range monitorsToBeUpdated.Iterator().C {
		monitorName := currentMonitor.(string)
		log.Debug("Running monitor: ", monitorName)
		runMonitor := preparedMonitors[monitorName]

		// Send the request to the Elasticsearch cluster
		path := fmt.Sprintf("/_opendistro/_alerting/monitors/%s", runMonitor.Id)
		res, err := es.client.PUT(path, runMonitor)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error posting monitor: %s", err)))
		}

		//TODO: validate we can push correctly and applied
		//NOTE: now, we can pushing but trigger is not going
		var monitorResponse model.UpdateMonitorResponse
		err = es.client.Bind(res.Body(), &monitorResponse)
		if err != nil {
			log.Fatal(errors.New(fmt.Sprintf("Error getting response: %s", err)))
		}
	}
}
