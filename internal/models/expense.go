// Package models defines the data structures used throughout the application.
package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Step represents the current step in the conversation flow
type Step int

const (
	StepStart Step = iota + 100 // Start from 100 to avoid conflicts with other steps
	StepVehicleType
	StepCategory
	StepOdometer
	StepPetrolPrice
	StepTotalPrice
	StepNotes
	StepComplete
	StepEditExpense
	StepEditCategory
	StepEditVehicleType
	StepEditOdometer
	StepEditPetrolPrice
	StepEditTotalPrice
	StepEditNotes
	StepDeleteExpense
	StepConfirmDelete
	StepSearchExpense
	StepNone
)

// User represents a Telegram user
type User struct {
	ID         int64     `db:"id"          json:"id"`
	TelegramID int64     `db:"telegram_id" json:"telegramId"`
	Username   string    `db:"username"    json:"username"`
	FirstName  string    `db:"first_name"  json:"firstName"`
	LastName   string    `db:"last_name"   json:"lastName"`
	CreatedAt  time.Time `db:"created_at"  json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at"  json:"updatedAt"`
}

// Category represents an expense category
type Category struct {
	ID        int64     `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Emoji     string    `db:"emoji"      json:"emoji"`
	Group     string    `db:"group"      json:"group"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Expense represents an expense record
type Expense struct {
	ID          int64          `db:"id"           json:"id"`
	UserID      int64          `db:"user_id"      json:"userId"`
	CategoryID  int64          `db:"category_id"  json:"categoryId"`
	VehicleType sql.NullString `db:"vehicle_type" json:"vehicleType"` // CAR/BIKE or NULL
	Odometer    float64        `db:"odometer"     json:"odometer"`    // Optional
	PetrolPrice float64        `db:"petrol_price" json:"petrolPrice"` // Optional
	TotalPrice  float64        `db:"total_price"  json:"totalPrice"`
	Notes       string         `db:"notes"        json:"notes,omitempty"`
	Timestamp   time.Time      `db:"timestamp"    json:"timestamp"`
	CreatedAt   time.Time      `db:"created_at"   json:"createdAt"`
	UpdatedAt   time.Time      `db:"updated_at"   json:"updatedAt"`
	DeletedAt   *time.Time     `db:"deleted_at"   json:"deletedAt,omitempty"`

	// Vector embeddings
	NotesEmbedding    Float32Vector `db:"notes_embedding"    json:"notesEmbedding,omitempty"`
	CategoryEmbedding Float32Vector `db:"category_embedding" json:"categoryEmbedding,omitempty"`

	Similarity float64 `db:"similarity" json:"similarity,omitempty"`

	// Joined fields from categories table
	CategoryName  string `db:"category_name"  json:"categoryName"`
	CategoryEmoji string `db:"category_emoji" json:"categoryEmoji"`
	CategoryGroup string `db:"category_group" json:"categoryGroup"`
}

// ExpenseStats represents expense statistics for a user
type ExpenseStats struct {
	TotalExpenses    int64     `db:"total_expenses"     json:"totalExpenses"`
	TotalSpent       float64   `db:"total_spent"        json:"totalSpent"`
	AvgExpense       float64   `db:"avg_expense"        json:"avgExpense"`
	MinExpense       float64   `db:"min_expense"        json:"minExpense"`
	MaxExpense       float64   `db:"max_expense"        json:"maxExpense"`
	FirstExpenseDate time.Time `db:"first_expense_date" json:"firstExpenseDate"`
	LastExpenseDate  time.Time `db:"last_expense_date"  json:"lastExpenseDate"`
}

// ExpenseEmbedding represents the vector embeddings for an expense
type ExpenseEmbedding struct {
	ID                int64     `db:"id"`
	NotesEmbedding    []float32 `db:"notes_embedding"`
	CategoryEmbedding []float32 `db:"category_embedding"`
	UpdatedAt         time.Time `db:"updated_at"`
}

// VehicleType represents the type of vehicle
type VehicleType string

const (
	VehicleTypeCar  VehicleType = "CAR"
	VehicleTypeBike VehicleType = "BIKE"
)

// CategoryType represents the type of expense
type CategoryType string

const (
	// Vehicle Expenses
	CategoryPetrol    CategoryType = "⛽ Petrol"
	CategoryService   CategoryType = "🔧 Service"
	CategoryRepairs   CategoryType = "🛠️ Repairs"
	CategoryInsurance CategoryType = "🚗 Insurance"
	CategoryParking   CategoryType = "🅿️ Parking"
	CategoryToll      CategoryType = "🛣️ Toll"
	CategoryCarLoan   CategoryType = "🚘 Car Loan EMI"

	// Home Expenses
	CategoryHomeLoan    CategoryType = "🏠 Home Loan EMI"
	CategoryElectricity CategoryType = "💡 Electricity"
	CategoryWater       CategoryType = "💧 Water"
	CategoryGas         CategoryType = "🔥 Gas"
	CategoryInternet    CategoryType = "📶 Internet"
	CategoryMobile      CategoryType = "📱 Mobile"
	CategoryCable       CategoryType = "📺 Cable/DTH"

	// Daily Living
	CategoryGrocery   CategoryType = "🛒 Grocery"
	CategoryDining    CategoryType = "🍽️ Dining"
	CategoryCoffee    CategoryType = "☕ Coffee/Tea"
	CategoryTransport CategoryType = "🚕 Transportation"
	CategoryShopping  CategoryType = "👕 Shopping"
	CategoryPersonal  CategoryType = "💇 Personal Care"

	// Entertainment
	CategoryMovies  CategoryType = "🎬 Movies"
	CategoryGaming  CategoryType = "🎮 Gaming"
	CategoryMusic   CategoryType = "🎵 Music"
	CategoryBooks   CategoryType = "📚 Books"
	CategoryHobbies CategoryType = "🎨 Hobbies"
	CategoryFitness CategoryType = "🏋️ Fitness"

	// Health & Medical
	CategoryMedicines CategoryType = "💊 Medicines"
	CategoryDoctor    CategoryType = "👨‍⚕️ Doctor"
	CategoryHospital  CategoryType = "🏥 Hospital"
	CategoryWellness  CategoryType = "🧘 Wellness"

	// Education
	CategoryCourses     CategoryType = "📖 Courses"
	CategoryStationery  CategoryType = "📝 Stationery"
	CategoryOnlineLearn CategoryType = "💻 Online Learning"

	// Travel
	CategoryFlights CategoryType = "✈️ Flights"
	CategoryHotels  CategoryType = "🏨 Hotels"
	CategoryTrains  CategoryType = "🚂 Trains"
	CategoryBuses   CategoryType = "🚌 Buses"

	// Investments & Savings
	CategorySavings     CategoryType = "💰 Savings"
	CategoryInvestments CategoryType = "📈 Investments"
	CategoryBankCharges CategoryType = "🏦 Bank Charges"

	// Gifts & Donations
	CategoryGifts        CategoryType = "🎁 Gifts"
	CategoryDonations    CategoryType = "🤝 Donations"
	CategoryCelebrations CategoryType = "🎉 Celebrations"

	// Other
	CategoryOther CategoryType = "📌 Other"
)

// CategoryGroup represents a group of related categories
type CategoryGroup struct {
	Name       string
	Emoji      string
	Categories []CategoryType
}

// GetCategoryGroups returns all category groups
func GetCategoryGroups() []CategoryGroup {
	return []CategoryGroup{
		{
			Name:  "Vehicle",
			Emoji: "🚗",
			Categories: []CategoryType{
				CategoryPetrol,
				CategoryService,
				CategoryRepairs,
				CategoryInsurance,
				CategoryParking,
				CategoryToll,
				CategoryCarLoan,
			},
		},
		{
			Name:  "Home",
			Emoji: "🏠",
			Categories: []CategoryType{
				CategoryHomeLoan,
				CategoryElectricity,
				CategoryWater,
				CategoryGas,
				CategoryInternet,
				CategoryMobile,
				CategoryCable,
			},
		},
		{
			Name:  "Daily Living",
			Emoji: "🏪",
			Categories: []CategoryType{
				CategoryGrocery,
				CategoryDining,
				CategoryCoffee,
				CategoryTransport,
				CategoryShopping,
				CategoryPersonal,
			},
		},
		{
			Name:  "Entertainment",
			Emoji: "🎮",
			Categories: []CategoryType{
				CategoryMovies,
				CategoryGaming,
				CategoryMusic,
				CategoryBooks,
				CategoryHobbies,
				CategoryFitness,
			},
		},
		{
			Name:  "Health",
			Emoji: "🏥",
			Categories: []CategoryType{
				CategoryMedicines,
				CategoryDoctor,
				CategoryHospital,
				CategoryWellness,
			},
		},
		{
			Name:  "Education",
			Emoji: "📚",
			Categories: []CategoryType{
				CategoryCourses,
				CategoryStationery,
				CategoryOnlineLearn,
			},
		},
		{
			Name:  "Travel",
			Emoji: "✈️",
			Categories: []CategoryType{
				CategoryFlights,
				CategoryHotels,
				CategoryTrains,
				CategoryBuses,
			},
		},
		{
			Name:  "Investments",
			Emoji: "💹",
			Categories: []CategoryType{
				CategorySavings,
				CategoryInvestments,
				CategoryBankCharges,
			},
		},
		{
			Name:  "Gifts",
			Emoji: "🎁",
			Categories: []CategoryType{
				CategoryGifts,
				CategoryDonations,
				CategoryCelebrations,
			},
		},
		{
			Name:  "Other",
			Emoji: "📌",
			Categories: []CategoryType{
				CategoryOther,
			},
		},
	}
}

// Operation represents the type of operation being performed
type Operation string

const (
	OperationCreate Operation = "CREATE"
	OperationEdit   Operation = "EDIT"
	OperationDelete Operation = "DELETE"
)

// Float32Vector is a custom type for scanning/valuing Postgres vector columns
// (pgvector stores as string like '[0.1,0.2,0.3]')
type Float32Vector []float32

// MarshalJSON for JSON encoding
func (v *Float32Vector) MarshalJSON() ([]byte, error) {
	return json.Marshal([]float32(*v))
}

// UnmarshalJSON for JSON decoding
func (v *Float32Vector) UnmarshalJSON(data []byte) error {
	var arr []float32
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*v = arr
	return nil
}

// Scan implements the sql.Scanner interface for Float32Vector
func (v *Float32Vector) Scan(src any) error {
	switch data := src.(type) {
	case string:
		return v.fromString(data)
	case []byte:
		return v.fromString(string(data))
	case nil:
		*v = nil
		return nil
	default:
		return fmt.Errorf("cannot scan %T into Float32Vector", src)
	}
}

// Value implements the driver.Valuer interface for Float32Vector
func (v *Float32Vector) Value() (driver.Value, error) {
	if v == nil {
		return nil, errors.New("Float32Vector is nil")
	}
	parts := make([]string, len(*v))
	for i, f := range *v {
		parts[i] = strconv.FormatFloat(float64(f), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]", nil
}

// fromString parses a pgvector string (e.g. '[0.1,0.2,0.3]')
func (v *Float32Vector) fromString(s string) error {
	s = strings.Trim(s, "[]")
	if s == "" {
		*v = Float32Vector{}
		return nil
	}
	parts := strings.Split(s, ",")
	vec := make(Float32Vector, len(parts))
	for i, p := range parts {
		f, err := strconv.ParseFloat(strings.TrimSpace(p), 32)
		if err != nil {
			return err
		}
		vec[i] = float32(f)
	}
	*v = vec
	return nil
}
