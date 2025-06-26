// Package bot implements the Telegram bot functionality for expense tracking.
package bot

import (
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/require"
)

func TestGetMainMenuKeyboard(t *testing.T) {
	t.Run("should create main menu keyboard", func(t *testing.T) {
		keyboard := GetMainMenuKeyboard()

		require.NotNil(t, keyboard)
		require.True(t, keyboard.OneTimeKeyboard)
		require.Len(t, keyboard.Keyboard, 3) // 3 rows

		// Check first row
		require.Len(t, keyboard.Keyboard[0], 2)
		require.Equal(t, "📝 Add Expense", keyboard.Keyboard[0][0].Text)
		require.Equal(t, "📋 List Expenses", keyboard.Keyboard[0][1].Text)

		// Check second row
		require.Len(t, keyboard.Keyboard[1], 2)
		require.Equal(t, "✏️ Edit Expense", keyboard.Keyboard[1][0].Text)
		require.Equal(t, "🗑️ Delete Expense", keyboard.Keyboard[1][1].Text)

		// Check third row
		require.Len(t, keyboard.Keyboard[2], 2)
		require.Equal(t, "📊 Reports", keyboard.Keyboard[2][0].Text)
		require.Equal(t, "📈 Dashboard", keyboard.Keyboard[2][1].Text)
	})
}

func TestGetCategoryKeyboard(t *testing.T) {
	t.Run("should create category keyboard", func(t *testing.T) {
		keyboard := GetCategoryKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 5) // 5 rows

		// Check first button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "🚗 Vehicle", firstButton.Text)
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "group_Vehicle", *firstButton.CallbackData)
	})
}

func TestGetCategoryGroupKeyboard(t *testing.T) {
	t.Run("should create category group keyboard with categories", func(t *testing.T) {
		categories := []*models.Category{
			{ID: 1, Name: "Petrol", Emoji: "⛽", Group: "Vehicle"},
			{ID: 2, Name: "Diesel", Emoji: "⛽", Group: "Vehicle"},
		}

		keyboard := GetCategoryGroupKeyboard(categories)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 2 categories + 1 back button

		// Check first category button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "⛽ Petrol", firstButton.Text)
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "category_Petrol", *firstButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back to Groups", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_groups", *backButton.CallbackData)
	})

	t.Run("should create category group keyboard with empty categories", func(t *testing.T) {
		categories := []*models.Category{}

		keyboard := GetCategoryGroupKeyboard(categories)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 1) // Only back button

		// Check back button
		backButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "⬅️ Back to Groups", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_groups", *backButton.CallbackData)
	})
}

func TestGetConfirmationKeyboard(t *testing.T) {
	t.Run("should create confirmation keyboard", func(t *testing.T) {
		keyboard := GetConfirmationKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 1) // 1 row with 2 buttons

		// Check confirm button
		confirmButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "✅ Yes", confirmButton.Text)
		require.NotNil(t, confirmButton.CallbackData)
		require.Equal(t, "confirm_delete", *confirmButton.CallbackData)

		// Check cancel button
		cancelButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "❌ No", cancelButton.Text)
		require.NotNil(t, cancelButton.CallbackData)
		require.Equal(t, "confirm_no", *cancelButton.CallbackData)
	})
}

