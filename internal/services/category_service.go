// Package services provides business logic services for the expense tracker bot.
package services

import (
	"context"
	"fmt"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/MitulShah1/expense-tracker-bot/internal/validation"
)

// CategoryService provides category-related business logic
type CategoryService struct {
	db        database.Storage
	logger    logger.Logger
	validator *validation.Validator
}

// NewCategoryService creates a new category service
func NewCategoryService(db database.Storage, logger logger.Logger) *CategoryService {
	return &CategoryService{
		db:        db,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// GetAllCategories retrieves all categories
func (s *CategoryService) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	// Get categories from database
	categories, err := s.db.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all categories", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get categories", err)
	}

	return categories, nil
}

// GetCategoryByName retrieves a category by name
func (s *CategoryService) GetCategoryByName(ctx context.Context, categoryName string) (*models.Category, error) {
	// Validate input
	if err := s.validator.ValidateCategoryName(categoryName); err != nil {
		return nil, err
	}

	// Get category from database
	category, err := s.db.GetCategoryByName(ctx, categoryName)
	if err != nil {
		s.logger.Error(ctx, "Failed to get category by name", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get category", err)
	}

	if category == nil {
		return nil, errors.NewNotFoundError("Category not found", fmt.Sprintf("Category '%s' not found", categoryName))
	}

	return category, nil
}

// GetCategoriesByGroup retrieves categories by group
func (s *CategoryService) GetCategoriesByGroup(ctx context.Context, groupName string) ([]*models.Category, error) {
	// Validate input
	if groupName == "" {
		return nil, errors.NewValidationError("Group name is required", "Group name cannot be empty")
	}

	if len(groupName) > 100 {
		return nil, errors.NewValidationError("Group name too long", "Group name must be 100 characters or less")
	}

	// Get categories from database
	categories, err := s.db.GetCategoriesByGroup(ctx, groupName)
	if err != nil {
		s.logger.Error(ctx, "Failed to get categories by group", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get categories", err)
	}

	return categories, nil
}

// GetCategoryGroups retrieves all category groups
func (s *CategoryService) GetCategoryGroups(ctx context.Context) ([]*models.CategoryGroup, error) {
	// Get all categories first
	categories, err := s.db.GetAllCategories(ctx)
	if err != nil {
		s.logger.Error(ctx, "Failed to get all categories for groups", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get categories", err)
	}

	// Group categories by their group
	groupMap := make(map[string]*models.CategoryGroup)
	for _, category := range categories {
		if _, exists := groupMap[category.Group]; !exists {
			groupMap[category.Group] = &models.CategoryGroup{
				Name:  category.Group,
				Emoji: getGroupEmoji(category.Group),
			}
		}
		// Add category to group (you might want to add a Categories field to CategoryGroup)
	}

	// Convert map to slice
	groups := make([]*models.CategoryGroup, 0, len(groupMap))
	for _, group := range groupMap {
		groups = append(groups, group)
	}

	return groups, nil
}

// getGroupEmoji returns the emoji for a category group
func getGroupEmoji(groupName string) string {
	emojiMap := map[string]string{
		"Vehicle":       "🚗",
		"Home":          "🏠",
		"Daily Living":  "🏪",
		"Entertainment": "🎬",
		"Health":        "🏥",
		"Education":     "📚",
		"Travel":        "✈️",
		"Investments":   "💹",
		"Gifts":         "🎁",
		"Other":         "📌",
	}

	if emoji, exists := emojiMap[groupName]; exists {
		return emoji
	}
	return "📌" // Default emoji
}
