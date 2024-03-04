package internal

import (
	"errors"
	"testing"

	clientMock "github.com/Trendyol/es-alert-cli/mocks/github.com/Trendyol/es-alert-cli/pkg/client"
	readerMock "github.com/Trendyol/es-alert-cli/mocks/github.com/Trendyol/es-alert-cli/pkg/reader"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Different_Monitor_Should_Return_Changed(t *testing.T) {
	// Arrange
	localMonitor := createTestMonitor("monitor1", "destination1")
	remoteMonitor := createTestMonitor("monitor1", "destination2")

	// Act
	changed := isMonitorChanged(localMonitor, remoteMonitor)

	// Assert
	assert.True(t, changed, "The monitors should be considered changed")
}

func Test_Same_Monitor_Should_Return_Same(t *testing.T) {
	// Arrange
	localMonitor := createTestMonitor("monitor1", "destination1")
	remoteMonitor := createTestMonitor("monitor1", "destination1")

	// Act
	changed := isMonitorChanged(localMonitor, remoteMonitor)

	// Assert
	assert.False(t, changed, "The monitors should be considered changed")
}

func Test_PrepareForCreate(t *testing.T) {
	localMonitors := map[string]model.Monitor{
		"monitor1": createTestMonitor("monitor1", "destination1"),
		"monitor2": createTestMonitor("monitor2", "destination2"),
	}
	destinations := map[string]model.Destination{
		"destination1": {ID: "dest1"},
		"destination2": {ID: "dest2"},
	}
	newMonitors := mapset.NewSet()
	newMonitors.Add("monitor1")

	expected := map[string]model.Monitor{
		"monitor1": createTestMonitor("monitor1", "dest1"),
	}

	createdMonitors := prepareForCreate(newMonitors, localMonitors, destinations)

	assert.Equal(t, expected, createdMonitors, "Prepared monitors for create should match expected")
}

func Test_IsMonitorChanged(t *testing.T) {
	localMonitor := createTestMonitor("monitor1", "destination1")
	remoteMonitor := createTestMonitor("monitor1", "destination2")

	// Monitors have different destinations, so they should be considered changed
	assert.True(t, isMonitorChanged(localMonitor, remoteMonitor), "Monitors should be considered changed")

	localMonitor = createTestMonitor("monitor1", "destination1")
	remoteMonitor = createTestMonitor("monitor1", "destination1")

	// Monitors have the same configuration, so they should not be considered changed
	assert.False(t, isMonitorChanged(localMonitor, remoteMonitor), "Monitors should not be considered changed")
}

func TestUpsert_ShouldUpdateMonitor(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockDestinations := map[string]model.Destination{"destination1": {ID: "destination1"}}
	mockRemoteMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 123},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}
	mockLocalMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 1234},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("FetchDestinations").Return(mockDestinations, nil)
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("UpdateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.NoError(t, err)
	assert.Equal(t, resp[0].ID, "monitor1")
}

func TestUpsert_ShouldCreateMonitor(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockDestinations := map[string]model.Destination{"destination1": {ID: "destination1"}}
	mockRemoteMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 123},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}
	mockLocalMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 1234},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("FetchDestinations").Return(mockDestinations, nil)
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet(), nil)
	mockClient.On("CreateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.NoError(t, err)
	assert.Equal(t, resp[0].ID, "monitor1")
}

