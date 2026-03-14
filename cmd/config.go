package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/config"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage axelix CLI configuration",
	}

	// config add <name> <url>
	configCmd.AddCommand(&cobra.Command{
		Use:   "add <name> <url>",
		Short: "Add or update a named service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name, url := args[0], args[1]
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			cfg.Services[name] = url
			if err := config.Save(cfg); err != nil {
				return err
			}
			output.NewPrinter(false).Success(fmt.Sprintf("service %q added (%s)", name, url))
			return nil
		},
	})

	// config remove <name>
	configCmd.AddCommand(&cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a named service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if _, ok := cfg.Services[name]; !ok {
				return fmt.Errorf("service %q not found", name)
			}
			delete(cfg.Services, name)
			if err := config.Save(cfg); err != nil {
				return err
			}
			output.NewPrinter(false).Success(fmt.Sprintf("service %q removed", name))
			return nil
		},
	})

	// config list
	configCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all configured services",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.Load()
			if err != nil {
				return err
			}
			if len(cfg.Services) == 0 {
				fmt.Println("No services configured. Run 'axelix config add <name> <url>'.")
				return nil
			}
			var rows [][]string
			for name, url := range cfg.Services {
				rows = append(rows, []string{name, url})
			}
			output.NewPrinter(false).Table([]string{"Name", "URL"}, rows)
			return nil
		},
	})

	return configCmd
}
