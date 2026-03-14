package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

var loggersCmd = &cobra.Command{
	Use:   "loggers",
	Short: "Manage application loggers",
}

var loggersListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all loggers",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetLoggers()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Logger", "Configured Level", "Effective Level"}
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
		printer.Table(headers, rows)
		return nil
	},
}

var loggerGetName string

var loggersGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single logger",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetLogger(loggerGetName)
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		var pairs [][2]string
		for k, v := range data {
			pairs = append(pairs, [2]string{k, fmt.Sprintf("%v", v)})
		}
		sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
		printer.KV(pairs)
		return nil
	},
}

var (
	loggerSetName  string
	loggerSetLevel string
)

var loggersSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a logger's level",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.SetLogLevel(loggerSetName, loggerSetLevel); err != nil {
			return err
		}
		printer.Success(fmt.Sprintf("Logger %q set to %s.", loggerSetName, loggerSetLevel))
		return nil
	},
}

func init() {
	loggersGetCmd.Flags().StringVar(&loggerGetName, "name", "", "Logger name")
	loggersGetCmd.MarkFlagRequired("name")

	loggersSetCmd.Flags().StringVar(&loggerSetName, "name", "", "Logger name")
	loggersSetCmd.Flags().StringVar(&loggerSetLevel, "level", "", "Log level (e.g. DEBUG, INFO, WARN, ERROR)")
	loggersSetCmd.MarkFlagRequired("name")
	loggersSetCmd.MarkFlagRequired("level")

	loggersCmd.AddCommand(loggersListCmd)
	loggersCmd.AddCommand(loggersGetCmd)
	loggersCmd.AddCommand(loggersSetCmd)
	rootCmd.AddCommand(loggersCmd)
}
