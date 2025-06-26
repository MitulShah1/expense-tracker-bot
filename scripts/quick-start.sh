#!/bin/bash

# 🚗 Expense Tracker Bot - Quick Start Script
# Simple one-command setup for new users

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🚗 Expense Tracker Bot - Quick Start"
echo "==================================="
echo ""

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo -e "${RED}Error: Please run this script from the project root directory${NC}"
    exit 1
fi

echo -e "${BLUE}This script will:${NC}"
echo "1. ✅ Check prerequisites (Go, Docker)"
echo "2. 📝 Create .env file with default settings"
echo "3. 🗄️  Start PostgreSQL database"
echo "4. 🔧 Install dependencies and run migrations"
echo "5. 🏗️  Build the application"
echo "6. 🧪 Run tests"
echo "7. 🚀 Provide next steps"
echo ""

read -p "Continue? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Setup cancelled."
    exit 0
fi

# Run the main setup script
echo ""
echo "Running setup script..."
./scripts/setup.sh

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}🎉 Quick start completed!${NC}"
    echo ""
    echo -e "${YELLOW}Next steps:${NC}"
    echo "1. Edit .env file and add your Telegram bot token"
    echo "2. Run: ./expense-tracker-bot"
    echo ""
    echo -e "${BLUE}Need help?${NC}"
    echo "- Read the README.md for detailed instructions"
    echo "- Check scripts/test_bot.sh for testing"
    echo "- Use 'make help' for available commands"
else
    echo ""
    echo -e "${RED}❌ Setup failed. Please check the errors above.${NC}"
    exit 1
fi 