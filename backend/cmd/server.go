package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/walterfan/lazy-rabbit-reminder/internal/api"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
	"github.com/walterfan/lazy-rabbit-reminder/internal/service"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/database"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/email"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP/HTTPS server",
	Long:  "Starts the web server with public/private routes and optional TLS support.",
	Run: func(command *cobra.Command, args []string) {

		logger := GetLogger()

		// Initialize database
		if err := database.InitDB(); err != nil {
			logger.Fatal("Failed to initialize database", zap.Error(err))
		}
		logger.Info("Database initialized, will init data")
		// Initialize data (prompts, default user, etc.)
		if err := database.InitData(); err != nil {
			logger.Fatal("Failed to initialize database data", zap.Error(err))
		}

		// Initialize email templates
		if err := email.InitEmailTemplates(); err != nil {
			logger.Fatal("Failed to initialize email templates", zap.Error(err))
		}

		defer database.CloseDB()

		// Initialize Redis and Auth services
		rdb := initRedis(logger)
		authService := initAuth(logger)

		logger.Info("Starting HTTP service...")
		webService := api.NewWebApiService(logger, rdb, authService)
		go webService.Run()

		logger.Info("Starting Job Manager...")
		tm := service.NewJobManager(logger, rdb)
		go tm.CheckTasks()

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		logger.Info("Server is running. Press Ctrl+C to stop.")
		<-signalChan
		logger.Info("Received shutdown signal, shutting down.")
	},
}

func initRedis(logger *zap.Logger) *redis.Client {

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	// Initialize the Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0, // Use default DB
	})

	ctx := context.Background()
	newsKey := "news:latest"
	timestamp := time.Now().Unix()

	err := rdb.ZAdd(ctx, newsKey, &redis.Z{
		Score:  float64(timestamp),
		Member: "Your news message here",
	}).Err()

	if err != nil {
		logger.Error("Error adding news to Redis",
			zap.Error(err),
		)
	}

	return rdb
}

func initAuth(logger *zap.Logger) *auth.AuthService {
	// For now, create a minimal auth service with mock implementations
	// In a real application, you would initialize proper database connections
	// and other dependencies here

	// Create mock implementations (you'll need to implement these properly later)
	userService := &auth.SimpleUserService{}
	passwordManager := auth.NewPasswordManager(10)                                                                      // cost parameter
	jwtManager, err := auth.NewJWTManager("./certs/private.pem", "./certs/public.pem", "lazy-rabbit-reminder", "users") // empty paths for now
	if err != nil {
		logger.Fatal("Failed to create JWT manager", zap.Error(err))
	}

	policyService := &auth.SimplePolicyService{}
	permissionEngine := auth.NewPermissionEngine(policyService)

	authService := auth.NewAuthService(userService, passwordManager, jwtManager, permissionEngine)

	logger.Info("Authentication service initialized")
	return authService
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
