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
	monitorMap, _, err := fileReader.ReadLocalYaml(tempFile)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Define the expected result based on the content of the temporary YAML file
	expectedMonitors := []model.Monitor{
		{
			Name:    "Monitor1",
			Enabled: true,
			Schedule: model.Schedule{
				Period: model.Period{
					Interval: 5,
					Unit:     "MINUTES",
				},
			},
			Inputs: []model.Input{
				{
					Search: model.Search{
						Indices: []string{"test-index"},
						Query: model.QueryParam{
							Query: model.InnerQuery{
								Bool: model.BoolParam{
									AdjustPureNegative: true,
									Boost:              1,
									Must: []model.MustParam{
										{
											Match: map[string]any{
												"x.level": map[string]any{
													"auto_generate_synonyms_phrase_query": true,
													"fuzzy_transpositions":                true,
													"query":                               "ERROR",
													"boost":                               1,
													"lenient":                             false,
													"max_expansions":                      50,
													"operator":                            "AND",
													"prefix_length":                       0,
													"zero_terms_query":                    "NONE",
												},
											},
											Range: nil,
										},
										{
											Match: map[string]any{
												"kubernetes.labels.release": map[string]any{
													"auto_generate_synonyms_phrase_query": true,
													"fuzzy_transpositions":                true,
													"query":                               "${QUERY}",
													"boost":                               1,
													"lenient":                             false,
													"max_expansions":                      50,
													"operator":                            "AND",
													"prefix_length":                       0,
													"zero_terms_query":                    "NONE",
												},
											},
											Range: nil,
										},
										{
											Match: nil,
											Range: map[string]any{
												"@timestamp": map[string]any{
													"boost":         1,
													"from":          "now-5m",
													"include_lower": true,
													"include_upper": false,
													"time_zone":     "+03:00",
													"to":            "now",
												},
											},
										},
									},
									MustNot: []model.MustParam{
										{
											Match: map[string]any{
												"x.message": map[string]any{
													"auto_generate_synonyms_phrase_query": true,
													"fuzzy_transpositions":                true,
													"query":                               "Generic Exception Occurred. org.springframework.web.HttpRequestMethodNotSupportedException",
													"boost":                               1,
													"lenient":                             false,
													"max_expansions":                      50,
													"operator":                            "AND",
													"prefix_length":                       0,
													"zero_terms_query":                    "NONE",
												},
											},
											Range: nil,
										},
									},
								},
							},
						},
					},
				},
			},
			Triggers: []model.Trigger{
				{
					ID:       "",
					Name:     "test-alert",
					Severity: "3",
					Condition: model.Condition{
						Script: model.Script{
							Source: "ctx.results[0].hits.total.value > 300",
							Lang:   "painless",
						},
					},
					Actions: []model.Action{
						{
							Name:            "test-alert",
							DestinationName: "",
							DestinationID:   "inventory-alerts",
							SubjectTemplate: model.Script{
								Source: "My Test Alert",
								Lang:   "mustache",
							},
							MessageTemplate: model.Script{
								Source: "{\n  \"title\":\"title\",\n  \"monitor\":{\n      \"name\":\"monitor\",\n      \"enabled\":\"true\"\n  },\n  \"trigger\":{\n      \"id\":\"id\",\n      \"name\":\"{{ctx.trigger.name}} \\n> *Cluster-based Logs:* *<https://my-test-kibana-url.com/app/discover#/?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-20m,to:now))&_a=(columns:!(_source),filters:!(('$state':(store:appState),meta:(alias:!n,disabled:!f,index:a0980db0-49f8-11ed-8c51-6f32aac10425,key:kubernetes.labels.release.keyword,negate:!f,params:(query:product-detail-api),type:phrase),query:(match_phrase:(kubernetes.labels.release.keyword:product-detail-api))),('$state':(store:appState),meta:(alias:!n,disabled:!f,index:a0980db0-49f8-11ed-8c51-6f32aac10425,key:x.level,negate:!f,params:(query:x.level),type:phrase),query:(match_phrase:(x.level:ERROR)))),index:a0980db0-49f8-11ed-8c51-6f32aac10425,interval:auto,query:(language:kuery,query:''),sort:!(!('@timestamp',desc)))|stage >* \\n> *Sample Message:* {{ctx.results.0.hits.hits.0._source.x.message}}\\n> *Agent Name:* {{ctx.results.0.hits.hits.0._source.x.agent-name}}\\n> *Cluster:* {{ctx.results.0.hits.hits.0._source.x_cluster}}\\n> *Severity:* {{ctx.results.0.hits.hits.0._source.x.level}}\\n> *Count:* {{ctx.results.0.hits.total.value}}\",\n      \"severity\":\"1\"\n  },\n  \"periodStart\":\"start\",\n  \"periodEnd\":\"end\"\n}",
								Lang:   "mustache",
							},
						},
					},
				},
			},

			// ... (other fields)
		},
	}
	monitor := monitorMap["Monitor1"]
	// Compare the actual result with the expected result
	if !reflect.DeepEqual(monitor, expectedMonitors[0]) {
		t.Errorf("Result mismatch. Expected:\n%v\nActual:\n%v", expectedMonitors[0], monitor)
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
          - 'test-index'
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
