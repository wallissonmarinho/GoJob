package services

import (
	"context"
	"fmt"
	"log"

	"github.com/wallissonmarinho/GoJob/internal/adapters/cli"
	"github.com/wallissonmarinho/GoJob/internal/core/ports"
)

// SyncService implements ports.Executor and orchestrates the sync operation
type SyncService struct {
	syncClient ports.SyncClient
	config     *cli.SyncConfig
	logger     *cli.Logger
}

// NewSyncService creates a new SyncService
func NewSyncService(syncClient ports.SyncClient, config *cli.SyncConfig) *SyncService {
	return &SyncService{
		syncClient: syncClient,
		config:     config,
		logger:     cli.NewLogger(config.Verbose),
	}
}

// Execute runs the sync operation
func (s *SyncService) Execute() error {
	s.logger.Info("Starting GoAnimes sync job...")
	s.logger.Debugf("URL: %s", s.config.URL)
	s.logger.Debugf("Timeout: %v", s.config.Timeout)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), s.config.Timeout)
	defer cancel()

	// Trigger sync
	statusCode, err := s.syncClient.TriggerSync(ctx, s.config.URL, s.config.APIKey)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Sync failed: %v", err))
		return err
	}

	s.logger.Success(fmt.Sprintf("Sync completed (HTTP %d)", statusCode))
	log.Printf("✅ Sync completed successfully\n")

	return nil
}
