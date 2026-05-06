# 🚀 CI/CD Setup Guide for GoJob

## GitHub Actions Workflow

O arquivo `.github/workflows/deploy.yml` automatiza:
1. **Build** - Cria imagem Docker on push para `main`
2. **Push** - Envia para GitHub Container Registry (GHCR)
3. **Deploy** - Atualiza o CronJob no k3s

## 📋 Configuração Necessária

### 1️⃣ Criar KUBECONFIG

Você precisa adicionar um secret com sua configuração do k3s:

```bash
# Em seu computador com k3s
cat ~/.kube/config | base64 | tr -d '\n' | pbcopy
```

Isso copia o conteúdo encoded em base64 para clipboard.

### 2️⃣ Adicionar Secret no GitHub

1. Vá para: **GitHub Repo → Settings → Secrets and variables → Actions**
2. Clique em **"New repository secret"**
3. Nome: `KUBECONFIG`
4. Cole o valor copiado acima

### 3️⃣ Criar ConfigMap e Secret no k3s

```bash
# Criar ConfigMap com URL do endpoint
kubectl create configmap gojob-config \
  --from-literal=sync-url=https://seu-endpoint/sync \
  -n default

# Criar Secret com API key
kubectl create secret generic gojob-secrets \
  --from-literal=api-key=seu-api-key-aqui \
  -n default
```

Ou aplique com manifesto:

```bash
# Edite deploy/k8s/cronjob.yaml com seus valores
kubectl apply -f deploy/k8s/cronjob.yaml
```

## 🔑 Variáveis de Ambiente

O workflow usa estas variáveis:

| Variável | Origem | Descrição |
|----------|--------|-----------|
| `GITHUB_SHA` | GitHub | Hash do commit |
| `REGISTRY` | Workflow env | GitHub Container Registry (ghcr.io) |
| `IMAGE_NAME` | GitHub | wallissonmarinho/GoJob |
| `KUBECONFIG` | Secret | Arquivo de config do k3s |

## 📦 Image Registry

As imagens são armazenadas em: `ghcr.io/wallissonmarinho/gojob`

Tags automáticas:
- `main` - branch
- `sha-xxxxx` - commit hash
- Semver se usar tags (v1.0.0)

### Tornar Imagem Pública (opcional)

```bash
# Após primeiro push automático:
# 1. GitHub → Seu perfil → Packages
# 2. Selecione gojob → Package settings
# 3. Change visibility → Public
```

## ✅ Fluxo Completo

1. Push código para `main`
   ↓
2. GitHub Actions:
   - Build imagem Docker
   - Push para ghcr.io
   ↓
3. Deploy job:
   - Atualiza manifesto K8s
   - Aplica ao cluster
   - Verifica status
   ↓
4. CronJob roda a cada 10 minutos

## 🐛 Troubleshooting

### Ver logs do workflow
GitHub → Actions → Clique no workflow mais recente

### Testar conexão k3s localmente
```bash
# Usar o mesmo KUBECONFIG que setou no secret
export KUBECONFIG=$(echo $KUBECONFIG_VALUE | base64 -d)
kubectl get nodes
```

### Verificar status do CronJob
```bash
kubectl get cronjob gojob-sync -n default
kubectl get jobs -n default -l app=gojob-sync
kubectl logs -n default -l app=gojob-sync --tail=20
```

## 📝 Próximos Passos

1. ✅ Configure o secret `KUBECONFIG` no GitHub
2. ✅ Crie ConfigMap e Secret no k3s
3. ✅ Faça um push para `main`
4. ✅ Monitore em GitHub → Actions

Ready to go! 🚀
