package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/walterfan/lazy-rabbit-reminder/internal/auth"
	"github.com/walterfan/lazy-rabbit-reminder/internal/secret"
	"github.com/walterfan/lazy-rabbit-reminder/internal/task"
	"github.com/walterfan/lazy-rabbit-reminder/pkg/metrics"
)

type StaticRoute struct {
	Path string
	Dir  string
}

type CommandResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type NewsItem struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// WebApiService holds necessary dependencies for the web service
type WebApiService struct {
	logger      *zap.Logger
	redisClient *redis.Client
	authService *auth.AuthService
}

// NewWebApiService creates a new instance of WebApiService with the required dependencies.
func NewWebApiService(logger *zap.Logger, redisClient *redis.Client, authService *auth.AuthService) *WebApiService {
	return &WebApiService{
		logger:      logger,
		redisClient: redisClient,
		authService: authService,
	}
}

// Run starts the HTTP/HTTPS server with routes and middleware
func (thiz *WebApiService) Run() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	// Paths to skip for logging and metrics
	skipPrefixes := []string{"/assets/", "/asserts/", "/.well-known/"}

	// Register Prometheus metrics and middleware (with skips)
	metrics.Register()
	r.Use(wrapMiddlewareWithSkips(metrics.MetricsMiddleware(), skipPrefixes))

	// Structured request logging middleware (with skips)
	r.Use(thiz.requestLoggerWithSkips(skipPrefixes))

	// Health and metrics
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/api/v1/news", func(c *gin.Context) {
		ctx := context.Background()
		newsKey := "news:latest"

		// Fetch all items sorted by timestamp desc
		items, err := thiz.redisClient.ZRevRangeWithScores(ctx, newsKey, 0, -1).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch news",
			})
			return
		}

		var newsList []NewsItem
		for _, item := range items {
			newsList = append(newsList, NewsItem{
				Message:   item.Member.(string),
				Timestamp: int64(item.Score),
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"news": newsList,
		})
	})

	r.GET("/api/v1/commands", func(ctx *gin.Context) {
		var commands []CommandResponse
		if err := viper.UnmarshalKey("commands", &commands); err != nil {
			thiz.logger.Error("Failed to unmarshal commands from config",
				zap.Error(err),
			)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load commands"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"commands": commands})
	})

	r.POST("/api/v1/commands", func(c *gin.Context) {
		type CommandRequest struct {
			Name       string `json:"name"`
			Parameters string `json:"parameters"`
		}

		var req CommandRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		handler, exists := commandHandlers[req.Name]
		if !exists {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Command '%s' not found", req.Name),
			})
			return
		}

		result, err := handler(req.Parameters)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": result,
		})
	})

	// Register auth routes
	thiz.setupAuthRoutes(r)

	// Register secret routes
	repo := secret.NewSecretRepository()
	secretService := secret.NewSecretService(repo)
	authMiddleware := auth.NewAuthMiddleware(thiz.authService)
	secret.RegisterRoutes(r, secretService, authMiddleware)

	// Register task routes
	taskRepo := task.NewTaskRepository()
	taskService := task.NewTaskService(taskRepo)
	task.RegisterRoutes(r, taskService, authMiddleware)

	// Default route for SPA
	r.NoRoute(func(ctx *gin.Context) {
		frontendPath := viper.GetString("server.webroot")
		ctx.File(frontendPath)
	})

	thiz.setupPublicRoutes(r)

	thiz.setupPrivateRoutes(r)

	thiz.setupCommands(r)

	// Run the server
	thiz.startServer(r)
}

