package backup_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/mountain-reverie/mikrotik-configuration-backup/internal/backup"
)

// mockSSHClient is a mock implementation of SSHClient for testing.
type mockSSHClient struct {
	connectFunc        func(ctx context.Context, config backup.Config) error
	executeCommandFunc func(ctx context.Context, cmd string) (string, error)
	closeFunc          func() error
}

func (m *mockSSHClient) Connect(ctx context.Context, config backup.Config) error {
	if m.connectFunc != nil {
		return m.connectFunc(ctx, config)
	}
	return nil
}

func (m *mockSSHClient) ExecuteCommand(ctx context.Context, cmd string) (string, error) {
	if m.executeCommandFunc != nil {
		return m.executeCommandFunc(ctx, cmd)
	}
	return "", nil
}

func (m *mockSSHClient) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestService_Execute_Success(t *testing.T) {
	t.Parallel()

	expectedConfig := "# MikroTik config\n/system identity set name=test\n"

	client := &mockSSHClient{
		executeCommandFunc: func(_ context.Context, cmd string) (string, error) {
			if cmd != "/export" {
				t.Errorf("unexpected command: %s", cmd)
			}
			return expectedConfig, nil
		},
	}

	service := backup.New(client)
	output := &bytes.Buffer{}

	config := backup.Config{
		Host:     "192.168.88.1",
		Port:     22,
		Username: "admin",
		Password: "password",
	}

	err := service.Execute(context.Background(), config, output)
	if err != nil {
		t.Fatalf("Execute() error = %v, want nil", err)
	}

	if got := output.String(); got != expectedConfig {
		t.Errorf("Execute() output = %q, want %q", got, expectedConfig)
	}
}

func TestService_Execute_ConnectionError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("connection failed")

	client := &mockSSHClient{
		connectFunc: func(_ context.Context, _ backup.Config) error {
			return expectedErr
		},
	}

	service := backup.New(client)
	output := &bytes.Buffer{}

	config := backup.Config{
		Host:     "192.168.88.1",
		Port:     22,
		Username: "admin",
		Password: "password",
	}

	err := service.Execute(context.Background(), config, output)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Execute() error does not wrap expected error")
	}
}

func TestService_Execute_CommandError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("command execution failed")

	client := &mockSSHClient{
		executeCommandFunc: func(_ context.Context, _ string) (string, error) {
			return "", expectedErr
		},
	}

	service := backup.New(client)
	output := &bytes.Buffer{}

	config := backup.Config{
		Host:     "192.168.88.1",
		Port:     22,
		Username: "admin",
	}

	err := service.Execute(context.Background(), config, output)
	if err == nil {
		t.Fatal("Execute() error = nil, want error")
	}
}
