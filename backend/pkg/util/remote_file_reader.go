package util

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

// RemoteFileReader defines an interface to read remote files.
type RemoteFileReader interface {
	ReadFile(path string) (string, error)
}

// SSHClient holds the SSH connection configuration.
type SSHClient struct {
	Host     string
	Port     int
	Username string
	Password string
	client   *ssh.Client
}

// NewSSHClient creates a new SSHClient and establishes a connection.
func NewSSHClient(host string, port int, username, password string) (*SSHClient, error) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ Use for testing only
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("Failed to dial: %s by %s", addr, username)
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &SSHClient{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		client:   client,
	}, nil
}

// ReadFile reads the content of a file on the remote server.
func (c *SSHClient) ReadFile(path string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	output, err := session.Output("cat " + path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(output), nil
}
