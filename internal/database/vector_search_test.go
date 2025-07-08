package database

import (
	"context"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestEmbedding creates a 1536-dimensional test embedding
func generateTestEmbedding() []float32 {
	embedding := make([]float32, 1536)
	for i := range embedding {
		embedding[i] = float32(i) / 1536.0
	}
	return embedding
}

func TestClient_SearchExpensesBySimilarity(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	tests := []struct {
		name                string
		userID              int64
		queryEmbedding      []float32
		similarityThreshold float32
		limit               int
		setupData           func(Storage) error
		expectError         bool
		errorType           errors.ErrorType
		expectedLen         int
	}{
		{
			name:                "successful similarity search",
			userID:              1,
			queryEmbedding:      generateTestEmbedding(),
			similarityThreshold: 0.1,
			limit:               10,
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12345, Username: "testuser"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expenses with embeddings
				expense1 := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Petrol expense",
					Timestamp:  time.Now(),
				}
				if err := client.CreateExpense(context.Background(), expense1); err != nil {
					return err
				}

				// Update with embedding
				embedding := generateTestEmbedding()
				return client.UpdateExpenseEmbedding(context.Background(), expense1.ID, embedding, embedding)
			},
			expectError: false,
			expectedLen: 1,
		},
		{
			name:                "no matching expenses",
			userID:              1,
			queryEmbedding:      generateTestEmbedding(),
			similarityThreshold: 0.9,
			limit:               10,
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12346, Username: "testuser2"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expense with different embedding
				expense1 := &models.Expense{
					UserID:     2,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Different expense",
					Timestamp:  time.Now(),
				}
				if err := client.CreateExpense(context.Background(), expense1); err != nil {
					return err
				}

				// Update with different embedding
				embedding := generateTestEmbedding()
				return client.UpdateExpenseEmbedding(context.Background(), expense1.ID, embedding, embedding)
			},
			expectError: false,
			expectedLen: 0,
		},
		{
			name:                "invalid embedding dimensions",
			userID:              1,
			queryEmbedding:      []float32{0.1, 0.2}, // Too short
			similarityThreshold: 0.1,
			limit:               10,
			setupData:           func(client Storage) error { return nil },
			expectError:         true,
			errorType:           errors.ErrorTypeDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database
			ctx := context.Background()
			logger := logger.NewMockLogger()
			client, err := NewClient(ctx, "postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable", logger)
			require.NoError(t, err)
			defer client.Close()

			// Setup test data
			if tt.setupData != nil {
				err = tt.setupData(client)
				require.NoError(t, err)
			}

			// Execute
			expenses, err := client.SearchExpensesBySimilarity(ctx, tt.userID, tt.queryEmbedding, tt.similarityThreshold, tt.limit)

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
		})
	}
}

func TestClient_FindSimilarExpenses(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	tests := []struct {
		name                string
		expenseID           int64
		similarityThreshold float32
		limit               int
		setupData           func(Storage) error
		expectError         bool
		errorType           errors.ErrorType
		expectedLen         int
	}{
		{
			name:                "successful similar expenses search",
			expenseID:           1,
			similarityThreshold: 0.1,
			limit:               5,
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12347, Username: "testuser3"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expenses
				expense1 := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Base expense",
					Timestamp:  time.Now(),
				}
				if err := client.CreateExpense(context.Background(), expense1); err != nil {
					return err
				}

				expense2 := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 150.0,
					Notes:      "Similar expense",
					Timestamp:  time.Now(),
				}
				if err := client.CreateExpense(context.Background(), expense2); err != nil {
					return err
				}

				// Update with similar embeddings
				embedding1 := generateTestEmbedding()
				embedding2 := generateTestEmbedding()
				// Make embedding2 slightly different
				embedding2[0] += 0.1
				if err := client.UpdateExpenseEmbedding(context.Background(), expense1.ID, embedding1, embedding1); err != nil {
					return err
				}
				return client.UpdateExpenseEmbedding(context.Background(), expense2.ID, embedding2, embedding2)
			},
			expectError: false,
			expectedLen: 1, // Should find expense2 as similar to expense1
		},
		{
			name:                "expense not found",
			expenseID:           999,
			similarityThreshold: 0.1,
			limit:               5,
			setupData:           func(client Storage) error { return nil },
			expectError:         true,
			errorType:           errors.ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database
			ctx := context.Background()
			logger := logger.NewMockLogger()
			client, err := NewClient(ctx, "postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable", logger)
			require.NoError(t, err)
			defer client.Close()

			// Setup test data
			if tt.setupData != nil {
				err = tt.setupData(client)
				require.NoError(t, err)
			}

			// Execute
			expenses, err := client.FindSimilarExpenses(ctx, tt.expenseID, tt.similarityThreshold, tt.limit)

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
		})
	}
}

