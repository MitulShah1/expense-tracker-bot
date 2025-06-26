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

// UserService provides user-related business logic
type UserService struct {
	db        database.Storage
	logger    logger.Logger
	validator *validation.Validator
}

// NewUserService creates a new user service
func NewUserService(db database.Storage, logger logger.Logger) *UserService {
	return &UserService{
		db:        db,
		logger:    logger,
		validator: validation.NewValidator(),
	}
}

// GetOrCreateUser gets an existing user or creates a new one
func (s *UserService) GetOrCreateUser(ctx context.Context, telegramID int64, username, firstName, lastName string) (*models.User, error) {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return nil, err
	}

	if username != "" {
		if err := s.validator.ValidateUsername(username); err != nil {
			return nil, err
		}
	}

	if firstName != "" {
		if err := s.validator.ValidateName(firstName, "first name"); err != nil {
			return nil, err
		}
	}

	if lastName != "" {
		if err := s.validator.ValidateName(lastName, "last name"); err != nil {
			return nil, err
		}
	}

	// Try to get existing user
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get user", err)
	}

	// If user exists, return it
	if user != nil {
		return user, nil
	}

	// Create new user
	newUser := &models.User{
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		LastName:   lastName,
	}

	if err := s.db.CreateUser(ctx, newUser); err != nil {
		s.logger.Error(ctx, "Failed to create user", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to create user", err)
	}

	s.logger.Info(ctx, "User created successfully",
		logger.Int("user_id", int(newUser.ID)),
		logger.Int("telegram_id", int(telegramID)))

	return newUser, nil
}

// GetUserByTelegramID retrieves a user by Telegram ID
func (s *UserService) GetUserByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	// Validate input
	if err := s.validator.ValidateTelegramID(telegramID); err != nil {
		return nil, err
	}

	// Get user from database
	user, err := s.db.GetUserByTelegramID(ctx, telegramID)
	if err != nil {
		s.logger.Error(ctx, "Failed to get user by Telegram ID", logger.ErrorField(err))
		return nil, errors.NewDatabaseError("Failed to get user", err)
	}

	if user == nil {
		return nil, errors.NewNotFoundError("User not found", fmt.Sprintf("User with Telegram ID %d not found", telegramID))
	}

	return user, nil
}
