package cmd

import (
	"github.com/Trendyol/es-alert-cli/pkg/graceful"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type cli struct {
	command *cobra.Command
	env     string
	debug   bool
}

var apiCmd = &cli{
	command: &cobra.Command{
		Use:   "cli",
		Short: "CLI App",
	},
	env:   "dev",
	debug: false,
}

func init() {
	rootCmd.AddCommand(apiCmd.command)

	apiCmd.command.Flags().StringVarP(&apiCmd.env, "env", "e", "dev", "select your env.")
	apiCmd.command.Flags().BoolVarP(&apiCmd.debug, "debug", "d", false, "enable debugging")

	apiCmd.command.RunE = func(cmd *cobra.Command, args []string) error {
		go func() {
			log.Printf("Cli app running in this scope...")
		}()

		graceful.Shutdown(2 * time.Second)

		return nil
	}
}
