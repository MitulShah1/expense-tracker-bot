package database

import (
	"context"
	"testing"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func TestUserStorage_Mock(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	t.Run("Create and get user", func(t *testing.T) {
		user := &models.User{
			TelegramID: 12345,
			Username:   "testuser",
			FirstName:  "Test",
			LastName:   "User",
		}
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotZero(t, user.ID)

		fetched, err := mock.GetUserByTelegramID(ctx, 12345)
		require.NoError(t, err)
		require.Equal(t, user.ID, fetched.ID)
		require.Equal(t, "testuser", fetched.Username)
	})

	t.Run("Update user on conflict", func(t *testing.T) {
		user := &models.User{
			TelegramID: 12345,
			Username:   "updateduser",
			FirstName:  "Updated",
			LastName:   "User",
		}
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)

		fetched, err := mock.GetUserByTelegramID(ctx, 12345)
		require.NoError(t, err)
		require.Equal(t, "updateduser", fetched.Username)
		require.Equal(t, "Updated", fetched.FirstName)
	})

	t.Run("User not found", func(t *testing.T) {
		fetched, err := mock.GetUserByTelegramID(ctx, 99999)
		require.Error(t, err)
		require.Nil(t, fetched)
	})
}