func TestUpsert_ShouldReturnErrorWhenFetchDestinationsReturnError(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockRemoteMonitors := map[string]model.Monitor{}
	mockLocalMonitors := map[string]model.Monitor{}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("FetchDestinations").Return(nil, errors.New("fetchDestinationError"))
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("UpdateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.Equal(t, err.Error(), "err while reading destinations: fetchDestinationError\n")
	assert.Equal(t, len(resp), 0)
}

func TestUpsert_ShouldReturnErrorIfFetchMonitorReturnsError(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockDestinations := map[string]model.Destination{"destination1": {ID: "destination1"}}
	mockRemoteMonitors := map[string]model.Monitor{}
	mockLocalMonitors := map[string]model.Monitor{}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("FetchDestinations").Return(mockDestinations, nil)
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet(), errors.New("fetchMonitorError"))
	mockClient.On("UpdateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.Equal(t, err.Error(), "err while reading remote monitors: fetchMonitorError\n")
	assert.Equal(t, len(resp), 0)
}

func TestUpsert_ShouldReturnErrorIfReadLocalMonitorReturnsError(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockDestinations := map[string]model.Destination{"destination1": {ID: "destination1"}}
	mockRemoteMonitors := map[string]model.Monitor{}
	mockLocalMonitors := map[string]model.Monitor{}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet(), errors.New("localMonitorError"))
	mockClient.On("FetchDestinations").Return(mockDestinations, nil)
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet(), nil)
	mockClient.On("UpdateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.Equal(t, err.Error(), "err while reading local files: localMonitorError\n")
	assert.Equal(t, len(resp), 0)
}

func TestUpsert_ShouldReturnEmptyIfNothingToUpdate(t *testing.T) {
	mockReader := &readerMock.FileReaderInterface{}
	mockClient := &clientMock.ElasticsearchAPIClientInterface{}
	mockDestinations := map[string]model.Destination{"destination1": {ID: "destination1"}}
	mockRemoteMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 123},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}
	mockLocalMonitors := map[string]model.Monitor{"monitor1": {ID: "monitor1", Inputs: []model.Input{{
		Search: model.Search{
			Indices: []string{"test"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: false,
						Boost:              0,
						Must: []model.MustParam{
							{
								Match: map[string]any{"test": 123},
								Range: nil,
							},
						},
						MustNot: nil,
					},
				},
			},
		},
	}}}}

	mockReader.On("ReadLocalYaml", "filename").Return(mockLocalMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("FetchDestinations").Return(mockDestinations, nil)
	mockClient.On("FetchMonitors").Return(mockRemoteMonitors, mapset.NewSet("monitor1"), nil)
	mockClient.On("UpdateMonitors", mock.Anything).Return([]model.UpdateMonitorResponse{
		{
			ID: "monitor1",
		},
	})

	mockMonitorService := MonitorService{
		reader: mockReader,
		client: mockClient,
	}

	// Test Upsert
	resp, err := mockMonitorService.Upsert("filename", false)
	assert.Nil(t, err)
	assert.Equal(t, len(resp), 0)
}

func createTestMonitor(monitorName string, destinationName string) model.Monitor {
	// Create test Schedule
	schedule := model.Schedule{
		Period: model.Period{
			Interval: 5,
			Unit:     "minutes",
		},
	}

	input := createTestInput()

	trigger := createTestTrigger(destinationName)

	monitor := model.Monitor{
		ID:       "monitor1",
		Type:     "type1",
		Name:     monitorName,
		Enabled:  true,
		Schedule: schedule,
		Inputs: []model.Input{
			input,
		},
		Triggers: []model.Trigger{
			trigger,
		},
	}

	return monitor
}

func createTestTrigger(destinationName string) model.Trigger {
	return model.Trigger{
		ID:       "trigger1",
		Name:     "trigger1",
		Severity: "high",
		Condition: model.Condition{
			Script: model.Script{
				Source: "source code",
				Lang:   "painless",
			},
		},
		Actions: createTestActions(destinationName),
	}
}

func createTestActions(destinationName string) []model.Action {
	return []model.Action{
		{
			Name:          "action1",
			DestinationID: destinationName,
			SubjectTemplate: model.Script{
				Source: "subject source",
				Lang:   "painless",
			},
			MessageTemplate: model.Script{
				Source: "message source",
				Lang:   "painless",
			},
		},
	}
}

func createTestInput() model.Input {
	return model.Input{
		Search: model.Search{
			Indices: []string{"index1", "index2"},
			Query: model.QueryParam{
				Query: model.InnerQuery{
					Bool: model.BoolParam{
						AdjustPureNegative: true,
						Boost:              1.0,
						Must: []model.MustParam{
							{
								Match: map[string]interface{}{
									"field1": "value1",
								},
							},
						},
					},
				},
			},
		},
	}
}
