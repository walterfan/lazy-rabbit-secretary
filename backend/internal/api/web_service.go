package api

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
}

// NewWebApiService creates a new instance of WebApiService with the required dependencies.
func NewWebApiService(logger *zap.Logger, redisClient *redis.Client) *WebApiService {
	return &WebApiService{
		logger:      logger,
		redisClient: redisClient,
	}
}

// Run starts the HTTP/HTTPS server with routes and middleware
func (thiz *WebApiService) Run() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

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
