// Package utils provides utility functions used throughout the expense tracker application.
package utils

import (
	"fmt"
	"strconv"
	"time"
)

// FormatCurrency formats a float64 as currency
func FormatCurrency(amount float64) string {
	return fmt.Sprintf("₹%.2f", amount)
}

// FormatDate formats a time.Time as a readable date
func FormatDate(t time.Time) string {
	return t.Format("02 Jan 2006")
}

// FormatDateTime formats a time.Time as a readable date and time
func FormatDateTime(t time.Time) string {
	return t.Format("02 Jan 2006 15:04:05")
}

// FormatMonth formats a time.Time as YYYY-MM
func FormatMonth(t time.Time) string {
	return t.Format("2006-01")
}

// CalculateAverage calculates the average of a slice of float64
func CalculateAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum / float64(len(numbers))
}

// CalculateTotal calculates the total of a slice of float64
func CalculateTotal(numbers []float64) float64 {
	var sum float64
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// GetMonthRange returns the start and end dates of a month
func GetMonthRange(year int, month time.Month) (time.Time, time.Time) {
	start := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0).Add(-time.Second)
	return start, end
}

// GetYearRange returns the start and end dates of a year
func GetYearRange(year int) (time.Time, time.Time) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0).Add(-time.Second)
	return start, end
}

// GetCurrentMonthRange returns the start and end dates of the current month
func GetCurrentMonthRange() (time.Time, time.Time) {
	now := time.Now()
	return GetMonthRange(now.Year(), now.Month())
}

// GetCurrentYearRange returns the start and end dates of the current year
func GetCurrentYearRange() (time.Time, time.Time) {
	return GetYearRange(time.Now().Year())
}

// IsValidAmount checks if a string can be parsed as a valid amount
func IsValidAmount(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// FormatAmount formats an amount string to ensure it's a valid number
func FormatAmount(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
