package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/walterfan/lazy-rabbit-secretary/internal/auth"
	"github.com/walterfan/lazy-rabbit-secretary/internal/book"
	"github.com/walterfan/lazy-rabbit-secretary/internal/daily"
	"github.com/walterfan/lazy-rabbit-secretary/internal/inbox"
	"github.com/walterfan/lazy-rabbit-secretary/internal/post"
	"github.com/walterfan/lazy-rabbit-secretary/internal/prompt"
	"github.com/walterfan/lazy-rabbit-secretary/internal/reminder"
	"github.com/walterfan/lazy-rabbit-secretary/internal/secret"
	"github.com/walterfan/lazy-rabbit-secretary/internal/task"
	"github.com/walterfan/lazy-rabbit-secretary/internal/wiki"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/database"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/metrics"
)

type StaticRoute struct {
	Path string
	Dir  string
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
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

	// API-specific detailed logging for /api/v1/* routes
	r.Use(thiz.apiLoggerMiddleware())

	// Health and metrics
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.GET("/api/v1/news", func(c *gin.Context) {
		ctx := context.Background()
		newsKey := "news:latest"

		if thiz.redisClient == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Redis client not initialized",
			})
			return
		}

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

	// Register reminder routes first (needed by task service)
	reminderRepo := reminder.NewReminderRepository()
	reminderService := reminder.NewReminderService(reminderRepo)
	reminder.RegisterRoutes(r, reminderService, authMiddleware)

	// Register task routes (with reminder service dependency)
	taskRepo := task.NewTaskRepository()
	taskService := task.NewTaskService(taskRepo, reminderService)
	task.RegisterRoutes(r, taskService, authMiddleware)

	// Register prompt routes
	promptRoutes := prompt.NewPromptRoutes(database.GetDB())
	promptRoutes.RegisterRoutes(r)

	bookRepo := book.NewBookRepository()
	bookService := book.NewBookService(bookRepo)
	book.RegisterRoutes(r, bookService, authMiddleware)

	// Register post routes
	postService := post.NewPostService(database.GetDB())
	post.RegisterRoutes(r, postService, authMiddleware)

	// Register wiki routes
	wikiService := wiki.NewWikiService(database.GetDB())
	wiki.RegisterRoutes(r, wikiService, authMiddleware)

	// Register GTD system routes
	inboxService := inbox.NewInboxService(database.GetDB())
	inbox.RegisterInboxRoutes(r, inboxService, authMiddleware)

	dailyService := daily.NewDailyService(database.GetDB())
	daily.RegisterDailyRoutes(r, dailyService, authMiddleware)

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

// apiLoggerMiddleware provides detailed logging for /api/v1/* routes including request/response bodies
func (thiz *WebApiService) apiLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Only log /api/v1/* routes
		if !strings.HasPrefix(path, "/api/v1/") {
			c.Next()
			return
		}

		start := time.Now()

		// Capture request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create custom response writer to capture response
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)
		status := c.Writer.Status()

		// Prepare log fields
		logFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// Add request body if present and not too large
		if len(requestBody) > 0 && len(requestBody) < 10240 { // 10KB limit
			// Mask sensitive fields in request body
			requestBodyStr := thiz.maskSensitiveData(string(requestBody))
			logFields = append(logFields, zap.String("request_body", requestBodyStr))
		}

		// Add response body if not too large
		responseBody := w.body.String()
		if len(responseBody) > 0 && len(responseBody) < 10240 { // 10KB limit
			// Mask sensitive fields in response body
			responseBodyStr := thiz.maskSensitiveData(responseBody)
			logFields = append(logFields, zap.String("response_body", responseBodyStr))
		}

		// Add content type headers
		if contentType := c.Request.Header.Get("Content-Type"); contentType != "" {
			logFields = append(logFields, zap.String("request_content_type", contentType))
		}
		if responseContentType := c.Writer.Header().Get("Content-Type"); responseContentType != "" {
			logFields = append(logFields, zap.String("response_content_type", responseContentType))
		}

		// Log based on status code
		if status >= 400 {
			thiz.logger.Error("API request completed with error", logFields...)
		} else {
			thiz.logger.Info("API request completed", logFields...)
		}
	}
}

// maskSensitiveData replaces sensitive information in JSON strings
func (thiz *WebApiService) maskSensitiveData(data string) string {
	// List of sensitive field names to mask
	sensitiveFields := []string{"password", "kek", "value", "secret", "token", "key", "auth"}

	result := data
	for _, field := range sensitiveFields {
		// Create regex pattern to match JSON field patterns like "field":"value" or "field": "value"
		pattern := fmt.Sprintf(`"%s"\s*:\s*"[^"]*"`, field)
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}

		// Replace the matched pattern with masked version
		replacement := fmt.Sprintf(`"%s":"***MASKED***"`, field)
		result = re.ReplaceAllString(result, replacement)
	}

	return result
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
