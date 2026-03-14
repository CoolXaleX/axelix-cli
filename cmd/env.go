package cmd

import (
	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newEnvCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	var pattern string
	cmd := &cobra.Command{
		Use:   "env",
		Short: "Show environment properties",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetEnv(pattern)
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			var rows [][]string
			for _, src := range feed.PropertySources {
				for _, prop := range src.Properties {
					val := ""
					if prop.Value != nil {
						val = *prop.Value
					}
					rows = append(rows, []string{src.Name, prop.Name, val})
				}
			}
			pr.Table([]string{"Source", "Property", "Value"}, rows)
			return nil
		},
	}
	cmd.Flags().StringVar(&pattern, "pattern", "", "Filter properties by pattern")
	return cmd
}
