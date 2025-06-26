package database

import (
	"context"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
)

// CategoryStorage defines operations for category management
type CategoryStorage interface {
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	GetCategoriesByGroup(ctx context.Context, group string) ([]*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
}

// GetAllCategories retrieves all categories
func (c *Client) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	var categories []*models.Category
	query := `SELECT * FROM categories ORDER BY "group", name`

	err := c.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoriesByGroup retrieves categories by group
func (c *Client) GetCategoriesByGroup(ctx context.Context, group string) ([]*models.Category, error) {
	var categories []*models.Category
	query := `SELECT * FROM categories WHERE "group" = $1 ORDER BY name`

	err := c.db.SelectContext(ctx, &categories, query, group)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// GetCategoryByName retrieves a category by name
func (c *Client) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	var category models.Category
	query := `SELECT * FROM categories WHERE name = $1`

	err := c.db.GetContext(ctx, &category, query, name)
	if err != nil {
		if isNoRows(err) {
			return nil, errNotFound
		}
		return nil, err
	}

	return &category, nil
}
