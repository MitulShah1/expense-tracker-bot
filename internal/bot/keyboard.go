// Package bot implements the Telegram bot functionality for expense tracking.
package bot

import (
	"fmt"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// GetMainMenuKeyboard returns the main menu keyboard
func GetMainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📝 Add Expense"),
			tgbotapi.NewKeyboardButton("📋 List Expenses"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("✏️ Edit Expense"),
			tgbotapi.NewKeyboardButton("🗑️ Delete Expense"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📊 Reports"),
			tgbotapi.NewKeyboardButton("📈 Dashboard"),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

// GetCategoryKeyboard returns the category group keyboard
func GetCategoryKeyboard() tgbotapi.InlineKeyboardMarkup {
	// This will be replaced with a dynamic version that fetches from DB
	// For now, return a static keyboard with the main groups
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚗 Vehicle", "group_Vehicle"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 Home", "group_Home"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏪 Daily Living", "group_Daily Living"),
			tgbotapi.NewInlineKeyboardButtonData("🎬 Entertainment", "group_Entertainment"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏥 Health", "group_Health"),
			tgbotapi.NewInlineKeyboardButtonData("📚 Education", "group_Education"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✈️ Travel", "group_Travel"),
			tgbotapi.NewInlineKeyboardButtonData("💰 Investments", "group_Investments"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎁 Gifts", "group_Gifts"),
			tgbotapi.NewInlineKeyboardButtonData("📌 Other", "group_Other"),
		),
	)
}

// GetCategoryGroupKeyboard returns the category group keyboard
func GetCategoryGroupKeyboard(categories []*models.Category) tgbotapi.InlineKeyboardMarkup {
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(categories)+1)

	for _, cat := range categories {
		buttonText := cat.Emoji + " " + cat.Name
		callbackData := "category_" + cat.Name
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		})
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Groups", "back_to_groups"),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// GetConfirmationKeyboard returns the confirmation keyboard
func GetConfirmationKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("✅ Yes", "confirm_delete"),
			tgbotapi.NewInlineKeyboardButtonData("❌ No", "confirm_no"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// GetReportKeyboard returns the report selection keyboard
func GetReportKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("📅 Monthly", "report_monthly"),
			tgbotapi.NewInlineKeyboardButtonData("📊 Yearly", "report_yearly"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("🏷️ By Category", "report_category"),
			tgbotapi.NewInlineKeyboardButtonData("🚗 By Vehicle", "report_vehicle"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "back_to_main"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// GetBudgetKeyboard returns the budget management keyboard
func GetBudgetKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("💰 Set Budget", "budget_set"),
			tgbotapi.NewInlineKeyboardButtonData("📊 View Budget", "budget_view"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("📈 Budget History", "budget_history"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Budget Settings", "budget_settings"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "back_to_main"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// GetReminderKeyboard returns the reminder management keyboard
func GetReminderKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("⏰ Set Reminder", "reminder_set"),
			tgbotapi.NewInlineKeyboardButtonData("📋 View Reminders", "reminder_view"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("🔔 Edit Reminder", "reminder_edit"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Delete Reminder", "reminder_delete"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "back_to_main"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// GetSettingsKeyboard returns the settings keyboard
func GetSettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("💱 Currency", "settings_currency"),
			tgbotapi.NewInlineKeyboardButtonData("📅 Date Format", "settings_date"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("🔔 Notifications", "settings_notifications"),
			tgbotapi.NewInlineKeyboardButtonData("🌍 Language", "settings_language"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("⬅️ Back", "back_to_main"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

// GetEditExpenseKeyboard returns the edit expense selection keyboard
func GetEditExpenseKeyboard(expenses []*models.Expense) tgbotapi.InlineKeyboardMarkup {
	maxExpenses := 10
	expensesToShow := expenses
	if len(expenses) > maxExpenses {
		expensesToShow = expenses[:maxExpenses]
	}
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(expensesToShow)+1)

	for _, expense := range expensesToShow {
		buttonText := fmt.Sprintf("%s - %s: ₹%.2f",
			expense.Timestamp.Format("02 Jan"),
			expense.CategoryName,
			expense.TotalPrice)
		callbackData := fmt.Sprintf("edit_%d", expense.ID)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		})
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Main Menu", "back_to_main"),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// GetDeleteExpenseKeyboard returns the delete expense selection keyboard
func GetDeleteExpenseKeyboard(expenses []*models.Expense) tgbotapi.InlineKeyboardMarkup {
	maxExpenses := 10
	expensesToShow := expenses
	if len(expenses) > maxExpenses {
		expensesToShow = expenses[:maxExpenses]
	}
	buttons := make([][]tgbotapi.InlineKeyboardButton, 0, len(expensesToShow)+1)

	for _, expense := range expensesToShow {
		buttonText := fmt.Sprintf("%s - %s: ₹%.2f",
			expense.Timestamp.Format("02 Jan"),
			expense.CategoryName,
			expense.TotalPrice)
		callbackData := fmt.Sprintf("delete_%d", expense.ID)
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(buttonText, callbackData),
		})
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("⬅️ Back to Main Menu", "back_to_main"),
	})
	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

// GetEditFieldKeyboard returns the edit field selection keyboard
func GetEditFieldKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("🏷️ Category", "edit_field_category"),
			tgbotapi.NewInlineKeyboardButtonData("🚗 Vehicle Type", "edit_field_vehicle"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("🔢 Odometer", "edit_field_odometer"),
			tgbotapi.NewInlineKeyboardButtonData("⛽ Petrol Price", "edit_field_petrol"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("💰 Total Price", "edit_field_total"),
			tgbotapi.NewInlineKeyboardButtonData("📝 Notes", "edit_field_notes"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("✅ Save Changes", "edit_save"),
			tgbotapi.NewInlineKeyboardButtonData("❌ Cancel", "edit_cancel"),
		},
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
