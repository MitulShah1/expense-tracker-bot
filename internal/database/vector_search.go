package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
)

// VectorSearchStorage defines operations for vector-based search
type VectorSearchStorage interface {
	SearchExpensesBySimilarity(ctx context.Context, userID int64, queryEmbedding []float32, similarityThreshold float32, limit int) ([]*models.Expense, error)
	FindSimilarExpenses(ctx context.Context, expenseID int64, similarityThreshold float32, limit int) ([]*models.Expense, error)
	UpdateExpenseEmbedding(ctx context.Context, expenseID int64, notesEmbedding, categoryEmbedding []float32) error
	GetExpenseEmbedding(ctx context.Context, expenseID int64) (*models.ExpenseEmbedding, error)
}

// ExpenseEmbedding represents the vector embeddings for an expense
type ExpenseEmbedding struct {
	ID                int64     `db:"id"`
	NotesEmbedding    []float32 `db:"notes_embedding"`
	CategoryEmbedding []float32 `db:"category_embedding"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// SearchExpensesBySimilarity searches for expenses using semantic similarity
func (c *Client) SearchExpensesBySimilarity(ctx context.Context, userID int64, queryEmbedding []float32, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	var expenses []*models.Expense

	// Validate query embedding
	if len(queryEmbedding) == 0 {
		return nil, errors.New("invalid query embedding: cannot be nil or empty")
	}

	// Convert []float32 to string representation for PostgreSQL vector
	embeddingStr := fmt.Sprintf("[%f", queryEmbedding[0])
	for i := 1; i < len(queryEmbedding); i++ {
		embeddingStr += fmt.Sprintf(",%f", queryEmbedding[i])
	}
	embeddingStr += "]"

	query := `
		SELECT 
			e.id, e.user_id, e.category_id, e.vehicle_type, e.odometer, 
			e.petrol_price, e.total_price, e.notes, e.timestamp, 
			e.created_at, e.updated_at, e.deleted_at,
			c.name as category_name, c.emoji as category_emoji, c."group" as category_group,
			1 - (e.notes_embedding <=> $2::vector) as similarity
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = $1 
			AND e.deleted_at IS NULL
			AND e.notes_embedding IS NOT NULL
			AND 1 - (e.notes_embedding <=> $2::vector) > $3
		ORDER BY e.notes_embedding <=> $2::vector
		LIMIT $4`

	err := c.db.SelectContext(ctx, &expenses, query, userID, embeddingStr, similarityThreshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search expenses by similarity: %w", err)
	}

	return expenses, nil
}

// FindSimilarExpenses finds expenses similar to a given expense
func (c *Client) FindSimilarExpenses(ctx context.Context, expenseID int64, similarityThreshold float32, limit int) ([]*models.Expense, error) {
	var expenses []*models.Expense

	query := `
		SELECT 
			e.id, e.user_id, e.category_id, e.vehicle_type, e.odometer, 
			e.petrol_price, e.total_price, e.notes, e.timestamp, 
			e.created_at, e.updated_at, e.deleted_at,
			c.name as category_name, c.emoji as category_emoji, c."group" as category_group,
			1 - (e.notes_embedding <=> target.notes_embedding) as similarity
		FROM expenses e
		JOIN categories c ON e.category_id = c.id,
		LATERAL (
			SELECT notes_embedding, user_id 
			FROM expenses 
			WHERE id = $1 AND deleted_at IS NULL
		) target
		WHERE e.user_id = target.user_id 
			AND e.id != $1
			AND e.deleted_at IS NULL
			AND e.notes_embedding IS NOT NULL
			AND 1 - (e.notes_embedding <=> target.notes_embedding) > $2
		ORDER BY e.notes_embedding <=> target.notes_embedding
		LIMIT $3`

	err := c.db.SelectContext(ctx, &expenses, query, expenseID, similarityThreshold, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar expenses: %w", err)
	}

	return expenses, nil
}

// UpdateExpenseEmbedding updates the vector embeddings for an expense
func (c *Client) UpdateExpenseEmbedding(ctx context.Context, expenseID int64, notesEmbedding, categoryEmbedding []float32) error {
	// Convert []float32 to string representation for PostgreSQL vector
	var notesEmbeddingStr, categoryEmbeddingStr any

	if len(notesEmbedding) > 0 {
		notesEmbeddingStr = fmt.Sprintf("[%f", notesEmbedding[0])
		for i := 1; i < len(notesEmbedding); i++ {
			notesEmbeddingStr = fmt.Sprintf("%s,%f", notesEmbeddingStr, notesEmbedding[i])
		}
		notesEmbeddingStr = fmt.Sprintf("%s]", notesEmbeddingStr)
	} else {
		notesEmbeddingStr = nil
	}

	if len(categoryEmbedding) > 0 {
		categoryEmbeddingStr = fmt.Sprintf("[%f", categoryEmbedding[0])
		for i := 1; i < len(categoryEmbedding); i++ {
			categoryEmbeddingStr = fmt.Sprintf("%s,%f", categoryEmbeddingStr, categoryEmbedding[i])
		}
		categoryEmbeddingStr = fmt.Sprintf("%s]", categoryEmbeddingStr)
	} else {
		categoryEmbeddingStr = nil
	}

	query := `
		UPDATE expenses 
		SET notes_embedding = $2::vector, category_embedding = $3::vector, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`

	result, err := c.db.ExecContext(ctx, query, expenseID, notesEmbeddingStr, categoryEmbeddingStr)
	if err != nil {
		return fmt.Errorf("failed to update expense embedding: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errNotFound
	}

	return nil
}

// GetExpenseEmbedding retrieves the vector embeddings for an expense
func (c *Client) GetExpenseEmbedding(ctx context.Context, expenseID int64) (*models.ExpenseEmbedding, error) {
	var embedding models.ExpenseEmbedding

	query := `
		SELECT id, notes_embedding, category_embedding, updated_at
		FROM expenses 
		WHERE id = $1 AND deleted_at IS NULL`

	err := c.db.GetContext(ctx, &embedding, query, expenseID)
	if err != nil {
		if isNoRows(err) {
			return nil, errNotFound
		}
		return nil, fmt.Errorf("failed to get expense embedding: %w", err)
	}

	return &embedding, nil
}
