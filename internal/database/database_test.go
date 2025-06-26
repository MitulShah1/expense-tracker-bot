package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func newTestUser(telegramID int64) *models.User {
	return &models.User{
		TelegramID: telegramID,
		Username:   "testuser",
		FirstName:  "Test",
		LastName:   "User",
	}
}

func newTestCategory(name, group string) *models.Category {
	return &models.Category{
		Name:  name,
		Emoji: "⛽",
		Group: group,
	}
}

func newTestExpense(userID, categoryID int64) *models.Expense {
	return &models.Expense{
		UserID:      userID,
		CategoryID:  categoryID,
		VehicleType: sql.NullString{String: "CAR", Valid: true},
		Odometer:    123.4,
		PetrolPrice: 100.5,
		TotalPrice:  500.0,
		Notes:       "Test expense",
		Timestamp:   time.Now(),
	}
}

func TestUserStorage(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	t.Run("Create and fetch user", func(t *testing.T) {
		user := newTestUser(12345)
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		require.NotZero(t, user.ID)
		fetched, err := mock.GetUserByTelegramID(ctx, 12345)
		require.NoError(t, err)
		require.Equal(t, user.ID, fetched.ID)
		require.Equal(t, "testuser", fetched.Username)
	})

	t.Run("Update user on conflict", func(t *testing.T) {
		user := newTestUser(12345)
		user.Username = "updateduser"
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		fetched, err := mock.GetUserByTelegramID(ctx, 12345)
		require.NoError(t, err)
		require.Equal(t, "updateduser", fetched.Username)
	})

	t.Run("User not found", func(t *testing.T) {
		fetched, err := mock.GetUserByTelegramID(ctx, 99999)
		require.Error(t, err)
		require.Nil(t, fetched)
		require.True(t, isNoRows(err))
	})
}

func TestCategoryStorage(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	cat1 := newTestCategory("Petrol", "Vehicle")
	cat2 := newTestCategory("Grocery", "Daily Living")
	mock.AddMockCategory(cat1)
	mock.AddMockCategory(cat2)

	t.Run("Get all categories", func(t *testing.T) {
		cats, err := mock.GetAllCategories(ctx)
		require.NoError(t, err)
		require.Len(t, cats, 2)
	})

	t.Run("Get categories by group", func(t *testing.T) {
		cats, err := mock.GetCategoriesByGroup(ctx, "Vehicle")
		require.NoError(t, err)
		require.Len(t, cats, 1)
		require.Equal(t, "Petrol", cats[0].Name)
	})

	t.Run("Get category by name", func(t *testing.T) {
		cat, err := mock.GetCategoryByName(ctx, "Grocery")
		require.NoError(t, err)
		require.Equal(t, "Grocery", cat.Name)
	})

	t.Run("Category not found", func(t *testing.T) {
		cat, err := mock.GetCategoryByName(ctx, "Nonexistent")
		require.Error(t, err)
		require.Nil(t, cat)
		require.True(t, isNoRows(err))
	})
}

func TestExpenseStorage(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	// Setup user and category
	user := newTestUser(11111)
	err := mock.CreateUser(ctx, user)
	require.NoError(t, err)
	cat := newTestCategory("Petrol", "Vehicle")
	mock.AddMockCategory(cat)

	t.Run("Create and fetch expense", func(t *testing.T) {
		exp := newTestExpense(user.ID, cat.ID)
		err := mock.CreateExpense(ctx, exp)
		require.NoError(t, err)
		require.NotZero(t, exp.ID)
		fetched, err := mock.GetExpenseByID(ctx, exp.ID)
		require.NoError(t, err)
		require.Equal(t, exp.ID, fetched.ID)
	})

	t.Run("Update expense", func(t *testing.T) {
		exp := newTestExpense(user.ID, cat.ID)
		err := mock.CreateExpense(ctx, exp)
		require.NoError(t, err)
		exp.Notes = "Updated note"
		err = mock.UpdateExpense(ctx, exp)
		require.NoError(t, err)
		fetched, err := mock.GetExpenseByID(ctx, exp.ID)
		require.NoError(t, err)
		require.Equal(t, "Updated note", fetched.Notes)
	})

	t.Run("Get expenses by user ID", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		err = mock.CreateExpense(ctx, newTestExpense(user.ID, cat.ID))
		require.NoError(t, err)
		err = mock.CreateExpense(ctx, newTestExpense(user.ID, cat.ID))
		require.NoError(t, err)
		exps, err := mock.GetExpensesByUserID(ctx, user.ID)
		require.NoError(t, err)
		require.Len(t, exps, 2)
	})

	t.Run("Get expenses by Telegram ID", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		err = mock.CreateExpense(ctx, newTestExpense(user.ID, cat.ID))
		require.NoError(t, err)
		exps, err := mock.GetExpensesByTelegramID(ctx, user.TelegramID)
		require.NoError(t, err)
		require.Len(t, exps, 1)
	})

	t.Run("Get expenses by date range", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		e1 := newTestExpense(user.ID, cat.ID)
		e1.Timestamp = time.Now().Add(-48 * time.Hour)
		err = mock.CreateExpense(ctx, e1)
		require.NoError(t, err)
		e2 := newTestExpense(user.ID, cat.ID)
		e2.Timestamp = time.Now().Add(-24 * time.Hour)
		err = mock.CreateExpense(ctx, e2)
		require.NoError(t, err)
		start := time.Now().Add(-36 * time.Hour)
		end := time.Now()
		exps, err := mock.GetExpensesByDateRange(ctx, user.ID, start, end)
		require.NoError(t, err)
		require.Len(t, exps, 1)
		require.Equal(t, e2.ID, exps[0].ID)
	})

	t.Run("Delete expense (soft)", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		e := newTestExpense(user.ID, cat.ID)
		err = mock.CreateExpense(ctx, e)
		require.NoError(t, err)
		err = mock.DeleteExpense(ctx, e.ID, user.ID)
		require.NoError(t, err)
		fetched, err := mock.GetExpenseByID(ctx, e.ID)
		require.Error(t, err)
		require.Nil(t, fetched)
	})

	t.Run("Expense not found", func(t *testing.T) {
		e, err := mock.GetExpenseByID(ctx, 99999)
		require.Error(t, err)
		require.Nil(t, e)
		require.True(t, isNoRows(err))
	})

	t.Run("Get expense stats", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		e1 := newTestExpense(user.ID, cat.ID)
		e1.TotalPrice = 100
		err = mock.CreateExpense(ctx, e1)
		require.NoError(t, err)
		e2 := newTestExpense(user.ID, cat.ID)
		e2.TotalPrice = 200
		err = mock.CreateExpense(ctx, e2)
		require.NoError(t, err)
		stats, err := mock.GetExpenseStats(ctx, user.ID)
		require.NoError(t, err)
		require.Equal(t, int64(2), stats.TotalExpenses)
		require.Equal(t, 300.0, stats.TotalSpent)
		require.Equal(t, 150.0, stats.AvgExpense)
		require.Equal(t, 100.0, stats.MinExpense)
		require.Equal(t, 200.0, stats.MaxExpense)
	})
}

func TestHelpers(t *testing.T) {
	require.True(t, isNoRows(sql.ErrNoRows))
	require.False(t, isNoRows(sql.ErrConnDone))
	require.EqualError(t, errNotFound, "not found")
}
