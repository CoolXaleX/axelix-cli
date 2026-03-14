package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage axelix CLI configuration",
}

var (
	cfgSetURL  string
	cfgSetUser string
	cfgSetPass string
)

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			cfg = &config.Config{}
		}
		if cfgSetURL != "" {
			cfg.URL = cfgSetURL
		}
		if cfgSetUser != "" {
			cfg.Username = cfgSetUser
		}
		if cfgSetPass != "" {
			cfg.Password = cfgSetPass
		}
		if err := config.Save(cfg); err != nil {
			return err
		}
		fmt.Println("✓ Configuration saved.")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		pass := cfg.Password
		if pass != "" {
			pass = "***"
		}
		fmt.Printf("URL:      %s\n", cfg.URL)
		fmt.Printf("Username: %s\n", cfg.Username)
		fmt.Printf("Password: %s\n", pass)
		return nil
	},
}

func init() {
	configSetCmd.Flags().StringVar(&cfgSetURL, "url", "", "SBS app URL")
	configSetCmd.Flags().StringVar(&cfgSetUser, "user", "", "Basic Auth username")
	configSetCmd.Flags().StringVar(&cfgSetPass, "password", "", "Basic Auth password")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
