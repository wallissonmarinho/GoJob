# Unit Tests

This directory contains unit tests for GoJob.

## Structure

```
unit/
├── adapters/          # Tests for adapters (HTTP, CLI, Commands)
│   ├── command_factory_test.go      # Tests for CommandFactory
│   └── http_sync_client_test.go     # Tests for HTTP SyncClient adapter
└── core/              # Tests for core business logic
    └── sync_service_test.go         # Tests for SyncService
```

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test package
go test ./tests/unit/core/...
go test ./tests/unit/adapters/...
```

## Test Coverage

- **Adapters**: HTTP client requests, command factory, CLI commands
- **Core Services**: Sync service execution, error handling, logging
- **Mocks**: MockSyncClient for testing without real HTTP calls
