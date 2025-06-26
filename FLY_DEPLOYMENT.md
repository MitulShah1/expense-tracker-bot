# 🚀 Fly.io Deployment Guide

This guide will help you deploy your Expense Tracker Bot to Fly.io with CI/CD automation.

## 📋 Prerequisites

1. **Fly.io Account**: Sign up at [fly.io](https://fly.io)
2. **Fly CLI**: Install the Fly CLI
3. **GitHub Repository**: Your code should be in a GitHub repository
4. **PostgreSQL Database**: You'll need a database (we'll use Fly Postgres)

## 🔧 Setup Steps

### Step 1: Install Fly CLI

```bash
# macOS
brew install flyctl

# Linux
curl -L https://fly.io/install.sh | sh

# Windows
iwr https://fly.io/install.ps1 -useb | iex
```

### Step 2: Login to Fly.io

```bash
fly auth login
```

### Step 3: Create a Fly.io App

```bash
# Create the app (this will create fly.toml)
fly launch --name expense-tracker-bot --region bom --no-deploy
```

### Step 4: Set Up PostgreSQL Database

```bash
# Create a PostgreSQL database
fly postgres create --name expense-tracker-db --region bom

# Attach the database to your app
fly postgres attach --app expense-tracker-bot expense-tracker-db
```

### Step 5: Configure Environment Variables

```bash
# Set your Telegram bot token
fly secrets set TELEGRAM_TOKEN=your_telegram_bot_token_here

# Set other environment variables
fly secrets set BOT_ID=expense-tracker
fly secrets set LOG_LEVEL=info
fly secrets set IS_DEV_MODE=false
```

### Step 6: Set Up GitHub Secrets

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Add the following secret:
   - **Name**: `FLY_API_TOKEN`
   - **Value**: Get this from Fly.io dashboard or run `fly tokens create deploy`

### Step 7: Deploy

```bash
# Deploy manually (for testing)
fly deploy

# Or use GitHub Actions for manual deployment
# Go to Actions tab in your GitHub repository and run the workflow manually
```

## 🔄 CI/CD Pipeline

The GitHub Actions workflow provides:

### Automatic (on every push/PR):
1. **🧪 Run Tests**: Execute all Go tests with PostgreSQL
2. **🔍 Lint Code**: Run golangci-lint for code quality
3. **🏗️ Build Application**: Compile the application
4. **📦 Upload Artifacts**: Store build artifacts for deployment

### Manual Deployment:
1. **🚀 Deploy**: Manual trigger to deploy to Fly.io
2. **✅ Health Check**: Verify deployment success
3. **📊 Notifications**: Deployment status and URLs

### How to Deploy Manually:

1. **Go to GitHub Repository** → **Actions** tab
2. **Select "CI/CD Pipeline"** workflow
3. **Click "Run workflow"** button
4. **Select "true"** for the deploy option
5. **Click "Run workflow"** to start deployment

This ensures that:
- ✅ Tests and linting run automatically on every change
- ✅ Deployment is controlled and intentional
- ✅ You can review changes before deploying
- ✅ Rollback is easier with manual control

## 📊 Monitoring

### Health Check Endpoints

- **Health**: `https://expense-tracker-bot.fly.dev/health`
- **Root**: `https://expense-tracker-bot.fly.dev/`
- **Metrics**: `https://expense-tracker-bot.fly.dev/metrics`

### Fly.io Dashboard

- **App Status**: [Fly.io Dashboard](https://fly.io/dashboard)
- **Logs**: `fly logs`
- **Status**: `fly status`

## 🛠️ Management Commands

### View Logs

```bash
# View recent logs
fly logs

# Follow logs in real-time
fly logs --follow

# View logs for specific app
fly logs --app expense-tracker-bot
```

### Scale Application

```bash
# Scale to 2 instances
fly scale count 2

# Scale with specific resources
fly scale vm shared-cpu-1x --memory 512
```

### Update Secrets

```bash
# Update Telegram token
fly secrets set TELEGRAM_TOKEN=new_token_here

# View current secrets
fly secrets list
```

### Database Management

```bash
# Connect to database
fly postgres connect --app expense-tracker-bot

# Create backup
fly postgres backup --app expense-tracker-bot

# View database status
fly postgres status --app expense-tracker-bot
```

## 🔧 Troubleshooting

### Common Issues

1. **Deployment Fails**
   ```bash
   # Check logs
   fly logs
   
   # Check status
   fly status
   
   # Redeploy
   fly deploy
   ```

2. **Database Connection Issues**
   ```bash
   # Check database status
   fly postgres status --app expense-tracker-bot
   
   # Reattach database
   fly postgres attach --app expense-tracker-bot expense-tracker-db
   ```

3. **Health Check Fails**
   ```bash
   # Check health endpoint
   curl https://expense-tracker-bot.fly.dev/health
   
   # Check app status
   fly status
   ```

## 📈 Scaling

### Auto-scaling Configuration

The `fly.toml` includes auto-scaling settings:

```toml
[http_service]
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 0
```

### Manual Scaling

```bash
# Scale to specific number of instances
fly scale count 3

# Scale with specific resources
fly scale vm shared-cpu-2x --memory 1024
```

## 🔒 Security

### Environment Variables

- All sensitive data is stored as Fly.io secrets
- Database credentials are automatically managed
- HTTPS is enforced by default

### Network Security

- Application runs in isolated containers
- Database is in private network
- All traffic is encrypted

## 💰 Cost Optimization

### Free Tier

- **3 shared-cpu-1x 256mb VMs**
- **3GB persistent volume storage**
- **160GB outbound data transfer**

### Cost Monitoring

```bash
# View current usage
fly billing show

# Set spending limits
fly billing set-credit-card
```

## 🚀 Production Checklist

- [ ] ✅ Fly.io account created
- [ ] ✅ Fly CLI installed and authenticated
- [ ] ✅ App created with `fly launch`
- [ ] ✅ PostgreSQL database created and attached
- [ ] ✅ Environment variables set as secrets
- [ ] ✅ GitHub secrets configured
- [ ] ✅ CI/CD pipeline working
- [ ] ✅ Health checks passing
- [ ] ✅ Monitoring set up
- [ ] ✅ Backup strategy configured

## 📞 Support

- **Fly.io Documentation**: [fly.io/docs](https://fly.io/docs)
- **Fly.io Community**: [community.fly.io](https://community.fly.io)
- **GitHub Issues**: Create an issue in your repository

---

**Happy deploying! 🚀** 