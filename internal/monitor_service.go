package internal

import (
	"errors"
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	mapset "github.com/deckarep/golang-set"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v3"
)

type MonitorService struct {
	reader *reader.FileReader
	client *client.ElasticsearchAPIClient
}

func NewMonitorService(reader *reader.FileReader, client *client.ElasticsearchAPIClient) *MonitorService {
	return &MonitorService{reader: reader, client: client}
}

func (m MonitorService) Upsert(filename string, deleteUntracked bool) error {

	// Get Destinations
	destinations, err := m.client.FetchDestinations()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %v\n", "err while reading destinations", err))
	}

	fmt.Println("Destinations fetched.")

	// Get Remote Monitors
	remoteMonitors, remoteMonitorSet, err := m.client.FetchMonitors()
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %v\n", "err while reading remote monitors", err))
	}

	fmt.Println("Monitors fetched.")

	// Get Local Monitors
	localMonitors, localMonitorSet, err := m.reader.ReadLocalYaml(filename)
	if err != nil {
		return errors.New(fmt.Sprintf("%s: %v\n", "err while reading local files", err))
	}

	fmt.Println("Local monitors read.")

	unTrackedMonitors := remoteMonitorSet.Difference(localMonitorSet)
	newMonitors := localMonitorSet.Difference(remoteMonitorSet)
	intersectedMonitors := remoteMonitorSet.Intersect(localMonitorSet)

	// Find modified monitors
	modifiedMonitors := findModifiedMonitors(intersectedMonitors, localMonitors, remoteMonitors)

	shouldDelete := deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		fmt.Println("All monitors are up-to-date with remote monitors")
		return nil
	}

	if shouldCreate {
		createMonitors(newMonitors, localMonitors, destinations, m.client)
	}

	if shouldUpdate {
		updateMonitors(modifiedMonitors, localMonitors, remoteMonitors, m.client)
	}

	return nil
}

func findModifiedMonitors(intersectedMonitors mapset.Set, localMonitors map[string]model.Monitor, remoteMonitors map[string]model.Monitor) mapset.Set {
	modifiedMonitors := mapset.NewSet()
	for intersectedMonitorName := range intersectedMonitors.Iterator().C {
		if isMonitorChanged(localMonitors[intersectedMonitorName.(string)], remoteMonitors[intersectedMonitorName.(string)]) {
			modifiedMonitors.Add(intersectedMonitorName)
		}
	}
	return modifiedMonitors
}

func updateMonitors(modifiedMonitors mapset.Set,
	localMonitors map[string]model.Monitor,
	remoteMonitors map[string]model.Monitor,
	esAPIClient *client.ElasticsearchAPIClient,
) {
	monitorsToBeUpdated := prepareForUpdate(modifiedMonitors, localMonitors, remoteMonitors)
	esAPIClient.UpdateMonitors(monitorsToBeUpdated)
	fmt.Println("Monitors updated.")
}

func createMonitors(newMonitors mapset.Set,
	localMonitors map[string]model.Monitor,
	destinations map[string]model.Destination,
	esAPIClient *client.ElasticsearchAPIClient,
) {
	monitorsToBeCreated := prepareForCreate(newMonitors, localMonitors, destinations)
	esAPIClient.CreateMonitors(monitorsToBeCreated)
	fmt.Println("Monitors created.")
}

func prepareForCreate(monitorSet mapset.Set, localMonitors map[string]model.Monitor, destinations map[string]model.Destination) map[string]model.Monitor {
	preparedMonitors := make(map[string]model.Monitor)
	for m := range monitorSet.Iterator().C {
		monitorName := m.(string)
		monitor := localMonitors[monitorName]

		for i, trigger := range localMonitors[monitorName].Triggers {
			monitor.Triggers[i].ID = trigger.ID

			for j, action := range trigger.Actions {
				monitor.Triggers[i].Actions[j].DestinationID = destinations[action.DestinationID].ID
			}
		}

		preparedMonitors[monitorName] = monitor
	}

	return preparedMonitors
}

func prepareForUpdate(monitorsToBeUpdated mapset.Set, localMonitors map[string]model.Monitor, remoteMonitors map[string]model.Monitor) map[string]model.Monitor {
	preparedMonitors := make(map[string]model.Monitor)

	for m := range monitorsToBeUpdated.Iterator().C {
		monitorName := m.(string)
		monitor := localMonitors[monitorName]
		monitor.ID = remoteMonitors[monitorName].ID

		for i, trigger := range remoteMonitors[monitorName].Triggers {
			monitor.Triggers[i].ID = trigger.ID

			for j, action := range trigger.Actions {
				monitor.Triggers[i].Actions[j].DestinationID = action.DestinationID
			}
		}

		preparedMonitors[monitorName] = monitor
	}

	return preparedMonitors
}

func isMonitorChanged(localMonitor model.Monitor, remoteMonitor model.Monitor) bool {
	localYaml, _ := yaml.Marshal(localMonitor)
	remoteYml, _ := yaml.Marshal(remoteMonitor)
	dmp := diffmatchpatch.New() // TODO: refactor this.
	diffs := dmp.DiffMain(string(remoteYml), string(localYaml), true)

	return len(diffs) > 1
}
