# 📜 Scripts Directory

This directory contains various setup and utility scripts for the Expense Tracker Bot.

## 🚀 Quick Start Scripts

### `quick-start.sh` - One-Command Setup

The easiest way to get started with the bot.

```bash
./scripts/quick-start.sh
```

**What it does:**

- ✅ Checks all prerequisites (Go, Docker)
- 📝 Creates and configures `.env` file
- 🗄️ Starts PostgreSQL database
- 🔧 Installs dependencies and runs migrations
- 🏗️ Builds the application
- 🧪 Runs tests
- 🚀 Provides next steps

**Best for:** New users who want the fastest setup experience.

### `setup.sh` - Complete Setup

Comprehensive setup script with detailed output.

```bash
./scripts/setup.sh
```

**What it does:**

- Everything from `quick-start.sh` plus:
- 🔍 Detailed error checking
- 📊 Database verification
- 🛠️ Migration tool installation
- 📋 Comprehensive status reporting

**Best for:** Users who want detailed setup information and troubleshooting.

## 🗄️ Database Scripts

### `setup_database.sh` - Database Only Setup

Sets up only the database components.

```bash
./scripts/setup_database.sh
```

**What it does:**

- 🗄️ Starts PostgreSQL container
- 🔧 Installs migration tool
- 📊 Runs database migrations
- ✅ Verifies database setup
- 📋 Shows database status

**Best for:** When you only need to set up or reset the database.

## 🧪 Testing Scripts

### `test_bot.sh` - Bot Testing

Interactive testing script for the bot.

```bash
./scripts/test_bot.sh
```

**What it does:**

- 🏗️ Builds the bot
- 🤖 Starts the bot in background
- 📝 Provides test scenarios
- 🛑 Graceful shutdown on Ctrl+C

**Best for:** Testing bot functionality after setup.

## 📋 Test Data

### `test_data.md` - Sample Data

Contains sample expense data for testing.

**What it includes:**

- 🚗 Vehicle expenses (petrol, service)
- 🍽️ Food expenses (groceries, restaurants)
- ⚡ Utility expenses (electricity, water)
- 📊 Test scenarios and flows

**Best for:** Manual testing and development.

## 🔧 Using with Makefile

You can also use these scripts through the Makefile:

```bash
# Quick setup
make setup-quick

# Complete setup
make setup

# Database setup only
make setup-db
```

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
   docker-compose down
   ./scripts/setup.sh
   ```

4. **Migration Tool Not Found**

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

### Reset Everything

To completely reset and start fresh:

```bash
# Stop and remove containers
docker-compose down -v

# Remove build artifacts
make clean

# Run setup again
./scripts/setup.sh
```

## 📝 Script Dependencies

All scripts require:

- **Go 1.21+** - For building and running the application
- **Docker & Docker Compose** - For PostgreSQL database
- **Git** - For cloning the repository
- **Bash** - For running the scripts

## 🔒 Security Notes

- Scripts create a `.env` file with default values
- **Always update `TELEGRAM_TOKEN`** in `.env` before running the bot
- Database credentials are set to default values for development
- Change passwords for production use
