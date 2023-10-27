package cmd

import (
	"fmt"
	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"os"
)

var fetch = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all alerts",
	Long:  `Long text example`,
	Run:   showDiff,
}

func fetchAlerts(cmd *cobra.Command, args []string) {
	destinations, err := destination.GetRemote(esClient)

	localMonitors, localMonitorSet, err := monitor.GetAllLocal(rootDir)

	if localMonitorSet.Cardinality() == 0 {
		log.Info("There are no monitors")
		os.Exit(1)
	}
	allRemoteMonitors, remoteMonitorsSet, err := monitor.GetAllRemote(esClient, destinations)
	check(err)
	allNewMonitors := localMonitorSet.Difference(remoteMonitorsSet)
	allCommonMonitors := remoteMonitorsSet.Intersect(localMonitorSet)

	changedMonitors := mapset.NewSet()
	allCommonMonitorsIt := allCommonMonitors.Iterator()
	for commonMonitor := range allCommonMonitorsIt.C {
		monitorName := commonMonitor.(string)
		if isMonitorChanged(localMonitors[monitorName], allRemoteMonitors[monitorName]) == true {
			changedMonitors.Add(commonMonitor)
		}
	}
	hasCreated := allNewMonitors.Cardinality() > 0
	//All New Monitors
	if hasCreated {
		log.Debug("New monitors to be publushed", allNewMonitors)
		fmt.Println("---------------------------------------------------------")
		fmt.Println(" These monitors are currently missing in alerting ")
		fmt.Println("---------------------------------------------------------")
		for newMonitor := range allNewMonitors.Iterator().C {
			monitorName := newMonitor.(string)
			localYaml, _ := yaml.Marshal(localMonitors[monitorName])
			color.Green(string(localYaml))
		}
	}
}

func init() {
	rootCmd.AddCommand(fetch)
}
