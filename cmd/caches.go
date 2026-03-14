package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newCachesCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	cachesCmd := &cobra.Command{Use: "caches", Short: "Manage application caches"}

	cachesCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all caches",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetCaches()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			var rows [][]string
			for _, mgr := range feed.CacheManagers {
				for _, c := range mgr.Caches {
					rows = append(rows, []string{
						mgr.Name, c.Name, c.Target,
						fmt.Sprintf("%v", c.Enabled),
						fmt.Sprintf("%v", c.ContainsStats),
					})
				}
			}
			pr.Table([]string{"Manager", "Cache", "Target", "Enabled", "HasStats"}, rows)
			return nil
		},
	})

	var getManager, getCache string
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single cache",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetSingleCache(getManager, getCache)
			if err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).JSON(data)
			return nil
		},
	}
	getCmd.Flags().StringVar(&getManager, "manager", "", "Cache manager name")
	getCmd.Flags().StringVar(&getCache, "cache", "", "Cache name")
	getCmd.MarkFlagRequired("manager")
	getCmd.MarkFlagRequired("cache")
	cachesCmd.AddCommand(getCmd)

	addToggle := func(use, short string, fn func(manager, cache string) error) {
		var manager, cache string
		c := &cobra.Command{
			Use:   use,
			Short: short,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := fn(manager, cache); err != nil {
					return err
				}
				output.NewPrinter(*jsonFlag).Success(fmt.Sprintf("Cache %sd.", use))
				return nil
			},
		}
		c.Flags().StringVar(&manager, "manager", "", "Cache manager name")
		c.Flags().StringVar(&cache, "cache", "", "Cache name (optional — omit to act on all caches in the manager)")
		c.MarkFlagRequired("manager")
		cachesCmd.AddCommand(c)
	}
	addToggle("enable", "Enable a cache or all caches in a manager", cl.EnableCache)
	addToggle("disable", "Disable a cache or all caches in a manager", cl.DisableCache)

	var clearAll bool
	var clearManager, clearCache, clearKey string
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear caches",
		Long: `Clear caches.

Scope is controlled by the combination of flags:
  --all                         clear every cache across all managers
  --manager M                   clear all caches in manager M
  --manager M --cache C         clear cache C in manager M
  --manager M --cache C --key K evict a single key from cache C`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			switch {
			case clearAll:
				err = cl.ClearAllCaches()
			case clearManager != "" && clearCache != "":
				err = cl.ClearCache(clearManager, clearCache, clearKey)
			case clearManager != "":
				err = cl.ClearManagerCaches(clearManager)
			default:
				return fmt.Errorf("specify --all to clear everything, or --manager (and optionally --cache / --key) to narrow the scope")
			}
			if err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Caches cleared.")
			return nil
		},
	}
	clearCmd.Flags().BoolVar(&clearAll, "all", false, "Clear every cache across all managers")
	clearCmd.Flags().StringVar(&clearManager, "manager", "", "Cache manager name")
	clearCmd.Flags().StringVar(&clearCache, "cache", "", "Cache name")
	clearCmd.Flags().StringVar(&clearKey, "key", "", "Specific key to evict (requires --cache)")
	cachesCmd.AddCommand(clearCmd)

	return cachesCmd
}
