package services

import (
	"context"
	errorsstd "errors"
	"fmt"
	"strings"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	errors "github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
)

// VectorService provides vector-based search and embedding functionality
type VectorService struct {
	db     database.Storage
	logger logger.Logger
}

// NewVectorService creates a new vector service
func NewVectorService(db database.Storage, logger logger.Logger) *VectorService {
	return &VectorService{
		db:     db,
		logger: logger,
	}
}

// SearchExpensesByQuery performs semantic search on expenses using natural language query
func (s *VectorService) SearchExpensesByQuery(ctx context.Context, telegramID int64, query string, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	// Get user by Telegram ID
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Generate embedding for the query
	queryEmbedding, err := s.generateEmbedding(ctx, query)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate query embedding", logger.ErrorField(err))
		return nil, errors.NewInternalError("Failed to process search query", err)
	}

	// Search for similar expenses
	expenses, err := s.db.SearchExpensesBySimilarity(ctx, user.ID, queryEmbedding, similarityThreshold, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to search expenses by similarity", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to search expenses", err)
	}

	return expenses, nil
}

// FindSimilarExpenses finds expenses similar to a given expense
func (s *VectorService) FindSimilarExpenses(ctx context.Context, expenseID int64, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	expenses, err := s.db.FindSimilarExpenses(ctx, expenseID, similarityThreshold, limit)
	if err != nil {
		s.logger.Error(ctx, "Failed to find similar expenses", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to find similar expenses", err)
	}

	return expenses, nil
}

// UpdateExpenseEmbeddings updates the vector embeddings for an expense
func (s *VectorService) UpdateExpenseEmbeddings(ctx context.Context, expenseID int64) error {
	// Get the expense
	expense, err := s.db.GetExpenseByID(ctx, expenseID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get expense for embedding update", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get expense", err)
	}

	if expense == nil {
		return errors.NewNotFoundError("Expense not found", fmt.Sprintf("Expense with ID %d not found", expenseID))
	}

	// Generate embeddings for notes and category
	var notesEmbedding, categoryEmbedding []float32

	if expense.Notes != "" {
		notesEmbedding, err = s.generateEmbedding(ctx, expense.Notes)
		if err != nil {
			s.logger.Error(ctx, "Failed to generate notes embedding", logger.ErrorField(err))
			return errors.NewInternalError("Failed to generate notes embedding", err)
		}
	}

	// Generate category embedding
	categoryText := fmt.Sprintf("%s %s", expense.CategoryName, expense.CategoryGroup)
	categoryEmbedding, err = s.generateEmbedding(ctx, categoryText)
	if err != nil {
		s.logger.Error(ctx, "Failed to generate category embedding", logger.ErrorField(err))
		return errors.NewInternalError("Failed to generate category embedding", err)
	}

	// Update embeddings in database
	err = s.db.UpdateExpenseEmbedding(ctx, expenseID, notesEmbedding, categoryEmbedding)
	if err != nil {
		s.logger.Error(ctx, "Failed to update expense embedding", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to update expense embedding", err)
	}

	s.logger.Info(ctx, "Expense embeddings updated successfully",
		logger.Int("expense_id", int(expenseID)))

	return nil
}

// BatchUpdateEmbeddings updates embeddings for all expenses without embeddings
func (s *VectorService) BatchUpdateEmbeddings(ctx context.Context, telegramID int64) error {
	// Get user by Telegram ID
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	// Get all expenses for the user
	expenses, err := s.db.GetExpensesByUserID(ctx, user.ID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get expenses for batch embedding update", logger.ErrorField(err))
		return errors.NewDatabaseError("Failed to get expenses", err)
	}

	updatedCount := 0
	for _, expense := range expenses {
		// Check if expense already has embeddings
		embedding, err := s.db.GetExpenseEmbedding(ctx, expense.ID)
		if err != nil {
			// Check if it's a not found error (no embedding exists yet)
			if appErr, ok := err.(*errors.AppError); ok && appErr.IsNotFoundError() {
				// This is expected - no embedding exists yet
			} else {
				s.logger.Error(ctx, "Failed to check expense embedding", logger.ErrorField(err))
				continue
			}
		}

		// Skip if embedding already exists
		if embedding != nil && len(embedding.NotesEmbedding) > 0 {
			continue
		}

		// Update embeddings for this expense
		if err := s.UpdateExpenseEmbeddings(ctx, expense.ID); err != nil {
			s.logger.Error(ctx, "Failed to update embeddings for expense",
				logger.Int("expense_id", int(expense.ID)),
				logger.ErrorField(err))
			continue
		}

		updatedCount++
	}

	s.logger.Info(ctx, "Batch embedding update completed",
		logger.Int("total_expenses", len(expenses)),
		logger.Int("updated_count", updatedCount))

	return nil
}

// generateEmbedding generates a vector embedding for the given text
// This is a placeholder implementation - in a real application, you would:
// 1. Use an embedding service (OpenAI, Cohere, etc.)
// 2. Cache embeddings to avoid repeated API calls
// 3. Handle rate limiting and errors
func (s *VectorService) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Normalize text
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errorsstd.New("empty text provided for embedding")
	}

	// This is a simplified placeholder implementation
	// In production, you would call an embedding API like:
	// - OpenAI's text-embedding-ada-002
	// - Cohere's embed-multilingual-v3.0
	// - Hugging Face's sentence-transformers

	// For now, we'll create a simple hash-based embedding
	// This is NOT suitable for production - just for demonstration
	embedding := make([]float32, 1536)
	hash := 0
	for _, char := range text {
		hash = (hash*31 + int(char)) % 1536
		embedding[hash] = float32(char) / 255.0
	}

	s.logger.Debug(ctx, "Generated embedding for text",
		logger.String("text", text[:myMin(len(text), 50)]),
		logger.Int("embedding_length", len(embedding)))

	return embedding, nil
}

// myMin returns the minimum of two integers
func myMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
