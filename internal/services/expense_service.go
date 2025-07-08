// Package services provides business logic services for the expense tracker bot.
package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/MitulShah1/expense-tracker-bot/internal/validation"
)

// ExpenseService provides expense-related business logic
type ExpenseService struct {
	db        database.Storage
	logger    logger.Logger
	validator *validation.Validator
}

// NewExpenseService creates a new expense service
func NewExpenseService(db database.Storage, logger logger.Logger) *ExpenseService {
	return &ExpenseService{
		db:        db,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// CreateExpense creates a new expense
func (s *ExpenseService) CreateExpense(ctx context.Context, expense *models.Expense, telegramID int64) error {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return err
	}

	if err := s.validator.ValidateAmount(expense.TotalPrice, "total price"); err != nil {
		return err
	}

	if expense.Odometer > 0 {
		if err := s.validator.ValidateOdometer(expense.Odometer); err != nil {
			return err
		}
	}

	if expense.PetrolPrice > 0 {
		if err := s.validator.ValidateAmount(expense.PetrolPrice, "petrol price"); err != nil {
			return err
		}
	}

	if err := s.validator.ValidateNotes(expense.Notes); err != nil {
		return err
	}

	if err := s.validator.ValidateCategoryName(expense.CategoryName); err != nil {
		return err
	}

	// Get or create user
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		// This should not happen as user should be created before calling this service
		return errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Get category by name
	category, err := s.db.GetCategoryByName(ctx, expense.CategoryName)
	if err != nil {
		s.logger.Error(ctx, "Failed to get category by name", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get category", err)
	}

	if category == nil {
		return errors.NewNotFoundError("Category not found", fmt.Sprintf("Category '%s' not found", expense.CategoryName))
	}

	// Set the expense details
	var vehicleType sql.NullString
	if category.Group == "Vehicle" && expense.VehicleType.Valid {
		if err := s.validator.ValidateVehicleType(expense.VehicleType.String); err != nil {
			return err
		}
		vehicleType = expense.VehicleType
	} else {
		vehicleType = sql.NullString{Valid: false}
	}

	// Create expense record
	expenseRecord := &models.Expense{
		UserID:      user.ID,
		CategoryID:  category.ID,
		VehicleType: vehicleType,
		Odometer:    expense.Odometer,
		PetrolPrice: expense.PetrolPrice,
		TotalPrice:  expense.TotalPrice,
		Notes:       expense.Notes,
		Timestamp:   expense.Timestamp,
	}

	// Save expense to database
	if err := s.db.CreateExpense(ctx, expenseRecord); err != nil {
		s.logger.Error(ctx, "Failed to create expense", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to create expense", err)
	}

	s.logger.Info(ctx, "Expense created successfully",
		logger.Int("user_id", int(user.ID)),
		logger.Int("expense_id", int(expenseRecord.ID)),
		logger.Float64("total_price", expense.TotalPrice))

	// Generate embeddings for the new expense
	vectorService := NewVectorService(s.db, s.logger)
	if err := vectorService.UpdateExpenseEmbeddings(ctx, expenseRecord.ID); err != nil {
		s.logger.Error(ctx, "Failed to generate embeddings for new expense", logger.ErrorField(err))
		// Don't fail the expense creation if embedding generation fails
		// The expense is still created successfully
	}

	return nil
}

// GetExpensesByTelegramID retrieves expenses for a user by Telegram ID
func (s *ExpenseService) GetExpensesByTelegramID(ctx context.Context, telegramID int64, limit, offset int) ([]*models.Expense, error) {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return nil, err
	}

	if err := s.validator.ValidatePagination(limit, offset); err != nil {
		return nil, err
	}

	// Get expenses from database
	expenses, err := s.db.GetExpensesByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get expenses by Telegram ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get expenses", err)
	}

	// Apply pagination
	if offset >= len(expenses) {
		return []*models.Expense{}, nil
	}

	end := offset + limit
	if end > len(expenses) {
		end = len(expenses)
	}

	return expenses[offset:end], nil
}

// GetExpenseByID retrieves an expense by ID
func (s *ExpenseService) GetExpenseByID(ctx context.Context, expenseID int64) (*models.Expense, error) {
	// Validate input
	if err := s.validator.ValidateExpenseID(expenseID); err != nil {
		return nil, err
	}

	// Get expense from database
	expense, err := s.db.GetExpenseByID(ctx, expenseID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get expense by ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get expense", err)
	}

	if expense == nil {
		return nil, errors.NewNotFoundError("Expense not found", fmt.Sprintf("Expense with ID %d not found", expenseID))
	}

	return expense, nil
}

// UpdateExpense updates an existing expense
func (s *ExpenseService) UpdateExpense(ctx context.Context, expense *models.Expense, telegramID int64) error {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return err
	}

	if err := s.validator.ValidateExpenseID(expense.ID); err != nil {
		return err
	}

	if err := s.validator.ValidateAmount(expense.TotalPrice, "total price"); err != nil {
		return err
	}

	if expense.Odometer > 0 {
		if err := s.validator.ValidateOdometer(expense.Odometer); err != nil {
			return err
		}
	}

	if expense.PetrolPrice > 0 {
		if err := s.validator.ValidateAmount(expense.PetrolPrice, "petrol price"); err != nil {
			return err
		}
	}

	if err := s.validator.ValidateNotes(expense.Notes); err != nil {
		return err
	}

	// Get user by Telegram ID
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Get existing expense
	existingExpense, err := s.db.GetExpenseByID(ctx, expense.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get existing expense", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get existing expense", err)
	}

	if existingExpense == nil {
		return errors.NewNotFoundError("Expense not found", fmt.Sprintf("Expense with ID %d not found", expense.ID))
	}

	// Check ownership
	if existingExpense.UserID != user.ID {
		return errors.NewUnauthorizedError("You can only edit your own expenses")
	}

	// Update expense in database
	if err := s.db.UpdateExpense(ctx, expense); err != nil {
		s.logger.Error(ctx, "Failed to update expense", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to update expense", err)
	}

	s.logger.Info(ctx, "Expense updated successfully",
		logger.Int("user_id", int(user.ID)),
		logger.Int("expense_id", int(expense.ID)),
		logger.Float64("total_price", expense.TotalPrice))

	return nil
}

