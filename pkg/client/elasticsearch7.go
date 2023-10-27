package client

import (
	"github.com/elastic/go-elasticsearch/v7"
)

type ElasticsearchAPIClient struct {
	client *elasticsearch.Client
}

func NewElasticsearchAPI(client []string) (*ElasticsearchAPIClient, error) {
	elasticClient, err := NewElasticClient(client)
	if err != nil {
		return nil, err
	}
	return &ElasticsearchAPIClient{client: elasticClient}, nil
}

func (es ElasticsearchAPIClient) FetchMonitors() (map[string]any, error) {

	// Since this is very simple call to match all maximum monitors which is 1000 for now
	query := []byte(`{"size": 1000, "query":{ "match_all": {}}}`)
	_ = query
	//TODO fill in the blanks
	/* resp, err := es.client.Search(http.MethodPost,
		"/_opendistro/_alerting/monitors/_search",
		query,
		getCommonHeaders(esClient))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Error retriving all the monitors")
	}
	*/
	testData := map[string]any{"alerts": "will", "be": "here"}
	return testData, nil
}
