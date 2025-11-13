// Package backup provides functionality for backing up MikroTik RouterOS configurations.
package backup

import (
	"context"
	"fmt"
	"io"
)

// Config holds the configuration for a backup operation.
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	KeyFile  string
}

// Service handles backup operations.
type Service struct {
	sshClient SSHClient
}

// SSHClient defines the interface for SSH operations.
type SSHClient interface {
	Connect(ctx context.Context, config Config) error
	ExecuteCommand(ctx context.Context, cmd string) (string, error)
	Close() error
}

// New creates a new backup service.
func New(client SSHClient) *Service {
	return &Service{
		sshClient: client,
	}
}

// Execute performs a backup operation.
func (s *Service) Execute(ctx context.Context, config Config, output io.Writer) error {
	if err := s.sshClient.Connect(ctx, config); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if closeErr := s.sshClient.Close(); closeErr != nil {
			// Log or handle close error if needed
			_ = closeErr
		}
	}()

	// Export configuration
	result, err := s.sshClient.ExecuteCommand(ctx, "/export")
	if err != nil {
		return fmt.Errorf("failed to export configuration: %w", err)
	}

	if _, err := output.Write([]byte(result)); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
