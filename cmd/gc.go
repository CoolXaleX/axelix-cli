package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newGcCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	gcCmd := &cobra.Command{Use: "gc", Short: "Manage garbage collection"}

	gcCmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Show GC log status",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetGCLogStatus()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
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
			pr.KV(pairs)
			return nil
		},
	})

	gcCmd.AddCommand(&cobra.Command{
		Use:   "log-file",
		Short: "Print GC log file content",
		RunE: func(cmd *cobra.Command, args []string) error {
			text, err := cl.GetGCLogFile()
			if err != nil {
				return err
			}
			fmt.Print(text)
			return nil
		},
	})

	gcCmd.AddCommand(&cobra.Command{
		Use:   "trigger",
		Short: "Trigger garbage collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.TriggerGC(); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("GC triggered.")
			return nil
		},
	})

	var logLevel string
	logEnableCmd := &cobra.Command{
		Use:   "log-enable",
		Short: "Enable GC logging at the given level",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.EnableGCLog(logLevel); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("GC logging enabled.")
			return nil
		},
	}
	logEnableCmd.Flags().StringVar(&logLevel, "level", "", "GC log level")
	logEnableCmd.MarkFlagRequired("level")
	gcCmd.AddCommand(logEnableCmd)

	gcCmd.AddCommand(&cobra.Command{
		Use:   "log-disable",
		Short: "Disable GC logging",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.DisableGCLog(); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("GC logging disabled.")
			return nil
		},
	})

	return gcCmd
}