// requestLoggerWithSkips logs request/response details unless path matches skip prefixes
func (thiz *WebApiService) requestLoggerWithSkips(skipPrefixes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, p := range skipPrefixes {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		start := time.Now()
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		fullPath := c.FullPath()
		if fullPath == "" {
			fullPath = path
		}

		thiz.logger.Info("http request",
			zap.String("method", c.Request.Method),
			zap.String("path", fullPath),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Int("response_size", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// wrapMiddlewareWithSkips wraps a gin middleware and skips it for matching path prefixes
func wrapMiddlewareWithSkips(mw gin.HandlerFunc, skipPrefixes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for _, p := range skipPrefixes {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}
		mw(c)
	}
}

func (thiz *WebApiService) startServer(r *gin.Engine) {
	port := viper.GetInt("server.port")
	tlsConfig := viper.Sub("server.tls")

	if tlsConfig != nil && tlsConfig.IsSet("cert_file") && tlsConfig.IsSet("key_file") {
		certFile := tlsConfig.GetString("cert_file")
		keyFile := tlsConfig.GetString("key_file")

		thiz.logger.Info("Starting HTTPS server",
			zap.Int("port", port),
			zap.String("cert_file", certFile),
			zap.String("key_file", keyFile),
		)

		if err := r.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile); err != nil {
			thiz.logger.Fatal("HTTPS server failed to start",
				zap.Error(err),
			)
		}
	} else {
		thiz.logger.Info("Starting HTTP server",
			zap.Int("port", port),
		)

		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			thiz.logger.Fatal("HTTP server failed to start",
				zap.Error(err),
			)
		}
	}
}
func (thiz *WebApiService) AuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != username || pass != password {
			thiz.logger.Warn("Unauthorized access attempt",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		thiz.logger.Info("Authorized access",
			zap.String("user", user),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}

func (thiz *WebApiService) setupPublicRoutes(r *gin.Engine) {
	var publicRoutes []StaticRoute
	if err := viper.UnmarshalKey("static_routes.public", &publicRoutes); err != nil {
		thiz.logger.Fatal("Error unmarshaling public routes",
			zap.Error(err),
		)
	}

	// Ensure /assets is served when server.webroot is configured but assets route is missing
	webroot := viper.GetString("server.webroot")
	if webroot != "" {
		distDir := filepath.Dir(webroot)
		assetsDir := filepath.Join(distDir, "assets")
		if info, err := os.Stat(assetsDir); err == nil && info.IsDir() {
			missing := true
			for _, route := range publicRoutes {
				if route.Path == "/assets" {
					missing = false
					break
				}
			}
			if missing {
				publicRoutes = append(publicRoutes, StaticRoute{Path: "/assets", Dir: assetsDir})
				thiz.logger.Info("Auto-registered assets static route",
					zap.String("path", "/assets"),
					zap.String("dir", assetsDir),
				)
			}
		}
	}

	for _, route := range publicRoutes {
		r.Static(route.Path, route.Dir)
		thiz.logger.Info("Registered public route",
			zap.String("path", route.Path),
			zap.String("dir", route.Dir),
		)
	}
}

func (thiz *WebApiService) setupPrivateRoutes(r *gin.Engine) {
	var privateRoutes []StaticRoute
	if err := viper.UnmarshalKey("static_routes.private", &privateRoutes); err != nil {
		thiz.logger.Fatal("Error unmarshaling private routes",
			zap.Error(err),
		)
	}

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")

	for _, route := range privateRoutes {
		authorized := r.Group(route.Path, thiz.AuthMiddleware(username, password))
		authorized.Static("/", route.Dir)
		thiz.logger.Info("Registered private route",
			zap.String("path", route.Path),
			zap.String("dir", route.Dir),
		)
	}
}

func (thiz *WebApiService) setupCommands(r *gin.Engine) {

	RegisterCommandHandler("read_over_ssh", ReadOverSSH)
}

func (thiz *WebApiService) setupAuthRoutes(r *gin.Engine) {
	// Create auth handlers and middleware
	authHandlers := auth.NewAuthHandlers(thiz.authService)
	authMiddleware := auth.NewAuthMiddleware(thiz.authService)

	// Register auth routes
	auth.RegisterRoutes(r, authHandlers, authMiddleware)

	thiz.logger.Info("Registered authentication routes")
}