// DeleteExpense deletes an expense
func (s *ExpenseService) DeleteExpense(ctx context.Context, expenseID, telegramID int64) error {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return err
	}

	if err := s.validator.ValidateExpenseID(expenseID); err != nil {
		return err
	}

	// Get user by Telegram ID
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Get existing expense
	existingExpense, err := s.db.GetExpenseByID(ctx, expenseID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get existing expense", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get existing expense", err)
	}

	if existingExpense == nil {
		return errors.NewNotFoundError("Expense not found", fmt.Sprintf("Expense with ID %d not found", expenseID))
	}

	// Check ownership
	if existingExpense.UserID != user.ID {
		return errors.NewUnauthorizedError("You can only delete your own expenses")
	}

	// Delete expense from database
	if err := s.db.DeleteExpense(ctx, expenseID, user.ID); err != nil {
		s.logger.Error(ctx, "Failed to delete expense", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to delete expense", err)
	}

	s.logger.Info(ctx, "Expense deleted successfully",
		logger.Int("user_id", int(user.ID)),
		logger.Int("expense_id", int(expenseID)))

	return nil
}

// GetExpenseStats retrieves expense statistics for a user
func (s *ExpenseService) GetExpenseStats(ctx context.Context, telegramID int64, startDate, endDate *time.Time) (*models.ExpenseStats, error) {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return nil, err
	}

	// Validate date range if provided
	if startDate != nil && endDate != nil {
		if err := s.validator.ValidateDateRange(*startDate, *endDate); err != nil {
			return nil, err
		}
	}

	// Get user by Telegram ID
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Get expenses for statistics calculation
	expenses, err := s.db.GetExpensesByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get expenses for statistics", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get expenses for statistics", err)
	}

	// Filter by date range if provided
	if startDate != nil && endDate != nil {
		filteredExpenses := make([]*models.Expense, 0)
		for _, expense := range expenses {
			if expense.Timestamp.After(*startDate) && expense.Timestamp.Before(*endDate) {
				filteredExpenses = append(filteredExpenses, expense)
			}
		}
		expenses = filteredExpenses
	}

	// Calculate statistics
	stats := &models.ExpenseStats{
		TotalExpenses: int64(len(expenses)),
	}

	if len(expenses) == 0 {
		return stats, nil
	}

	var totalSpent float64
	var minExpense, maxExpense float64
	var firstExpenseDate, lastExpenseDate time.Time

	for i, expense := range expenses {
		totalSpent += expense.TotalPrice

		if i == 0 {
			minExpense = expense.TotalPrice
			maxExpense = expense.TotalPrice
			firstExpenseDate = expense.Timestamp
			lastExpenseDate = expense.Timestamp
		} else {
			if expense.TotalPrice < minExpense {
				minExpense = expense.TotalPrice
			}
			if expense.TotalPrice > maxExpense {
				maxExpense = expense.TotalPrice
			}
			if expense.Timestamp.Before(firstExpenseDate) {
				firstExpenseDate = expense.Timestamp
			}
			if expense.Timestamp.After(lastExpenseDate) {
				lastExpenseDate = expense.Timestamp
			}
		}
	}

	stats.TotalSpent = totalSpent
	stats.AvgExpense = totalSpent / float64(len(expenses))
	stats.MinExpense = minExpense
	stats.MaxExpense = maxExpense
	stats.FirstExpenseDate = firstExpenseDate
	stats.LastExpenseDate = lastExpenseDate

	return stats, nil
}
