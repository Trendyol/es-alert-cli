package cmd

import (
	"github.com/labstack/gommon/log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "es-alert-cli",
	Short: "ES Alert CLI",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf(err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(upsertCmd)
}
