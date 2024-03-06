package cmd

import (
	"os"

	"github.com/labstack/gommon/log"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "es-alert-cli",
	Short: "ES Alert CLI",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.AddCommand(upsertCmd)
}
