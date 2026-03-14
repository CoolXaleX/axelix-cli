package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cachesCmd = &cobra.Command{
	Use:   "caches",
	Short: "Manage application caches",
}

var cachesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all caches",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetCaches()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Manager", "Cache", "Target", "Enabled", "HasStats"}
		var rows [][]string
		for _, mgr := range feed.CacheManagers {
			for _, c := range mgr.Caches {
				rows = append(rows, []string{
					mgr.Name,
					c.Name,
					c.Target,
					fmt.Sprintf("%v", c.Enabled),
					fmt.Sprintf("%v", c.ContainsStats),
				})
			}
		}
		printer.Table(headers, rows)
		return nil
	},
}

var (
	cacheGetManager string
	cacheGetCache   string
)

var cachesGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetSingleCache(cacheGetManager, cacheGetCache)
		if err != nil {
			return err
		}
		printer.JSON(data)
		return nil
	},
}

var (
	cacheEnableManager string
	cacheEnableCache   string
)

var cachesEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a cache or all caches in a manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.EnableCache(cacheEnableManager, cacheEnableCache); err != nil {
			return err
		}
		printer.Success("Cache enabled.")
		return nil
	},
}

var (
	cacheDisableManager string
	cacheDisableCache   string
)

var cachesDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a cache or all caches in a manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DisableCache(cacheDisableManager, cacheDisableCache); err != nil {
			return err
		}
		printer.Success("Cache disabled.")
		return nil
	},
}

var (
	cacheClearManager string
	cacheClearCache   string
	cacheClearKey     string
)

var cachesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear caches",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.ClearCaches(cacheClearManager, cacheClearCache, cacheClearKey); err != nil {
			return err
		}
		printer.Success("Caches cleared.")
		return nil
	},
}

func init() {
	cachesGetCmd.Flags().StringVar(&cacheGetManager, "manager", "", "Cache manager name")
	cachesGetCmd.Flags().StringVar(&cacheGetCache, "cache", "", "Cache name")
	cachesGetCmd.MarkFlagRequired("manager")
	cachesGetCmd.MarkFlagRequired("cache")

	cachesEnableCmd.Flags().StringVar(&cacheEnableManager, "manager", "", "Cache manager name")
	cachesEnableCmd.Flags().StringVar(&cacheEnableCache, "cache", "", "Cache name (optional)")
	cachesEnableCmd.MarkFlagRequired("manager")

	cachesDisableCmd.Flags().StringVar(&cacheDisableManager, "manager", "", "Cache manager name")
	cachesDisableCmd.Flags().StringVar(&cacheDisableCache, "cache", "", "Cache name (optional)")
	cachesDisableCmd.MarkFlagRequired("manager")

	cachesClearCmd.Flags().StringVar(&cacheClearManager, "manager", "", "Cache manager name (optional)")
	cachesClearCmd.Flags().StringVar(&cacheClearCache, "cache", "", "Cache name (optional)")
	cachesClearCmd.Flags().StringVar(&cacheClearKey, "key", "", "Cache key (optional)")

	cachesCmd.AddCommand(cachesListCmd)
	cachesCmd.AddCommand(cachesGetCmd)
	cachesCmd.AddCommand(cachesEnableCmd)
	cachesCmd.AddCommand(cachesDisableCmd)
	cachesCmd.AddCommand(cachesClearCmd)
	rootCmd.AddCommand(cachesCmd)
}
