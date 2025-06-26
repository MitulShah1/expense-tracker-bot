# Database Migrations

This directory contains PostgreSQL migration scripts for the Expense Tracker Bot.

## Migration Files

### 001_init.sql

- Creates the initial database structure
- Sets up tables: `users`, `categories`, `expenses`, `recurring_expenses`
- Creates indexes for performance
- Sets up triggers for automatic `updated_at` timestamps

### 002_seed_categories.sql

- Seeds the `categories` table with default expense categories
- Includes all category groups: Vehicle, Home, Daily Living, Entertainment, Health, Education, Travel, Investments, Gifts, Other

### 003_add_budgets_table.sql

- Adds budget tracking functionality
- Creates `budgets` and `budget_limits` tables
- Enables setting spending limits by category and time period

### 004_add_views.sql

- Creates useful database views for reporting and analytics
- Includes views for expense summaries, monthly breakdowns, and user statistics

## Running Migrations

### Option 1: Manual Execution

```bash
# Connect to your PostgreSQL database
psql -h localhost -U username -d expense_tracker

# Run migrations in order
\i migrations/001_init.sql
\i migrations/002_seed_categories.sql
\i migrations/003_add_budgets_table.sql
\i migrations/004_add_views.sql
```

### Option 2: Using a Migration Tool

We recommend using [golang-migrate](https://github.com/golang-migrate/migrate) for automated migrations:

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgres://username:password@localhost:5432/expense_tracker?sslmode=disable" up
```

### Option 3: Using Docker

```bash
# Run PostgreSQL in Docker
docker run --name postgres-expense-tracker \
  -e POSTGRES_DB=expense_tracker \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  -d postgres:15

# Run migrations
migrate -path migrations -database "postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable" up
```

## Database Schema

### Core Tables

#### users

- `id`: Primary key
- `telegram_id`: Unique Telegram user ID
- `username`, `first_name`, `last_name`: User information
- `created_at`, `updated_at`: Timestamps

#### categories

- `id`: Primary key
- `name`: Category name (unique)
- `emoji`: Category emoji
- `group`: Category group (Vehicle, Home, etc.)

#### expenses

- `id`: Primary key
- `user_id`: Foreign key to users
- `category_id`: Foreign key to categories
- `vehicle_type`: CAR/BIKE (optional)
- `odometer`, `petrol_price`: Optional vehicle data
- `total_price`: Required expense amount
- `notes`: Optional notes
- `timestamp`: When the expense occurred
- `deleted_at`: Soft delete timestamp

### Optional Tables

#### recurring_expenses

- For setting up recurring payments
- Supports daily, weekly, monthly, yearly intervals

#### budgets

- For budget tracking and limits
- Supports different time periods

## Views

The migration creates several useful views:

- `expense_summary`: Summary by user and category
- `monthly_expense_summary`: Monthly breakdowns
- `category_expense_breakdown`: Category-wise statistics
- `user_expense_stats`: User statistics
- `recent_expenses`: Recent expenses (last 30 days)

## Rollback

To rollback migrations using golang-migrate:

```bash
# Rollback one step
migrate -path migrations -database "postgres://username:password@localhost:5432/expense_tracker?sslmode=disable" down 1

# Rollback all
migrate -path migrations -database "postgres://username:password@localhost:5432/expense_tracker?sslmode=disable" down
```

## Environment Variables

Make sure to set these environment variables:

```env
DATABASE_URL=postgres://username:password@localhost:5432/expense_tracker?sslmode=disable
```

## Notes

- All tables use `TIMESTAMPTZ` for timezone-aware timestamps
- Soft deletes are implemented using `deleted_at` columns
- Automatic `updated_at` triggers are set up
- Indexes are created for optimal query performance
- Foreign key constraints ensure data integrity 