#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🚀 Starting Expense Tracker Bot Test Suite"

# Check if .env file exists
if [ ! -f .env ]; then
    echo -e "${RED}Error: .env file not found${NC}"
    echo "Please create a .env file with the following variables:"
    echo "TELEGRAM_TOKEN=your_bot_token"
    echo "DATABASE_URL=postgres://username:password@localhost:5432/expense_tracker"
    echo "BOT_ID=expense-tracker"
    echo "LOG_LEVEL=debug"
    echo "IS_DEV_MODE=true"
    exit 1
fi

# Build the bot
echo "📦 Building the bot..."
go build -o expense-tracker-bot cmd/main.go

# Start the bot in the background
echo "🤖 Starting the bot..."
./expense-tracker-bot &
BOT_PID=$!

# Wait for bot to start
sleep 2

echo -e "${GREEN}✅ Bot is running with PID: $BOT_PID${NC}"
echo "
📝 Test Cases to Verify:

1. Basic Commands:
   - /start - Should show welcome message
   - /help - Should show help message
   - /cancel - Should cancel current operation

2. Add Expense Flow:
   - /add
   - Select category group (e.g., Vehicle)
   - Select category (e.g., Petrol)
   - Enter vehicle type (CAR/BIKE)
   - Enter odometer reading
   - Enter petrol price
   - Enter total price
   - Add optional notes

3. List Expenses:
   - /list - Should show placeholder message (database not connected)

4. Edit Expense:
   - /edit - Should show placeholder message (database not connected)

5. Delete Expense:
   - /delete - Should show placeholder message (database not connected)

6. Reports:
   - /report - Should show placeholder message (database not connected)
   - /dashboard - Should show placeholder message (database not connected)

7. Error Handling:
   - Invalid inputs
   - Rate limiting
   - Network issues

Note: Database features are currently placeholders. Connect PostgreSQL to enable full functionality.

To stop the bot, press Ctrl+C
"

# Wait for user to press Ctrl+C
trap "kill $BOT_PID; echo -e '\n${GREEN}✅ Bot stopped${NC}'; exit 0" INT
wait 