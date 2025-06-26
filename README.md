# 🚗 Expense Tracker Bot

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-12+-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://www.docker.com/)
[![Go Report Card](https://goreportcard.com/badge/github.com/MitulShah1/expense-tracker-bot)](https://goreportcard.com/report/github.com/MitulShah1/expense-tracker-bot)

> A robust, enterprise-grade Telegram bot for tracking personal expenses with
> PostgreSQL backend, built with Go. Features comprehensive expense management,
> category organization, vehicle expense tracking, and detailed analytics.

[🚀 Quick Start](#-quick-start) • [🌟 Features](#-features) • [🏗️ Architecture](#️-architecture) • [🧪 Testing](#-testing) • [🔧 Development](#-development) • [🚀 Deployment](#-deployment)

## 🌟 Features

### 💰 Core Functionality

- **📝 Expense Tracking**: Add, edit, delete, and list expenses with ease
- **📂 Category Management**: Organized expense categories with emojis
- **🚗 Vehicle Expenses**: Special handling for fuel, service, and maintenance costs
- **📊 Reports & Analytics**: Generate comprehensive expense reports and statistics
- **📈 Dashboard**: Visual overview of spending patterns and trends

### 🏢 Enterprise Features

- **🛡️ Robust Error Handling**: Custom error types with proper error categorization
- **✅ Input Validation**: Comprehensive validation for all user inputs
- **🏗️ Service Layer Architecture**: Clean separation of business logic
- **🔧 Middleware Support**: Cross-cutting concerns like logging, rate limiting, and metrics
- **🧪 Comprehensive Testing**: Unit tests with mocking and test coverage
- **📝 Structured Logging**: JSON logging with context and correlation IDs
- **🗄️ Database Abstraction**: Interface-based database layer with PostgreSQL support
- **⚙️ Configuration Management**: Environment-based configuration with validation
- **🚦 Rate Limiting**: Built-in rate limiting to prevent abuse
- **📊 Metrics Collection**: Application metrics for monitoring

## 🏗️ Architecture

### High-Level Architecture

```text
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Telegram API  │    │   Bot Handler   │    │  Service Layer  │
│                 │◄──►│                 │◄──►│                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Middleware    │    │   Validation    │
                       │                 │    │                 │
                       └─────────────────┘    └─────────────────┘
                                │                       │
                                ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Database Layer │    │   Error Handler │
                       │                 │    │                 │
                       └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   PostgreSQL    │
                       │                 │
                       └─────────────────┘
```

## 🚀 Quick Start

### 🚀 One-Command Setup (Recommended)

For the fastest setup experience, use our automated setup script:

```bash
# Clone the repository
git clone https://github.com/MitulShah1/expense-tracker-bot.git
cd expense-tracker-bot

# Run the quick start script
./scripts/quick-start.sh
```

This script will automatically:
- ✅ Check all prerequisites (Go, Docker)
- 📝 Create and configure your `.env` file
- 🗄️ Start PostgreSQL database with Docker
- 🔧 Install dependencies and run migrations
- 🏗️ Build the application
- 🧪 Run tests
- 🚀 Provide next steps

### 📋 Prerequisites

- [Go 1.21](https://golang.org/dl/) or higher
- [Docker](https://www.docker.com/get-started) and Docker Compose
- [Git](https://git-scm.com/) (for cloning)

### ⚙️ Manual Setup

If you prefer manual setup or the automated script doesn't work:

#### Step 1: Clone the Repository

```bash
git clone https://github.com/MitulShah1/expense-tracker-bot.git
cd expense-tracker-bot
```

#### Step 2: Run Complete Setup Script

```bash
# Make scripts executable
chmod +x scripts/*.sh

# Run the complete setup
./scripts/setup.sh
```

#### Step 3: Manual Environment Configuration

1. Create the environment file:

   ```bash
   cp .env.example .env  # if .env.example exists
   # or create manually:
   ```

2. Configure your environment variables:

   ```bash
   # Telegram Bot Configuration
   TELEGRAM_TOKEN=your_telegram_bot_token
   BOT_ID=your_bot_id
   
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

   > **Note**: You'll need to create a Telegram bot first. Visit [@BotFather](https://t.me/botfather) on Telegram to create your bot and get the token.

#### Step 4: Manual Database Setup

1. Start PostgreSQL using Docker Compose:

   ```bash
   docker-compose up -d postgres
   ```

2. Run database migrations:

   ```bash
   make migrate
   ```

#### Step 5: Manual Application Setup

1. Install dependencies:

   ```bash
   make deps
   ```

2. Build the application:

   ```bash
   make build
   ```

3. Run the application:

   ```bash
   make run
   ```

### 🎯 After Setup

Once setup is complete:

1. **Update your bot token**: Edit `.env` file and replace `your_telegram_bot_token_here` with your actual bot token
2. **Start the bot**: `./expense-tracker-bot` or `make run`
3. **Test the bot**: Use `./scripts/test_bot.sh` to test functionality
4. **Access pgAdmin**: Visit http://localhost:8080 (admin@expense-tracker.com / admin)

### 🔧 Available Scripts

- `./scripts/quick-start.sh` - One-command setup with confirmation
- `./scripts/setup.sh` - Complete automated setup
- `./scripts/setup_database.sh` - Database-only setup
- `./scripts/test_bot.sh` - Test the bot functionality

## 🧪 Testing

### 🏃‍♂️ Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test package
go test ./internal/services/...
```

### 📊 Test Coverage

The application includes comprehensive unit tests with:

- ✅ Service layer testing with mocked dependencies
- ✅ Input validation testing
- ✅ Error handling testing
- ✅ Database operation testing

## 🔧 Development

### 🎯 Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run all quality checks
make all
```

### 🚀 Adding New Features

1. **Service Layer**: Add business logic in `internal/services/`
2. **Validation**: Add input validation in `internal/validation/`
3. **Database**: Add database operations in `internal/database/`
4. **Tests**: Add corresponding tests in the same package
5. **Documentation**: Update this README and add inline documentation

### ⚠️ Error Handling

The application uses custom error types for better error categorization:

- `ValidationError`: Input validation errors
- `NotFoundError`: Resource not found errors
- `DatabaseError`: Database operation errors
- `TelegramError`: Telegram API errors
- `UnauthorizedError`: Authorization errors
- `RateLimitError`: Rate limiting errors

## 📊 Monitoring & Observability

### 📝 Logging

The application uses structured JSON logging with:

- 🔗 Request correlation IDs
- 👤 User context
- ⚡ Performance metrics
- ❌ Error details

### 📈 Metrics

Built-in metrics collection for:

- 📊 Request counts
- ❌ Error rates
- ⏱️ Response times
- 👥 Active users

### 🏥 Health Checks

Database connectivity and application health monitoring.

## 🔒 Security

### 🛡️ Input Validation

- ✅ Comprehensive validation for all user inputs
- 🛡️ SQL injection prevention through parameterized queries
- 🚦 Rate limiting to prevent abuse
- 🔐 User authorization checks

### 🔐 Data Protection

- 🔒 Secure database connections
- ⚙️ Environment-based configuration
- 🚫 No sensitive data in logs

## 🚀 Deployment

### 🐳 Docker Deployment

```bash
# Build Docker image
docker build -t expense-tracker-bot .

# Run with Docker Compose
docker-compose up -d
```

### 🏭 Production Considerations

- ⚙️ Use environment-specific configurations
- 📝 Set up proper logging aggregation
- 🗄️ Configure database connection pooling
- 📊 Implement monitoring and alerting
- 💾 Set up backup strategies

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Contribution Guidelines

1. 🍴 Fork the repository
2. 🌿 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. ✏️ Make your changes
4. 🧪 Add tests for new functionality
5. ✅ Ensure all tests pass (`make test`)
6. 📤 Commit your changes (`git commit -m 'Add some amazing feature'`)
7. 🚀 Push to the branch (`git push origin feature/amazing-feature`)
8. 📋 Open a Pull Request

### 📏 Code Standards

- 📚 Follow Go coding standards and conventions
- 🧪 Add comprehensive tests for new functionality
- 📖 Update documentation for any new features
- 💬 Use meaningful commit messages
- 🔍 Ensure code passes linting (`make lint`)

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

Need help? Here are your options:

- 🐛 [Create an issue](https://github.com/MitulShah1/expense-tracker-bot/issues) in the GitHub repository
- 📚 Check the documentation in the `docs/` directory
- 🧪 Review the code examples in the test files
- 💬 Join our community discussions

## 🙏 Acknowledgments

- Built with [Go](https://golang.org/) and [PostgreSQL](https://www.postgresql.org/)
- Uses [Telegram Bot API](https://core.telegram.org/bots/api) for messaging
- Inspired by the need for better personal expense tracking

---

## Made with ❤️ using Go and PostgreSQL

[⬆️ Back to top](#-expense-tracker-bot)
