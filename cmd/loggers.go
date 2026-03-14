package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newLoggersCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	loggersCmd := &cobra.Command{Use: "loggers", Short: "Manage application loggers"}

	loggersCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all loggers",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetLoggers()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			names := make([]string, 0, len(feed.Loggers))
			for name := range feed.Loggers {
				names = append(names, name)
			}
			sort.Strings(names)
			rows := make([][]string, 0, len(names))
			for _, name := range names {
				lvl := feed.Loggers[name]
				configured := ""
				if lvl.ConfiguredLevel != nil {
					configured = *lvl.ConfiguredLevel
				}
				rows = append(rows, []string{name, configured, lvl.EffectiveLevel})
			}
			pr.Table([]string{"Logger", "Configured Level", "Effective Level"}, rows)
			return nil
		},
	})

	var getName string
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single logger",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetLogger(getName)
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
				return nil
			}
			var pairs [][2]string
			for k, v := range data {
				pairs = append(pairs, [2]string{k, fmt.Sprintf("%v", v)})
			}
			sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
			pr.KV(pairs)
			return nil
		},
	}
	getCmd.Flags().StringVar(&getName, "name", "", "Logger name")
	getCmd.MarkFlagRequired("name")
	loggersCmd.AddCommand(getCmd)

	var setName, setLevel string
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a logger's level",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.SetLogLevel(setName, setLevel); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success(fmt.Sprintf("Logger %q set to %s.", setName, setLevel))
			return nil
		},
	}
	setCmd.Flags().StringVar(&setName, "name", "", "Logger name")
	setCmd.Flags().StringVar(&setLevel, "level", "", "Log level (e.g. DEBUG, INFO, WARN, ERROR)")
	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("level")
	loggersCmd.AddCommand(setCmd)

	return loggersCmd
}
