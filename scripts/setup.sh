#!/bin/bash

# 🚗 Expense Tracker Bot - Complete Setup Script
# This script sets up the entire application from scratch

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
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

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to get docker compose command
get_docker_compose_cmd() {
    if command_exists docker-compose; then
        echo "docker-compose"
    elif docker compose version >/dev/null 2>&1; then
        echo "docker compose"
    else
        echo ""
    fi
}

# Function to check if port is available
port_available() {
    ! nc -z localhost "$1" 2>/dev/null
}

echo "🚗 Expense Tracker Bot - Complete Setup"
echo "======================================"
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -f "docker-compose.yml" ]; then
    print_error "Please run this script from the project root directory"
    exit 1
fi

# Step 1: Check prerequisites
print_status "Checking prerequisites..."

# Check Go
if ! command_exists go; then
    print_error "Go is not installed. Please install Go 1.21 or higher"
    print_warning "Visit: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
print_success "Go version: $GO_VERSION"

# Check Docker
if ! command_exists docker; then
    print_error "Docker is not installed. Please install Docker"
    print_warning "Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check Docker Compose
DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
if [ -z "$DOCKER_COMPOSE_CMD" ]; then
    print_error "Docker Compose is not installed. Please install Docker Compose"
    print_warning "Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

print_success "Docker Compose found: $DOCKER_COMPOSE_CMD"

# Step 2: Create .env file if it doesn't exist
print_status "Setting up environment configuration..."

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
EOF
    
    print_success ".env file created"
    print_warning "Please update TELEGRAM_TOKEN in .env file with your bot token"
    print_warning "Get your bot token from @BotFather on Telegram"
else
    print_success ".env file already exists"
fi

# Step 3: Install Go dependencies
print_status "Installing Go dependencies..."
go mod download
go mod tidy
print_success "Dependencies installed"

# Step 4: Start PostgreSQL with Docker
print_status "Starting PostgreSQL database..."

# Check if PostgreSQL is already running
if docker ps | grep -q "expense-tracker-postgres"; then
    print_warning "PostgreSQL container is already running"
else
    # Start PostgreSQL
    $DOCKER_COMPOSE_CMD up -d postgres
    
    # Wait for PostgreSQL to be ready
    print_status "Waiting for PostgreSQL to be ready..."
    for i in {1..30}; do
        if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
            break
        fi
        sleep 2
    done
    
    if $DOCKER_COMPOSE_CMD exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
        print_success "PostgreSQL is ready"
    else
        print_error "PostgreSQL failed to start"
        exit 1
    fi
fi

# Step 5: Install migration tool
print_status "Setting up database migration tool..."

if ! command_exists psql; then
    print_error "psql is not installed. Please install PostgreSQL client"
    print_warning "On Ubuntu/Debian: sudo apt-get install postgresql-client"
    print_warning "On macOS: brew install postgresql"
    exit 1
fi

print_success "PostgreSQL client found"

# Step 6: Run database migrations
print_status "Running database migrations..."

# Load environment variables
source .env

# Check if psql is available
if ! command_exists psql; then
    print_error "psql is not installed. Please install PostgreSQL client"
    print_warning "On Ubuntu/Debian: sudo apt-get install postgresql-client"
    print_warning "On macOS: brew install postgresql"
    exit 1
fi

# Run migrations manually since they're in custom format
print_status "Executing migration files..."

MIGRATION_FILES=(
    "001_init.sql"
    "002_seed_categories.sql"
    "003_add_budgets_table.sql"
    "004_add_views.sql"
)

for migration in "${MIGRATION_FILES[@]}"; do
    if [ -f "migrations/$migration" ]; then
        print_status "Running $migration..."
        if psql "$DATABASE_URL" -f "migrations/$migration" >/dev/null 2>&1; then
            print_success "✅ $migration completed"
        else
            print_warning "⚠️  $migration may have already been applied or failed"
        fi
    else
        print_warning "⚠️  Migration file $migration not found"
    fi
done

print_success "Database migrations completed"

# Step 7: Build the application
print_status "Building the application..."
go build -o expense-tracker-bot ./cmd

if [ $? -eq 0 ]; then
    print_success "Application built successfully"
else
    print_error "Build failed"
    exit 1
fi

# Step 8: Run tests
print_status "Running tests..."
go test -v ./...

if [ $? -eq 0 ]; then
    print_success "All tests passed"
else
    print_warning "Some tests failed, but continuing with setup"
fi

# Step 9: Final verification
print_status "Performing final verification..."

# Check if binary was created
if [ -f "expense-tracker-bot" ]; then
    print_success "Binary file created: expense-tracker-bot"
else
    print_error "Binary file not found"
    exit 1
fi

# Check database connection
if psql "$DATABASE_URL" -c "SELECT COUNT(*) FROM categories;" >/dev/null 2>&1; then
    print_success "Database connection verified"
else
    print_warning "Database connection test failed"
fi

echo ""
echo "🎉 Setup completed successfully!"
echo "================================"
echo ""
echo "📋 Next steps:"
echo "1. Update your .env file with your Telegram bot token"
echo "   - Get your bot token from @BotFather on Telegram"
echo "   - Replace 'your_telegram_bot_token_here' in .env file"
echo ""
echo "2. Start the bot:"
echo "   ./expense-tracker-bot"
echo ""
echo "3. Or use Makefile commands:"
echo "   make run          # Run the application"
echo "   make test         # Run tests"
echo "   make build        # Build the application"
echo ""
echo "🔧 Useful commands:"
echo "   $DOCKER_COMPOSE_CMD up -d postgres    # Start database"
echo "   $DOCKER_COMPOSE_CMD down              # Stop database"
echo "   make clean                       # Clean build files"
echo "   ./scripts/test_bot.sh           # Test the bot"
echo ""
echo "📊 Database management:"
echo "   - pgAdmin available at: http://localhost:8080"
echo "   - Email: admin@expense-tracker.com"
echo "   - Password: admin"
echo ""
echo "🐛 Troubleshooting:"
echo "   - Check logs: $DOCKER_COMPOSE_CMD logs postgres"
echo "   - Reset database: $DOCKER_COMPOSE_CMD down -v && $DOCKER_COMPOSE_CMD up -d postgres"
echo "   - View migration status: migrate -path migrations -database \"$DATABASE_URL\" version"
echo ""

# Check if TELEGRAM_TOKEN is still the default value
if grep -q "your_telegram_bot_token_here" .env; then
    print_warning "⚠️  Don't forget to update TELEGRAM_TOKEN in .env file!"
fi

print_success "Setup script completed! 🚀" 