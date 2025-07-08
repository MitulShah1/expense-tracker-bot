# 🚀 Setup Guide - Expense Tracker Bot

This guide provides multiple ways to set up your Expense Tracker Bot, from the simplest one-command setup to detailed manual configuration.

## 🎯 Quick Start (Recommended)

For the fastest setup experience:

```bash
# 1. Clone the repository
git clone https://github.com/MitulShah1/expense-tracker-bot.git
cd expense-tracker-bot

# 2. Run the quick start script
./scripts/quick-start.sh
```

This will automatically:

- ✅ Check all prerequisites (Go, Docker)
- 📝 Create and configure your `.env` file
- 🗄️ Start PostgreSQL database with Docker
- 🔧 Install dependencies and run migrations
- 🏗️ Build the application
- 🧪 Run tests
- 🚀 Provide next steps

## 🔧 Alternative Setup Methods

### Method 1: Complete Setup Script

For detailed setup with comprehensive output:

```bash
./scripts/setup.sh
```

### Method 2: Using Makefile

```bash
# Quick setup with confirmation
make setup-quick

# Complete setup
make setup

# Database setup only
make setup-db
```

### Method 3: Manual Setup

If you prefer step-by-step control:

```bash
# 1. Install dependencies
make deps

# 2. Start database
docker compose -f docker/docker-compose.yml up -d postgres

# 3. Run database setup
./scripts/setup_database.sh

# 4. Build application
make build

# 5. Run application
make run
```

## 📋 Prerequisites

Before running any setup script, ensure you have:

- **Go 1.21+** - [Download here](https://golang.org/dl/)
- **Docker & Docker Compose** - [Install here](https://docs.docker.com/get-docker/) (see docker/ directory)
- **Git** - [Install here](https://git-scm.com/)

## ⚙️ Configuration

### Environment Variables

The setup scripts will create a `.env` file with default values:

```env
# Telegram Bot Configuration
TELEGRAM_TOKEN=your_telegram_bot_token_here
BOT_ID=expense-tracker

# Database Configuration
DATABASE_URL=postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable

# Application Configuration
LOG_LEVEL=info
IS_DEV_MODE=true

# Database Connection Pool (optional)
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=5m
```

### Getting Your Telegram Bot Token

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Copy the token and update `TELEGRAM_TOKEN` in your `.env` file

## 🎯 After Setup

Once setup is complete:

1. **Update your bot token**: Edit `.env` file and replace `your_telegram_bot_token_here`
2. **Start the bot**: `./expense-tracker-bot` or `make run`
3. **Test the bot**: Run `./expense-tracker-bot` to start the bot and test functionality
4. **Access pgAdmin**: Visit http://localhost:8080 (admin@expense-tracker.com / admin)

## 🔧 Available Scripts

| Script | Purpose | Best For |
|--------|---------|----------|
| `./scripts/quick-start.sh` | One-command setup with confirmation | New users |
| `./scripts/setup.sh` | Complete automated setup | Detailed setup info |
| `./scripts/setup_database.sh` | Database-only setup | Database management |


## 🛠️ Troubleshooting

### Common Issues

1. **Permission Denied**

   ```bash
   chmod +x scripts/*.sh
   ```

2. **Docker Not Running**

   ```bash
   sudo systemctl start docker
   ```

3. **Port Already in Use**

   ```bash
   docker compose -f docker/docker-compose.yml down
   ./scripts/setup.sh
   ```

4. **psql Not Found**

   ```bash
   # Ubuntu/Debian
   sudo apt-get install postgresql-client
   
   # macOS
   brew install postgresql
   ```

### Reset Everything

To completely reset and start fresh:

```bash
# Stop and remove containers
docker compose -f docker/docker-compose.yml down -v

# Remove build artifacts
make clean

# Run setup again
./scripts/setup.sh
```

## 📊 Database Management

### Accessing the Database

- **pgAdmin**: http://localhost:8080
  - Email: admin@expense-tracker.com
  - Password: admin

- **Direct Connection**:

  ```bash
  psql postgres://postgres:password@localhost:5432/expense_tracker
  ```

### Database Commands

```bash
# Start database
docker compose -f docker/docker-compose.yml up -d postgres

# Stop database
docker compose -f docker/docker-compose.yml down

# View logs
docker compose -f docker/docker-compose.yml logs postgres

# Reset database
docker compose -f docker/docker-compose.yml down -v && docker compose -f docker/docker-compose.yml up -d postgres
```

## 🧪 Testing

### Run Tests

```bash
# All tests
make test

# Tests with coverage
make test-coverage

# Test the bot
./scripts/test_bot.sh
```

### Sample Test Data

See `scripts/test_data.md` for sample expenses and test scenarios.

## 🚀 Production Deployment

For production deployment:

1. **Update environment variables**:
   - Set `IS_DEV_MODE=false`
   - Use strong database passwords
   - Configure proper logging levels

2. **Security considerations**:
   - Change default database credentials
   - Use SSL for database connections
   - Set up proper backup strategies

3. **Monitoring**:
   - Set up application monitoring
   - Configure log aggregation
   - Implement health checks

## 📚 Additional Resources

- **Main README**: [README.md](README.md)
- **Scripts Documentation**: [scripts/README.md](scripts/README.md)
- **Database Migrations**: [migrations/README.md](migrations/README.md)
- **Code of Conduct**: [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## 🆘 Getting Help

If you encounter issues:

1. Check the troubleshooting section above
2. Review the logs: `docker compose -f docker/docker-compose.yml logs postgres`
3. Check the test data: `scripts/test_data.md`
4. Create an issue on GitHub
5. Review the documentation in the `docs/` directory

---

**Happy coding! 🚀**