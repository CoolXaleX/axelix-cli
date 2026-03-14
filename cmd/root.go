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
	flagURL      string
	flagUser     string
	flagPassword string
	flagJSON     bool
	printer      *output.Printer
	apiClient    *client.Client
)

var rootCmd = &cobra.Command{
	Use:           "axelix",
	Short:         "CLI for Axelix SBS — direct Spring Boot monitoring",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip for config subcommands.
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}
		cfg := config.Resolve(flagURL, flagUser, flagPassword)
		if cfg.URL == "" {
			return fmt.Errorf("no URL — use --url or run 'axelix config set --url http://...'")
		}
		apiClient = client.New(cfg.URL, cfg.Username, cfg.Password)
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
	rootCmd.PersistentFlags().StringVar(&flagURL, "url", "", "SBS app URL (e.g. http://localhost:8080)")
	rootCmd.PersistentFlags().StringVar(&flagUser, "user", "", "Basic Auth username")
	rootCmd.PersistentFlags().StringVar(&flagPassword, "password", "", "Basic Auth password")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output raw JSON")
}
