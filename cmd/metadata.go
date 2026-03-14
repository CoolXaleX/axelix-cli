package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var metadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Show instance metadata",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetMetadata()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		keys := make([]string, 0, len(data))
		for k := range data {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var pairs [][2]string
		for _, k := range keys {
			pairs = append(pairs, [2]string{k, fmt.Sprintf("%v", data[k])})
		}
		printer.KV(pairs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(metadataCmd)
}
