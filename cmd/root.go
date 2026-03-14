package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/config"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

var (
	flagURL     string
	flagService string
	flagJSON    bool
	printer     *output.Printer
	apiClient   *client.Client
)

var rootCmd = &cobra.Command{
	Use:           "axelix",
	Short:         "CLI for Axelix SBS — direct Spring Boot monitoring",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// config subcommands manage their own state — no client needed.
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}
		url, err := config.Resolve(flagURL, flagService)
		if err != nil {
			return err
		}
		apiClient = client.New(url)
		printer = output.NewPrinter(flagJSON)
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "SBS app URL (overrides saved service)")
	rootCmd.PersistentFlags().StringVar(&flagService, "service", "", "Named service from config (e.g. --service prod)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output raw JSON")
}
