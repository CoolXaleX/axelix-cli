package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/models"
)

var scheduledTasksCmd = &cobra.Command{
	Use:   "scheduled-tasks",
	Short: "Manage scheduled tasks",
}

var scheduledTasksListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all scheduled tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		feed, err := apiClient.GetScheduledTasks()
		if err != nil {
			return err
		}
		if printer.IsJSON() {
			printer.JSON(feed)
			return nil
		}
		headers := []string{"Type", "Target/Trigger", "Schedule", "Enabled", "Next Exec", "Last Status"}
		var rows [][]string

		for _, t := range feed.Cron {
			rows = append(rows, taskRow("cron", t.Runnable.Target, t.Expression, t.Enabled, t.NextExecution, t.LastExecution))
		}
		for _, t := range feed.FixedDelay {
			schedule := fmt.Sprintf("delay %vms", t.Interval)
			rows = append(rows, taskRow("fixedDelay", t.Runnable.Target, schedule, t.Enabled, t.NextExecution, t.LastExecution))
		}
		for _, t := range feed.FixedRate {
			schedule := fmt.Sprintf("rate %vms", t.Interval)
			rows = append(rows, taskRow("fixedRate", t.Runnable.Target, schedule, t.Enabled, t.NextExecution, t.LastExecution))
		}
		for _, t := range feed.Custom {
			rows = append(rows, taskRow("custom", t.Runnable.Target, t.Trigger, t.Enabled, t.NextExecution, t.LastExecution))
		}
		printer.Table(headers, rows)
		return nil
	},
}

func taskRow(taskType, target, schedule string, enabled bool, next *models.NextExecution, last *models.LastExecution) []string {
	nextTime := ""
	if next != nil {
		nextTime = next.Time
	}
	lastStatus := ""
	if last != nil {
		lastStatus = last.Status
	}
	return []string{taskType, target, schedule, strconv.FormatBool(enabled), nextTime, lastStatus}
}

var (
	taskToggleTrigger string
	taskToggleForce   bool
)

var scheduledTasksEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable a scheduled task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.EnableScheduledTask(taskToggleTrigger, taskToggleForce); err != nil {
			return err
		}
		printer.Success("Scheduled task enabled.")
		return nil
	},
}

var scheduledTasksDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable a scheduled task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.DisableScheduledTask(taskToggleTrigger, taskToggleForce); err != nil {
			return err
		}
		printer.Success("Scheduled task disabled.")
		return nil
	},
}

var taskExecuteTrigger string

var scheduledTasksExecuteCmd = &cobra.Command{
	Use:   "execute",
	Short: "Immediately execute a scheduled task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.ExecuteScheduledTask(taskExecuteTrigger); err != nil {
			return err
		}
		printer.Success("Scheduled task executed.")
		return nil
	},
}

var (
	taskCronTrigger string
	taskCronExpr    string
)

var scheduledTasksSetCronCmd = &cobra.Command{
	Use:   "set-cron",
	Short: "Set the cron expression of a scheduled task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.SetTaskCronExpression(taskCronTrigger, taskCronExpr); err != nil {
			return err
		}
		printer.Success("Cron expression updated.")
		return nil
	},
}

var (
	taskIntervalTrigger   string
	taskIntervalMs        int64
)

var scheduledTasksSetIntervalCmd = &cobra.Command{
	Use:   "set-interval",
	Short: "Set the interval of a scheduled task",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := apiClient.SetTaskInterval(taskIntervalTrigger, taskIntervalMs); err != nil {
			return err
		}
		printer.Success("Task interval updated.")
		return nil
	},
}

func init() {
	scheduledTasksEnableCmd.Flags().StringVar(&taskToggleTrigger, "trigger", "", "Task trigger identifier")
	scheduledTasksEnableCmd.Flags().BoolVar(&taskToggleForce, "force", false, "Force enable")
	scheduledTasksEnableCmd.MarkFlagRequired("trigger")

	scheduledTasksDisableCmd.Flags().StringVar(&taskToggleTrigger, "trigger", "", "Task trigger identifier")
	scheduledTasksDisableCmd.Flags().BoolVar(&taskToggleForce, "force", false, "Force disable")
	scheduledTasksDisableCmd.MarkFlagRequired("trigger")

	scheduledTasksExecuteCmd.Flags().StringVar(&taskExecuteTrigger, "trigger", "", "Task trigger identifier")
	scheduledTasksExecuteCmd.MarkFlagRequired("trigger")

	scheduledTasksSetCronCmd.Flags().StringVar(&taskCronTrigger, "trigger", "", "Task trigger identifier")
	scheduledTasksSetCronCmd.Flags().StringVar(&taskCronExpr, "cron", "", "Cron expression")
	scheduledTasksSetCronCmd.MarkFlagRequired("trigger")
	scheduledTasksSetCronCmd.MarkFlagRequired("cron")

	scheduledTasksSetIntervalCmd.Flags().StringVar(&taskIntervalTrigger, "trigger", "", "Task trigger identifier")
	scheduledTasksSetIntervalCmd.Flags().Int64Var(&taskIntervalMs, "interval", 0, "Interval in milliseconds")
	scheduledTasksSetIntervalCmd.MarkFlagRequired("trigger")
	scheduledTasksSetIntervalCmd.MarkFlagRequired("interval")

	scheduledTasksCmd.AddCommand(scheduledTasksListCmd)
	scheduledTasksCmd.AddCommand(scheduledTasksEnableCmd)
	scheduledTasksCmd.AddCommand(scheduledTasksDisableCmd)
	scheduledTasksCmd.AddCommand(scheduledTasksExecuteCmd)
	scheduledTasksCmd.AddCommand(scheduledTasksSetCronCmd)
	scheduledTasksCmd.AddCommand(scheduledTasksSetIntervalCmd)
	rootCmd.AddCommand(scheduledTasksCmd)
}
