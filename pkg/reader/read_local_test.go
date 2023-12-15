package reader

import (
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReadLocalYaml(t *testing.T) {
	// Create a temporary YAML file for testing
	tempFile := createTempYAMLFile(t)

	// Ensure the temporary file is removed after the test
	defer os.Remove(tempFile)

	// Create an instance of your FileReader
	fileReader := &FileReader{}

	// Call the function with the temporary file path
	config, _, err := fileReader.ReadLocalYaml(tempFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Define the expected result based on the content of the temporary YAML file
	expectedConfig := []model.Monitor{
		{
			Name:    "Monitor1",
			Enabled: true,
			Schedule: model.Schedule{
				Period: model.Period{
					Interval: 5,
					Unit:     "MINUTES",
				},
			},
			// ... (other fields)
		},
		{
			Name:    "Monitor2",
			Enabled: false,
			Schedule: model.Schedule{
				Period: model.Period{
					Interval: 10,
					Unit:     "MINUTES",
				},
			},
			// ... (other fields)
		},
	}

	// Compare the actual result with the expected result
	if !reflect.DeepEqual(config, expectedConfig) {
		t.Errorf("Result mismatch. Expected:\n%v\nActual:\n%v", expectedConfig, config)
	}
}

// Helper function to create a temporary YAML file for testing
func createTempYAMLFile(t *testing.T) string {
	tmpfile, err := ioutil.TempFile("", "test-*.yaml")
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
	return `- name: ${MONITOR}
  enabled: true
  schedule:
    period:
      interval: 5
      unit: MINUTES
    cron: null
  inputs:
    - search:
        indices:
          - 'indexing-offer-log-*'
        query:
          query:
            bool:
              adjust_pure_negative: true
              boost: 1
              must:
                - match:
                    ${LEVEL}:
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
    - name: ${ALERT}
      severity: "3"
      condition: ctx.results[0].hits.total.value > ${COUNT}
      actions:
        - name: ${ALERT}
          destinationId: inventory-alerts
          subject: ${ALERT}
          message: |-
            Monitor {{ctx.monitor.name}} just entered alert status. Please investigate the issue. 
            ❗ *${ALERT}* is triggered by more than ${COUNT} errors. <!channel>
            > *Cluster-based Logs:* *<https://${CLUSTER}/app/kibana#/discover?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-20m,to:now))&_a=(columns:!(x.message,x.msg),filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'${INDEX_PATTERN_ID}',key:kubernetes.labels.release,negate:!f,params:(query:${QUERY}),type:phrase),query:(match_phrase:(kubernetes.labels.release:${QUERY}))),('$state':(store:appState),meta:(alias:!n,disabled:!f,index:'${INDEX_PATTERN_ID}',key:${LEVEL},negate:!f,params:(query:ERROR),type:phrase),query:(match_phrase:(${LEVEL}:ERROR)))),index:'${INDEX_PATTERN_ID}',interval:auto,query:(language:kuery,query:''),sort:!(!('@timestamp',desc)))|$DC>*
            > *Sample Message:* {{ctx.results.0.hits.hits.0._source.x.message}}
            > *Agent Name:* {{ctx.results.0.hits.hits.0._source.x.agent-name}}
            > *Cluster:* {{ctx.results.0.hits.hits.0._source.x.cluster}}
            > *Severity:* {{ctx.results.0.hits.hits.0._source.${LEVEL}}}
            > *Count:* {{ctx.results.0.hits.total.value}}
`
}
