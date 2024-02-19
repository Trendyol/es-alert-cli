package internal

import (
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
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
