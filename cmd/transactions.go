package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newTransactionsCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	txCmd := &cobra.Command{Use: "transactions", Short: "Manage transaction monitoring"}

	txCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List monitored transaction entrypoints",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetTransactions()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			rows := make([][]string, 0, len(feed.Entrypoints))
			for _, ep := range feed.Entrypoints {
				rows = append(rows, []string{
					ep.ClassName,
					ep.MethodName,
					strconv.Itoa(len(ep.Executions)),
					fmt.Sprintf("%d", ep.ExecutionStats.AverageDurationMs),
					fmt.Sprintf("%d", ep.ExecutionStats.MaxDurationMs),
					fmt.Sprintf("%d", ep.ExecutionStats.MedianDurationMs),
				})
			}
			pr.Table([]string{"Class", "Method", "Executions", "Avg ms", "Max ms", "Median ms"}, rows)
			return nil
		},
	})

	txCmd.AddCommand(&cobra.Command{
		Use:   "clear",
		Short: "Clear all recorded transaction data",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.ClearTransactions(); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Transaction data cleared.")
			return nil
		},
	})

	return txCmd
}
