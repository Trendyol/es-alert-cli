package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/spf13/cobra"
)

var fetch = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch all alerts",
	Long:  `Long text example`,
	Run:   fetchAlerts,
}

func fetchAlerts(cmd *cobra.Command, args []string) {
	// TODO make sure args are validated & mapped appropriately
	esAPIClient, err := client.NewElasticsearchAPI(args)
	if err != nil {
		fmt.Println("we have an error", err)
		return
	}
	monitors, err := esAPIClient.FetchMonitors()
	if err != nil {
		fmt.Println("Client error", err)
		return
	}
	fmt.Println(monitors)
}

func init() {
	rootCmd.AddCommand(fetch)
}
