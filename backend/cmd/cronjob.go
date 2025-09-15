package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walterfan/lazy-rabbit-reminder/internal/service"
)

var cronjob = &cobra.Command{
	Use:   "check-tasks",
	Short: "Check scheduled tasks and manage task expiry",
	Long:  "Runs cron jobs and monitors task expiries using Redis.",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize JobManager and start checking tasks
		tm := &service.JobManager{}
		tm.CheckTasks()
	},
}

func init() {
	rootCmd.AddCommand(cronjob)
}
