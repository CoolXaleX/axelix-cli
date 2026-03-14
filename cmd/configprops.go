package cmd

import (
	"github.com/spf13/cobra"
)

var configpropsCmd = &cobra.Command{
	Use:   "configprops",
	Short: "Show @ConfigurationProperties beans",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetConfigProps()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Prefix", "Bean", "Property", "Value"}
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
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configpropsCmd)
}
