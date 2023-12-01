package cmd

import (
	"fmt"
	"github.com/Trendyol/es-alert-cli/pkg/graceful"
	"github.com/spf13/cobra"
	"log"
	"time"
)

type cli struct {
	command            *cobra.Command
	helpCommand        *cobra.Command
	env                string
	monitoringFilename string
	deleteUntracked    bool
	debug              bool
}

var cliCmd = &cli{
	command: &cobra.Command{
		Use:   "cli",
		Short: "cli app to run",
	},
	helpCommand: &cobra.Command{
		Use:   "help",
		Short: "help for es-alert-cli",
		Run: func(cmd *cobra.Command, args []string) {
			printHelp()
		}},
	env:                "dev",
	monitoringFilename: "monitoring.yaml",
	debug:              false,
}

func init() {
	rootCmd.AddCommand(cliCmd.command)
	rootCmd.AddCommand(cliCmd.helpCommand)

	cliCmd.command.Flags().StringVarP(&cliCmd.env, "env", "e", "dev", "select your env.")
	cliCmd.command.Flags().StringVarP(&cliCmd.monitoringFilename, "filename", "n", "monitoring.yaml", "select your file name")
	cliCmd.command.Flags().BoolVarP(&cliCmd.debug, "debug", "d", false, "enable debugging")

	cliCmd.command.RunE = func(cmd *cobra.Command, args []string) error {
		go func() {
			log.Printf("Cli app running in this scope...")
		}()

		graceful.Shutdown(2 * time.Second)

		return nil
	}
}

func printHelp() {
	fmt.Println("Welcome to es-alert-cli!")
	fmt.Println("This is a CLI application for inserting your kibana alerts as a code.")

	fmt.Println("\nAvailable commands:")
	fmt.Println("cli (Example usage: -- cli)")
	fmt.Println("help (Example usage: -- help)")

	fmt.Println("\nFor more information, visit: https://docs.google.com/document/d/1GLngKFtxt6XqmRTDCj2zGZODSSMZRdNts7fKpYDthGo")
}
