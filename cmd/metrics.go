package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newMetricsCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	metricsCmd := &cobra.Command{Use: "metrics", Short: "Inspect application metrics"}

	metricsCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all metric groups",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetMetrics()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
				return nil
			}
			groups, ok := data["metricsGroups"].([]any)
			if !ok {
				pr.JSON(data)
				return nil
			}
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
			pr.Table([]string{"Group", "Metric", "Description"}, rows)
			return nil
		},
	})

	var getName, getTag string
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single metric",
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := cl.GetMetric(getName, getTag)
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(data)
				return nil
			}
			var pairs [][2]string
			for _, key := range []string{"name", "description", "baseUnit"} {
				if v, ok := data[key].(string); ok {
					pairs = append(pairs, [2]string{key, v})
				}
			}
			if measurements, ok := data["measurements"].([]any); ok {
				for _, m := range measurements {
					mMap, ok := m.(map[string]any)
					if !ok {
						continue
					}
					stat, _ := mMap["statistic"].(string)
					pairs = append(pairs, [2]string{stat, fmt.Sprintf("%v", mMap["value"])})
				}
			}
			pr.KV(pairs)
			return nil
		},
	}
	getCmd.Flags().StringVar(&getName, "name", "", "Metric name")
	getCmd.Flags().StringVar(&getTag, "tag", "", "Tag filter in key:value format")
	getCmd.MarkFlagRequired("name")
	metricsCmd.AddCommand(getCmd)

	return metricsCmd
}
