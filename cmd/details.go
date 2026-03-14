package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newDetailsCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "details",
		Short: "Show instance details",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetDetails()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
				return nil
			}
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
					pr.KV(pairs)
				default:
					fmt.Printf("  %v\n", v)
				}
			}
			return nil
		},
	}
}
