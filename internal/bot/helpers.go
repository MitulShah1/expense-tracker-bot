package bot

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/MitulShah1/expense-tracker-bot/pkg/utils"
)

// prepareExpenseSelection fetches expenses for a user, limits them to maxExpenses,
// creates a deep copy, and stores them in user state. Returns the expenses to show.
func (b *Bot) prepareExpenseSelection(ctx context.Context, userID int64, maxExpenses int) ([]*models.Expense, error) {
	// Get expenses from service by Telegram ID
	expenses, err := b.expenseService.GetExpensesByTelegramID(ctx, userID, 100, 0) // Get all expenses
	if err != nil {
		return nil, err
	}

	if len(expenses) == 0 {
		return nil, nil
	}

	// Limit expenses to show
	expensesToShow := expenses
	if len(expenses) > maxExpenses {
		expensesToShow = expenses[:maxExpenses]
	}

	// Create a deep copy to ensure independence
	expensesCopy := make([]*models.Expense, len(expensesToShow))
	copy(expensesCopy, expensesToShow)

	// Store in user state
	state := b.getState(userID)
	if state == nil {
		state = models.NewUserState()
		b.setState(userID, state)
	}
	state.ExpenseSelection = expensesCopy

	// Debug logging
	b.logger.Info(ctx, "Expense selection prepared",
		logger.Int("total_expenses", len(expenses)),
		logger.Int("expenses_to_show", len(expensesToShow)),
		logger.Int("stored_in_state", len(state.ExpenseSelection)))

	return expensesToShow, nil
}

// parseFloatOrReply attempts to parse a string as float64 and sends an error message if it fails
func (b *Bot) parseFloatOrReply(ctx context.Context, chatID int64, text, fieldName string) (float64, error) {
	value, err := strconv.ParseFloat(text, 64)
	if err != nil {
		errorMsg := fmt.Sprintf("Please enter a valid number for the %s.", fieldName)
		return 0, b.sendMessage(ctx, chatID, errorMsg)
	}
	return value, nil
}

// checkExpenseOwnership verifies that the user owns the expense and returns an error if not
func (b *Bot) checkExpenseOwnership(ctx context.Context, userID int64, expense *models.Expense) error {
	if expense == nil {
		return errors.New("expense not found")
	}

	// Get user by Telegram ID
	user, err := b.userService.GetUserByTelegramID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("could not find your user record")
	}

	if expense.UserID != user.ID {
		return errors.New("you can only modify your own expenses")
	}

	return nil
}

// buildExpenseListMessage builds a formatted message showing expenses grouped by category
func (b *Bot) buildExpenseListMessage(expenses []*models.Expense) string {
	if len(expenses) == 0 {
		return "No expenses found."
	}

	// Group expenses by category
	categoryExpenses := make(map[string][]*models.Expense)
	for _, expense := range expenses {
		categoryExpenses[expense.CategoryName] = append(categoryExpenses[expense.CategoryName], expense)
	}

	// Sort categories for consistent ordering
	var categories []string
	for category := range categoryExpenses {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	// Build message
	var sb strings.Builder
	sb.WriteString("Your expenses:\n\n")

	for _, category := range categories {
		categoryExpenses := categoryExpenses[category]
		sb.WriteString(fmt.Sprintf("📊 %s:\n", category))
		var prices []float64
		for _, expense := range categoryExpenses {
			sb.WriteString(fmt.Sprintf("• %s: %s\n", utils.FormatDate(expense.Timestamp), utils.FormatCurrency(expense.TotalPrice)))
			prices = append(prices, expense.TotalPrice)
		}
		total := utils.CalculateTotal(prices)
		sb.WriteString(fmt.Sprintf("Total: %s\n\n", utils.FormatCurrency(total)))
	}

	return sb.String()
}

// buildReportMessage builds a formatted expense report message
func (b *Bot) buildReportMessage(expenses []*models.Expense) string {
	if len(expenses) == 0 {
		return "No expenses found to generate report."
	}

	// Calculate report metrics
	var totalExpense float64
	categoryTotals := make(map[string]float64)
	monthlyTotals := make(map[string]float64)

	for _, expense := range expenses {
		totalExpense += expense.TotalPrice
		categoryTotals[expense.CategoryName] += expense.TotalPrice

		// Group by month
		monthKey := expense.Timestamp.Format("January 2006")
		monthlyTotals[monthKey] += expense.TotalPrice
	}

	// Build report message
	var sb strings.Builder
	sb.WriteString("📊 Expense Report\n\n")

	// Overall total
	sb.WriteString(fmt.Sprintf("💰 Total Expenses: %s\n\n", utils.FormatCurrency(totalExpense)))

	// Category breakdown
	sb.WriteString("📈 Category Breakdown:\n")
	for category, total := range categoryTotals {
		percentage := (total / totalExpense) * 100
		sb.WriteString(fmt.Sprintf("• %s: %s (%.1f%%)\n",
			category,
			utils.FormatCurrency(total),
			percentage))
	}
	sb.WriteString("\n")

	// Monthly breakdown
	sb.WriteString("📅 Monthly Breakdown:\n")
	for month, total := range monthlyTotals {
		sb.WriteString(fmt.Sprintf("• %s: %s\n", month, utils.FormatCurrency(total)))
	}

	return sb.String()
}

// buildDashboardMessage builds a formatted dashboard message
func (b *Bot) buildDashboardMessage(expenses []*models.Expense) string {
	if len(expenses) == 0 {
		return "No expenses found to show dashboard."
	}

	// Sort expenses by timestamp (most recent first)
	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Timestamp.After(expenses[j].Timestamp)
	})

	// Get last 5 expenses
	recentExpenses := expenses
	if len(expenses) > 5 {
		recentExpenses = expenses[:5]
	}

	// Calculate dashboard metrics
	var totalExpense float64
	var totalFuelExpense float64
	var totalDistance float64
	var avgFuelEfficiency float64

	// Calculate metrics
	for _, expense := range expenses {
		totalExpense += expense.TotalPrice
		if expense.CategoryName == "⛽ Petrol" {
			totalFuelExpense += expense.TotalPrice
			if expense.Odometer > 0 {
				totalDistance += expense.Odometer
			}
		}
	}

	// Calculate average fuel efficiency if possible
	if totalFuelExpense > 0 && totalDistance > 0 {
		avgFuelEfficiency = totalDistance / (totalFuelExpense / 100) // km per 100 rupees
	}

	// Build dashboard message
	var sb strings.Builder
	sb.WriteString("📱 Expense Dashboard\n\n")

	// Overall metrics
	sb.WriteString("📊 Overall Metrics:\n")
	sb.WriteString(fmt.Sprintf("• Total Expenses: %s\n", utils.FormatCurrency(totalExpense)))
	sb.WriteString(fmt.Sprintf("• Total Fuel Expenses: %s\n", utils.FormatCurrency(totalFuelExpense)))
	if avgFuelEfficiency > 0 {
		sb.WriteString(fmt.Sprintf("• Average Fuel Efficiency: %.1f km/₹100\n", avgFuelEfficiency))
	}
	sb.WriteString("\n")

	// Recent expenses
	sb.WriteString("🕒 Recent Expenses:\n")
	for _, expense := range recentExpenses {
		sb.WriteString(fmt.Sprintf("• %s - %s: %s\n",
			utils.FormatDate(expense.Timestamp),
			expense.CategoryName,
			utils.FormatCurrency(expense.TotalPrice)))
	}

	return sb.String()
}
