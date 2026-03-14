package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Inspect application metrics",
}

var metricsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all metric groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetMetrics()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		// Try to parse the metricsGroups structure.
		groups, ok := data["metricsGroups"].([]any)
		if !ok {
			printer.JSON(data)
			return nil
		}
		headers := []string{"Group", "Metric", "Description"}
		var rows [][]string
		for _, g := range groups {
			gMap, ok := g.(map[string]any)
			if !ok {
				continue
			}
			groupName, _ := gMap["groupName"].(string)
			metrics, _ := gMap["metrics"].([]any)
			for _, m := range metrics {
				mMap, ok := m.(map[string]any)
				if !ok {
					continue
				}
				name, _ := mMap["metricName"].(string)
				desc, _ := mMap["description"].(string)
				rows = append(rows, []string{groupName, name, desc})
			}
		}
		printer.Table(headers, rows)
		return nil
	},
}

var (
	metricGetName string
	metricGetTag  string
)

var metricsGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a single metric",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := apiClient.GetMetric(metricGetName, metricGetTag)
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(data)
			return nil
		}
		// Print name and description.
		var pairs [][2]string
		if name, ok := data["name"].(string); ok {
			pairs = append(pairs, [2]string{"name", name})
		}
		if desc, ok := data["description"].(string); ok {
			pairs = append(pairs, [2]string{"description", desc})
		}
		if baseUnit, ok := data["baseUnit"].(string); ok {
			pairs = append(pairs, [2]string{"baseUnit", baseUnit})
		}
		// Print measurements.
		if measurements, ok := data["measurements"].([]any); ok {
			for _, m := range measurements {
				mMap, ok := m.(map[string]any)
				if !ok {
					continue
				}
				stat, _ := mMap["statistic"].(string)
				val := fmt.Sprintf("%v", mMap["value"])
				pairs = append(pairs, [2]string{stat, val})
			}
		}
		printer.KV(pairs)
		return nil
	},
}

func init() {
	metricsGetCmd.Flags().StringVar(&metricGetName, "name", "", "Metric name")
	metricsGetCmd.Flags().StringVar(&metricGetTag, "tag", "", "Tag filter in key:value format")
	metricsGetCmd.MarkFlagRequired("name")

	metricsCmd.AddCommand(metricsListCmd)
	metricsCmd.AddCommand(metricsGetCmd)
	rootCmd.AddCommand(metricsCmd)
}
