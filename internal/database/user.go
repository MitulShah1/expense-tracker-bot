package database

import (
	"context"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
)

// UserStorage defines operations for user management
type UserStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
}

// CreateUser creates a new user
func (c *Client) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (telegram_id) DO UPDATE SET
			username = EXCLUDED.username,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			updated_at = now()
		RETURNING id, created_at, updated_at`

	return c.db.QueryRowxContext(ctx, query,
		user.TelegramID, user.Username, user.FirstName, user.LastName).
		StructScan(user)
}

// GetUserByTelegramID retrieves a user by Telegram ID
func (c *Client) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE telegram_id = $1`

	err := c.db.GetContext(ctx, &user, query, telegramID)
	if err != nil {
		if isNoRows(err) {
			return nil, errNotFound
		}
		return nil, err
	}

	return &user, nil
}
