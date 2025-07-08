package services

import (
	"context"
	"testing"

	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVectorService is a mock implementation of VectorService
type MockVectorService struct {
	mock.Mock
}

func (m *MockVectorService) SearchExpensesByQuery(ctx context.Context, telegramID int64, query string, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	args := m.Called(ctx, telegramID, query, similarityThreshold, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockVectorService) FindSimilarExpenses(ctx context.Context, expenseID int64, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	args := m.Called(ctx, expenseID, similarityThreshold, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockVectorService) UpdateExpenseEmbeddings(ctx context.Context, expenseID int64) error {
	args := m.Called(ctx, expenseID)
	return args.Error(0)
}

func (m *MockVectorService) BatchUpdateEmbeddings(ctx context.Context, telegramID int64) error {
	args := m.Called(ctx, telegramID)
	return args.Error(0)
}

func TestVectorService_SearchExpensesByQuery(t *testing.T) {
	tests := []struct {
		name                string
		telegramID          int64
		query               string
		similarityThreshold float32
		limit               int
		setupMock           func(*MockStorage)
		expectError         bool
		errorType           errors.ErrorType
		expectedLen         int
	}{
		{
			name:                "successful semantic search",
			telegramID:          12345,
			query:               "petrol",
			similarityThreshold: 0.1,
			limit:               10,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				expenses := []*models.Expense{
					{ID: 1, Notes: "Petrol", CategoryName: "⛽ Petrol", TotalPrice: 100.0},
					{ID: 2, Notes: "Fuel", CategoryName: "⛽ Petrol", TotalPrice: 200.0},
				}

				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("SearchExpensesBySimilarity", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), float32(0.1), 10).Return(expenses, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:                "user not found",
			telegramID:          12345,
			query:               "petrol",
			similarityThreshold: 0.1,
			limit:               10,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(nil, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
		{
			name:                "database error during search",
			telegramID:          12345,
			query:               "petrol",
			similarityThreshold: 0.1,
			limit:               10,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("SearchExpensesBySimilarity", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), float32(0.1), 10).Return(nil, assert.AnError)
			},
			expectError: true,
			errorType:   errors.ErrorTypeDatabase,
		},
		{
			name:                "empty query",
			telegramID:          12345,
			query:               "",
			similarityThreshold: 0.1,
			limit:               10,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()
			service := NewVectorService(mockDB, logger)

			// Execute
			expenses, err := service.SearchExpensesByQuery(context.Background(), tt.telegramID, tt.query, tt.similarityThreshold, tt.limit)

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

func TestVectorService_FindSimilarExpenses(t *testing.T) {
	tests := []struct {
		name                string
		expenseID           int64
		similarityThreshold float32
		limit               int
		setupMock           func(*MockStorage)
		expectError         bool
		errorType           errors.ErrorType
		expectedLen         int
	}{
		{
			name:                "successful similar expenses search",
			expenseID:           1,
			similarityThreshold: 0.8,
			limit:               5,
			setupMock: func(mockDB *MockStorage) {
				expenses := []*models.Expense{
					{ID: 2, Notes: "Similar expense", CategoryName: "⛽ Petrol", TotalPrice: 150.0},
					{ID: 3, Notes: "Another similar", CategoryName: "⛽ Petrol", TotalPrice: 180.0},
				}
				mockDB.On("FindSimilarExpenses", mock.Anything, int64(1), float32(0.8), 5).Return(expenses, nil)
			},
			expectError: false,
			expectedLen: 2,
		},
		{
			name:                "database error during search",
			expenseID:           1,
			similarityThreshold: 0.8,
			limit:               5,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("FindSimilarExpenses", mock.Anything, int64(1), float32(0.8), 5).Return(nil, assert.AnError)
			},
			expectError: true,
			errorType:   errors.ErrorTypeDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()
			service := NewVectorService(mockDB, logger)

			// Execute
			expenses, err := service.FindSimilarExpenses(context.Background(), tt.expenseID, tt.similarityThreshold, tt.limit)

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

func TestVectorService_UpdateExpenseEmbeddings(t *testing.T) {
	tests := []struct {
		name        string
		expenseID   int64
		setupMock   func(*MockStorage)
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:      "successful embedding update",
			expenseID: 1,
			setupMock: func(mockDB *MockStorage) {
				expense := &models.Expense{
					ID:            1,
					Notes:         "Petrol expense",
					CategoryName:  "⛽ Petrol",
					CategoryGroup: "Vehicle",
				}
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(expense, nil)
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(nil)
			},
			expectError: false,
		},
		{
			name:      "expense not found",
			expenseID: 999,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("GetExpenseByID", mock.Anything, int64(999)).Return(nil, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
		{
			name:      "database error during update",
			expenseID: 1,
			setupMock: func(mockDB *MockStorage) {
				expense := &models.Expense{
					ID:            1,
					Notes:         "Petrol expense",
					CategoryName:  "⛽ Petrol",
					CategoryGroup: "Vehicle",
				}
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(expense, nil)
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(assert.AnError)
			},
			expectError: true,
			errorType:   errors.ErrorTypeDatabase,
		},
		{
			name:      "expense with empty notes",
			expenseID: 1,
			setupMock: func(mockDB *MockStorage) {
				expense := &models.Expense{
					ID:            1,
					Notes:         "",
					CategoryName:  "⛽ Petrol",
					CategoryGroup: "Vehicle",
				}
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(expense, nil)
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()
			service := NewVectorService(mockDB, logger)

			// Execute
			err := service.UpdateExpenseEmbeddings(context.Background(), tt.expenseID)

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

func TestVectorService_BatchUpdateEmbeddings(t *testing.T) {
	tests := []struct {
		name        string
		telegramID  int64
		setupMock   func(*MockStorage)
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:       "successful batch update",
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				expenses := []*models.Expense{
					{ID: 1, Notes: "Petrol", CategoryName: "⛽ Petrol", CategoryGroup: "Vehicle"},
					{ID: 2, Notes: "Service", CategoryName: "🔧 Service", CategoryGroup: "Vehicle"},
				}

				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetExpensesByUserID", mock.Anything, int64(1)).Return(expenses, nil)
				mockDB.On("GetExpenseEmbedding", mock.Anything, int64(1)).Return(nil, &errors.AppError{Type: errors.ErrorTypeNotFound})
				mockDB.On("GetExpenseEmbedding", mock.Anything, int64(2)).Return(nil, &errors.AppError{Type: errors.ErrorTypeNotFound})
				mockDB.On("GetExpenseByID", mock.Anything, int64(1)).Return(expenses[0], nil)
				mockDB.On("GetExpenseByID", mock.Anything, int64(2)).Return(expenses[1], nil)
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, int64(1), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(nil)
				mockDB.On("UpdateExpenseEmbedding", mock.Anything, int64(2), mock.AnythingOfType("[]float32"), mock.AnythingOfType("[]float32")).Return(nil)
			},
			expectError: false,
		},
		{
			name:       "user not found",
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(nil, nil)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
		{
			name:       "database error getting expenses",
			telegramID: 12345,
			setupMock: func(mockDB *MockStorage) {
				user := &models.User{ID: 1, TelegramID: 12345}
				mockDB.On("GetUserByTelegramID", mock.Anything, int64(12345)).Return(user, nil)
				mockDB.On("GetExpensesByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)
			},
			expectError: true,
			errorType:   errors.ErrorTypeDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			logger := logger.NewMockLogger()
			service := NewVectorService(mockDB, logger)

			// Execute
			err := service.BatchUpdateEmbeddings(context.Background(), tt.telegramID)

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

func TestVectorService_GenerateEmbedding(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		expectError bool
	}{
		{
			name:        "successful embedding generation",
			text:        "petrol",
			expectError: false,
		},
		{
			name:        "empty text",
			text:        "",
			expectError: true,
		},
		{
			name:        "whitespace only text",
			text:        "   ",
			expectError: true,
		},
		{
			name:        "long text",
			text:        "This is a very long text that should still generate an embedding",
			expectError: false,
		},
		{
			name:        "special characters",
			text:        "Petrol & Gas ⛽",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockDB := &MockStorage{}
			logger := logger.NewMockLogger()
			service := NewVectorService(mockDB, logger)

			// Execute
			embedding, err := service.generateEmbedding(context.Background(), tt.text)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, embedding)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, embedding)
				assert.Len(t, embedding, 1536) // Should be 1536-dimensional
				// Check that at least some values are non-zero
				hasNonZero := false
				for _, val := range embedding {
					if val != 0 {
						hasNonZero = true
						break
					}
				}
				assert.True(t, hasNonZero, "Embedding should have some non-zero values")
			}
		})
	}
}

func TestVectorService_GenerateEmbedding_Consistency(t *testing.T) {
	// Test that the same text generates the same embedding
	mockDB := &MockStorage{}
	logger := logger.NewMockLogger()
	service := NewVectorService(mockDB, logger)

	text := "petrol"
	embedding1, err1 := service.generateEmbedding(context.Background(), text)
	embedding2, err2 := service.generateEmbedding(context.Background(), text)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, embedding1, embedding2, "Same text should generate same embedding")
}

func TestVectorService_GenerateEmbedding_DifferentTexts(t *testing.T) {
	// Test that different texts generate different embeddings
	mockDB := &MockStorage{}
	logger := logger.NewMockLogger()
	service := NewVectorService(mockDB, logger)

	text1 := "petrol"
	text2 := "diesel"
	embedding1, err1 := service.generateEmbedding(context.Background(), text1)
	embedding2, err2 := service.generateEmbedding(context.Background(), text2)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, embedding1, embedding2, "Different texts should generate different embeddings")
}
