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
		printer.Table([]string{"Manager", "Cache", "Target", "Enabled", "HasStats"}, rows)
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
	cacheClearAll     bool
	cacheClearManager string
	cacheClearCache   string
	cacheClearKey     string
)

var cachesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear caches",
	Long: `Clear caches.

Scope is controlled by the combination of flags:
  --all                         clear every cache across all managers
  --manager M                   clear all caches in manager M
  --manager M --cache C         clear cache C in manager M
  --manager M --cache C --key K evict a single key from cache C`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch {
		case cacheClearAll:
			if err := apiClient.ClearAllCaches(); err != nil {
				return err
			}
		case cacheClearManager != "" && cacheClearCache != "":
			if err := apiClient.ClearCache(cacheClearManager, cacheClearCache, cacheClearKey); err != nil {
				return err
			}
		case cacheClearManager != "":
			if err := apiClient.ClearManagerCaches(cacheClearManager); err != nil {
				return err
			}
		default:
			return fmt.Errorf("specify --all to clear everything, or --manager (and optionally --cache / --key) to narrow the scope")
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
	cachesEnableCmd.Flags().StringVar(&cacheEnableCache, "cache", "", "Cache name (optional — omit to act on all caches in the manager)")
	cachesEnableCmd.MarkFlagRequired("manager")

	cachesDisableCmd.Flags().StringVar(&cacheDisableManager, "manager", "", "Cache manager name")
	cachesDisableCmd.Flags().StringVar(&cacheDisableCache, "cache", "", "Cache name (optional — omit to act on all caches in the manager)")
	cachesDisableCmd.MarkFlagRequired("manager")

	cachesClearCmd.Flags().BoolVar(&cacheClearAll, "all", false, "Clear every cache across all managers")
	cachesClearCmd.Flags().StringVar(&cacheClearManager, "manager", "", "Cache manager name")
	cachesClearCmd.Flags().StringVar(&cacheClearCache, "cache", "", "Cache name")
	cachesClearCmd.Flags().StringVar(&cacheClearKey, "key", "", "Specific key to evict (requires --cache)")

	cachesCmd.AddCommand(cachesListCmd, cachesGetCmd, cachesEnableCmd, cachesDisableCmd, cachesClearCmd)
	rootCmd.AddCommand(cachesCmd)
}
