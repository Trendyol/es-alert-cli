package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/model"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	mapset "github.com/deckarep/golang-set"
	"github.com/labstack/gommon/log"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var fetch = &cobra.Command{
	Use:   "upsert",
	Short: "Upsert all alerts",
	Long:  `Upsert command will push your monitoring yaml to remote if any change exists.`,
	Run:   upsertAlerts,
}

func upsertAlerts(cmd *cobra.Command, args []string) {
	// TODO make sure args are validated & mapped appropriately
	esAPIClient, err := client.NewElasticsearchAPI(args[0])
	fileReader, err := reader.NewFileReader()
	if err != nil {
		fmt.Println("we have an error", err)
		return
	}

	//Get Monitors(Local and Remote)
	remoteMonitors, remoteMonitorSet, err := esAPIClient.FetchMonitors()
	if err != nil {
		fmt.Println("error while read remote monitors", err)
		return
	}

	localMonitors, localMonitorSet, err := fileReader.ReadLocalYaml(cliCmd.monitoringFilename) //TODO: read yaml name from args like our old one
	if err != nil {
		fmt.Println("error while read local file", err)
		return
	}

	unTrackedMonitors := remoteMonitorSet.Difference(localMonitorSet)
	newMonitors := localMonitorSet.Difference(remoteMonitorSet)
	intersectedMonitors := remoteMonitorSet.Intersect(localMonitorSet)

	//Find modified monitors
	modifiedMonitors := mapset.NewSet()
	intersectedMonitorsIt := intersectedMonitors.Iterator()
	for intersectedMonitorName := range intersectedMonitorsIt.C {
		if isMonitorChanged(localMonitors[intersectedMonitorName.(string)], remoteMonitors[intersectedMonitorName.(string)]) {
			modifiedMonitors.Add(intersectedMonitorName)
		}
	}

	//monitorsToBeUpdated := newMonitors.Union(modifiedMonitors)
	shouldDelete := cliCmd.deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		log.Info("All monitors are up-to-date with remote monitors")
		return
	}

	//TODO: destination name should be set to destination id.
	//TODO: compare local with remote and push updated monitors

	fmt.Println(localMonitors)

	//TODO compare(localMonitors, remoteMonitors)
	if err != nil {
		fmt.Println("Client error", err)
		return
	}
	fmt.Println(remoteMonitors)
}

func init() {
	rootCmd.AddCommand(fetch)
}

func isMonitorChanged(localMonitor model.Monitor, remoteMonitor model.Monitor) bool {
	localYaml, _ := yaml.Marshal(localMonitor)
	remoteYml, _ := yaml.Marshal(remoteMonitor)
	var dmp = diffmatchpatch.New() //TODO: refactor this.
	diffs := dmp.DiffMain(string(remoteYml), string(localYaml), true)
	if len(diffs) > 1 {
		return true
	}
	return false
}
