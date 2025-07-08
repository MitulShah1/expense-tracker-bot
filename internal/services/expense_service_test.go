// Package services provides business logic services for the expense tracker bot.
package services

import (
	"context"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockStorage is a mock implementation of database.Storage
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStorage) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	args := m.Called(ctx, telegramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockStorage) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockStorage) GetCategoriesByGroup(ctx context.Context, group string) ([]*models.Category, error) {
	args := m.Called(ctx, group)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *MockStorage) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *MockStorage) CreateExpense(ctx context.Context, expense *models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockStorage) GetExpensesByTelegramID(ctx context.Context, telegramID int64) ([]*models.Expense, error) {
	args := m.Called(ctx, telegramID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockStorage) GetExpenseByID(ctx context.Context, expenseID int64) (*models.Expense, error) {
	args := m.Called(ctx, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}

func (m *MockStorage) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	args := m.Called(ctx, expense)
	return args.Error(0)
}

func (m *MockStorage) DeleteExpense(ctx context.Context, expenseID, userID int64) error {
	args := m.Called(ctx, expenseID, userID)
	return args.Error(0)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockStorage) GetDB() *sqlx.DB {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*sqlx.DB)
}

func (m *MockStorage) GetExpenseStats(ctx context.Context, userID int64) (*models.ExpenseStats, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExpenseStats), args.Error(1)
}

func (m *MockStorage) GetExpensesByDateRange(ctx context.Context, userID int64, startDate, endDate time.Time) ([]*models.Expense, error) {
	args := m.Called(ctx, userID, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockStorage) GetExpensesByUserID(ctx context.Context, userID int64) ([]*models.Expense, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

// VectorSearchStorage stubs
func (m *MockStorage) SearchExpensesBySimilarity(ctx context.Context, userID int64, queryEmbedding []float32, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	args := m.Called(ctx, userID, queryEmbedding, similarityThreshold, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockStorage) FindSimilarExpenses(ctx context.Context, expenseID int64, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	args := m.Called(ctx, expenseID, similarityThreshold, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockStorage) UpdateExpenseEmbedding(ctx context.Context, expenseID int64, notesEmbedding, categoryEmbedding []float32) error {
	args := m.Called(ctx, expenseID, notesEmbedding, categoryEmbedding)
	return args.Error(0)
}

func (m *MockStorage) GetExpenseEmbedding(ctx context.Context, expenseID int64) (*models.ExpenseEmbedding, error) {
	args := m.Called(ctx, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExpenseEmbedding), args.Error(1)
}

func TestExpenseService_CreateExpense(t *testing.T) {
	tests := []struct {
		name        string
		expense     *models.Expense
		telegramID  int64
		setupMock   func(*MockStorage)
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name: "successful expense creation",
			expense: &models.Expense{
				TotalPrice:   100.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Test expense",
				Timestamp:    time.Now(),
			},
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				category := &models.Category{ID: 1, Name: "⛽ Petrol", Group: "Vehicle"}
				expenseRecord := &models.Expense{ID: 1, UserID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol", Notes: "Test expense"}

				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetCategoryByName", mock.Anything, "⛽ Petrol").Return(category, nil)
				mockDB.On("CreateExpense", mock.Anything, mock.AnythingOfType("*models.Expense")).Return(nil)
				// General mock for all GetExpenseByID calls
				mockDB.On("GetExpenseByID", mock.Anything, mock.AnythingOfType("int64")).Return(expenseRecord, nil)
				// Mock UpdateExpenseEmbedding
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, mock.AnythingOfType("int64"), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "invalid telegram ID",
			expense: &models.Expense{
				TotalPrice:   100.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Test expense",
				Timestamp:    time.Now(),
			},
			telegramID:  -1,
			setupMock:   func(mockDB *MockStorage) {},
			expectError: true,
			errorType:   errors.ErrorTypeValidation,
		},
		{
			name: "invalid total price",
			expense: &models.Expense{
				TotalPrice:   -100.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Test expense",
				Timestamp:    time.Now(),
			},
			telegramID:  12345,
			setupMock:   func(mockDB *MockStorage) {},
			expectError: true,
			errorType:   errors.ErrorTypeValidation,
		},
		{
			name: "user not found",
			expense: &models.Expense{
				TotalPrice:   100.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Test expense",
				Timestamp:    time.Now(),
			},
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(nil, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
		{
			name: "category not found",
			expense: &models.Expense{
				TotalPrice:   100.0,
				CategoryName: "Invalid Category",
				Notes:        "Test expense",
				Timestamp:    time.Now(),
			},
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetCategoryByName", mock.Anything, "Invalid Category").Return(nil, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()

			service := NewExpenseService(mockDB, logger)

			// Execute
			err := service.CreateExpense(context.Background(), tt.expense, tt.telegramID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if appErr, ok := err.(*errors.AppError); ok {
					assert.Equal(t, tt.errorType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestExpenseService_GetExpensesByTelegramID(t *testing.T) {
	tests := []struct {
		name        string
		telegramID  int64
		limit       int
		offset      int
		setupMock   func(*MockStorage)
		expectError bool
		errorType   errors.ErrorType
		expectedLen int
	}{
		{
			name:       "successful retrieval",
			telegramID: 12345,
			limit:      10,
			offset:     0,
			setupMock: func(mockDB *MockStorage) {
				expenses := []*models.Expense{
					{ID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol"},
					{ID: 2, TotalPrice: 200.0, CategoryName: "🍽️ Dining"},
				}
				mockDB.On("GetExpensesByTelegramID", mock.Anything, int64(12345)).Return(expenses, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:        "invalid telegram ID",
			telegramID:  -1,
			limit:       10,
			offset:      0,
			setupMock:   func(mockDB *MockStorage) {},
			expectError: true,
			errorType:   errors.ErrorTypeValidation,
		},
		{
			name:        "invalid pagination",
			telegramID:  12345,
			limit:       -1,
			offset:      0,
			setupMock:   func(mockDB *MockStorage) {},
			expectError: true,
			errorType:   errors.ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()

			service := NewExpenseService(mockDB, logger)

			// Execute
			expenses, err := service.GetExpensesByTelegramID(context.Background(), tt.telegramID, tt.limit, tt.offset)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if appErr, ok := err.(*errors.AppError); ok {
					assert.Equal(t, tt.errorType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, expenses, tt.expectedLen)
			}

			mockDB.AssertExpectations(t)
		})
	}
}

func TestExpenseService_UpdateExpense(t *testing.T) {
	tests := []struct {
		name        string
		expense     *models.Expense
		telegramID  int64
		setupMock   func(*MockStorage)
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name: "successful expense update",
			expense: &models.Expense{
				ID:           1,
				TotalPrice:   150.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Updated expense",
			},
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				existingExpense := &models.Expense{ID: 1, UserID: 1, TotalPrice: 100.0}

				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(existingExpense, nil)
				mockDB.On("UpdateExpense", mock.Anything, mock.AnythingOfType("*models.Expense")).Return(nil)
			},
			expectError: false,
		},
		{
			name: "unauthorized update",
			expense: &models.Expense{
				ID:           1,
				TotalPrice:   150.0,
				CategoryName: "⛽ Petrol",
				Notes:        "Updated expense",
			},
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				existingExpense := &models.Expense{ID: 1, UserID: 2, TotalPrice: 100.0} // Different user

				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(existingExpense, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()

			service := NewExpenseService(mockDB, logger)

			// Execute
			err := service.UpdateExpense(context.Background(), tt.expense, tt.telegramID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if appErr, ok := err.(*errors.AppError); ok {
					assert.Equal(t, tt.errorType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
