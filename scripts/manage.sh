#!/bin/bash

# Expense Tracker Bot - Unified Management Script
# Usage:
#   ./scripts/manage.sh [quick|full|db|help]

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}
print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}
print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}
print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

command_exists() {
    command -v "$1" >/dev/null 2>&1
}
get_docker_compose_cmd() {
    if command_exists docker-compose; then
        echo "docker-compose -f docker/docker-compose.yml"
    elif docker compose version >/dev/null 2>&1; then
        echo "docker compose -f docker/docker-compose.yml"
    else
        echo ""
    fi
}

show_help() {
    echo "\nUsage: $0 [quick|full|db|help]"
    echo "  quick   - One-command quick setup (recommended for new users)"
    echo "  full    - Complete setup with detailed checks"
    echo "  db      - Database-only setup/reset"
    echo "  help    - Show this help message\n"
}

# --- Shared Steps ---
check_prerequisites() {
    print_status "Checking prerequisites..."
    if ! command_exists go; then
        print_error "Go is not installed. Please install Go 1.21 or higher."
        exit 1
    fi
    print_success "Go version: $(go version | awk '{print $3}')"
    if ! command_exists docker; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi
    DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
    if [ -z "$DOCKER_COMPOSE_CMD" ]; then
        print_error "Docker Compose is not installed."
        exit 1
    fi
    print_success "Docker Compose found: $DOCKER_COMPOSE_CMD"
}

create_env_file() {
    if [ ! -f ".env" ]; then
        print_warning ".env file not found. Creating from template..."
        cat > .env << EOF
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
# Vector Search (optional)
VECTOR_SEARCH_ENABLED=true
SIMILARITY_THRESHOLD=0.7
MAX_SEARCH_RESULTS=10
EOF
        print_success ".env file created. Please update TELEGRAM_TOKEN."
    else
        print_success ".env file already exists."
    fi
}

install_go_deps() {
    print_status "Installing Go dependencies..."
    go mod download && go mod tidy
    print_success "Dependencies installed."
}

start_postgres() {
    print_status "Starting PostgreSQL database..."
    if docker ps | grep -q "expense-tracker-postgres"; then
        print_warning "PostgreSQL container is already running."
    else
        $DOCKER_COMPOSE_CMD up -d postgres
        print_status "Waiting for PostgreSQL to be ready..."
        for i in {1..30}; do
            if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
                break
            fi
            sleep 2
        done
        if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
            print_success "PostgreSQL is ready."
        else
            print_error "PostgreSQL failed to start."
            exit 1
        fi
    fi
}

run_migrations() {
    print_status "Running database migrations..."
    source .env
    if ! command_exists psql; then
        print_error "psql is not installed. Please install PostgreSQL client."
        exit 1
    fi
    MIGRATION_FILES=(
        "001_init.sql"
        "002_seed_categories.sql"
        "003_add_budgets_table.sql"
        "004_add_views.sql"
        "005_add_pgvector.sql"
    )
    for migration in "${MIGRATION_FILES[@]}"; do
        if [ -f "migrations/$migration" ]; then
            print_status "Running $migration..."
            if psql "$DATABASE_URL" -f "migrations/$migration" >/dev/null 2>&1; then
                print_success "✅ $migration completed."
            else
                print_warning "⚠️  $migration may have already been applied or failed."
            fi
        else
            print_warning "⚠️  Migration file $migration not found."
        fi
    done
    print_success "Database migrations completed."
}

build_app() {
    print_status "Building the application..."
    go build -o expense-tracker-bot ./cmd
    if [ $? -eq 0 ]; then
        print_success "Application built successfully."
    else
        print_error "Build failed."
        exit 1
    fi
}

run_tests() {
    print_status "Running tests..."
    go test -v ./...
    if [ $? -eq 0 ]; then
        print_success "All tests passed."
    else
        print_warning "Some tests failed, but continuing with setup."
    fi
}

# --- Modes ---
quick_start() {
    echo -e "${BLUE}Expense Tracker Bot - Quick Start${NC}\n==================================="
    check_prerequisites
    create_env_file
    install_go_deps
    start_postgres
    run_migrations
    build_app
    run_tests
    echo -e "\n${GREEN}🎉 Quick start completed!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Edit .env file and add your Telegram bot token\n2. Run: ./expense-tracker-bot\n"
}

full_setup() {
    echo -e "${BLUE}Expense Tracker Bot - Complete Setup${NC}\n======================================"
    check_prerequisites
    create_env_file
    install_go_deps
    start_postgres
    run_migrations
    build_app
    run_tests
    echo -e "\n${GREEN}🎉 Setup completed successfully!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Update your .env file with your Telegram bot token\n2. Start the bot: ./expense-tracker-bot\n"
}

db_only() {
    echo -e "${BLUE}Expense Tracker Bot - Database Only Setup${NC}\n=========================================="
    check_prerequisites
    create_env_file
    start_postgres
    run_migrations
    echo -e "\n${GREEN}🎉 Database setup completed successfully!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Update your .env file with the correct TELEGRAM_TOKEN\n2. Build and run the bot: go build -o expense-tracker-bot cmd/main.go\n3. Start the bot: ./expense-tracker-bot\n"
}

# --- Main ---
MODE="$1"
case "$MODE" in
    quick)
        quick_start
        ;;
    full)
        full_setup
        ;;
    db)
        db_only
        ;;
    help|--help|-h|"")
        show_help
        ;;
    *)
        print_error "Unknown command: $MODE"
        show_help
        exit 1
        ;;
esac 