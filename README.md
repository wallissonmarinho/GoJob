# GoJob — HTTP Sync Trigger

**Universal command-line tool in Go with Hexagonal Architecture to trigger HTTP sync endpoints via POST request, designed to run as a Kubernetes CronJob.**

**Currently used by:** GoAnimes sync endpoint (first use case, but not limited to it)

---

## 🚀 Quick Start

### Build

```bash
go build -o gojob ./cmd/gojob
```

### Run (Local)

```bash
# Show help
./gojob --help
./gojob sync --help

# Basic sync (uses env vars or defaults)
./gojob sync

# With explicit URL and API key
./gojob sync \
  --url <your-sync-url> \
  --api-key your-api-key

# Verbose mode
./gojob sync --verbose \
  --url <your-sync-url> \
  --api-key your-api-key

# Custom timeout (30s default)
./gojob sync --timeout 60 --api-key your-api-key
```

---

## 📋 `sync` Command Flags

| Flag | Short | Type | Required | Description |
|------|-------|------|----------|-------------|
| `--url` | `-u` | string | yes | HTTP sync endpoint URL |
| `--api-key` | `-k` | string | yes | API key for authentication |
| `--timeout` | | int | no | Request timeout in seconds (default: 30) |
| `--verbose` | | bool | no | Enable verbose logging |

### Environment Variables (Recommended)

- `SYNC_URL` - HTTP sync endpoint URL
- `API_KEY` - API key for authentication (recommended for security)

---

## 🐳 Docker

### Build Image

```bash
docker build -t gojob:latest .
```

### Run Container

```bash
docker run --rm \
  -e SYNC_URL=<your-sync-url> \
  -e API_KEY=your-api-key \
  gojob:latest sync --verbose
```

---

## ☸️ Kubernetes / K3s CronJob

### Prerequisites

1. **Build and push image to registry:**

```bash
docker build -t your-registry/gojob:latest .
docker push your-registry/gojob:latest
```

2. **Update deployment manifests** with:
   - Your container registry URL
   - Sync endpoint URL
   - API key for authentication

### Deploy CronJob

```bash
# Apply manifests (creates ConfigMap, Secret, and CronJob)
kubectl apply -f deploy/k8s/cronjob.yaml

# Verify
kubectl get cronjobs
kubectl get pods -l app=gojob-sync
kubectl logs -l app=gojob-sync --tail=50
```

### Configuration via K8s

Edit `deploy/k8s/cronjob.yaml`:

```yaml
# ConfigMap - update sync URL
data:
  sync-url: "<your-sync-url>"

# Secret - update API key
stringData:
  api-key: "your-actual-api-key"
```

Then reapply:

```bash
kubectl apply -f deploy/k8s/cronjob.yaml
```

### Customizing CronJob Schedule

Edit the `schedule` field in `cronjob.yaml`:

```yaml
spec:
  schedule: "*/10 * * * *"  # Every 10 minutes
```

Common schedules:
- `0 * * * *` - Every hour
- `0 */6 * * *` - Every 6 hours
- `0 0 * * *` - Daily at midnight
- `*/5 * * * *` - Every 5 minutes

---

## 🏛️ Architecture

This project follows **Hexagonal Architecture (Ports & Adapters)**:

```
┌──────────────────────────────┐
│     CLI Input Adapter        │
└────────────────┬─────────────┘
                 │
                 ↓
       ┌─────────────────────┐
       │  SyncService        │ ← Use Case (core)
       │  (Orchestration)    │
       └────────────┬────────┘
                    │
           ┌────────┘
           ↓
      HTTP Client Adapter
      (SyncClient)
```

**Benefits:**
- ✅ **Testable**: Mock HTTP client without making real requests
- ✅ **Decoupled**: Replace HTTP client with another transport? New adapter, done
- ✅ **Flexible**: Add new input adapters (REST API, scheduler, webhook, etc.)

### Package Structure

```
cmd/gojob/              # CLI entry point
├── main.go             # Application bootstrap
└── commands.go         # Command factory & handlers

internal/
├── core/               # Business logic (use cases)
│   ├── services/       # SyncService (orchestration)
│   └── ports/          # Interfaces (SyncClient port)
└── adapters/           # External adapters
    └── http/           # HTTP client implementation
```

---

## 🧪 Testing

```bash
go test ./...
```

---

## 📝 License

MIT
