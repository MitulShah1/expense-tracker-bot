# 📜 Scripts Directory

This directory contains the unified management script for the Expense Tracker Bot.

## 🚀 Main Management Script

### `manage.sh` - Unified Setup & Management

The single entry point for all setup and management tasks.

```bash
./scripts/manage.sh [quick|full|db|help]
```

**Commands:**

- `quick` — One-command quick setup (recommended for new users)
- `full` — Complete setup with detailed checks
- `db` — Database-only setup/reset
- `help` — Show usage instructions

**Examples:**

- Quick start (recommended):

  ```bash
  ./scripts/manage.sh quick
  ```

- Full setup with detailed checks:

  ```bash
  ./scripts/manage.sh full
  ```

- Database-only setup/reset:

  ```bash
  ./scripts/manage.sh db
  ```

- Show help:

  ```bash
  ./scripts/manage.sh help
  ```

## 🔧 Using with Makefile

You can also use these scripts through the Makefile (update Makefile to use manage.sh):

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
   docker-compose -f docker/docker-compose.yml down
   ./scripts/manage.sh full
   ```

4. **Migration Tool Not Found**

   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

### Reset Everything

To completely reset and start fresh:

```bash
# Stop and remove containers
docker-compose -f docker/docker-compose.yml down -v

# Remove build artifacts
make clean

# Run setup again
./scripts/manage.sh full
```

## 📝 Script Dependencies

All scripts require:

- **Go 1.21+** - For building and running the application
- **Docker & Docker Compose** - For PostgreSQL database (see docker/ directory)
- **Git** - For cloning the repository
- **Bash** - For running the scripts

## 🔒 Security Notes

- Scripts create a `.env` file with default values
- **Always update `TELEGRAM_TOKEN`** in `.env` before running the bot
- Database credentials are set to default values for development
- Change passwords for production use
