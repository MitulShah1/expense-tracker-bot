package bot

import (
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/stretchr/testify/assert"
)

// parseTestDate is a helper function to parse dates for testing
func parseTestDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

// createTestBot creates a minimal bot instance for testing helper functions
func createTestBot() *Bot {
	return &Bot{
		logger: logger.NewMockLogger(),
		states: make(map[int64]*models.UserState),
	}
}

func TestBuildExpenseListMessage(t *testing.T) {
	tests := []struct {
		name           string
		expenses       []*models.Expense
		expectedResult string
	}{
		{
			name:           "empty_expenses",
			expenses:       []*models.Expense{},
			expectedResult: "No expenses found.",
		},
		{
			name: "single_expense",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
			},
			expectedResult: "Your expenses:\n\n📊 ⛽ Petrol:\n• 01 Jan 2024: ₹100.00\nTotal: ₹100.00\n\n",
		},
		{
			name: "multiple_expenses_same_category",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   200.0,
					Timestamp:    parseTestDate("2024-01-02"),
				},
			},
			expectedResult: "Your expenses:\n\n📊 ⛽ Petrol:\n• 01 Jan 2024: ₹100.00\n• 02 Jan 2024: ₹200.00\nTotal: ₹300.00\n\n",
		},
		{
			name: "multiple_categories",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
				{
					CategoryName: "🔧 Service",
					TotalPrice:   200.0,
					Timestamp:    parseTestDate("2024-01-02"),
				},
			},
			expectedResult: "Your expenses:\n\n📊 ⛽ Petrol:\n• 01 Jan 2024: ₹100.00\nTotal: ₹100.00\n\n📊 🔧 Service:\n• 02 Jan 2024: ₹200.00\nTotal: ₹200.00\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			bot := createTestBot()

			// Execute
			result := bot.buildExpenseListMessage(tt.expenses)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestBuildReportMessage(t *testing.T) {
	tests := []struct {
		name           string
		expenses       []*models.Expense
		expectedResult string
	}{
		{
			name:           "empty_expenses",
			expenses:       []*models.Expense{},
			expectedResult: "No expenses found to generate report.",
		},
		{
			name: "single_expense",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
			},
			expectedResult: "📊 Expense Report\n\n💰 Total Expenses: ₹100.00\n\n📈 Category Breakdown:\n• ⛽ Petrol: ₹100.00 (100.0%)\n\n📅 Monthly Breakdown:\n• January 2024: ₹100.00\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			bot := createTestBot()

			// Execute
			result := bot.buildReportMessage(tt.expenses)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}

	// Test multiple expenses with order-independent assertions
	t.Run("multiple_expenses_different_months", func(t *testing.T) {
		bot := createTestBot()
		expenses := []*models.Expense{
			{
				CategoryName: "⛽ Petrol",
				TotalPrice:   100.0,
				Timestamp:    parseTestDate("2024-01-01"),
			},
			{
				CategoryName: "🔧 Service",
				TotalPrice:   200.0,
				Timestamp:    parseTestDate("2024-02-01"),
			},
		}

		result := bot.buildReportMessage(expenses)

		// Check overall structure
		assert.Contains(t, result, "📊 Expense Report")
		assert.Contains(t, result, "💰 Total Expenses: ₹300.00")
		assert.Contains(t, result, "📈 Category Breakdown:")
		assert.Contains(t, result, "📅 Monthly Breakdown:")

		// Check category breakdown (order independent)
		assert.Contains(t, result, "⛽ Petrol: ₹100.00 (33.3%)")
		assert.Contains(t, result, "🔧 Service: ₹200.00 (66.7%)")

		// Check monthly breakdown (order independent)
		assert.Contains(t, result, "January 2024: ₹100.00")
		assert.Contains(t, result, "February 2024: ₹200.00")
	})

	t.Run("multiple_expenses_same_month", func(t *testing.T) {
		bot := createTestBot()
		expenses := []*models.Expense{
			{
				CategoryName: "⛽ Petrol",
				TotalPrice:   100.0,
				Timestamp:    parseTestDate("2024-01-01"),
			},
			{
				CategoryName: "🔧 Service",
				TotalPrice:   200.0,
				Timestamp:    parseTestDate("2024-01-15"),
			},
		}

		result := bot.buildReportMessage(expenses)

		// Check overall structure
		assert.Contains(t, result, "📊 Expense Report")
		assert.Contains(t, result, "💰 Total Expenses: ₹300.00")
		assert.Contains(t, result, "📈 Category Breakdown:")
		assert.Contains(t, result, "📅 Monthly Breakdown:")

		// Check category breakdown (order independent)
		assert.Contains(t, result, "⛽ Petrol: ₹100.00 (33.3%)")
		assert.Contains(t, result, "🔧 Service: ₹200.00 (66.7%)")

		// Check monthly breakdown
		assert.Contains(t, result, "January 2024: ₹300.00")
	})
}

