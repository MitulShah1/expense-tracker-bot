// Package database provides PostgreSQL database operations for the expense tracker bot.
package database

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/jmoiron/sqlx"
)

// MockStorage implements Storage interface for testing
type MockStorage struct {
	mu         sync.RWMutex
	users      map[int64]*models.User
	categories []*models.Category
	expenses   map[int64]*models.Expense
	nextID     int64
}

// NewMockStorage creates a new mock storage instance
func NewMockStorage() Storage {
	return &MockStorage{
		users:      make(map[int64]*models.User),
		categories: make([]*models.Category, 0),
		expenses:   make(map[int64]*models.Expense),
		nextID:     1,
	}
}

// Close closes the mock storage (no-op)
func (m *MockStorage) Close() error {
	return nil
}

// GetDB returns nil for mock storage
func (m *MockStorage) GetDB() *sqlx.DB {
	return nil
}

// User Operations

// CreateUser creates a new user in mock storage
func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if user already exists
	if existingUser, exists := m.users[user.TelegramID]; exists {
		// Update existing user
		existingUser.Username = user.Username
		existingUser.FirstName = user.FirstName
		existingUser.LastName = user.LastName
		existingUser.UpdatedAt = time.Now()
		*user = *existingUser
		return nil
	}

	// Create new user
	user.ID = m.nextID
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users[user.TelegramID] = user
	m.nextID++
	return nil
}

// GetUserByTelegramID retrieves a user by Telegram ID from mock storage
func (m *MockStorage) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if user, exists := m.users[telegramID]; exists {
		return user, nil
	}
	return nil, sql.ErrNoRows
}

// Category Operations

// GetAllCategories retrieves all categories from mock storage
func (m *MockStorage) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*models.Category, len(m.categories))
	copy(result, m.categories)
	return result, nil
}

// GetCategoriesByGroup retrieves categories by group from mock storage
func (m *MockStorage) GetCategoriesByGroup(ctx context.Context, group string) ([]*models.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*models.Category
	for _, category := range m.categories {
		if category.Group == group {
			result = append(result, category)
		}
	}
	return result, nil
}

// GetCategoryByName retrieves a category by name from mock storage
func (m *MockStorage) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, category := range m.categories {
		if category.Name == name {
			return category, nil
		}
	}
	return nil, sql.ErrNoRows
}

// Expense Operations

// CreateExpense creates a new expense in mock storage
func (m *MockStorage) CreateExpense(ctx context.Context, expense *models.Expense) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	expense.ID = m.nextID
	expense.CreatedAt = time.Now()
	expense.UpdatedAt = time.Now()
	m.expenses[expense.ID] = expense
	m.nextID++
	return nil
}

// GetExpensesByUserID retrieves all expenses for a user from mock storage
func (m *MockStorage) GetExpensesByUserID(ctx context.Context, userID int64) ([]*models.Expense, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*models.Expense
	for _, expense := range m.expenses {
		if expense.UserID == userID && expense.DeletedAt == nil {
			result = append(result, expense)
		}
	}
	return result, nil
}

// GetExpensesByTelegramID retrieves all expenses for a user by their Telegram ID from mock storage
func (m *MockStorage) GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*models.Expense, error) {
	// First get the user by Telegram ID
	user, err := m.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, sql.ErrNoRows
	}

	// Then get expenses by internal user ID
	return m.GetExpensesByUserID(ctx, user.ID)
}

// GetExpenseByID retrieves an expense by ID from mock storage
func (m *MockStorage) GetExpenseByID(ctx context.Context, id int64) (*models.Expense, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if expense, exists := m.expenses[id]; exists && expense.DeletedAt == nil {
		return expense, nil
	}
	return nil, sql.ErrNoRows
}

// UpdateExpense updates an existing expense in mock storage
func (m *MockStorage) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existingExpense, exists := m.expenses[expense.ID]; exists && existingExpense.DeletedAt == nil {
		existingExpense.CategoryID = expense.CategoryID
		existingExpense.VehicleType = expense.VehicleType
		existingExpense.Odometer = expense.Odometer
		existingExpense.PetrolPrice = expense.PetrolPrice
		existingExpense.TotalPrice = expense.TotalPrice
		existingExpense.Notes = expense.Notes
		existingExpense.Timestamp = expense.Timestamp
		existingExpense.UpdatedAt = time.Now()
		*expense = *existingExpense
		return nil
	}
	return sql.ErrNoRows
}

// DeleteExpense soft deletes an expense in mock storage
func (m *MockStorage) DeleteExpense(ctx context.Context, id, userID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if expense, exists := m.expenses[id]; exists && expense.UserID == userID && expense.DeletedAt == nil {
		now := time.Now()
		expense.DeletedAt = &now
		expense.UpdatedAt = now
		return nil
	}
	return sql.ErrNoRows
}

// GetExpenseStats retrieves expense statistics for a user from mock storage
func (m *MockStorage) GetExpenseStats(ctx context.Context, userID int64) (*models.ExpenseStats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var totalExpenses int64
	var totalSpent float64
	var minExpense, maxExpense float64
	var firstExpenseDate, lastExpenseDate time.Time
	var expenseCount int

	for _, expense := range m.expenses {
		if expense.UserID == userID && expense.DeletedAt == nil {
			totalExpenses++
			totalSpent += expense.TotalPrice

			if expenseCount == 0 {
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
			expenseCount++
		}
	}

	avgExpense := float64(0)
	if totalExpenses > 0 {
		avgExpense = totalSpent / float64(totalExpenses)
	}

	return &models.ExpenseStats{
		TotalExpenses:    totalExpenses,
		TotalSpent:       totalSpent,
		AvgExpense:       avgExpense,
		MinExpense:       minExpense,
		MaxExpense:       maxExpense,
		FirstExpenseDate: firstExpenseDate,
		LastExpenseDate:  lastExpenseDate,
	}, nil
}

// GetExpensesByDateRange retrieves expenses within a date range from mock storage
func (m *MockStorage) GetExpensesByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.Expense, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*models.Expense
	for _, expense := range m.expenses {
		if expense.UserID == userID && expense.DeletedAt == nil &&
			expense.Timestamp.After(startDate) && expense.Timestamp.Before(endDate) {
			result = append(result, expense)
		}
	}
	return result, nil
}

// Helper methods for testing

// AddMockCategory adds a category to mock storage for testing
func (m *MockStorage) AddMockCategory(category *models.Category) {
	m.mu.Lock()
	defer m.mu.Unlock()

	category.ID = m.nextID
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	m.categories = append(m.categories, category)
	m.nextID++
}

// ClearMockData clears all mock data for testing
func (m *MockStorage) ClearMockData() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.users = make(map[int64]*models.User)
	m.categories = make([]*models.Category, 0)
	m.expenses = make(map[int64]*models.Expense)
	m.nextID = 1
}
