package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	"github.com/spf13/cobra"
)

var fetch = &cobra.Command{
	Use:   "upsert",
	Short: "Upsert all alerts",
	Long:  `Long text example`,
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
	monitors, err := esAPIClient.FetchMonitors()
	localMonitors, err := fileReader.ReadLocalYaml("test_monitoring.yaml")

	fmt.Println(localMonitors)

	//TODO compare(localMonitors, monitors)
	if err != nil {
		fmt.Println("Client error", err)
		return
	}
	fmt.Println(monitors)
}

func init() {
	rootCmd.AddCommand(fetch)
}
