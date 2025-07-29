package cmd

import (
	"github.com/spf13/cobra"
	"github.com/walterfan/lazy-rabbit-reminder/internal/api"
	"github.com/walterfan/lazy-rabbit-reminder/internal/task"
 	"os"
	"os/signal"
	"syscall"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP/HTTPS server",
	Long:  "Starts the web server with public/private routes and optional TLS support.",
	Run: func(command *cobra.Command, args []string) {
		logger := GetLogger()

		logger.Info("Starting HTTP service...")
		webService := api.NewWebApiService(logger, rdb)
		go webService.Run()

		logger.Info("Starting Task service...")
		tm := task.NewTaskManager(logger, rdb)
		go tm.CheckTasks()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		logger.Info("Server is running. Press Ctrl+C to stop.")
		<-signalChan
		logger.Info("Received shutdown signal, shutting down.")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
