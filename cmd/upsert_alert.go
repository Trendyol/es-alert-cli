package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/internal"

	"github.com/Trendyol/es-alert-cli/pkg/errs"

	"github.com/Trendyol/es-alert-cli/pkg/client"
	"github.com/Trendyol/es-alert-cli/pkg/reader"
	"github.com/spf13/cobra"
)

var upsertCmd = &cobra.Command{
	Use:   "upsert",
	Short: "Upsert all monitors",
	Long:  `Upsert command will push your monitoring yaml to remote if any change exists.`,
	Run:   upsertMonitors,
}

func init() {
	upsertCmd.Flags().StringP("cluster", "c", "", "select your cluster ip to update.")
	upsertCmd.Flags().StringP("filename", "n", "", "select your monitoring file name.")
}

func upsertMonitors(cmd *cobra.Command, args []string) {
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

	monitorService := internal.NewMonitorService(fileReader, esAPIClient)
	err = monitorService.Upsert(filename, cliCmd.deleteUntracked)

	if err != nil {
		fmt.Printf("Upsert operation failed, err: %s", err.Error())
		return
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
