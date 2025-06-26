package database

import (
	"context"
	"testing"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func TestCategoryStorage_Mock(t *testing.T) {
	mock := NewMockStorage().(*MockStorage)
	mock.ClearMockData()
	ctx := context.Background()

	cat1 := &models.Category{Name: "Petrol", Emoji: "⛽", Group: "Vehicle"}
	cat2 := &models.Category{Name: "Grocery", Emoji: "🛒", Group: "Daily Living"}
	cat3 := &models.Category{Name: "Service", Emoji: "🔧", Group: "Vehicle"}
	mock.AddMockCategory(cat1)
	mock.AddMockCategory(cat2)
	mock.AddMockCategory(cat3)

	t.Run("GetAllCategories returns all", func(t *testing.T) {
		cats, err := mock.GetAllCategories(ctx)
		require.NoError(t, err)
		require.Len(t, cats, 3)
	})

	t.Run("GetCategoriesByGroup returns correct group", func(t *testing.T) {
		vehicleCats, err := mock.GetCategoriesByGroup(ctx, "Vehicle")
		require.NoError(t, err)
		require.Len(t, vehicleCats, 2)
		groups := map[string]bool{vehicleCats[0].Group: true, vehicleCats[1].Group: true}
		require.True(t, groups["Vehicle"])
	})

	t.Run("GetCategoryByName returns correct category", func(t *testing.T) {
		cat, err := mock.GetCategoryByName(ctx, "Grocery")
		require.NoError(t, err)
		require.Equal(t, "Grocery", cat.Name)
		cat, err = mock.GetCategoryByName(ctx, "Petrol")
		require.NoError(t, err)
		require.Equal(t, "Petrol", cat.Name)
	})

	t.Run("GetCategoryByName returns nil for missing", func(t *testing.T) {
		cat, err := mock.GetCategoryByName(ctx, "Nonexistent")
		require.Error(t, err)
		require.Nil(t, cat)
	})
}
