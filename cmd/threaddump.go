package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var threadDumpCmd = &cobra.Command{
	Use:   "thread-dump",
	Short: "Manage thread dump",
}

var threadDumpGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a thread dump",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetThreadDump()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Name", "ID", "State", "Daemon", "Priority", "Blocked"}
		rows := make([][]string, 0, len(feed.Threads))
		for _, t := range feed.Threads {
			rows = append(rows, []string{
				t.ThreadName,
				strconv.FormatInt(t.ThreadID, 10),
				fmt.Sprintf("%v", t.ThreadState),
				strconv.FormatBool(t.Daemon),
				strconv.Itoa(t.Priority),
				strconv.FormatInt(t.BlockedTime, 10),
			})
		}
		printer.Table(headers, rows)
		return nil
	},
}

var threadDumpEnableCmd = &cobra.Command{
	Use:   "enable-contention",
	Short: "Enable thread contention monitoring",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.EnableThreadContention(); err != nil {
			return err
		}
		printer.Success("Thread contention monitoring enabled.")
		return nil
	},
}

var threadDumpDisableCmd = &cobra.Command{
	Use:   "disable-contention",
	Short: "Disable thread contention monitoring",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DisableThreadContention(); err != nil {
			return err
		}
		printer.Success("Thread contention monitoring disabled.")
		return nil
	},
}

func init() {
	threadDumpCmd.AddCommand(threadDumpGetCmd)
	threadDumpCmd.AddCommand(threadDumpEnableCmd)
	threadDumpCmd.AddCommand(threadDumpDisableCmd)
	rootCmd.AddCommand(threadDumpCmd)
}
