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
# 1. Copy .env.example to .env and update with your values
cp .env.example .env
# Edit .env with your SYNC_URL and API_KEY

# 2. Show help
./gojob --help
./gojob sync --help

# 3. Basic sync (uses env vars from .env)
./gojob sync

# 4. With explicit URL and API key
export SYNC_URL=<your-sync-url>
export API_KEY=your-api-key
./gojob sync

# 5. Verbose mode
./gojob sync --verbose

# 6. Custom timeout (30s default)
./gojob sync --timeout 60
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

## 📝 Environment Variables Reference

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `SYNC_URL` | string | yes | - | HTTP sync endpoint URL (e.g., `https://api.example.com/sync`) |
| `API_KEY` | string | yes | - | API key for Bearer token authentication |
| `TIMEOUT` | duration | no | `30s` | Request timeout (e.g., `45s`, `1m`) |
| `VERBOSE` | bool | no | `false` | Enable verbose logging (`true`, `1`, or `yes`) |

### Example `.env` File

```bash
# HTTP Sync Endpoint
SYNC_URL=https://goanimes.example.com/admin/sync

# API Authentication
API_KEY=your-secret-api-key-here

# Optional: Custom timeout
TIMEOUT=45s

# Optional: Verbose logging
VERBOSE=true
```

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

2. **Create Secret with credentials:**

```bash
kubectl create secret generic gojob-env \
  --from-literal=SYNC_URL=https://your-endpoint/sync \
  --from-literal=API_KEY=your-api-key \
  -n default
```

### Deploy CronJob

```bash
# Update image reference in cronjob.yaml
sed -i 's|gojob:latest|your-registry/gojob:latest|g' deploy/k8s/cronjob.yaml

# Apply manifests
kubectl apply -f deploy/k8s/cronjob.yaml

# Verify
kubectl get cronjobs
kubectl get pods -l app=gojob-sync
kubectl logs -l app=gojob-sync --tail=50
```

### Configuration via K8s

The CronJob uses a consolidated Secret `gojob-env`:

```yaml
envFrom:
- secretRef:
    name: gojob-env
```

Update the secret:

```bash
# Edit existing secret
kubectl edit secret gojob-env -n default

# Or recreate:
kubectl delete secret gojob-env -n default
kubectl create secret generic gojob-env \
  --from-literal=SYNC_URL=https://new-endpoint/sync \
  --from-literal=API_KEY=new-api-key
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
