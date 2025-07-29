package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var rootCmd = &cobra.Command{
	Use:   "lazy-rabbit-reminder",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines...`,
}

var logger *zap.Logger
var rdb *redis.Client

func initLogger() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	defer logger.Sync()
}

func initConfig() {
	// Get config flag value
	configFile, _ := rootCmd.Flags().GetString("config")

	viper.SetConfigType("yaml")

	if configFile != "" {
		// Use provided config file
		viper.SetConfigFile(configFile)
	} else {
		// Fallback to default config
		viper.SetConfigName("config")
		viper.AddConfigPath("./config")
	}

	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Error reading config file",
			zap.Error(err),
		)
	}
	logger.Info("Configuration loaded successfully",
		zap.String("configFile", viper.ConfigFileUsed()),
	)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func GetRedisClient() *redis.Client {
	return rdb
}

func initRedis() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file",
			zap.Error(err),
		)
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Initialize the Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0, // Use default DB
	})

	ctx := context.Background()
	newsKey := "news:latest"
	timestamp := time.Now().Unix()

	err = rdb.ZAdd(ctx, newsKey, &redis.Z{
		Score:  float64(timestamp),
		Member: "Your news message here",
	}).Err()

	if err != nil {
		logger.Fatal("Error adding news to Redis",
			zap.Error(err),
		)
	}
}

func init() {
	rootCmd.PersistentFlags().String("config", "", "config file (optional) instead of using default 'config/config.yaml'")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogger()
		initConfig()
		initRedis()
	}
}
