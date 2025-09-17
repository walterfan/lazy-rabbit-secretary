package cmd

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/walterfan/lazy-rabbit-reminder/internal/jobs"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/database"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/log"
)

var (
	jobName  string
	listJobs bool
	runOnce  bool
)

var cronjob = &cobra.Command{
	Use:   "job [job-name]",
	Short: "Execute job handlers from command line",
	Long: `Execute specific job handlers or manage cron jobs.

Available job handlers:
  - checkTask: Check tasks for reminder generation
  - remindTask: Process due reminders and send notifications  
  - writeBlog: Generate daily technical blog content
  - generateCalendar: Generate daily calendar content

Examples:
  # List all available job handlers
  ./lazy-rabbit-reminder job --list

  # Execute a specific job handler once
  ./lazy-rabbit-reminder job checkTask
  ./lazy-rabbit-reminder job --name writeBlog

  # Start the full cron scheduler (runs continuously)
  ./lazy-rabbit-reminder job --scheduler`,
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize unified logger first
		err := log.InitLogger()
		if err != nil {
			fmt.Printf("Failed to initialize unified logger: %v\n", err)
			os.Exit(1)
		}
		sugar := log.GetLogger()

		// Initialize database (after logger is configured)
		if err := database.InitDB(); err != nil {
			sugar.Fatalf("Failed to initialize database: %v", err)
		}
		db := database.GetDB()

		// Initialize Redis (optional, can be nil for some handlers)
		var rdb *redis.Client
		// You can initialize Redis here if needed for your handlers

		// Initialize JobManager with proper dependencies
		jm := jobs.NewJobManager(sugar.Desugar(), rdb, db)


		// Handle list jobs flag
		if listJobs {
			listAvailableJobs()
			return
		}

		// Determine job name from args or flag
		targetJob := jobName
		if len(args) > 0 {
			targetJob = args[0]
		}

		// Handle scheduler mode (runs continuously)
		if targetJob == "scheduler" || (!runOnce && targetJob == "") {
			sugar.Info("Starting cron scheduler...")
			jm.CheckTasks() // This starts the continuous scheduler
			return
		}

		// Execute specific job handler
		if targetJob != "" {
			executeJobHandler(jm, targetJob, sugar)
		} else {
			cmd.Help()
		}
	},
}

// executeJobHandler executes a specific job handler by name
func executeJobHandler(jm *jobs.JobManager, jobName string, logger *zap.SugaredLogger) {
	// Map of available job handlers
	availableJobs := map[string]string{
		"checkTask":        "Check tasks for reminder generation",
		"remindTask":       "Process due reminders and send notifications",
		"writeBlog":        "Generate daily technical blog content",
		"generateCalendar": "Generate daily calendar content",
	}

	description, exists := availableJobs[jobName]
	if !exists {
		logger.Errorf("Unknown job handler: %s", jobName)
		logger.Info("Available job handlers:")
		for name, desc := range availableJobs {
			logger.Infof("  %s: %s", name, desc)
		}
		os.Exit(1)
	}

	logger.Infof("Executing job handler: %s (%s)", jobName, description)

	// Execute the job handler directly
	err := jm.ExecuteFunction(jobName, "")
	if err != nil {
		logger.Errorf("Failed to execute job handler %s: %v", jobName, err)
		os.Exit(1)
	}

	logger.Infof("Successfully completed job handler: %s", jobName)
}

// listAvailableJobs displays all available job handlers
func listAvailableJobs() {
	fmt.Println("Available job handlers:")
	fmt.Println()

	jobs := []struct {
		Name        string
		Description string
		Example     string
	}{
		{"checkTask", "Check tasks for reminder generation", "./lazy-rabbit-reminder job checkTask"},
		{"remindTask", "Process due reminders and send notifications", "./lazy-rabbit-reminder job remindTask"},
		{"writeBlog", "Generate daily technical blog content", "./lazy-rabbit-reminder job writeBlog"},
		{"generateCalendar", "Generate daily calendar content", "./lazy-rabbit-reminder job generateCalendar"},
		{"scheduler", "Start continuous cron scheduler", "./lazy-rabbit-reminder job scheduler"},
	}

	for _, job := range jobs {
		fmt.Printf("  %-18s %s\n", job.Name, job.Description)
		fmt.Printf("  %-18s Example: %s\n", "", job.Example)
		fmt.Println()
	}
}

func init() {
	// Add flags
	cronjob.Flags().StringVarP(&jobName, "name", "n", "", "Name of the job handler to execute")
	cronjob.Flags().BoolVarP(&listJobs, "list", "l", false, "List all available job handlers")
	cronjob.Flags().BoolVarP(&runOnce, "once", "o", true, "Run job once and exit (default: true)")

	rootCmd.AddCommand(cronjob)
}
