package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/config"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage axelix CLI configuration",
}

// config add <name> <url>
var configAddCmd = &cobra.Command{
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
		// If this is the first service, make it current automatically.
		if cfg.Current == "" {
			cfg.Current = name
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		p := output.NewPrinter(false)
		p.Success(fmt.Sprintf("service %q added (%s)", name, url))
		if cfg.Current == name {
			p.Line(fmt.Sprintf("  → now using %q", name))
		}
		return nil
	},
}

// config remove <name>
var configRemoveCmd = &cobra.Command{
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
		if cfg.Current == name {
			cfg.Current = ""
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		output.NewPrinter(false).Success(fmt.Sprintf("service %q removed", name))
		return nil
	},
}

// config use <name>
var configUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Switch the active service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		if _, ok := cfg.Services[name]; !ok {
			return fmt.Errorf("service %q not found — run 'axelix config list' to see available services", name)
		}
		cfg.Current = name
		if err := config.Save(cfg); err != nil {
			return err
		}
		output.NewPrinter(false).Success(fmt.Sprintf("now using %q (%s)", name, cfg.Services[name]))
		return nil
	},
}

// config list
var configListCmd = &cobra.Command{
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
		p := output.NewPrinter(false)
		var rows [][]string
		for name, url := range cfg.Services {
			active := ""
			if name == cfg.Current {
				active = "✓"
			}
			rows = append(rows, []string{active, name, url})
		}
		p.Table([]string{"", "Name", "URL"}, rows)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configListCmd)
	rootCmd.AddCommand(configCmd)
}
