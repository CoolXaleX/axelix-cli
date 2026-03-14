package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var gcCmd = &cobra.Command{
	Use:   "gc",
	Short: "Manage garbage collection",
}

var gcStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show GC log status",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetGCLogStatus()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var pairs [][2]string
		for _, k := range keys {
			pairs = append(pairs, [2]string{k, fmt.Sprintf("%v", data[k])})
		}
		printer.KV(pairs)
		return nil
	},
}

var gcLogFileCmd = &cobra.Command{
	Use:   "log-file",
	Short: "Print GC log file content",
	RunE: func(cmd *cobra.Command, args []string) error {
		text, err := apiClient.GetGCLogFile()
		if err != nil {
			return err
		}
		fmt.Print(text)
		return nil
	},
}

var gcTriggerCmd = &cobra.Command{
	Use:   "trigger",
	Short: "Trigger garbage collection",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.TriggerGC(); err != nil {
			return err
		}
		printer.Success("GC triggered.")
		return nil
	},
}

var gcLogEnableLevel string

var gcLogEnableCmd = &cobra.Command{
	Use:   "log-enable",
	Short: "Enable GC logging at the given level",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.EnableGCLog(gcLogEnableLevel); err != nil {
			return err
		}
		printer.Success("GC logging enabled.")
		return nil
	},
}

var gcLogDisableCmd = &cobra.Command{
	Use:   "log-disable",
	Short: "Disable GC logging",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DisableGCLog(); err != nil {
			return err
		}
		printer.Success("GC logging disabled.")
		return nil
	},
}

func init() {
	gcLogEnableCmd.Flags().StringVar(&gcLogEnableLevel, "level", "", "GC log level")
	gcLogEnableCmd.MarkFlagRequired("level")

	gcCmd.AddCommand(gcStatusCmd)
	gcCmd.AddCommand(gcLogFileCmd)
	gcCmd.AddCommand(gcTriggerCmd)
	gcCmd.AddCommand(gcLogEnableCmd)
	gcCmd.AddCommand(gcLogDisableCmd)
	rootCmd.AddCommand(gcCmd)
}