func TestClient_UpdateExpenseEmbedding(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	tests := []struct {
		name              string
		expenseID         int64
		notesEmbedding    []float32
		categoryEmbedding []float32
		setupData         func(Storage) error
		expectError       bool
		errorType         errors.ErrorType
	}{
		{
			name:              "successful embedding update",
			expenseID:         1,
			notesEmbedding:    generateTestEmbedding(),
			categoryEmbedding: generateTestEmbedding(),
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12348, Username: "testuser4"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expense
				expense := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Test expense",
					Timestamp:  time.Now(),
				}
				return client.CreateExpense(context.Background(), expense)
			},
			expectError: false,
		},
		{
			name:              "expense not found",
			expenseID:         999,
			notesEmbedding:    generateTestEmbedding(),
			categoryEmbedding: generateTestEmbedding(),
			setupData:         func(client Storage) error { return nil },
			expectError:       true,
			errorType:         errors.ErrorTypeNotFound,
		},
		{
			name:              "invalid embedding dimensions",
			expenseID:         1,
			notesEmbedding:    []float32{0.1, 0.2}, // Too short
			categoryEmbedding: []float32{0.2, 0.3}, // Too short
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12349, Username: "testuser5"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expense
				expense := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Test expense",
					Timestamp:  time.Now(),
				}
				return client.CreateExpense(context.Background(), expense)
			},
			expectError: true,
			errorType:   errors.ErrorTypeDatabase,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database
			ctx := context.Background()
			logger := logger.NewMockLogger()
			client, err := NewClient(ctx, "postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable", logger)
			require.NoError(t, err)
			defer client.Close()

			// Setup test data
			if tt.setupData != nil {
				err = tt.setupData(client)
				require.NoError(t, err)
			}

			// Execute
			err = client.UpdateExpenseEmbedding(ctx, tt.expenseID, tt.notesEmbedding, tt.categoryEmbedding)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if appErr, ok := err.(*errors.AppError); ok {
					assert.Equal(t, tt.errorType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)

				// Verify the embedding was actually updated
				expense, err := client.GetExpenseByID(ctx, tt.expenseID)
				assert.NoError(t, err)
				assert.NotNil(t, expense)
				// Note: We can't easily check the actual vector values in the test
				// but we can verify the expense was found
			}
		})
	}
}

func TestClient_GetExpenseEmbedding(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping database test in short mode")
	}

	tests := []struct {
		name        string
		expenseID   int64
		setupData   func(Storage) error
		expectError bool
		errorType   errors.ErrorType
	}{
		{
			name:      "successful embedding retrieval",
			expenseID: 1,
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12350, Username: "testuser6"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expense
				expense := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Test expense",
					Timestamp:  time.Now(),
				}
				if err := client.CreateExpense(context.Background(), expense); err != nil {
					return err
				}

				// Update with embedding
				embedding := generateTestEmbedding()
				return client.UpdateExpenseEmbedding(context.Background(), expense.ID, embedding, embedding)
			},
			expectError: false,
		},
		{
			name:        "expense not found",
			expenseID:   999,
			setupData:   func(client Storage) error { return nil },
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
		{
			name:      "expense without embedding",
			expenseID: 1,
			setupData: func(client Storage) error {
				// Create test user
				user := &models.User{TelegramID: 12351, Username: "testuser7"}
				if err := client.CreateUser(context.Background(), user); err != nil {
					return err
				}

				// Create test expense without embedding
				expense := &models.Expense{
					UserID:     1,
					CategoryID: 1, // Assuming category 1 exists
					TotalPrice: 100.0,
					Notes:      "Test expense",
					Timestamp:  time.Now(),
				}
				return client.CreateExpense(context.Background(), expense)
			},
			expectError: true,
			errorType:   errors.ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database
			ctx := context.Background()
			logger := logger.NewMockLogger()
			client, err := NewClient(ctx, "postgres://postgres:password@localhost:5432/expense_tracker?sslmode=disable", logger)
			require.NoError(t, err)
			defer client.Close()

			// Setup test data
			if tt.setupData != nil {
				err = tt.setupData(client)
				require.NoError(t, err)
			}

			// Execute
			embedding, err := client.GetExpenseEmbedding(ctx, tt.expenseID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				if appErr, ok := err.(*errors.AppError); ok {
					assert.Equal(t, tt.errorType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, embedding)
				assert.NotNil(t, embedding.NotesEmbedding)
				assert.NotNil(t, embedding.CategoryEmbedding)
			}
		})
	}
}
