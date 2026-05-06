package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wallissonmarinho/GoJob/internal/adapters/cli"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
	"github.com/wallissonmarinho/GoJob/internal/core/services"
)

// MockSyncClient is a mock implementation of ports.SyncClient for testing
type MockSyncClient struct {
	triggerSyncFunc func(ctx context.Context, url string, apiKey string) (int, error)
}

func (m *MockSyncClient) TriggerSync(ctx context.Context, url string, apiKey string) (int, error) {
	if m.triggerSyncFunc != nil {
		return m.triggerSyncFunc(ctx, url, apiKey)
	}
	return 200, nil
}

func TestSyncService_Execute_Success(t *testing.T) {
	mockClient := &MockSyncClient{
		triggerSyncFunc: func(ctx context.Context, url string, apiKey string) (int, error) {
			return 202, nil
		},
	}

	config := cli.NewSyncConfig(
		"http://localhost:8080/admin/sync",
		"test-api-key",
		10*time.Second,
		false,
	)

	service := services.NewSyncService(mockClient, config)
	err := service.Execute()

	assert.NoError(t, err)
}

func TestSyncService_Execute_HTTPError(t *testing.T) {
	mockClient := &MockSyncClient{
		triggerSyncFunc: func(ctx context.Context, url string, apiKey string) (int, error) {
			return 500, nil
		},
	}

	config := cli.NewSyncConfig(
		"http://localhost:8080/admin/sync",
		"test-api-key",
		10*time.Second,
		false,
	)

	service := services.NewSyncService(mockClient, config)
	err := service.Execute()

	// Service itself handles the error gracefully
	assert.NoError(t, err)
}

func TestSyncService_WithVerboseLogging(t *testing.T) {
	mockClient := &MockSyncClient{
		triggerSyncFunc: func(ctx context.Context, url string, apiKey string) (int, error) {
			return 202, nil
		},
	}

	config := cli.NewSyncConfig(
		"http://localhost:8080/admin/sync",
		"test-api-key",
		10*time.Second,
		true, // verbose enabled
	)

	service := services.NewSyncService(mockClient, config)
	err := service.Execute()

	assert.NoError(t, err)
}

var _ ports.SyncClient = (*MockSyncClient)(nil)
