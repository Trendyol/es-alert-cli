package internal

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/errs"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/gommon/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v3"
)

type MonitorService struct {
	reader reader.FileReaderInterface
	client client.ElasticsearchAPIClientInterface
}

type MonitorServiceInterface interface {
	Upsert(filename string, deleteUntracked bool) error
}

func NewMonitorService(cluster string) (*MonitorService, error) {
	esAPIClient, err := client.NewElasticsearchAPI(cluster)
	if errs.LogError(err, "Error creating Elasticsearch API client") {
		return nil, err
	}
	log.Info("Elastic api client created.")

	fileReader, err := reader.NewFileReader()
	if errs.LogError(err, "err while creating file reader") {
		return nil, err
	}
	return &MonitorService{reader: fileReader, client: esAPIClient}, nil
}

func (m MonitorService) Upsert(filename string, deleteUntracked bool) ([]model.UpdateMonitorResponse, error) {

	// Get Destinations
	destinations, err := m.client.FetchDestinations()
	if err != nil {
		return nil, fmt.Errorf("%s: %v\n", "err while reading destinations", err)
	}

	log.Info("Destinations fetched.")

	// Get Remote Monitors
	remoteMonitors, remoteMonitorSet, err := m.client.FetchMonitors()
	if err != nil {
		return nil, fmt.Errorf("%s: %v\n", "err while reading remote monitors", err)
	}

	log.Info("Monitors fetched.")

	// Get Local Monitors
	localMonitors, localMonitorSet, err := m.reader.ReadLocalYaml(filename)
	if err != nil {
		return nil, fmt.Errorf("%s: %v\n", "err while reading local files", err)
	}

	log.Info("Local monitors read.")

	unTrackedMonitors := remoteMonitorSet.Difference(localMonitorSet)
	newMonitors := localMonitorSet.Difference(remoteMonitorSet)
	intersectedMonitors := remoteMonitorSet.Intersect(localMonitorSet)

	// Find modified monitors
	modifiedMonitors := findModifiedMonitors(intersectedMonitors, localMonitors, remoteMonitors)

	shouldDelete := deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		log.Info("All monitors are up-to-date with remote monitors")
		return []model.UpdateMonitorResponse{}, nil
	}

	if shouldCreate {
		return m.createMonitors(newMonitors, localMonitors, destinations), nil
	}

	if shouldUpdate {
		return m.updateMonitors(modifiedMonitors, localMonitors, remoteMonitors), nil
	}

	return []model.UpdateMonitorResponse{}, nil
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

func (m MonitorService) updateMonitors(modifiedMonitors mapset.Set, localMonitors map[string]model.Monitor, remoteMonitors map[string]model.Monitor) []model.UpdateMonitorResponse {
	monitorsToBeUpdated := prepareForUpdate(modifiedMonitors, localMonitors, remoteMonitors)
	return m.client.UpdateMonitors(monitorsToBeUpdated)
}

func (m MonitorService) createMonitors(newMonitors mapset.Set, localMonitors map[string]model.Monitor, destinations map[string]model.Destination) []model.UpdateMonitorResponse {
	monitorsToBeCreated := prepareForCreate(newMonitors, localMonitors, destinations)
	return m.client.CreateMonitors(monitorsToBeCreated)
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
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(remoteYml), string(localYaml), true)

	return len(diffs) > 1
}
