package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/axelixlabs/axelix-cli/internal/client"
	"github.com/axelixlabs/axelix-cli/internal/models"
	"github.com/axelixlabs/axelix-cli/internal/output"
)

func newScheduledTasksCmd(cl *client.Client, jsonFlag *bool) *cobra.Command {
	tasksCmd := &cobra.Command{Use: "scheduled-tasks", Short: "Manage scheduled tasks"}

	tasksCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all scheduled tasks",
		RunE: func(cmd *cobra.Command, args []string) error {
			feed, err := cl.GetScheduledTasks()
			if err != nil {
				return err
			}
			pr := output.NewPrinter(*jsonFlag)
			if pr.IsJSON() {
				pr.JSON(feed)
				return nil
			}
			headers := []string{"Type", "Target/Trigger", "Schedule", "Enabled", "Next Exec", "Last Status"}
			var rows [][]string
			for _, t := range feed.Cron {
				rows = append(rows, taskRow("cron", t.Runnable.Target, t.Expression, t.Enabled, t.NextExecution, t.LastExecution))
			}
			for _, t := range feed.FixedDelay {
				rows = append(rows, taskRow("fixedDelay", t.Runnable.Target, fmt.Sprintf("delay %vms", t.Interval), t.Enabled, t.NextExecution, t.LastExecution))
			}
			for _, t := range feed.FixedRate {
				rows = append(rows, taskRow("fixedRate", t.Runnable.Target, fmt.Sprintf("rate %vms", t.Interval), t.Enabled, t.NextExecution, t.LastExecution))
			}
			for _, t := range feed.Custom {
				rows = append(rows, taskRow("custom", t.Runnable.Target, t.Trigger, t.Enabled, t.NextExecution, t.LastExecution))
			}
			pr.Table(headers, rows)
			return nil
		},
	})

	addToggleCmd := func(use, short string, fn func(trigger string, force bool) error) {
		var trigger string
		var force bool
		c := &cobra.Command{
			Use:   use,
			Short: short,
			RunE: func(cmd *cobra.Command, args []string) error {
				if err := fn(trigger, force); err != nil {
					return err
				}
				output.NewPrinter(*jsonFlag).Success(fmt.Sprintf("Scheduled task %sd.", use))
				return nil
			},
		}
		c.Flags().StringVar(&trigger, "trigger", "", "Task trigger identifier")
		c.Flags().BoolVar(&force, "force", false, fmt.Sprintf("Force %s", use))
		c.MarkFlagRequired("trigger")
		tasksCmd.AddCommand(c)
	}
	addToggleCmd("enable", "Enable a scheduled task", cl.EnableScheduledTask)
	addToggleCmd("disable", "Disable a scheduled task", cl.DisableScheduledTask)

	var execTrigger string
	execCmd := &cobra.Command{
		Use:   "execute",
		Short: "Immediately execute a scheduled task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.ExecuteScheduledTask(execTrigger); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Scheduled task executed.")
			return nil
		},
	}
	execCmd.Flags().StringVar(&execTrigger, "trigger", "", "Task trigger identifier")
	execCmd.MarkFlagRequired("trigger")
	tasksCmd.AddCommand(execCmd)

	var cronTrigger, cronExpr string
	cronCmd := &cobra.Command{
		Use:   "set-cron",
		Short: "Set the cron expression of a scheduled task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.SetTaskCronExpression(cronTrigger, cronExpr); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Cron expression updated.")
			return nil
		},
	}
	cronCmd.Flags().StringVar(&cronTrigger, "trigger", "", "Task trigger identifier")
	cronCmd.Flags().StringVar(&cronExpr, "cron", "", "Cron expression")
	cronCmd.MarkFlagRequired("trigger")
	cronCmd.MarkFlagRequired("cron")
	tasksCmd.AddCommand(cronCmd)

	var intervalTrigger string
	var intervalMs int64
	intervalCmd := &cobra.Command{
		Use:   "set-interval",
		Short: "Set the interval of a scheduled task",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cl.SetTaskInterval(intervalTrigger, intervalMs); err != nil {
				return err
			}
			output.NewPrinter(*jsonFlag).Success("Task interval updated.")
			return nil
		},
	}
	intervalCmd.Flags().StringVar(&intervalTrigger, "trigger", "", "Task trigger identifier")
	intervalCmd.Flags().Int64Var(&intervalMs, "interval", 0, "Interval in milliseconds")
	intervalCmd.MarkFlagRequired("trigger")
	intervalCmd.MarkFlagRequired("interval")
	tasksCmd.AddCommand(intervalCmd)

	return tasksCmd
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
