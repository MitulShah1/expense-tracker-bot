package bot

import (
	"context"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/MitulShah1/expense-tracker-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

// MockVectorService is a mock implementation of VectorServiceInterface
type MockVectorService struct {
	mock.Mock
}

func (m *MockVectorService) SearchExpensesByQuery(ctx context.Context, userID int64, query string, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	args := m.Called(ctx, userID, query, similarityThreshold, limit)
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

// MockExpenseService is a mock implementation that embeds the real service
type MockExpenseService struct {
	*services.ExpenseService
	mock.Mock
}

func NewMockExpenseService(db database.Storage, logger logger.Logger) *services.ExpenseService {
	return services.NewExpenseService(db, logger)
}

func (m *MockExpenseService) GetExpensesByTelegramID(ctx context.Context, telegramID int64, limit, offset int) ([]*models.Expense, error) {
	args := m.Called(ctx, telegramID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Expense), args.Error(1)
}

func (m *MockExpenseService) GetExpenseByID(ctx context.Context, expenseID int64) (*models.Expense, error) {
	args := m.Called(ctx, expenseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Expense), args.Error(1)
}

// MockBotAPI mocks the tgbotapi.BotAPI for testing
// Only implements the Send method needed for handler tests
// All other methods can panic if called
type MockBotAPI struct {
	mock.Mock
}

func (m *MockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.Called(c)
	return tgbotapi.Message{}, nil
}

func (m *MockBotAPI) GetUpdatesChan(u tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel {
	ch := make(chan tgbotapi.Update)
	close(ch)
	return ch
}

func TestBot_handleSearchCommand(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(*MockStorage, *MockVectorService)
		expectError bool
	}{
		{
			name:   "successful search command",
			userID: 12345,
			setupMock: func(mockDB *MockStorage, mockVector *MockVectorService) {
				expenses := []*models.Expense{
					{ID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol", Notes: "fuel", Timestamp: time.Now()},
				}
				mockDB.On("GetExpensesByTelegramID", mock.Anything, int64(12345)).Return(expenses, nil)
			},
			expectError: false,
		},
		{
			name:   "no expenses found",
			userID: 12345,
			setupMock: func(mockDB *MockStorage, mockVector *MockVectorService) {
				mockDB.On("GetExpensesByTelegramID", mock.Anything, int64(12345)).Return([]*models.Expense{}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockStorage{}
			mockVector := &MockVectorService{}
			tt.setupMock(mockDB, mockVector)

			// Create a mock logger
			mockLogger := &logger.MockLogger{}

			mockAPI := &MockBotAPI{}
			mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)
			bot := &Bot{
				db:             mockDB,
				logger:         mockLogger,
				states:         make(map[int64]*models.UserState),
				expenseService: NewMockExpenseService(mockDB, mockLogger),
				vectorService:  mockVector,
				api:            mockAPI,
			}

			// Create a mock message
			message := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: tt.userID},
				From: &tgbotapi.User{ID: tt.userID},
			}

			// Test the method
			err := bot.handleSearchCommand(context.Background(), message)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockVector.AssertExpectations(t)
		})
	}
}

func TestBot_handleSearchQuery(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		query       string
		setupMock   func(*MockStorage, *MockVectorService)
		expectError bool
	}{
		{
			name:   "successful search query",
			userID: 12345,
			query:  "petrol",
			setupMock: func(mockDB *MockStorage, mockVector *MockVectorService) {
				mockVector.On("SearchExpensesByQuery", mock.Anything, int64(12345), "petrol", float32(0.1), 10).Return([]*models.Expense{{ID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol", Notes: "fuel", Timestamp: time.Now()}}, nil)
			},
			expectError: false,
		},
		{
			name:   "empty query",
			userID: 12345,
			query:  "",
			setupMock: func(mockDB *MockStorage, mockVector *MockVectorService) {
				// No mock setup needed for empty query
			},
			expectError: false,
		},
		{
			name:   "no results found",
			userID: 12345,
			query:  "nonexistent",
			setupMock: func(mockDB *MockStorage, mockVector *MockVectorService) {
				mockVector.On("SearchExpensesByQuery", mock.Anything, int64(12345), "nonexistent", float32(0.1), 10).Return([]*models.Expense{}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockStorage{}
			mockVector := &MockVectorService{}
			mockLogger := &logger.MockLogger{}
			mockAPI := &MockBotAPI{}
			mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)

			// Set up the vectorService mock expectation for this test
			if tt.query != "" {
				mockVector.On("SearchExpensesByQuery", mock.Anything, tt.userID, tt.query, float32(0.1), 10).Return([]*models.Expense{{ID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol", Notes: "fuel", Timestamp: time.Now()}}, nil)
			}

			bot := &Bot{
				db:             mockDB,
				logger:         mockLogger,
				states:         make(map[int64]*models.UserState),
				expenseService: NewMockExpenseService(mockDB, mockLogger),
				vectorService:  mockVector,
				api:            mockAPI,
			}

			message := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: tt.userID},
				From: &tgbotapi.User{ID: tt.userID},
				Text: tt.query,
			}

			err := bot.handleSearchQuery(context.Background(), message)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
			mockVector.AssertExpectations(t)
			mockAPI.AssertExpectations(t)
		})
	}
}

func TestBot_handleListCommand(t *testing.T) {
	tests := []struct {
		name        string
		userID      int64
		setupMock   func(*MockStorage)
		expectError bool
	}{
		{
			name:   "successful list command",
			userID: 12345,
			setupMock: func(mockDB *MockStorage) {
				expenses := []*models.Expense{
					{ID: 1, TotalPrice: 100.0, CategoryName: "⛽ Petrol", Notes: "fuel", Timestamp: time.Now()},
					{ID: 2, TotalPrice: 50.0, CategoryName: "🍕 Food", Notes: "lunch", Timestamp: time.Now()},
				}
				mockDB.On("GetExpensesByTelegramID", mock.Anything, int64(12345)).Return(expenses, nil)
			},
			expectError: false,
		},
		{
			name:   "no expenses found",
			userID: 12345,
			setupMock: func(mockDB *MockStorage) {
				mockDB.On("GetExpensesByTelegramID", mock.Anything, int64(12345)).Return([]*models.Expense{}, nil)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := &MockStorage{}
			tt.setupMock(mockDB)

			// Create a mock logger
			mockLogger := &logger.MockLogger{}

			mockAPI := &MockBotAPI{}
			mockAPI.On("Send", mock.Anything).Return(tgbotapi.Message{}, nil)
			bot := &Bot{
				db:             mockDB,
				logger:         mockLogger,
				expenseService: NewMockExpenseService(mockDB, mockLogger),
				api:            mockAPI,
			}

			// Create a mock message
			message := &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: tt.userID},
				From: &tgbotapi.User{ID: tt.userID},
			}

			// Test the method
			err := bot.handleListCommand(context.Background(), message)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockDB.AssertExpectations(t)
		})
	}
}
