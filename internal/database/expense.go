package database

import (
	"context"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
)

// ExpenseStorage defines operations for expense management
type ExpenseStorage interface {
	CreateExpense(ctx context.Context, expense *models.Expense) error
	GetExpensesByUserID(ctx context.Context, userID int64) ([]*models.Expense, error)
	GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*models.Expense, error)
	GetExpenseByID(ctx context.Context, id int64) (*models.Expense, error)
	UpdateExpense(ctx context.Context, expense *models.Expense) error
	DeleteExpense(ctx context.Context, id, userID int64) error
	GetExpenseStats(ctx context.Context, userID int64) (*models.ExpenseStats, error)
	GetExpensesByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.Expense, error)
}

// CreateExpense creates a new expense
func (c *Client) CreateExpense(ctx context.Context, expense *models.Expense) error {
	query := `
		INSERT INTO expenses (user_id, category_id, vehicle_type, odometer, petrol_price, total_price, notes, timestamp)
		VALUES ($1, $2, CASE WHEN $3 = '' THEN NULL ELSE $3 END, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	return c.db.QueryRowxContext(ctx, query,
		expense.UserID, expense.CategoryID, expense.VehicleType, expense.Odometer,
		expense.PetrolPrice, expense.TotalPrice, expense.Notes, expense.Timestamp).
		StructScan(expense)
}

// GetExpensesByUserID retrieves all expenses for a user
func (c *Client) GetExpensesByUserID(ctx context.Context, userID int64) ([]*models.Expense, error) {
	var expenses []*models.Expense
	query := `
		SELECT e.*, c.name as category_name, c.emoji as category_emoji, c."group" as category_group
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = $1 AND e.deleted_at IS NULL
		ORDER BY e.timestamp DESC`

	err := c.db.SelectContext(ctx, &expenses, query, userID)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetExpenseByID retrieves an expense by ID
func (c *Client) GetExpenseByID(ctx context.Context, id int64) (*models.Expense, error) {
	var expense models.Expense
	query := `
		SELECT e.*, c.name as category_name, c.emoji as category_emoji, c."group" as category_group
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.id = $1 AND e.deleted_at IS NULL`

	err := c.db.GetContext(ctx, &expense, query, id)
	if err != nil {
		if isNoRows(err) {
			return nil, errNotFound
		}
		return nil, err
	}

	return &expense, nil
}

// UpdateExpense updates an existing expense
func (c *Client) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	query := `
		UPDATE expenses 
		SET category_id = $1, vehicle_type = CASE WHEN $2 = '' THEN NULL ELSE $2 END, odometer = $3, petrol_price = $4, 
		    total_price = $5, notes = $6, timestamp = $7, updated_at = now()
		WHERE id = $8 AND user_id = $9 AND deleted_at IS NULL
		RETURNING updated_at`

	return c.db.QueryRowxContext(ctx, query,
		expense.CategoryID, expense.VehicleType, expense.Odometer, expense.PetrolPrice,
		expense.TotalPrice, expense.Notes, expense.Timestamp, expense.ID, expense.UserID).
		Scan(&expense.UpdatedAt)
}

// DeleteExpense soft deletes an expense
func (c *Client) DeleteExpense(ctx context.Context, id, userID int64) error {
	query := `
		UPDATE expenses 
		SET deleted_at = now(), updated_at = now()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	result, err := c.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}

// GetExpenseStats retrieves expense statistics for a user
func (c *Client) GetExpenseStats(ctx context.Context, userID int64) (*models.ExpenseStats, error) {
	var stats models.ExpenseStats
	query := `
		SELECT 
			COUNT(*) as total_expenses,
			COALESCE(SUM(total_price), 0) as total_spent,
			COALESCE(AVG(total_price), 0) as avg_expense,
			COALESCE(MIN(total_price), 0) as min_expense,
			COALESCE(MAX(total_price), 0) as max_expense,
			COALESCE(MIN(timestamp), '1970-01-01'::timestamptz) as first_expense_date,
			COALESCE(MAX(timestamp), '1970-01-01'::timestamptz) as last_expense_date
		FROM expenses 
		WHERE user_id = $1 AND deleted_at IS NULL`

	err := c.db.GetContext(ctx, &stats, query, userID)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetExpensesByDateRange retrieves expenses within a date range
func (c *Client) GetExpensesByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.Expense, error) {
	var expenses []*models.Expense
	query := `
		SELECT e.*, c.name as category_name, c.emoji as category_emoji, c."group" as category_group
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = $1 AND e.deleted_at IS NULL 
		  AND e.timestamp >= $2 AND e.timestamp <= $3
		ORDER BY e.timestamp DESC`

	err := c.db.SelectContext(ctx, &expenses, query, userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return expenses, nil
}

// GetExpensesByTelegramID retrieves all expenses for a user by their Telegram ID
func (c *Client) GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*models.Expense, error) {
	// First get the user by Telegram ID
	user, err := c.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errNotFound
	}

	// Then get expenses by internal user ID
	return c.GetExpensesByUserID(ctx, user.ID)
}
