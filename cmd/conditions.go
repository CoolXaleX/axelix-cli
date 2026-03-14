package cmd

import (
	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newConditionsCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	var negative bool
	cmd := &cobra.Command{
		Use:   "conditions",
		Short: "Show auto-configuration conditions",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetConditions()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			headers := []string{"Class", "Method", "Condition", "Message"}
			var rows [][]string
			if negative {
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
			pr.Table(headers, rows)
			return nil
		},
	}
	cmd.Flags().BoolVar(&negative, "negative", false, "Show negative matches instead of positive")
	return cmd
}
