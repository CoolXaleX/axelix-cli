package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/config"
)

func Execute() {
	var flagJSON bool

	rootCmd := &cobra.Command{
		Use:           "axelix",
		Short:         "CLI for Axelix SBS — direct Spring Boot monitoring",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output raw JSON")

	// config subcommands never need a service connection.
	rootCmd.AddCommand(newConfigCmd())

	// Add one subcommand per named service from ~/.axelix/config.json.
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "✗ failed to load config:", err)
		os.Exit(1)
	}

	for name, url := range cfg.Services {
		name, url := name, url // capture loop variables
		cl := client.New(url)

		svcCmd := &cobra.Command{
			Use:   name,
			Short: url,
		}
		svcCmd.AddCommand(
			newBeansCmd(cl, &flagJSON),
			newCachesCmd(cl, &flagJSON),
			newConditionsCmd(cl, &flagJSON),
			newConfigPropsCmd(cl, &flagJSON),
			newDetailsCmd(cl, &flagJSON),
			newEnvCmd(cl, &flagJSON),
			newGcCmd(cl, &flagJSON),
			newHeapDumpCmd(cl),
			newLoggersCmd(cl, &flagJSON),
			newMetadataCmd(cl, &flagJSON),
			newMetricsCmd(cl, &flagJSON),
			newScheduledTasksCmd(cl, &flagJSON),
			newThreadDumpCmd(cl, &flagJSON),
			newTransactionsCmd(cl, &flagJSON),
		)
		rootCmd.AddCommand(svcCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "✗", err)
		os.Exit(1)
	}
}
