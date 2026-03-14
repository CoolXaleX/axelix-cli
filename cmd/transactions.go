package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage transaction monitoring",
}

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List monitored transaction entrypoints",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetTransactions()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Class", "Method", "Executions", "Avg ms", "Max ms", "Median ms"}
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
		printer.Table(headers, rows)
		return nil
	},
}

var transactionsClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all recorded transaction data",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.ClearTransactions(); err != nil {
			return err
		}
		printer.Success("Transaction data cleared.")
		return nil
	},
}

func init() {
	transactionsCmd.AddCommand(transactionsListCmd)
	transactionsCmd.AddCommand(transactionsClearCmd)
	rootCmd.AddCommand(transactionsCmd)
}
