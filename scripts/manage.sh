#!/bin/bash

# Expense Tracker Bot - Unified Management Script
# Usage:
#   ./scripts/manage.sh [quick|full|db|start|stop|restart|logs|status|help]

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
    echo "\nUsage: $0 [quick|full|db|start|stop|restart|logs|status|help]"
    echo "  quick    - One-command quick setup (recommended for new users)"
    echo "  full     - Complete setup with detailed checks"
    echo "  db       - Database-only setup/reset"
    echo "  start    - Start all services (postgres, pgadmin, app)"
    echo "  stop     - Stop all services"
    echo "  restart  - Restart all services"
    echo "  logs     - Show logs from all services"
    echo "  status   - Show status of all services"
    echo "  help     - Show this help message\n"
}

# --- Shared Steps ---
check_prerequisites() {
    print_status "Checking prerequisites..."
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
POSTGRES_DB=expense_tracker
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
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
# pgAdmin Configuration (optional)
PGADMIN_DEFAULT_EMAIL=admin@expense-tracker.com
PGADMIN_DEFAULT_PASSWORD=admin
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

start_services() {
    print_status "Starting all services with Docker Compose..."
    $DOCKER_COMPOSE_CMD up -d
    if [ $? -eq 0 ]; then
        print_success "All services started successfully."
        print_status "Waiting for services to be ready..."
        sleep 10
        show_service_status
    else
        print_error "Failed to start services."
        exit 1
    fi
}

stop_services() {
    print_status "Stopping all services..."
    $DOCKER_COMPOSE_CMD down
    print_success "All services stopped."
}

restart_services() {
    print_status "Restarting all services..."
    $DOCKER_COMPOSE_CMD restart
    print_success "All services restarted."
}

show_service_status() {
    print_status "Service Status:"
    echo ""
    $DOCKER_COMPOSE_CMD ps
    echo ""
    print_status "Service URLs:"
    echo "  📊 pgAdmin: http://localhost:8080 (admin@expense-tracker.com / admin)"
    echo "  🤖 Bot API: http://localhost:8081/health"
    echo "  🗄️  Database: localhost:5432"
}

show_logs() {
    print_status "Showing logs from all services..."
    $DOCKER_COMPOSE_CMD logs -f
}

run_migrations() {
    print_status "Running database migrations..."
    # Wait for postgres to be ready
    print_status "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
            break
        fi
        sleep 2
    done
    
    if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
        print_success "PostgreSQL is ready."
        
        # Run migrations
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
                if $DOCKER_COMPOSE_CMD exec -T postgres psql -U postgres -d expense_tracker -f "/docker-entrypoint-initdb.d/$migration" >/dev/null 2>&1; then
                    print_success "✅ $migration completed."
                else
                    print_warning "⚠️  $migration may have already been applied or failed."
                fi
            else
                print_warning "⚠️  Migration file $migration not found."
            fi
        done
        print_success "Database migrations completed."
    else
        print_error "PostgreSQL failed to start."
        exit 1
    fi
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
    start_services
    run_migrations
    build_app
    run_tests
    echo -e "\n${GREEN}🎉 Quick start completed!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Edit .env file and add your Telegram bot token\n2. The bot should be running at http://localhost:8081\n3. Access pgAdmin at http://localhost:8080\n"
}

full_setup() {
    echo -e "${BLUE}Expense Tracker Bot - Complete Setup${NC}\n======================================"
    check_prerequisites
    create_env_file
    install_go_deps
    start_services
    run_migrations
    build_app
    run_tests
    echo -e "\n${GREEN}🎉 Setup completed successfully!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Update your .env file with your Telegram bot token\n2. The bot should be running at http://localhost:8081\n3. Access pgAdmin at http://localhost:8080\n"
}

db_only() {
    echo -e "${BLUE}Expense Tracker Bot - Database Only Setup${NC}\n=========================================="
    check_prerequisites
    create_env_file
    start_services
    run_migrations
    echo -e "\n${GREEN}🎉 Database setup completed successfully!${NC}"
    echo -e "\n${YELLOW}Next steps:${NC}\n1. Update your .env file with the correct TELEGRAM_TOKEN\n2. The bot should be running at http://localhost:8081\n3. Access pgAdmin at http://localhost:8080\n"
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
    start)
        check_prerequisites
        start_services
        ;;
    stop)
        check_prerequisites
        stop_services
        ;;
    restart)
        check_prerequisites
        restart_services
        ;;
    logs)
        check_prerequisites
        show_logs
        ;;
    status)
        check_prerequisites
        show_service_status
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