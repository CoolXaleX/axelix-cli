package cmd

import (
	"github.com/spf13/cobra"
)

var conditionsNegative bool

var conditionsCmd = &cobra.Command{
	Use:   "conditions",
	Short: "Show auto-configuration conditions",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetConditions()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Class", "Method", "Condition", "Message"}
		var rows [][]string
		if conditionsNegative {
			for _, nc := range feed.NegativeMatches {
				method := ""
				if nc.MethodName != nil {
					method = *nc.MethodName
				}
				for _, m := range nc.NotMatched {
					rows = append(rows, []string{nc.ClassName, method, m.Condition, m.Message})
				}
			}
		} else {
			for _, pc := range feed.PositiveMatches {
				method := ""
				if pc.MethodName != nil {
					method = *pc.MethodName
				}
				for _, m := range pc.Matched {
					rows = append(rows, []string{pc.ClassName, method, m.Condition, m.Message})
				}
			}
		}
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	conditionsCmd.Flags().BoolVar(&conditionsNegative, "negative", false, "Show negative matches instead of positive")
	rootCmd.AddCommand(conditionsCmd)
}
