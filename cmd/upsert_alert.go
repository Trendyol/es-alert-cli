package cmd

import (
	"fmt"
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
	cluster, err := cmd.Flags().GetString("cluster")
	if err != nil {
		fmt.Println("error getting cluster parameter", err)
		return
	}
	fmt.Printf("Cli will connect to %s\n", cluster)

	filename, err := cmd.Flags().GetString("filename")
	if err != nil {
		fmt.Println("error getting filename parameter", err)
		return
	}

	esAPIClient, err := client.NewElasticsearchAPI(cluster)
	fmt.Println("Elastic api client created.")

	fileReader, err := reader.NewFileReader()
	if err != nil {
		fmt.Println("we have an error", err)
		return
	}

	//Get Destinations
	destinations, err := esAPIClient.FetchDestinations()
	if err != nil {
		fmt.Println("error while read destinations", err)
		return
	}
	fmt.Println("Destinations fetched.")

	//Get Remote Monitors
	remoteMonitors, remoteMonitorSet, err := esAPIClient.FetchMonitors()
	if err != nil {
		fmt.Println("error while read remote monitors", err)
		return
	}
	fmt.Println("Monitors fetched.")

	//Get Local Monitors
	localMonitors, localMonitorSet, err := fileReader.ReadLocalYaml(filename)
	if err != nil {
		fmt.Println("error while read local file", err)
		return
	}
	fmt.Println("Local monitors read.")

	unTrackedMonitors := remoteMonitorSet.Difference(localMonitorSet)
	newMonitors := localMonitorSet.Difference(remoteMonitorSet)
	intersectedMonitors := remoteMonitorSet.Intersect(localMonitorSet)

	//Find modified monitors
	modifiedMonitors := mapset.NewSet()
	for intersectedMonitorName := range intersectedMonitors.Iterator().C {
		if isMonitorChanged(localMonitors[intersectedMonitorName.(string)], remoteMonitors[intersectedMonitorName.(string)]) {
			modifiedMonitors.Add(intersectedMonitorName)
		}
	}

	shouldDelete := cliCmd.deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		fmt.Println("All monitors are up-to-date with remote monitors")
		return
	}

	if shouldCreate {
		monitorsToBeCreated := prepareForCreate(newMonitors, localMonitors, destinations)
		esAPIClient.CreateMonitors(monitorsToBeCreated)
		fmt.Println("Monitors created.")
	}

	if shouldUpdate {
		monitorsToBeUpdated := prepareForUpdate(modifiedMonitors, localMonitors, remoteMonitors)
		esAPIClient.UpdateMonitors(monitorsToBeUpdated)
		fmt.Println("Monitors updated.")
	}

	fmt.Println("All process completed.")
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
	var dmp = diffmatchpatch.New() //TODO: refactor this.
	diffs := dmp.DiffMain(string(remoteYml), string(localYaml), true)

	return len(diffs) > 1
}
