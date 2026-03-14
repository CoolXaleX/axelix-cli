package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Show instance details",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetDetails()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		// Print each top-level section as a group of KV pairs.
		sections := make([]string, 0, len(data))
		for k := range data {
			sections = append(sections, k)
		}
		sort.Strings(sections)
		for _, section := range sections {
			fmt.Printf("\n[%s]\n", section)
			switch v := data[section].(type) {
			case map[string]any:
				var pairs [][2]string
				for key, val := range v {
					pairs = append(pairs, [2]string{key, fmt.Sprintf("%v", val)})
				}
				sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
				printer.KV(pairs)
			default:
				fmt.Printf("  %v\n", v)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)
}
