package cmd

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "lazy-rabbit-reminder",
	Short: "A brief description of your application",
	Long:  `A longer description that spans multiple lines...`,
}

var logger *zap.Logger

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

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		// .env file is optional, so we just log a warning instead of failing
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	rootCmd.PersistentFlags().String("config", "", "config file (optional) instead of using default 'config/config.yaml'")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		initLogger()
		initConfig()
	}
}
