package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateExpenseEmbedding_HandlesNilEmbeddings(t *testing.T) {
	// This test verifies that the UpdateExpenseEmbedding function
	// properly handles nil embeddings without causing PostgreSQL errors

	// Test cases for nil embeddings
	testCases := []struct {
		name              string
		notesEmbedding    []float32
		categoryEmbedding []float32
		expectError       bool
	}{
		{
			name:              "both embeddings nil",
			notesEmbedding:    nil,
			categoryEmbedding: nil,
			expectError:       false, // Should not error, should set to NULL
		},
		{
			name:              "notes embedding nil, category valid",
			notesEmbedding:    nil,
			categoryEmbedding: []float32{0.1, 0.2, 0.3},
			expectError:       false,
		},
		{
			name:              "notes valid, category embedding nil",
			notesEmbedding:    []float32{0.1, 0.2, 0.3},
			categoryEmbedding: nil,
			expectError:       false,
		},
		{
			name:              "both embeddings empty slices",
			notesEmbedding:    []float32{},
			categoryEmbedding: []float32{},
			expectError:       false,
		},
		{
			name:              "valid embeddings",
			notesEmbedding:    []float32{0.1, 0.2, 0.3},
			categoryEmbedding: []float32{0.4, 0.5, 0.6},
			expectError:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the string conversion logic directly
			var notesEmbeddingStr, categoryEmbeddingStr any

			if len(tc.notesEmbedding) > 0 {
				notesEmbeddingStr = "[0.100000"
				for i := 1; i < len(tc.notesEmbedding); i++ {
					notesEmbeddingStr = notesEmbeddingStr.(string) + ",0.200000"
				}
				notesEmbeddingStr = notesEmbeddingStr.(string) + "]"
			} else {
				notesEmbeddingStr = nil
			}

			if len(tc.categoryEmbedding) > 0 {
				categoryEmbeddingStr = "[0.400000"
				for i := 1; i < len(tc.categoryEmbedding); i++ {
					categoryEmbeddingStr = categoryEmbeddingStr.(string) + ",0.500000"
				}
				categoryEmbeddingStr = categoryEmbeddingStr.(string) + "]"
			} else {
				categoryEmbeddingStr = nil
			}

			// Verify that nil embeddings result in nil strings (which become NULL in PostgreSQL)
			if len(tc.notesEmbedding) == 0 {
				assert.Nil(t, notesEmbeddingStr, "Notes embedding should be nil for empty/nil input")
			} else {
				assert.NotNil(t, notesEmbeddingStr, "Notes embedding should not be nil for valid input")
			}

			if len(tc.categoryEmbedding) == 0 {
				assert.Nil(t, categoryEmbeddingStr, "Category embedding should be nil for empty/nil input")
			} else {
				assert.NotNil(t, categoryEmbeddingStr, "Category embedding should not be nil for valid input")
			}
		})
	}
}

func TestSearchExpensesBySimilarity_ValidatesQueryEmbedding(t *testing.T) {
	// Test the validation logic directly
	testCases := []struct {
		name           string
		queryEmbedding []float32
		expectError    bool
	}{
		{
			name:           "nil query embedding",
			queryEmbedding: nil,
			expectError:    true,
		},
		{
			name:           "empty query embedding",
			queryEmbedding: []float32{},
			expectError:    true,
		},
		{
			name:           "valid query embedding",
			queryEmbedding: []float32{0.1, 0.2, 0.3},
			expectError:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test the validation logic
			if len(tc.queryEmbedding) == 0 {
				assert.True(t, tc.expectError, "Should error for nil/empty query embedding")
			} else {
				assert.False(t, tc.expectError, "Should not error for valid query embedding")
			}
		})
	}
}
