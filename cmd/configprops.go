package cmd

import (
	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newConfigPropsCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "configprops",
		Short: "Show @ConfigurationProperties beans",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetConfigProps()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			var rows [][]string
			for _, b := range feed.Beans {
				for _, kv := range b.Properties {
					val := ""
					if kv.Value != nil {
						val = *kv.Value
					}
					rows = append(rows, []string{b.Prefix, b.BeanName, kv.Key, val})
				}
			}
			pr.Table([]string{"Prefix", "Bean", "Property", "Value"}, rows)
			return nil
		},
	}
}
