package command

import (

	"fmt"
	"net/http"


	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"

	"go.uber.org/zap"
)

type CommandService struct {
	logger *zap.Logger
}

type CommandResponse struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CommandRequest struct {
	Name       string `json:"name"`
	Question   string `json:"question"`
	Parameters string `json:"parameters"`
}

var commandHandlers = make(map[string]func(string) (string, error))

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity
	},
}

func NewCommandService(logger *zap.Logger) *CommandService {
	// Register command handlers
	//RegisterCommandHandler("read_over_ssh", ReadOverSSH)

	return &CommandService{
		logger: logger,
	}
}

func RegisterCommandHandler(name string, webCmd func(string) (string, error)) {
	commandHandlers[name] = webCmd
}

func (s *CommandService) GetCommands(c *gin.Context) {
	var commands []CommandResponse
	if err := viper.UnmarshalKey("commands", &commands); err != nil {
		s.logger.Error("Failed to unmarshal commands from config",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load commands"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"commands": commands})
}

func (s *CommandService) ExecuteCommand(c *gin.Context) {
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
}

func (s *CommandService) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		s.logger.Error("Failed to upgrade connection", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	var req CommandRequest
	if err := conn.ReadJSON(&req); err != nil {
		s.logger.Error("Invalid request", zap.Error(err))
		conn.WriteJSON(gin.H{"error": "Invalid request"})
		return
	}

	s.logger.Info("Received WebSocket command request",
		zap.String("name", req.Name),
		zap.String("question", req.Question))

	handler, exists := commandHandlers[req.Name]
	if !exists {
		conn.WriteJSON(gin.H{
			"error": fmt.Sprintf("Command '%s' not found", req.Name),
		})
		return
	}

	// Execute command and stream result
	result, err := handler(req.Parameters)
	if err != nil {
		conn.WriteJSON(gin.H{
			"error": err.Error(),
		})
		return
	}

	// Send result as WebSocket message
	conn.WriteJSON(gin.H{
		"result": result,
	})
}

