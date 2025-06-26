package database

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func newExpense(userID, categoryID int64, ts time.Time) *models.Expense {
	return &models.Expense{
		UserID:      userID,
		CategoryID:  categoryID,
		VehicleType: sql.NullString{String: "CAR", Valid: true},
		Odometer:    100,
		PetrolPrice: 90,
		TotalPrice:  500,
		Notes:       "test",
		Timestamp:   ts,
	}
}

func TestExpenseStorage_Mock(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	// Setup user and category
	user := &models.User{TelegramID: 123, Username: "u", FirstName: "f", LastName: "l"}
	err := mock.CreateUser(ctx, user)
	require.NoError(t, err)
	cat := &models.Category{Name: "Petrol", Emoji: "⛽", Group: "Vehicle"}
	mock.AddMockCategory(cat)

	t.Run("Create and get by ID", func(t *testing.T) {
		e := newExpense(user.ID, cat.ID, time.Now())
		err := mock.CreateExpense(ctx, e)
		require.NoError(t, err)
		fetched, err := mock.GetExpenseByID(ctx, e.ID)
		require.NoError(t, err)
		require.Equal(t, e.ID, fetched.ID)
	})

	t.Run("Update expense", func(t *testing.T) {
		e := newExpense(user.ID, cat.ID, time.Now())
		err := mock.CreateExpense(ctx, e)
		require.NoError(t, err)
		e.Notes = "updated"
		err = mock.UpdateExpense(ctx, e)
		require.NoError(t, err)
		fetched, err := mock.GetExpenseByID(ctx, e.ID)
		require.NoError(t, err)
		require.Equal(t, "updated", fetched.Notes)
	})

	t.Run("Get by user ID", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		err = mock.CreateExpense(ctx, newExpense(user.ID, cat.ID, time.Now()))
		require.NoError(t, err)
		err = mock.CreateExpense(ctx, newExpense(user.ID, cat.ID, time.Now()))
		require.NoError(t, err)
		exps, err := mock.GetExpensesByUserID(ctx, user.ID)
		require.NoError(t, err)
		require.Len(t, exps, 2)
	})

	t.Run("Get by telegram ID", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		err = mock.CreateExpense(ctx, newExpense(user.ID, cat.ID, time.Now()))
		require.NoError(t, err)
		exps, err := mock.GetExpensesByTelegramID(ctx, user.TelegramID)
		require.NoError(t, err)
		require.Len(t, exps, 1)
	})

	t.Run("Get by date range", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		e1 := newExpense(user.ID, cat.ID, time.Now().Add(-48*time.Hour))
		err = mock.CreateExpense(ctx, e1)
		require.NoError(t, err)
		e2 := newExpense(user.ID, cat.ID, time.Now().Add(-24*time.Hour))
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
		e := newExpense(user.ID, cat.ID, time.Now())
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
	})

	t.Run("Get expense stats", func(t *testing.T) {
		mock.ClearMockData()
		err := mock.CreateUser(ctx, user)
		require.NoError(t, err)
		mock.AddMockCategory(cat)
		e1 := newExpense(user.ID, cat.ID, time.Now())
		e1.TotalPrice = 100
		err = mock.CreateExpense(ctx, e1)
		require.NoError(t, err)
		e2 := newExpense(user.ID, cat.ID, time.Now())
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
