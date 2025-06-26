# 🚀 Quick Deployment Guide

## 🔄 CI/CD Workflow

### Automatic Actions (on every push/PR):

- ✅ **Tests**: All Go tests run automatically
- ✅ **Linting**: Code quality checks with golangci-lint
- ✅ **Build**: Application compilation
- ✅ **Artifacts**: Build artifacts stored for deployment

### Manual Deployment:

- 🚀 **Deploy**: Manual trigger required for deployment

## 📋 How to Deploy

### Method 1: GitHub Actions (Recommended)

1. **Push your changes** to the `main` branch
2. **Go to GitHub** → Your repository → **Actions** tab
3. **Select "CI/CD Pipeline"** workflow
4. **Click "Run workflow"** button
5. **Set deploy option to "true"**
6. **Click "Run workflow"** to start deployment

### Method 2: Direct Fly.io Deployment

```bash
# Deploy directly to Fly.io
fly deploy
```

## 🔍 Pre-Deployment Checklist

Before deploying, ensure:

- [ ] ✅ All tests pass locally: `make test`
- [ ] ✅ Code is linted: `make lint`
- [ ] ✅ Environment variables are set: `fly secrets list`
- [ ] ✅ Database is running: `fly postgres status`
- [ ] ✅ Telegram bot token is configured

## 📊 Post-Deployment Verification

After deployment, check:

```bash
# Check app status
fly status

# View logs
fly logs

# Test health endpoint
curl https://expense-tracker-bot.fly.dev/health

# Check app URL
open https://expense-tracker-bot.fly.dev
```

## 🛠️ Troubleshooting

### Deployment Fails:

1. **Check GitHub Actions logs** for test/lint failures
2. **Verify Fly.io secrets**: `fly secrets list`
3. **Check database status**: `fly postgres status`
4. **View deployment logs**: `fly logs`

### App Not Responding:

1. **Check health endpoint**: `/health`
2. **View app logs**: `fly logs`
3. **Check app status**: `fly status`
4. **Restart app**: `fly apps restart`

## 🔄 Rollback

If deployment fails:

```bash
# List deployments
fly releases

# Rollback to previous version
fly deploy --image-label v1

# Or restart the app
fly apps restart
```

## 📈 Monitoring

### Health Endpoints:

- **Health**: `https://expense-tracker-bot.fly.dev/health`
- **Root**: `https://expense-tracker-bot.fly.dev/`
- **Metrics**: `https://expense-tracker-bot.fly.dev/metrics`

### Fly.io Commands:

```bash
# View app status
fly status

# View logs
fly logs

# Scale app
fly scale count 2

# Update secrets
fly secrets set TELEGRAM_TOKEN=new_token
```

---

**Happy deploying! 🚀**