package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newBeansCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "beans",
		Short: "List Spring beans",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetBeans()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			rows := make([][]string, 0, len(feed.Beans))
			for _, b := range feed.Beans {
				rows = append(rows, []string{
					b.BeanName, b.Scope, b.ClassName, b.ProxyType,
					fmt.Sprintf("%v", b.IsPrimary),
					fmt.Sprintf("%v", b.IsLazyInit),
				})
			}
			pr.Table([]string{"Name", "Scope", "Class", "ProxyType", "Primary", "Lazy"}, rows)
			return nil
		},
	}
}
