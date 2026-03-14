package cmd

import (
	"github.com/spf13/cobra"
)

var envPattern string

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Show environment properties",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetEnv(envPattern)
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Source", "Property", "Value"}
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
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	envCmd.Flags().StringVar(&envPattern, "pattern", "", "Filter properties by pattern")
	rootCmd.AddCommand(envCmd)
}
