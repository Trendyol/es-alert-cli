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
	esAPIClient, err := client.NewElasticsearchAPI(args[0])
	fileReader, err := reader.NewFileReader()
	if err != nil {
		fmt.Println("we have an error", err)
		return
	}

	//Get Remote Monitors
	remoteMonitors, remoteMonitorSet, err := esAPIClient.FetchMonitors()
	if err != nil {
		fmt.Println("error while read remote monitors", err)
		return
	}

	//TODO: diğer takımıdn değişiklikleri intersected'a gelir.
	//TODO: dosya içeriğini kontrol etmektense son gelenin isminden dümdüz update atabiliriz.
	//TODO: takımların folder altında monitor'leri bulunacak şekilde yapmamız düzenli yapabilir bunu da
	//TODO: mesela bi alert'ı mars'a bi alert'ı moon'a farklı threshold'larla nasıl aktaracağız?

	//Get Local Monitors
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
	for intersectedMonitorName := range intersectedMonitors.Iterator().C {
		if isMonitorChanged(localMonitors[intersectedMonitorName.(string)], remoteMonitors[intersectedMonitorName.(string)]) {
			modifiedMonitors.Add(intersectedMonitorName)
		}
	}

	shouldDelete := cliCmd.deleteUntracked && unTrackedMonitors.Cardinality() > 0
	shouldUpdate := modifiedMonitors.Cardinality() > 0
	shouldCreate := newMonitors.Cardinality() > 0
	if !shouldCreate && !shouldUpdate && !shouldDelete {
		log.Info("All monitors are up-to-date with remote monitors")
		return
	}

	monitorsToBeUpdated := newMonitors.Union(modifiedMonitors)
	preparedMonitors := prepareMonitors(monitorsToBeUpdated, localMonitors, remoteMonitors)

	/*
		//TODO: push created monitors
		if shouldCreate {
			esAPIClient.PushMonitors(monitorsToBeUpdated, preparedMonitors)
		}*/

	//NOTE: in progress
	if shouldUpdate {
		esAPIClient.PushMonitors(monitorsToBeUpdated, preparedMonitors)
	}

	//TODO: continue
	fmt.Println(monitorsToBeUpdated)
	/*
		fmt.Println(remoteMonitors)*/
}

func prepareMonitors(monitorsToBeUpdated mapset.Set, localMonitors map[string]model.Monitor, remoteMonitors map[string]model.Monitor) map[string]model.Monitor {
	preparedMonitors := make(map[string]model.Monitor)

	for m := range monitorsToBeUpdated.Iterator().C {
		monitorName := m.(string)
		monitor := localMonitors[monitorName]
		monitor.Id = remoteMonitors[monitorName].Id

		for i, trigger := range remoteMonitors[monitorName].Triggers {
			monitor.Triggers[i].Id = trigger.Id

			for j, action := range trigger.Actions {
				monitor.Triggers[i].Actions[j].Id = action.Id
				monitor.Triggers[i].Actions[j].DestinationId = action.DestinationId
			}
		}

		preparedMonitors[monitorName] = monitor
	}

	return preparedMonitors
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
