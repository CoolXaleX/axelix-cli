package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var beansCmd = &cobra.Command{
	Use:   "beans",
	Short: "List Spring beans",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetBeans()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Name", "Scope", "Class", "ProxyType", "Primary", "Lazy"}
		rows := make([][]string, 0, len(feed.Beans))
		for _, b := range feed.Beans {
			rows = append(rows, []string{
				b.BeanName,
				b.Scope,
				b.ClassName,
				b.ProxyType,
				fmt.Sprintf("%v", b.IsPrimary),
				fmt.Sprintf("%v", b.IsLazyInit),
			})
		}
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(beansCmd)
}
