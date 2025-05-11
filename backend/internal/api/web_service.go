package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// WebServiceConfig holds necessary dependencies for the web service
type WebServiceConfig struct {
	Logger *zap.Logger
}

// Run starts the HTTP/HTTPS server with routes and middleware
func (c *WebServiceConfig) Run() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "pong"})
	})

	// Default route for SPA
	r.NoRoute(func(ctx *gin.Context) {
		frontendPath := viper.GetString("server.webroot")
		ctx.File(frontendPath)
	})

	// Setup public routes
	var publicRoutes []StaticRoute
	if err := viper.UnmarshalKey("static_routes.public", &publicRoutes); err != nil {
		c.Logger.Fatal("Error unmarshaling public routes",
			zap.Error(err),
		)
	}
	for _, route := range publicRoutes {
		r.Static(route.Path, route.Dir)
		c.Logger.Info("Registered public route",
			zap.String("path", route.Path),
			zap.String("dir", route.Dir),
		)
	}

	// Setup private routes
	var privateRoutes []StaticRoute
	if err := viper.UnmarshalKey("static_routes.private", &privateRoutes); err != nil {
		c.Logger.Fatal("Error unmarshaling private routes",
			zap.Error(err),
		)
	}

	username := viper.GetString("auth.username")
	password := viper.GetString("auth.password")

	for _, route := range privateRoutes {
		authorized := r.Group(route.Path, AuthMiddleware(username, password))
		authorized.Static("/", route.Dir)
		c.Logger.Info("Registered private route",
			zap.String("path", route.Path),
			zap.String("dir", route.Dir),
		)
	}

	// Run the server
	port := viper.GetInt("server.port")
	tlsConfig := viper.Sub("server.tls")

	if tlsConfig != nil && tlsConfig.IsSet("cert_file") && tlsConfig.IsSet("key_file") {
		certFile := tlsConfig.GetString("cert_file")
		keyFile := tlsConfig.GetString("key_file")

		c.Logger.Info("Starting HTTPS server",
			zap.Int("port", port),
			zap.String("cert_file", certFile),
			zap.String("key_file", keyFile),
		)

		if err := r.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile); err != nil {
			c.Logger.Fatal("HTTPS server failed to start",
				zap.Error(err),
			)
		}
	} else {
		c.Logger.Info("Starting HTTP server",
			zap.Int("port", port),
		)

		if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
			c.Logger.Fatal("HTTP server failed to start",
				zap.Error(err),
			)
		}
	}
}