package cmd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)
type StaticRoute struct {
	Path string
	Dir  string
}


// AuthMiddleware applies HTTP Basic Authentication for protected routes
func AuthMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != username || pass != password {
			logger.Warn("Unauthorized access attempt",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
			)
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		logger.Info("Authorized access",
			zap.String("user", user),
			zap.String("path", c.Request.URL.Path),
		)
		c.Next()
	}
}

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP/HTTPS server",
	Long:  "Starts the web server with public/private routes and optional TLS support.",
	Run: func(cmd *cobra.Command, args []string) {

		r := gin.Default()

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message":"pong"})
		})

		// Default route for SPA
		r.NoRoute(func(c *gin.Context) {
			frontendPath := viper.GetString("server.webroot")
			c.File(frontendPath)
		})

		// Setup public routes
		var publicRoutes []StaticRoute
		if err := viper.UnmarshalKey("static_routes.public", &publicRoutes); err != nil {
			logger.Fatal("Error unmarshaling public routes",
				zap.Error(err),
			)
		}
		for _, route := range publicRoutes {
			r.Static(route.Path, route.Dir)
			logger.Info("Registered public route",
				zap.String("path", route.Path),
				zap.String("dir", route.Dir),
			)
		}

		// Setup private routes
		var privateRoutes []StaticRoute
		if err := viper.UnmarshalKey("static_routes.private", &privateRoutes); err != nil {
			logger.Fatal("Error unmarshaling private routes",
				zap.Error(err),
			)
		}

		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		for _, route := range privateRoutes {
			authorized := r.Group(route.Path, AuthMiddleware(username, password))
			authorized.Static("/", route.Dir)
			logger.Info("Registered private route",
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

			logger.Info("Starting HTTPS server",
				zap.Int("port", port),
				zap.String("cert_file", certFile),
				zap.String("key_file", keyFile),
			)

			if err := r.RunTLS(fmt.Sprintf(":%d", port), certFile, keyFile); err != nil {
				logger.Fatal("HTTPS server failed to start",
					zap.Error(err),
				)
			}
		} else {
			logger.Info("Starting HTTP server",
				zap.Int("port", port),
			)

			if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
				logger.Fatal("HTTP server failed to start",
					zap.Error(err),
				)
			}
		}
	},
}


func init() {
	rootCmd.AddCommand(serverCmd)
}