func TestGetReportKeyboard(t *testing.T) {
	t.Run("should create report keyboard", func(t *testing.T) {
		keyboard := GetReportKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 3 rows

		// Check first row
		monthlyButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "📅 Monthly", monthlyButton.Text)
		require.NotNil(t, monthlyButton.CallbackData)
		require.Equal(t, "report_monthly", *monthlyButton.CallbackData)

		yearlyButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "📊 Yearly", yearlyButton.Text)
		require.NotNil(t, yearlyButton.CallbackData)
		require.Equal(t, "report_yearly", *yearlyButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetBudgetKeyboard(t *testing.T) {
	t.Run("should create budget keyboard", func(t *testing.T) {
		keyboard := GetBudgetKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 3 rows

		// Check first row
		setBudgetButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "💰 Set Budget", setBudgetButton.Text)
		require.NotNil(t, setBudgetButton.CallbackData)
		require.Equal(t, "budget_set", *setBudgetButton.CallbackData)

		viewBudgetButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "📊 View Budget", viewBudgetButton.Text)
		require.NotNil(t, viewBudgetButton.CallbackData)
		require.Equal(t, "budget_view", *viewBudgetButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetReminderKeyboard(t *testing.T) {
	t.Run("should create reminder keyboard", func(t *testing.T) {
		keyboard := GetReminderKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 3 rows

		// Check first row
		setReminderButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "⏰ Set Reminder", setReminderButton.Text)
		require.NotNil(t, setReminderButton.CallbackData)
		require.Equal(t, "reminder_set", *setReminderButton.CallbackData)

		viewReminderButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "📋 View Reminders", viewReminderButton.Text)
		require.NotNil(t, viewReminderButton.CallbackData)
		require.Equal(t, "reminder_view", *viewReminderButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetSettingsKeyboard(t *testing.T) {
	t.Run("should create settings keyboard", func(t *testing.T) {
		keyboard := GetSettingsKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 3 rows

		// Check first row
		currencyButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "💱 Currency", currencyButton.Text)
		require.NotNil(t, currencyButton.CallbackData)
		require.Equal(t, "settings_currency", *currencyButton.CallbackData)

		dateFormatButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "📅 Date Format", dateFormatButton.Text)
		require.NotNil(t, dateFormatButton.CallbackData)
		require.Equal(t, "settings_date", *dateFormatButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetEditExpenseKeyboard(t *testing.T) {
	t.Run("should create edit expense keyboard with expenses", func(t *testing.T) {
		expenses := []*models.Expense{
			{ID: 1, CategoryName: "Fuel", TotalPrice: 100.0, Timestamp: time.Now()},
			{ID: 2, CategoryName: "Food", TotalPrice: 50.0, Timestamp: time.Now()},
		}

		keyboard := GetEditExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 2 expenses + 1 back button

		// Check first expense button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.Contains(t, firstButton.Text, "Fuel")
		require.Contains(t, firstButton.Text, "₹100.00")
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "edit_1", *firstButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back to Main Menu", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})

	t.Run("should create edit expense keyboard with more than 10 expenses", func(t *testing.T) {
		expenses := make([]*models.Expense, 15)
		for i := 0; i < 15; i++ {
			expenses[i] = &models.Expense{
				ID:           int64(i + 1),
				CategoryName: "Fuel",
				TotalPrice:   float64(i+1) * 10.0,
				Timestamp:    time.Now(),
			}
		}

		keyboard := GetEditExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 11) // 10 expenses + 1 back button

		// Check first expense button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "edit_1", *firstButton.CallbackData)
	})

	t.Run("should create edit expense keyboard with no expenses", func(t *testing.T) {
		expenses := []*models.Expense{}

		keyboard := GetEditExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 1) // Only back button

		// Check back button
		backButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "⬅️ Back to Main Menu", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetDeleteExpenseKeyboard(t *testing.T) {
	t.Run("should create delete expense keyboard with expenses", func(t *testing.T) {
		expenses := []*models.Expense{
			{ID: 1, CategoryName: "Fuel", TotalPrice: 100.0, Timestamp: time.Now()},
			{ID: 2, CategoryName: "Food", TotalPrice: 50.0, Timestamp: time.Now()},
		}

		keyboard := GetDeleteExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 3) // 2 expenses + 1 back button

		// Check first expense button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.Contains(t, firstButton.Text, "Fuel")
		require.Contains(t, firstButton.Text, "₹100.00")
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "delete_1", *firstButton.CallbackData)

		// Check back button
		backButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "⬅️ Back to Main Menu", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})

	t.Run("should create delete expense keyboard with more than 10 expenses", func(t *testing.T) {
		expenses := make([]*models.Expense, 15)
		for i := 0; i < 15; i++ {
			expenses[i] = &models.Expense{
				ID:           int64(i + 1),
				CategoryName: "Fuel",
				TotalPrice:   float64(i+1) * 10.0,
				Timestamp:    time.Now(),
			}
		}

		keyboard := GetDeleteExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 11) // 10 expenses + 1 back button

		// Check first expense button
		firstButton := keyboard.InlineKeyboard[0][0]
		require.NotNil(t, firstButton.CallbackData)
		require.Equal(t, "delete_1", *firstButton.CallbackData)
	})

	t.Run("should create delete expense keyboard with no expenses", func(t *testing.T) {
		expenses := []*models.Expense{}

		keyboard := GetDeleteExpenseKeyboard(expenses)
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 1) // Only back button

		// Check back button
		backButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "⬅️ Back to Main Menu", backButton.Text)
		require.NotNil(t, backButton.CallbackData)
		require.Equal(t, "back_to_main", *backButton.CallbackData)
	})
}

func TestGetEditFieldKeyboard(t *testing.T) {
	t.Run("should create edit field keyboard", func(t *testing.T) {
		keyboard := GetEditFieldKeyboard()
		require.NotNil(t, keyboard)
		require.Len(t, keyboard.InlineKeyboard, 4) // 4 rows

		// Check first row
		categoryButton := keyboard.InlineKeyboard[0][0]
		require.Equal(t, "🏷️ Category", categoryButton.Text)
		require.NotNil(t, categoryButton.CallbackData)
		require.Equal(t, "edit_field_category", *categoryButton.CallbackData)

		vehicleButton := keyboard.InlineKeyboard[0][1]
		require.Equal(t, "🚗 Vehicle Type", vehicleButton.Text)
		require.NotNil(t, vehicleButton.CallbackData)
		require.Equal(t, "edit_field_vehicle", *vehicleButton.CallbackData)

		// Check second row
		odometerButton := keyboard.InlineKeyboard[1][0]
		require.Equal(t, "🔢 Odometer", odometerButton.Text)
		require.NotNil(t, odometerButton.CallbackData)
		require.Equal(t, "edit_field_odometer", *odometerButton.CallbackData)

		petrolButton := keyboard.InlineKeyboard[1][1]
		require.Equal(t, "⛽ Petrol Price", petrolButton.Text)
		require.NotNil(t, petrolButton.CallbackData)
		require.Equal(t, "edit_field_petrol", *petrolButton.CallbackData)

		// Check third row
		totalButton := keyboard.InlineKeyboard[2][0]
		require.Equal(t, "💰 Total Price", totalButton.Text)
		require.NotNil(t, totalButton.CallbackData)
		require.Equal(t, "edit_field_total", *totalButton.CallbackData)

		notesButton := keyboard.InlineKeyboard[2][1]
		require.Equal(t, "📝 Notes", notesButton.Text)
		require.NotNil(t, notesButton.CallbackData)
		require.Equal(t, "edit_field_notes", *notesButton.CallbackData)

		// Check fourth row
		saveButton := keyboard.InlineKeyboard[3][0]
		require.Equal(t, "✅ Save Changes", saveButton.Text)
		require.NotNil(t, saveButton.CallbackData)
		require.Equal(t, "edit_save", *saveButton.CallbackData)

		cancelButton := keyboard.InlineKeyboard[3][1]
		require.Equal(t, "❌ Cancel", cancelButton.Text)
		require.NotNil(t, cancelButton.CallbackData)
		require.Equal(t, "edit_cancel", *cancelButton.CallbackData)
	})
}
