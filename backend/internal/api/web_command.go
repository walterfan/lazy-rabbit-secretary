package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
	"github.com/walterfan/lazy-rabbit-secretary/pkg/util"
)

type SSHParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	FilePath string `json:"filePath"`
}

var commandHandlers = make(map[string]func(string) (string, error))

func RegisterCommandHandler(name string, webCmd func(string) (string, error)) {
	commandHandlers[name] = webCmd
}

func ReadOverSSH(parameters string) (string, error) {
	var params SSHParams = SSHParams{}

	// Parse the JSON parameters
	if err := json.Unmarshal([]byte(parameters), &params); err != nil {
		log.GetLogger().Errorf("Error parsing parameters: %v", err)
		//return "", fmt.Errorf("invalid parameters: %w", err)
	}

	portStr := os.Getenv("SSH_PORT")
	port := 22 // default SSH port if not set or invalid
	var err error
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.GetLogger().Errorf("Invalid SSH_PORT value: %v", err)
			return "", fmt.Errorf("invalid SSH_PORT: %w", err)
		}
	}

	params = SSHParams{
		Host:     os.Getenv("SSH_HOST"),
		Port:     port,
		Username: os.Getenv("SSH_USERNAME"),
		Password: os.Getenv("SSH_PASSWORD"),
		FilePath: "/root/.aws/credentials",
	}

	// Establish SSH connection
	client, err := util.NewSSHClient(params.Host, params.Port, params.Username, params.Password)
	if err != nil {
		log.GetLogger().Errorf("Error connecting: %v", err)
		return "", fmt.Errorf("ssh connection failed: %w", err)
	}

	// Read remote file
	content, err := client.ReadFile(params.FilePath)
	if err != nil {
		log.GetLogger().Errorf("Error reading file: %v", err)
		return "", fmt.Errorf("file read failed: %w", err)
	}

	fmt.Println("Remote file content:")
	fmt.Println(content)
	return content, nil
}
