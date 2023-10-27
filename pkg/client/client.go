package client

import (
	"math"

	"github.com/elastic/go-elasticsearch/v7"
)

func NewElasticClient(url []string) (*elasticsearch.Client, error) {
	//TODO configuration
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		MaxRetries: math.MaxInt,
		Addresses:  url,
		Transport:  newTransport(),
		//	CompressRequestBody:   config.Elasticsearch.CompressionEnabled,
		//	DiscoverNodesOnStart:  !config.Elasticsearch.DisableDiscoverNodesOnStart,
		//	DiscoverNodesInterval: *config.Elasticsearch.DiscoverNodesInterval,
	})
	if err != nil {
		return nil, err
	}
	return es, nil
}
