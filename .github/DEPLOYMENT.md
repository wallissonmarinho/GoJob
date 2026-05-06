# 🚀 CI/CD Setup Guide for GoJob (Oracle K3s Deployment)

## GitHub Actions Workflow: `oracle-deploy.yml`

The workflow automates:
1. **Test** — `go test ./...` validation
2. **Image** — Multi-arch Docker build (amd64 + arm64) and push to GHCR
3. **Deploy** — SSH to VM and apply K3s CronJob manifests

---

## 📋 GitHub Configuration

### Environment: `prd`

Create a **prd** environment in your repo:
- **Settings → Environments → New environment**
- Name: `prd`
- Add environment variables and secrets (see below)

### Variables (Non-sensitive)

Add these as **Repository Variables** or in the **prd** environment:

| Variable | Required | Example | Description |
|----------|----------|---------|-------------|
| `OCI_VM_HOST` | yes | `your-vm.example.com` | Hostname/IP of your k3s VM |
| `OCI_VM_USER` | yes | `ubuntu` | SSH user on the VM |
| `OCI_DEPLOY_ROOT` | no | `/opt/go` | Root directory for deployments |
| `SYNC_URL` | yes | `https://goanimes.example.com/admin/sync` | HTTP sync endpoint URL |
| `API_KEY` | yes | `your-api-key` | API key for authentication |
| `GOJOB_ENV_B64` | no | (base64) | Additional env vars as base64 |

### Secrets (Sensitive)

Add these as **Repository Secrets** or in the **prd** environment:

| Secret | Required | Description |
|--------|----------|-------------|
| `OCI_SSH_PRIVATE_KEY` | yes | SSH private key for connecting to VM (PEM format) |

---

## 🔧 Setup Instructions

### 1️⃣ Generate SSH Key Pair (if needed)

```bash
ssh-keygen -t ed25519 -f ~/.ssh/gojob_deploy -C "gojob-deploy"
```

### 2️⃣ Add SSH Public Key to VM

On your k3s VM:

```bash
# As the user who will run the deploy
cat >> ~/.ssh/authorized_keys < ~/.ssh/gojob_deploy.pub

chmod 600 ~/.ssh/authorized_keys
```

### 3️⃣ Add Secrets to GitHub

**Repository Settings → Secrets and variables → Actions**

Only one secret is required:

#### Add `OCI_SSH_PRIVATE_KEY`

```bash
# Copy the private key
cat ~/.ssh/id_ed25519 | pbcopy
```

In GitHub:
- Name: `OCI_SSH_PRIVATE_KEY`
- Value: (paste the private key)

### 4️⃣ Add Variables to GitHub

**Repository Settings → Environments → prd** (or as repo variables)

```
OCI_VM_HOST=your-vm.example.com
OCI_VM_USER=ubuntu
OCI_DEPLOY_ROOT=/opt/go
SYNC_URL=https://goanimes.example.com/admin/sync
API_KEY=your-actual-api-key
```

### 5️⃣ Setup K3s VM

Ensure kubeconfig is accessible:

```bash
# Check if k3s is running
sudo k3s kubectl get nodes

# Make kubeconfig accessible to deploy user
sudo cp /etc/rancher/k3s/k3s.yaml ~/.kube/config
sudo chown ubuntu:ubuntu ~/.kube/config
chmod 600 ~/.kube/config

# Test access
kubectl cluster-info
```

---

## 🎯 Deployment Flow

```
git push to main
    ↓
GitHub Actions:
  1. Test job — runs go test
  2. Image job — builds multi-arch Docker image
  3. Deploy job (prd environment):
     - Stamps image tag in cronjob.yaml
     - Packs K8s manifests
     - Copies via SCP to VM
     - SSHs to VM and executes deploy script:
       * Creates .env file from vars/secrets
       * Applies CronJob manifest
       * Creates gojob-env Secret from .env
     ↓
K3s CronJob updated
```

---

## 📝 Alternative: Base64 Environment Block

Instead of individual variables, you can provide:

```bash
# Create .env with all variables
cat > .env << EOF
SYNC_URL=https://your-endpoint/sync
API_KEY=your-api-key
EOF

# Encode as base64
cat .env | base64 | tr -d '\n'
# Output: U1lOQ19VUkw9...
```

Then in GitHub, add as `GOJOB_ENV_B64` variable. The deploy script will decode and merge.

---

## ✅ Verification

After deployment:

```bash
# On the VM
kubectl get cronjobs
kubectl get cronjob gojob-sync -n default -o yaml
kubectl get pods -n default -l app=gojob-sync
kubectl logs -n default -l app=gojob-sync --tail=50

# Check the secret was created
kubectl get secret gojob-env -n default
kubectl get secret gojob-env -n default -o jsonpath='{.data.SYNC_URL}' | base64 -d
```

---

## 🐛 Troubleshooting

### "kubectl: command not found"

On the VM:
```bash
# Add k3s to PATH
export PATH=$PATH:/usr/local/bin
echo "export PATH=$PATH:/usr/local/bin" >> ~/.profile
```

### "KUBECONFIG invalid"

```bash
# Check if file exists and is readable
ls -la /etc/rancher/k3s/k3s.yaml
ls -la ~/.kube/config

# Ensure KUBECONFIG is exported
echo $KUBECONFIG

# Test cluster access
kubectl cluster-info --request-timeout=10s
```

### SSH Connection Fails

```bash
# Test SSH locally
ssh -i ~/.ssh/gojob_deploy ubuntu@your-vm.example.com kubectl cluster-info

# Check GitHub Secrets (OCI_SSH_PRIVATE_KEY format)
# Should be PEM format, no additional wrapping
```

### CronJob Not Updating

```bash
# Check manifest was applied
kubectl get cronjob gojob-sync -n default -o yaml | grep -A2 image

# Check Secret exists
kubectl get secret gojob-env -n default

# Inspect the latest Job
kubectl get jobs -n default -l app=gojob-sync --sort-by=.metadata.creationTimestamp | tail -1
```

---

## 📚 See Also

- [GoJob README](../../README.md) — General project documentation
- [GoAnimes oracle-deploy.yml](https://github.com/wallissonmarinho/GoAnimes/blob/main/.github/workflows/oracle-deploy.yml) — Reference implementation

Ready to go! 🚀