func TestBuildDashboardMessage(t *testing.T) {
	tests := []struct {
		name           string
		expenses       []*models.Expense
		expectedResult string
	}{
		{
			name:           "empty_expenses",
			expenses:       []*models.Expense{},
			expectedResult: "No expenses found to show dashboard.",
		},
		{
			name: "fuel_expenses_with_efficiency",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Odometer:     50.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
			},
			expectedResult: "📱 Expense Dashboard\n\n📊 Overall Metrics:\n• Total Expenses: ₹100.00\n• Total Fuel Expenses: ₹100.00\n• Average Fuel Efficiency: 50.0 km/₹100\n\n🕒 Recent Expenses:\n• 01 Jan 2024 - ⛽ Petrol: ₹100.00\n",
		},
		{
			name: "mixed_expenses_no_fuel_efficiency",
			expenses: []*models.Expense{
				{
					CategoryName: "🔧 Service",
					TotalPrice:   200.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Odometer:     0, // No odometer reading
					Timestamp:    parseTestDate("2024-01-02"),
				},
			},
			expectedResult: "📱 Expense Dashboard\n\n📊 Overall Metrics:\n• Total Expenses: ₹300.00\n• Total Fuel Expenses: ₹100.00\n\n🕒 Recent Expenses:\n• 02 Jan 2024 - ⛽ Petrol: ₹100.00\n• 01 Jan 2024 - 🔧 Service: ₹200.00\n",
		},
		{
			name: "multiple_fuel_expenses_with_efficiency",
			expenses: []*models.Expense{
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   100.0,
					Odometer:     50.0,
					Timestamp:    parseTestDate("2024-01-01"),
				},
				{
					CategoryName: "⛽ Petrol",
					TotalPrice:   200.0,
					Odometer:     100.0,
					Timestamp:    parseTestDate("2024-01-02"),
				},
			},
			expectedResult: "📱 Expense Dashboard\n\n📊 Overall Metrics:\n• Total Expenses: ₹300.00\n• Total Fuel Expenses: ₹300.00\n• Average Fuel Efficiency: 50.0 km/₹100\n\n🕒 Recent Expenses:\n• 02 Jan 2024 - ⛽ Petrol: ₹200.00\n• 01 Jan 2024 - ⛽ Petrol: ₹100.00\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			bot := createTestBot()

			// Execute
			result := bot.buildDashboardMessage(tt.expenses)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// Test utility functions
func TestParseTestDate(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		expected time.Time
	}{
		{
			name:     "valid_date",
			dateStr:  "2024-01-01",
			expected: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "valid_date_february",
			dateStr:  "2024-02-15",
			expected: time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTestDate(tt.dateStr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test edge cases for message building functions
func TestBuildExpenseListMessageEdgeCases(t *testing.T) {
	t.Run("nil_expenses", func(t *testing.T) {
		bot := createTestBot()
		result := bot.buildExpenseListMessage(nil)
		assert.Equal(t, "No expenses found.", result)
	})

	t.Run("expense_with_nil_timestamp", func(t *testing.T) {
		bot := createTestBot()
		expenses := []*models.Expense{
			{
				CategoryName: "⛽ Petrol",
				TotalPrice:   100.0,
				Timestamp:    time.Time{}, // Zero time
			},
		}
		result := bot.buildExpenseListMessage(expenses)
		assert.Contains(t, result, "⛽ Petrol")
		assert.Contains(t, result, "₹100.00")
	})
}

func TestBuildReportMessageEdgeCases(t *testing.T) {
	t.Run("nil_expenses", func(t *testing.T) {
		bot := createTestBot()
		result := bot.buildReportMessage(nil)
		assert.Equal(t, "No expenses found to generate report.", result)
	})

	t.Run("expense_with_zero_price", func(t *testing.T) {
		bot := createTestBot()
		expenses := []*models.Expense{
			{
				CategoryName: "⛽ Petrol",
				TotalPrice:   0.0,
				Timestamp:    parseTestDate("2024-01-01"),
			},
		}
		result := bot.buildReportMessage(expenses)
		assert.Contains(t, result, "Total Expenses: ₹0.00")
		assert.Contains(t, result, "⛽ Petrol: ₹0.00 (NaN%)")
	})
}

func TestBuildDashboardMessageEdgeCases(t *testing.T) {
	t.Run("nil_expenses", func(t *testing.T) {
		bot := createTestBot()
		result := bot.buildDashboardMessage(nil)
		assert.Equal(t, "No expenses found to show dashboard.", result)
	})

	t.Run("expense_with_negative_odometer", func(t *testing.T) {
		bot := createTestBot()
		expenses := []*models.Expense{
			{
				CategoryName: "⛽ Petrol",
				TotalPrice:   100.0,
				Odometer:     -10.0, // Negative odometer
				Timestamp:    parseTestDate("2024-01-01"),
			},
		}
		result := bot.buildDashboardMessage(expenses)
		assert.Contains(t, result, "Total Expenses: ₹100.00")
		assert.Contains(t, result, "Total Fuel Expenses: ₹100.00")
		// Should not show fuel efficiency with negative odometer
		assert.NotContains(t, result, "Average Fuel Efficiency")
	})
}
