package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/internal"
	"github.com/labstack/gommon/log"

	"github.com/Trendyol/es-alert-cli/pkg/errs"

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
	cluster, filename, ok := getFlags(cmd)
	if !ok {
		return
	}

	log.Info(fmt.Sprintf("Cli will connect to %s\n", cluster))

	monitorService, err := internal.NewMonitorService(cluster)
	if err != nil {
		log.Info(fmt.Sprintf(", err: %s", err.Error()))
		return
	}

	response, err := monitorService.Upsert(filename, cliCmd.deleteUntracked)
	if err != nil {
		log.Info(fmt.Sprintf("Upsert operation failed, err: %s", err.Error()))
		return
	}

	for _, monitorResponse := range response {
		log.Info(fmt.Sprintf("monitor with id: %s is created", monitorResponse.ID))
	}

	log.Info("All process completed.")
}

func getFlags(cmd *cobra.Command) (string, string, bool) {
	cluster, err := cmd.Flags().GetString("cluster")
	if errs.LogError(err, "err getting cluster parameter") {
		return "", "", false
	}
	filename, err := cmd.Flags().GetString("filename")
	if errs.LogError(err, "err getting filename parameter") {
		return "", "", false
	}
	return cluster, filename, true
}
