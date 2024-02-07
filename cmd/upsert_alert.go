package cmd

import (
	"fmt"

	"github.com/Trendyol/es-alert-cli/pkg/errs"

	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	mapset "github.com/deckarep/golang-set"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var upsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Upsert all alerts",
	Long:  `Upsert command will push your monitoring yaml to remote if any change exists.`,
	Run:   upsertAlerts,
}

func init() {
	upsertCmd.Flags().StringP("cluster", "c", "", "select your cluster ip to update.")
	upsertCmd.Flags().StringP("filename", "n", "", "select your monitoring file name.")
}

func upsertAlerts(cmd *cobra.Command, args []string) {
	cluster, filename, ok := getFlagVariables(cmd)
	if !ok {
		return
	}
	fmt.Printf("Cli will connect to %s\n", cluster)
	esAPIClient, err := client.NewElasticsearchAPI(cluster)
	if errs.HandleError(err, "Error creating Elasticsearch API client") {
		return
	}
	fmt.Println("Elastic api client created.")

	fileReader, err := reader.NewFileReader()
	if errs.HandleError(err, "err while creating file reader") {
		return
	}

	// Get Destinations
	destinations, err := esAPIClient.FetchDestinations()
	if errs.HandleError(err, "err while read destinations") {
		return
	}

	fmt.Println("Destinations fetched.")

	// Get Remote Monitors
	remoteMonitors, remoteMonitorSet, err := esAPIClient.FetchMonitors()
	if errs.HandleError(err, "err while read remote monitors") {
		return
	}

	fmt.Println("Monitors fetched.")

	// Get Local Monitors
	localMonitors, localMonitorSet, err := fileReader.ReadLocalYaml(filename)
	if errs.HandleError(err, "err while read local file") {
		return
	}
	fmt.Println("Local monitors read.")

	unTrackedMonitors := remoteMonitorSet.Difference(localMonitorSet)
	newMonitors := localMonitorSet.Difference(remoteMonitorSet)
	intersectedMonitors := remoteMonitorSet.Intersect(localMonitorSet)

	// Find modified monitors
	modifiedMonitors := findModifiedMonitors(intersectedMonitors, localMonitors, remoteMonitors)

	shouldDelete := cliCmd.deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		fmt.Println("All monitors are up-to-date with remote monitors")
		return
	}

	if shouldCreate {
		createMonitors(newMonitors, localMonitors, destinations, esAPIClient)
	}

	if shouldUpdate {
		updateMonitors(modifiedMonitors, localMonitors, remoteMonitors, esAPIClient)
	}

	fmt.Println("All process completed.")
}

func getFlagVariables(cmd *cobra.Command) (string, string, bool) {
	cluster, err := cmd.Flags().GetString("cluster")
	if errs.HandleError(err, "err getting cluster parameter") {
		return "", "", false
	}
	filename, err := cmd.Flags().GetString("filename")
	if errs.HandleError(err, "err getting filename parameter") {
		return "", "", false
	}
	return cluster, filename, true
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
