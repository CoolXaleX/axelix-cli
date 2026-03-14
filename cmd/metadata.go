package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newMetadataCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "metadata",
		Short: "Show instance metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetMetadata()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
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
			pr.KV(pairs)
			return nil
		},
	}
}
