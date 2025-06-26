// Package validation provides input validation utilities for the expense tracker bot.
package validation

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
)

// Validator provides validation methods
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateTelegramID validates a Telegram user ID
func (v *Validator) ValidateTelegramID(telegramID int64) error {
	if telegramID <= 0 {
		return errors.NewValidationError("Invalid Telegram ID", "Telegram ID must be positive")
	}
	return nil
}

// ValidateUsername validates a username
func (v *Validator) ValidateUsername(username string) error {
	if username == "" {
		return errors.NewValidationError("Username is required", "Username cannot be empty")
	}

	if len(username) > 32 {
		return errors.NewValidationError("Username too long", "Username must be 32 characters or less")
	}

	// Telegram username pattern: 5-32 characters, alphanumeric and underscore
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]{5,32}$`)
	if !usernameRegex.MatchString(username) {
		return errors.NewValidationError("Invalid username format", "Username must be 5-32 characters, alphanumeric and underscore only")
	}

	return nil
}

// ValidateName validates a person's name
func (v *Validator) ValidateName(name, fieldName string) error {
	if name == "" {
		return errors.NewValidationError(fieldName+" is required", fieldName+" cannot be empty")
	}

	if len(name) > 64 {
		return errors.NewValidationError(fieldName+" too long", fieldName+" must be 64 characters or less")
	}

	// Allow letters, spaces, hyphens, and apostrophes
	nameRegex := regexp.MustCompile(`^[a-zA-Z\s\-']+$`)
	if !nameRegex.MatchString(name) {
		return errors.NewValidationError("Invalid "+fieldName+" format", fieldName+" can only contain letters, spaces, hyphens, and apostrophes")
	}

	return nil
}

// ValidateAmount validates a monetary amount
func (v *Validator) ValidateAmount(amount float64, fieldName string) error {
	if amount <= 0 {
		return errors.NewValidationError("Invalid "+fieldName, fieldName+" must be greater than 0")
	}

	if amount > 999999999.99 {
		return errors.NewValidationError(fieldName+" too large", fieldName+" cannot exceed 999,999,999.99")
	}

	return nil
}

// ValidateAmountString validates a monetary amount as a string
func (v *Validator) ValidateAmountString(amountStr, fieldName string) (float64, error) {
	if amountStr == "" {
		return 0, errors.NewValidationError(fieldName+" is required", fieldName+" cannot be empty")
	}

	// Remove any currency symbols and whitespace
	cleanAmount := strings.TrimSpace(strings.ReplaceAll(amountStr, "₹", ""))

	amount, err := strconv.ParseFloat(cleanAmount, 64)
	if err != nil {
		return 0, errors.NewValidationError("Invalid "+fieldName+" format", fieldName+" must be a valid number")
	}

	if err := v.ValidateAmount(amount, fieldName); err != nil {
		return 0, err
	}

	return amount, nil
}

// ValidateOdometer validates an odometer reading
func (v *Validator) ValidateOdometer(odometer float64) error {
	if odometer < 0 {
		return errors.NewValidationError("Invalid odometer reading", "Odometer reading cannot be negative")
	}

	if odometer > 999999.9 {
		return errors.NewValidationError("Odometer reading too large", "Odometer reading cannot exceed 999,999.9 km")
	}

	return nil
}

// ValidateOdometerString validates an odometer reading as a string
func (v *Validator) ValidateOdometerString(odometerStr string) (float64, error) {
	if odometerStr == "" {
		return 0, errors.NewValidationError("Odometer reading is required", "Odometer reading cannot be empty")
	}

	odometer, err := strconv.ParseFloat(odometerStr, 64)
	if err != nil {
		return 0, errors.NewValidationError("Invalid odometer format", "Odometer reading must be a valid number")
	}

	if err := v.ValidateOdometer(odometer); err != nil {
		return 0, err
	}

	return odometer, nil
}

// ValidateNotes validates expense notes
func (v *Validator) ValidateNotes(notes string) error {
	if len(notes) > 500 {
		return errors.NewValidationError("Notes too long", "Notes must be 500 characters or less")
	}

	return nil
}

// ValidateCategoryName validates a category name
func (v *Validator) ValidateCategoryName(categoryName string) error {
	if categoryName == "" {
		return errors.NewValidationError("Category name is required", "Category name cannot be empty")
	}

	if len(categoryName) > 100 {
		return errors.NewValidationError("Category name too long", "Category name must be 100 characters or less")
	}

	return nil
}

// ValidateVehicleType validates a vehicle type
func (v *Validator) ValidateVehicleType(vehicleType string) error {
	validTypes := []string{"CAR", "BIKE"}

	for _, validType := range validTypes {
		if vehicleType == validType {
			return nil
		}
	}

	return errors.NewValidationError("Invalid vehicle type", "Vehicle type must be one of: "+strings.Join(validTypes, ", "))
}

// ValidateDate validates a date
func (v *Validator) ValidateDate(date time.Time, fieldName string) error {
	if date.IsZero() {
		return errors.NewValidationError("Invalid "+fieldName, fieldName+" cannot be zero")
	}

	// Check if date is not in the future
	if date.After(time.Now()) {
		return errors.NewValidationError("Invalid "+fieldName, fieldName+" cannot be in the future")
	}

	// Check if date is not too far in the past (e.g., more than 10 years)
	tenYearsAgo := time.Now().AddDate(-10, 0, 0)
	if date.Before(tenYearsAgo) {
		return errors.NewValidationError("Invalid "+fieldName, fieldName+" cannot be more than 10 years ago")
	}

	return nil
}

// ValidateExpenseID validates an expense ID
func (v *Validator) ValidateExpenseID(expenseID int64) error {
	if expenseID <= 0 {
		return errors.NewValidationError("Invalid expense ID", "Expense ID must be positive")
	}

	return nil
}

// ValidateUserID validates a user ID
func (v *Validator) ValidateUserID(userID int64) error {
	if userID <= 0 {
		return errors.NewValidationError("Invalid user ID", "User ID must be positive")
	}

	return nil
}

// ValidateCategoryID validates a category ID
func (v *Validator) ValidateCategoryID(categoryID int64) error {
	if categoryID <= 0 {
		return errors.NewValidationError("Invalid category ID", "Category ID must be positive")
	}

	return nil
}

// ValidatePagination validates pagination parameters
func (v *Validator) ValidatePagination(limit, offset int) error {
	if limit <= 0 {
		return errors.NewValidationError("Invalid limit", "Limit must be positive")
	}

	if limit > 100 {
		return errors.NewValidationError("Limit too large", "Limit cannot exceed 100")
	}

	if offset < 0 {
		return errors.NewValidationError("Invalid offset", "Offset cannot be negative")
	}

	return nil
}

// ValidateDateRange validates a date range
func (v *Validator) ValidateDateRange(startDate, endDate time.Time) error {
	if err := v.ValidateDate(startDate, "start date"); err != nil {
		return err
	}

	if err := v.ValidateDate(endDate, "end date"); err != nil {
		return err
	}

	if startDate.After(endDate) {
		return errors.NewValidationError("Invalid date range", "Start date cannot be after end date")
	}

	return nil
}
