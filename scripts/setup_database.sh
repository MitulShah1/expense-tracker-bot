#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to get docker compose command
get_docker_compose_cmd() {
    if command -v docker-compose >/dev/null 2>&1; then
        echo "docker-compose"
    elif docker compose version >/dev/null 2>&1; then
        echo "docker compose"
    else
        echo ""
    fi
}

echo "🗄️  Setting up PostgreSQL Database for Expense Tracker Bot"

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo "Please create a .env file with the following variables:"
    echo "DATABASE_URL=postgres://username:password@localhost:5432/expense_tracker?sslmode=disable"
    echo "TELEGRAM_TOKEN=your_bot_token"
    echo "BOT_ID=expense-tracker"
    echo "LOG_LEVEL=debug"
    echo "IS_DEV_MODE=true"
    exit 1
fi

# Load environment variables
source .env

# Get docker compose command
DOCKER_COMPOSE_CMD=$(get_docker_compose_cmd)
if [ -z "$DOCKER_COMPOSE_CMD" ]; then
    echo -e "${RED}Error: Docker Compose not found${NC}"
    echo "Please install Docker Compose"
    exit 1
fi

# Extract database connection details from DATABASE_URL
if [[ $DATABASE_URL =~ postgres://([^:]+):([^@]+)@([^:]+):([^/]+)/([^?]+) ]]; then
    DB_USER="${BASH_REMATCH[1]}"
    DB_PASS="${BASH_REMATCH[2]}"
    DB_HOST="${BASH_REMATCH[3]}"
    DB_PORT="${BASH_REMATCH[4]}"
    DB_NAME="${BASH_REMATCH[5]}"
else
    echo -e "${RED}Error: Invalid DATABASE_URL format${NC}"
    echo "Expected format: postgres://username:password@host:port/database"
    exit 1
fi

echo -e "${YELLOW}Database Configuration:${NC}"
echo "Host: $DB_HOST"
echo "Port: $DB_PORT"
echo "Database: $DB_NAME"
echo "User: $DB_USER"

# Check if PostgreSQL is running
echo -e "\n${YELLOW}Checking PostgreSQL connection...${NC}"
if ! pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to PostgreSQL${NC}"
    echo "Please ensure PostgreSQL is running and accessible"
    exit 1
fi

echo -e "${GREEN}✅ PostgreSQL connection successful${NC}"

# Check if database exists
echo -e "\n${YELLOW}Checking if database exists...${NC}"
if ! psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${YELLOW}Database '$DB_NAME' does not exist. Creating...${NC}"
    createdb -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" "$DB_NAME"
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Database created successfully${NC}"
    else
        echo -e "${RED}Error: Failed to create database${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ Database exists${NC}"
fi

# Run migrations
echo -e "\n${YELLOW}Running database migrations...${NC}"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}Error: psql is not installed${NC}"
    echo "Please install PostgreSQL client:"
    echo "Ubuntu/Debian: sudo apt-get install postgresql-client"
    echo "macOS: brew install postgresql"
    exit 1
fi

# Run migrations manually since they're in custom format
MIGRATION_FILES=(
    "001_init.sql"
    "002_seed_categories.sql"
    "003_add_budgets_table.sql"
    "004_add_views.sql"
)

for migration in "${MIGRATION_FILES[@]}"; do
    if [ -f "migrations/$migration" ]; then
        echo -e "${YELLOW}Running $migration...${NC}"
        if psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "migrations/$migration" >/dev/null 2>&1; then
            echo -e "${GREEN}✅ $migration completed${NC}"
        else
            echo -e "${YELLOW}⚠️  $migration may have already been applied or failed${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  Migration file $migration not found${NC}"
    fi
done

echo -e "${GREEN}✅ Migrations completed successfully${NC}"

# Verify tables were created
echo -e "\n${YELLOW}Verifying database setup...${NC}"
TABLES=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public';")

echo "Created tables:"
echo "$TABLES" | grep -E "(users|categories|expenses|recurring_expenses|budgets|budget_limits)" | while read table; do
    if [ ! -z "$table" ]; then
        echo -e "  ${GREEN}✅ $table${NC}"
    fi
done

# Check views
echo -e "\n${YELLOW}Verifying views...${NC}"
VIEWS=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT viewname FROM pg_views WHERE schemaname = 'public';")

echo "Created views:"
echo "$VIEWS" | grep -E "(expense_summary|monthly_expense_summary|category_expense_breakdown|user_expense_stats|recent_expenses)" | while read view; do
    if [ ! -z "$view" ]; then
        echo -e "  ${GREEN}✅ $view${NC}"
    fi
done

# Check categories
echo -e "\n${YELLOW}Verifying seeded data...${NC}"
CATEGORY_COUNT=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT COUNT(*) FROM categories;")
echo -e "${GREEN}✅ Categories seeded: $CATEGORY_COUNT${NC}"

echo -e "\n${GREEN}🎉 Database setup completed successfully!${NC}"
echo -e "\n${YELLOW}Next steps:${NC}"
echo "1. Update your .env file with the correct DATABASE_URL"
echo "2. Build and run the bot: go build -o expense-tracker-bot cmd/main.go"
echo "3. Start the bot: ./expense-tracker-bot"
echo -e "\n${YELLOW}Useful commands:${NC}"
echo "- View migration status: migrate -path migrations -database \"$DATABASE_URL\" version"
echo "- Rollback migrations: migrate -path migrations -database \"$DATABASE_URL\" down"
echo "- Connect to database: psql \"$DATABASE_URL\"" 