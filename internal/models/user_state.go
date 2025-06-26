package models

import "time"

// UserState represents the current state of a user in the conversation
type UserState struct {
	Step             Step
	VehicleType      string
	Category         string
	Odometer         float64
	PetrolPrice      float64
	TotalPrice       float64
	Notes            string
	EditMode         bool
	EditID           string
	DeleteExpense    *Expense   // Store the expense being deleted
	TempExpense      *Expense   // Temporary expense for adding/editing
	ExpenseSelection []*Expense // List of expenses shown for edit/delete
	LastActivity     time.Time  // Last activity timestamp
	CreatedAt        time.Time  // When the state was created
	UpdatedAt        time.Time  // When the state was last updated
}

// NewUserState creates a new user state
func NewUserState() *UserState {
	now := time.Now()
	return &UserState{
		Step:         StepStart,
		LastActivity: now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// UpdateActivity updates the last activity timestamp
func (s *UserState) UpdateActivity() {
	now := time.Now()
	s.LastActivity = now
	s.UpdatedAt = now
}
