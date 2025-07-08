// Package bot provides the Telegram bot logic.
package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/MitulShah1/expense-tracker-bot/internal/models"
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

// handleSearchCommand handles the /search command for semantic search
func (b *Bot) handleSearchCommand(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	userID := message.From.ID

	// Check if user has any expenses
	expenses, err := b.expenseService.GetExpensesByTelegramID(ctx, userID, 1, 0)
	if err != nil {
		return b.sendError(ctx, chatID, err)
	}

	if len(expenses) == 0 {
		return b.sendMessage(ctx, chatID, "You don't have any expenses yet. Add some expenses first to use semantic search.")
	}

	// Set user state to search mode
	state := b.getState(userID)
	if state == nil {
		state = models.NewUserState()
		b.setState(userID, state)
	}
	state.Step = models.StepSearchExpense

	// Send search instructions
	searchInstructions := `🔍 Semantic Search\n\nYou can search for expenses using natural language. Examples:\n• Find all fuel expenses from last month\n• Show me expensive car repairs\n• Find expenses related to maintenance\n\nType your search query:`

	return b.sendMessage(ctx, chatID, searchInstructions)
}

// handleSearchQuery handles the search query input
func (b *Bot) handleSearchQuery(ctx context.Context, message *tgbotapi.Message) error {
	chatID := message.Chat.ID
	userID := message.From.ID
	query := strings.TrimSpace(message.Text)

	if query == "" {
		return b.sendMessage(ctx, chatID, "Please enter a search query.")
	}

	// Perform semantic search with lower threshold for placeholder embeddings
	expenses, err := b.vectorService.SearchExpensesByQuery(ctx, userID, query, 0.1, 10)
	if err != nil {
		return b.sendError(ctx, chatID, err)
	}

	if len(expenses) == 0 {
		return b.sendMessage(ctx, chatID, fmt.Sprintf("No expenses found matching: %q\n\nTry a different search term or be more specific.", query))
	}

	// Build search results message
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Search Results for: %q\n\n", query))
	sb.WriteString(fmt.Sprintf("Found %d matching expenses:\n\n", len(expenses)))

	for i, expense := range expenses {
		sb.WriteString(fmt.Sprintf("%d. %s - %s: ₹%.2f\n",
			i+1,
			expense.Timestamp.Format("02 Jan 2006"),
			expense.CategoryName,
			expense.TotalPrice))
		if expense.Notes != "" {
			sb.WriteString(fmt.Sprintf("   Notes: %s\n", expense.Notes))
		}
		sb.WriteString("\n")
	}

	// Reset user state
	state := b.getState(userID)
	if state != nil {
		state.Step = models.StepNone
	}

	return b.sendMessage(ctx, chatID, sb.String())
}
