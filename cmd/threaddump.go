package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newThreadDumpCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	threadCmd := &cobra.Command{Use: "thread-dump", Short: "Manage thread dump"}

	threadCmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a thread dump",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetThreadDump()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			rows := make([][]string, 0, len(feed.Threads))
			for _, t := range feed.Threads {
				rows = append(rows, []string{
					t.ThreadName,
					strconv.FormatInt(t.ThreadID, 10),
					t.ThreadState,
					strconv.FormatBool(t.Daemon),
					strconv.Itoa(t.Priority),
					strconv.FormatInt(t.BlockedTime, 10),
				})
			}
			pr.Table([]string{"Name", "ID", "State", "Daemon", "Priority", "Blocked"}, rows)
			return nil
		},
	})

	threadCmd.AddCommand(&cobra.Command{
		Use:   "enable-contention",
		Short: "Enable thread contention monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.EnableThreadContention(); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Thread contention monitoring enabled.")
			return nil
		},
	})

	threadCmd.AddCommand(&cobra.Command{
		Use:   "disable-contention",
		Short: "Disable thread contention monitoring",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.DisableThreadContention(); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Thread contention monitoring disabled.")
			return nil
		},
	})

	return threadCmd
}
