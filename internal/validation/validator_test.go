package validation

import (
	"strings"
	"testing"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/errors"
)

func TestNewValidator(t *testing.T) {
	validator := NewValidator()
	if validator == nil {
		t.Error("NewValidator() should not return nil")
	}
}

func TestValidateTelegramID(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		telegramID int64
		wantErr    bool
	}{
		{"valid positive ID", 123456789, false},
		{"valid large ID", 999999999999, false},
		{"zero ID", 0, true},
		{"negative ID", -1, true},
		{"large negative ID", -123456789, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateTelegramID(tt.telegramID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTelegramID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateTelegramID() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateUsername(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{"valid username", "john_doe", false},
		{"valid username with numbers", "user123", false},
		{"valid username minimum length", "abcde", false},
		{"valid username maximum length", "a1234567890123456789012345678901", false},
		{"empty username", "", true},
		{"too short username", "abcd", true},
		{"too long username", "a12345678901234567890123456789012", true},
		{"username with special chars", "user@name", true},
		{"username with spaces", "user name", true},
		{"username with hyphens", "user-name", true},
		{"username starting with number", "1user", false},
		{"username with uppercase", "User123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUsername(tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateUsername() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateName(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		nameValue string
		fieldName string
		wantErr   bool
	}{
		{"valid first name", "John", "first name", false},
		{"valid last name", "Doe", "last name", false},
		{"valid name with space", "John Doe", "full name", false},
		{"valid name with hyphen", "Jean-Pierre", "first name", false},
		{"valid name with apostrophe", "O'Connor", "last name", false},
		{"valid name maximum length", strings.Repeat("A", 64), "name", false},
		{"empty name", "", "name", true},
		{"too long name", strings.Repeat("A", 65), "name", true},
		{"name with numbers", "John123", "name", true},
		{"name with special chars", "John@Doe", "name", true},
		{"name with underscore", "John_Doe", "name", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateName(tt.nameValue, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateName() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		amount    float64
		fieldName string
		wantErr   bool
	}{
		{"valid amount", 100.50, "amount", false},
		{"valid small amount", 0.01, "amount", false},
		{"valid large amount", 999999999.99, "amount", false},
		{"zero amount", 0, "amount", true},
		{"negative amount", -10.50, "amount", true},
		{"too large amount", 1000000000.00, "amount", true},
		{"valid amount with custom field", 50.25, "price", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAmount(tt.amount, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateAmount() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateAmountString(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		amountStr  string
		fieldName  string
		wantAmount float64
		wantErr    bool
	}{
		{"valid amount string", "100.50", "amount", 100.50, false},
		{"valid amount with rupee symbol", "₹50.25", "amount", 50.25, false},
		{"valid amount with spaces", " 75.00 ", "amount", 75.00, false},
		{"valid amount with rupee and spaces", " ₹ 25.50 ", "amount", 25.50, false},
		{"valid integer amount", "100", "amount", 100.0, false},
		{"empty string", "", "amount", 0, true},
		{"invalid number", "abc", "amount", 0, true},
		{"negative amount", "-10.50", "amount", 0, true},
		{"zero amount", "0", "amount", 0, true},
		{"too large amount", "1000000000", "amount", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := validator.ValidateAmountString(tt.amountStr, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAmountString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && amount != tt.wantAmount {
				t.Errorf("ValidateAmountString() amount = %v, want %v", amount, tt.wantAmount)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateAmountString() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateOdometer(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name     string
		odometer float64
		wantErr  bool
	}{
		{"valid odometer", 50000.5, false},
		{"zero odometer", 0, false},
		{"valid large odometer", 999999.9, false},
		{"negative odometer", -100, true},
		{"too large odometer", 1000000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateOdometer(tt.odometer)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOdometer() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateOdometer() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateOdometerString(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name         string
		odometerStr  string
		wantOdometer float64
		wantErr      bool
	}{
		{"valid odometer string", "50000.5", 50000.5, false},
		{"valid integer odometer", "75000", 75000.0, false},
		{"zero odometer", "0", 0.0, false},
		{"empty string", "", 0, true},
		{"invalid number", "abc", 0, true},
		{"negative odometer", "-100", 0, true},
		{"too large odometer", "1000000", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			odometer, err := validator.ValidateOdometerString(tt.odometerStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOdometerString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && odometer != tt.wantOdometer {
				t.Errorf("ValidateOdometerString() odometer = %v, want %v", odometer, tt.wantOdometer)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateOdometerString() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateNotes(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		notes   string
		wantErr bool
	}{
		{"valid notes", "This is a valid note", false},
		{"empty notes", "", false},
		{"notes with special chars", "Note with @#$% symbols", false},
		{"notes with numbers", "Note with 123 numbers", false},
		{"notes at limit", string(make([]rune, 500)), false},
		{"notes over limit", string(make([]rune, 501)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateNotes(tt.notes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNotes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateNotes() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateCategoryName(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name         string
		categoryName string
		wantErr      bool
	}{
		{"valid category name", "Food", false},
		{"valid category with spaces", "Transport Fuel", false},
		{"valid category at limit", string(make([]rune, 100)), false},
		{"empty category name", "", true},
		{"category name over limit", string(make([]rune, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCategoryName(tt.categoryName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategoryName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateCategoryName() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateVehicleType(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name        string
		vehicleType string
		wantErr     bool
	}{
		{"valid car type", "CAR", false},
		{"valid bike type", "BIKE", false},
		{"invalid type", "TRUCK", true},
		{"empty type", "", true},
		{"lowercase car", "car", true},
		{"lowercase bike", "bike", true},
		{"mixed case", "Car", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateVehicleType(tt.vehicleType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVehicleType() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateVehicleType() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	validator := NewValidator()
	now := time.Now()

	tests := []struct {
		name      string
		date      time.Time
		fieldName string
		wantErr   bool
	}{
		{"valid past date", now.AddDate(0, 0, -1), "date", false},
		{"valid date 5 years ago", now.AddDate(-5, 0, 0), "date", false},
		{"valid date 9 years ago", now.AddDate(-9, 0, 0), "date", false},
		{"zero date", time.Time{}, "date", true},
		{"future date", now.AddDate(0, 0, 1), "date", true},
		{"date too far in past", now.AddDate(-11, 0, 0), "date", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDate(tt.date, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateDate() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateExpenseID(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name      string
		expenseID int64
		wantErr   bool
	}{
		{"valid expense ID", 1, false},
		{"valid large expense ID", 999999999, false},
		{"zero expense ID", 0, true},
		{"negative expense ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateExpenseID(tt.expenseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateExpenseID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateExpenseID() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateUserID(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		userID  int64
		wantErr bool
	}{
		{"valid user ID", 1, false},
		{"valid large user ID", 999999999, false},
		{"zero user ID", 0, true},
		{"negative user ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateUserID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateUserID() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateCategoryID(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name       string
		categoryID int64
		wantErr    bool
	}{
		{"valid category ID", 1, false},
		{"valid large category ID", 999999999, false},
		{"zero category ID", 0, true},
		{"negative category ID", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCategoryID(tt.categoryID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCategoryID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateCategoryID() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidatePagination(t *testing.T) {
	validator := NewValidator()

	tests := []struct {
		name    string
		limit   int
		offset  int
		wantErr bool
	}{
		{"valid pagination", 10, 0, false},
		{"valid pagination with offset", 20, 50, false},
		{"valid maximum limit", 100, 0, false},
		{"zero limit", 0, 0, true},
		{"negative limit", -1, 0, true},
		{"limit too large", 101, 0, true},
		{"negative offset", 10, -1, true},
		{"negative limit and offset", -1, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePagination(tt.limit, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePagination() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidatePagination() should return validation error, got %T", err)
				}
			}
		})
	}
}

func TestValidateDateRange(t *testing.T) {
	validator := NewValidator()
	now := time.Now()

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		wantErr   bool
	}{
		{"valid date range", now.AddDate(0, 0, -10), now.AddDate(0, 0, -1), false},
		{"same start and end date", now.AddDate(0, 0, -5), now.AddDate(0, 0, -5), false},
		{"start date after end date", now.AddDate(0, 0, -1), now.AddDate(0, 0, -10), true},
		{"invalid start date (future)", now.AddDate(0, 0, 1), now.AddDate(0, 0, 10), true},
		{"invalid end date (future)", now.AddDate(0, 0, -10), now.AddDate(0, 0, 1), true},
		{"invalid start date (too old)", now.AddDate(-11, 0, 0), now.AddDate(0, 0, -1), true},
		{"invalid end date (too old)", now.AddDate(0, 0, -10), now.AddDate(-11, 0, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateDateRange(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if appErr, ok := err.(*errors.AppError); ok && !appErr.IsValidationError() {
					t.Errorf("ValidateDateRange() should return validation error, got %T", err)
				}
			}
		})
	}
}
