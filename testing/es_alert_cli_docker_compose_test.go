package testing

import (
	"bytes"
	"context"
	"fmt"
	"github.com/Trendyol/es-alert-cli/cmd"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

func TestEsAlertCli(t *testing.T) {
	ctx := context.Background()
	compose, err := tc.NewDockerCompose("docker-compose.yml")
	require.NoError(t, err, "NewDockerComposeAPI()")
	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	err = compose.
		WaitForService("opendistro", wait.ForLog("[opendistro] Node started")).
		Up(ctx, tc.Wait(true))

	esContainer, err := compose.ServiceContainer(ctx, "opendistro")
	if err != nil {
		println(err)
	}
	kibanaContainer, err := compose.ServiceContainer(ctx, "kibana")
	if err != nil {
		println(err)
	}

	elasticEndpoint, err := esContainer.Endpoint(ctx, "")
	if err != nil {
		t.Errorf("Error getting the Elasticsearch endpoint: %s", err)
	}
	elasticEndpoint = "http://" + elasticEndpoint

	println(elasticEndpoint)
	kibanaEndpoint, err := kibanaContainer.Endpoint(ctx, "")
	if err != nil {
		t.Errorf("Error getting the Kibana endpoint: %s", err)
	}
	kibanaEndpoint = "http://" + kibanaEndpoint
	println(kibanaEndpoint)

	elasticClient, err := client.NewElasticsearchAPI(elasticEndpoint, &client.BasicAuth{
		Username: "admin",
		Password: "admin",
	})
	createIndex(*elasticClient, t)

	// Create a temporary YAML file for testing
	tempFile := createTempYAMLFile(t)

	// Ensure the temporary file is removed after the test
	defer os.Remove(tempFile)

	actual := new(bytes.Buffer)
	cmd.RootCmd.SetOut(actual)
	cmd.RootCmd.SetErr(actual)
	cmd.RootCmd.SetArgs([]string{"upsert", "-c", elasticEndpoint, "-n", tempFile})

	//when
	err = cmd.RootCmd.Execute()
	if err != nil {
		println(err)
	}
	if err != nil {
		t.Errorf("Error creating elastic client %s", err)
	}

	time.Sleep(5000)
	monitors, monitorSet, err := elasticClient.FetchMonitors()
	if err != nil {
		t.Errorf("Error fething monitors: %s", err)
	}
	//then
	assert.Equal(t, 1, len(monitors), "actual is not expected")
	assert.Equal(t, 5, len(monitorSet.String()), "actual is not expected")
}

func createIndex(es client.ElasticsearchAPIClient, t *testing.T) {
	res, err := es.Client.Put("/created-index", nil)
	if err != nil {
		println(fmt.Errorf("err while creating index: %s", err))
	}
	assert.Equal(t, res.StatusCode(), 200, "index not created")
}

func createTempYAMLFile(t *testing.T) string {
	tmpfile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}

	defer tmpfile.Close()

	if _, err := tmpfile.Write([]byte(createMonitorYamlContent())); err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	return tmpfile.Name()
}

func createMonitorYamlContent() string {
	return `- name: Monitor1
  enabled: true
  schedule:
    period:
      interval: 5
      unit: MINUTES
    cron: null
  inputs:
    - search:
        indices:
          - 'created-index'
        query:
          query:
            bool:
              adjust_pure_negative: true
              boost: 1
              must:
                - match:
                    x.level:
                      auto_generate_synonyms_phrase_query: true
                      boost: 1
                      fuzzy_transpositions: true
                      lenient: false
                      max_expansions: 50
                      operator: AND
                      prefix_length: 0
                      query: ERROR
                      zero_terms_query: NONE
                - match:
                    kubernetes.labels.release:
                      auto_generate_synonyms_phrase_query: true
                      boost: 1
                      fuzzy_transpositions: true
                      lenient: false
                      max_expansions: 50
                      operator: AND
                      prefix_length: 0
                      query: ${QUERY}
                      zero_terms_query: NONE
                - range:
                    '@timestamp':
                      boost: 1
                      from: now-5m
                      include_lower: true
                      include_upper: false
                      time_zone: "+03:00"
                      to: now
              must_not:
                - match:
                    x.message:
                      auto_generate_synonyms_phrase_query: true
                      boost: 1
                      fuzzy_transpositions: true
                      lenient: false
                      max_expansions: 50
                      operator: AND
                      prefix_length: 0
                      query: Generic Exception Occurred. org.springframework.web.HttpRequestMethodNotSupportedException
                      zero_terms_query: NONE
  triggers:
    - name: test-alert
      severity: "3"
      condition:
        script:
          source: ctx.results[0].hits.total.value > 300
          lang: painless
      actions:
        - name: test-alert
          destinationID: inventory-alerts
          subject:
            source: My Test Alert
            lang: mustache
          message: 
            source: |-
              {
                "title":"title",
                "monitor":{
                    "name":"monitor",
                    "enabled":"true"
                },
                "trigger":{
                    "id":"id",
                    "name":"{{ctx.trigger.name}} \n> *Cluster-based Logs:* *<https://my-test-kibana-url.com/app/discover#/?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-20m,to:now))&_a=(columns:!(_source),filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:a0980db0-49f8-11ed-8c51-6f32aac10425,key:kubernetes.labels.release.keyword,negate:!f,params:(query:product-detail-api),type:phrase),query:(match_phrase:(kubernetes.labels.release.keyword:product-detail-api))),('$state':(store:appState),meta:(alias:!n,disabled:!f,index:a0980db0-49f8-11ed-8c51-6f32aac10425,key:x.level,negate:!f,params:(query:x.level),type:phrase),query:(match_phrase:(x.level:ERROR)))),index:a0980db0-49f8-11ed-8c51-6f32aac10425,interval:auto,query:(language:kuery,query:''),sort:!(!('@timestamp',desc)))|stage >* \n> *Sample Message:* {{ctx.results.0.hits.hits.0._source.x.message}}\n> *Agent Name:* {{ctx.results.0.hits.hits.0._source.x.agent-name}}\n> *Cluster:* {{ctx.results.0.hits.hits.0._source.x_cluster}}\n> *Severity:* {{ctx.results.0.hits.hits.0._source.x.level}}\n> *Count:* {{ctx.results.0.hits.total.value}}",
                    "severity":"1"
                },
                "periodStart":"start",
                "periodEnd":"end"
              }
            lang: mustache`
}
