package ports

import "context"

// SyncClient is a port that defines how to trigger sync on an external service
type SyncClient interface {
	TriggerSync(ctx context.Context, url string, apiKey string) (statusCode int, err error)
}
