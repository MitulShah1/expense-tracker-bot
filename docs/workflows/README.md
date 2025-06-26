# Telegram Bot Workflows

This directory contains detailed workflow diagrams for the Expense Tracker Telegram Bot. Each diagram illustrates the step-by-step process for different bot commands and features.

## Available Workflows

1. [Add Expense Flow](./add-expense.md)
2. [List Expenses Flow](./list-expenses.md)
3. [Edit Expense Flow](./edit-expense.md)
4. [Delete Expense Flow](./delete-expense.md)
5. [Report Flow](./report.md)
6. [Dashboard Flow](./dashboard.md)
7. [Error Handling Flow](./error-handling.md)

## How to Read the Diagrams

Each workflow is documented using Mermaid diagrams, which can be viewed directly on GitHub. The diagrams show:

- Step-by-step processes
- Decision points
- Error handling
- Validation steps
- Success/failure paths

## Category Structure

The bot uses a hierarchical category structure:

### Vehicle Group

- Fuel
- Oil
- Battery
- Tires

### Service Group

- Maintenance
- Repair
- Cleaning
- Washing

### Other Group

- Miscellaneous
- Insurance
- Tax
- Other

## Command Overview

- `/start` - Initialize bot and show welcome message
- `/add` - Add new expense
- `/list` - List all expenses
- `/edit` - Edit existing expense
- `/delete` - Delete expense
- `/report` - Generate expense report
- `/dashboard` - Show expense dashboard
- `/help` - Show help message
- `/cancel` - Cancel current operation