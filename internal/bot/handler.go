package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

// handleListCommand handles the /list command
func (b *Bot) handleListCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Get expenses from service by Telegram ID
	expenses, err := b.expenseService.GetExpensesByTelegramID(ctx, message.From.ID, 100, 0) // Get all expenses
	if err != nil {
		return b.sendError(ctx, message.Chat.ID, err)
	}

	// Build and send message using helper
	messageText := b.buildExpenseListMessage(expenses)
	return b.sendMessage(ctx, message.Chat.ID, messageText)
}

// handleEditCommand handles the /edit command
func (b *Bot) handleEditCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Use helper to prepare expense selection
	expenses, err := b.prepareExpenseSelection(ctx, message.From.ID, 10)
	if err != nil {
		return b.sendError(ctx, message.Chat.ID, err)
	}

	if len(expenses) == 0 {
		return b.sendMessage(ctx, message.Chat.ID, "No expenses found to edit.")
	}

	// Send message with inline keyboard for expense selection
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select an expense to edit:")
	msg.ReplyMarkup = GetEditExpenseKeyboard(expenses)
	_, err = b.api.Send(msg)
	return err
}

// handleDeleteCommand handles the /delete command
func (b *Bot) handleDeleteCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Use helper to prepare expense selection
	expenses, err := b.prepareExpenseSelection(ctx, message.From.ID, 10)
	if err != nil {
		return b.sendError(ctx, message.Chat.ID, err)
	}

	if len(expenses) == 0 {
		return b.sendMessage(ctx, message.Chat.ID, "No expenses found to delete.")
	}

	// Send message with inline keyboard for expense selection
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select an expense to delete:")
	msg.ReplyMarkup = GetDeleteExpenseKeyboard(expenses)
	_, err = b.api.Send(msg)
	return err
}

// handleReportCommand handles the /report command
func (b *Bot) handleReportCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Get expenses from service by Telegram ID
	expenses, err := b.expenseService.GetExpensesByTelegramID(ctx, message.From.ID, 100, 0) // Get all expenses
	if err != nil {
		b.logger.Error(ctx, "Failed to get expenses for report", zap.Error(err))
		b.incrementMetric(&b.metrics.errorCount)
		return b.sendError(ctx, message.Chat.ID, err)
	}

	// Build and send message using helper
	messageText := b.buildReportMessage(expenses)
	return b.sendMessage(ctx, message.Chat.ID, messageText)
}

// handleDashboardCommand handles the /dashboard command
func (b *Bot) handleDashboardCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Get recent expenses from service
	expenses, err := b.expenseService.GetExpensesByTelegramID(ctx, message.From.ID, 10, 0) // Get recent 10 expenses
	if err != nil {
		b.incrementMetric(&b.metrics.errorCount)
		return b.sendError(ctx, message.Chat.ID, err)
	}

	// Build and send message using helper
	messageText := b.buildDashboardMessage(expenses)
	return b.sendMessage(ctx, message.Chat.ID, messageText)
}

// handleAddCommand handles the /add command
func (b *Bot) handleAddCommand(ctx context.Context, message *tgbotapi.Message) error {
	// Send category group selection keyboard
	msg := tgbotapi.NewMessage(message.Chat.ID, "Select a category group:")
	msg.ReplyMarkup = GetCategoryKeyboard()
	_, err := b.api.Send(msg)
	return err
}
