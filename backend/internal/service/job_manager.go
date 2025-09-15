package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Define a Task structure with the 'function' and 'command' fields
type Task struct {
	Name       string `yaml:"name"`
	Schedule   string `yaml:"schedule"`
	Function   string `yaml:"function"`
	Parameters string `yaml:"parameters"`
	Deadline   string `yaml:"deadline"`
}

type Config struct {
	Tasks []Task `yaml:"tasks"`
}

// JobManager struct to encapsulate task-related functions
type JobManager struct {
	config *Config
	logger *zap.SugaredLogger
	ctx    context.Context
	rdb    *redis.Client
}

func NewJobManager(logger *zap.Logger, redisClient *redis.Client) *JobManager {

	return &JobManager{
		config: nil,
		logger: logger.Sugar(),
		ctx:    context.Background(),
		rdb:    redisClient,
	}
}

// Load configuration from the YAML file
func (tm *JobManager) loadConfig() error {

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}

	tm.config = &config
	return nil
}

// Define the function to checkTask the URL
func (tm *JobManager) checkTask(url string) {
	// Perform a basic HTTP GET request to check if the URL is reachable
	resp, err := http.Get(url)
	if err != nil {
		tm.logger.Infof("Error checking URL %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	tm.logger.Infof("Checked URL %s with status code: %d", url, resp.StatusCode)
}

// Function to map the function name to the actual Go function
func (tm *JobManager) executeFunction(functionName, param string) {
	// Remove the "()" from the function name
	functionName = strings.TrimSuffix(functionName, "()")

	tm.logger.Infof("Executing function '%s' with parameter '%s'", functionName, param)
	// Map the function name to the corresponding function
	if functionName == "checkTask" {
		tm.checkTask(param)
	} else {
		tm.logger.Infof("No predefined function: %s", functionName)
	}

	handler, exists := JobHandlers[functionName]
	if !exists {
		tm.logger.Infof("No plugin found for: %s", functionName)
		return
	}

	if err := handler.Execute(param); err != nil {
		tm.logger.Infof("Error executing plugin %s: %v", functionName, err)
	}
}

// Read unfinished tasks from Redis
func (tm *JobManager) readUnfinishedTasks() ([]string, error) {
	keys, err := tm.rdb.Keys(tm.ctx, "task:*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Publish task expiry event to Redis
func (tm *JobManager) publishTaskExpiryEvent(taskID string) error {
	return tm.rdb.Publish(tm.ctx, "task_expiry_channel", taskID).Err()
}

// Check task expiry and publish events
func (tm *JobManager) checkTaskExpiry() {
	tasks, err := tm.readUnfinishedTasks()
	if err != nil {
		tm.logger.Infof("Failed to read unfinished tasks: %v", err)
		return
	}

	for _, taskKey := range tasks {
		// Assume task expiration time is stored in Redis with key "task:<taskID>:expiry"
		expiryTime, err := tm.rdb.Get(tm.ctx, taskKey).Int64()
		if err != nil {
			tm.logger.Infof("Failed to get expiry time for task %s: %v", taskKey, err)
			continue
		}

		tm.logger.Infof("task %s expiry time is %d vs. %d", taskKey, expiryTime, time.Now().Unix())

		if time.Now().Unix() >= expiryTime {
			taskID := strings.TrimPrefix(taskKey, "task:")
			if err := tm.publishTaskExpiryEvent(taskID); err != nil {
				tm.logger.Infof("Failed to publish expiry event for task %s: %v", taskID, err)
			}
		}
	}
}

// Check tasks and set expiry times
func (tm *JobManager) CheckTasks() {
	// Load configuration from the YAML file
	if err := tm.loadConfig(); err != nil {
		tm.logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize the cron job scheduler
	c := cron.New(cron.WithSeconds()) // cron jobs with seconds precision

	// Add tasks from the YAML file to the cron scheduler
	for _, task := range tm.config.Tasks {
		// If the task has a function, we execute it
		if task.Function != "" && task.Schedule != "" {
			// Parse the function and its parameters from the string
			functionDesc := task.Function
			functionDesc = strings.TrimSpace(functionDesc)

			// Extract the URL or parameter from the function string
			functionName := functionDesc
			param := ""
			if strings.Contains(functionDesc, "(") && strings.Contains(functionDesc, ")") {
				startIdx := strings.Index(functionDesc, "(") + 1
				endIdx := strings.Index(functionDesc, ")")
				if startIdx < endIdx {
					param = functionDesc[startIdx:endIdx]
					functionName = strings.TrimSpace(functionDesc[:startIdx-1])
				}
			}

			// Add the task to the cron scheduler
			_, err := c.AddFunc(task.Schedule, func(functionName, param string) func() {
				return func() {
					tm.executeFunction(functionName, param)
				}
			}(functionName, param))

			if err != nil {
				tm.logger.Infof("Failed to add task %s: %v", task.Name, err)
			} else {
				tm.logger.Infof("Scheduled function task %s with schedule %s", task.Name, task.Schedule)
			}
		}

		// If the task has a command, we execute it
		if task.Function != "" && task.Schedule != "" {
			// Add the task to the cron scheduler
			_, err := c.AddFunc(task.Schedule, func(function string, parameters string) func() {
				return func() {
					tm.executeFunction(function, parameters)
				}
			}(task.Function, task.Parameters))

			if err != nil {
				tm.logger.Infof("Failed to add task %s: %v", task.Name, err)
			} else {
				tm.logger.Infof("Scheduled command task %s with schedule %s", task.Name, task.Schedule)
			}
		}

		if task.Deadline != "" {
			// Parse the deadline string to a time.Time object
			deadlineTime, err := time.Parse(time.RFC3339, task.Deadline)
			if err != nil {
				tm.logger.Infof("Failed to parse deadline for task %s: %v", task.Name, err)
				continue
			}

			// Convert the deadline time to a Unix timestamp
			expiryTime := deadlineTime.Unix()

			// Set the Redis key with the expiry time
			key := fmt.Sprintf("task:%s:expiry", task.Name)
			key = strings.ReplaceAll(key, " ", "_")

			err = tm.rdb.Set(tm.ctx, key, expiryTime, 0).Err()
			if err != nil {
				tm.logger.Infof("Failed to set expiry time for task %s: %v", task.Name, err)
				continue
			}

			tm.logger.Infof("Set expiry time for task %s: %d", task.Name, expiryTime)
		}
	}

	// Add the task expiry check cron job
	_, err := c.AddFunc("@every 1m", tm.checkTaskExpiry)
	if err != nil {
		tm.logger.Fatalf("Failed to add task expiry check cron job: %v", err)
	} else {
		tm.logger.Infof("Scheduled task expiry check cron job")
	}

	// Start the cron scheduler
	c.Start()

}
