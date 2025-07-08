package bot

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MitulShah1/expense-tracker-bot/internal/database"
	"github.com/MitulShah1/expense-tracker-bot/internal/logger"
	"github.com/MitulShah1/expense-tracker-bot/internal/models"
	"github.com/MitulShah1/expense-tracker-bot/internal/services"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// BotAPIInterface defines the methods used from tgbotapi.BotAPI
// This allows us to mock the API in tests
type BotAPIInterface interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetUpdatesChan(u tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
}

// Bot represents the Telegram bot
type Bot struct {
	api    BotAPIInterface
	db     database.Storage
	logger logger.Logger
	// Services
	expenseService  *services.ExpenseService
	categoryService *services.CategoryService
	userService     *services.UserService
	vectorService   services.VectorServiceInterface
	states          map[int64]*models.UserState
	// Add new fields for state management
	stateTimeout  time.Duration
	stateMutex    sync.RWMutex
	cleanupTicker *time.Ticker
	// Add rate limiting
	rateLimiter *rate.Limiter
	// Add metrics
	metrics struct {
		messageCount   int64
		commandCount   int64
		errorCount     int64
		expenseCount   int64
		activeUsers    int64
		lastUpdateTime time.Time
	}
	metricsMutex sync.RWMutex
}

// NewBot creates a new bot instance
func NewBot(ctx context.Context, token string, dbClient database.Storage, logger logger.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	// Initialize services
	expenseService := services.NewExpenseService(dbClient, logger)
	categoryService := services.NewCategoryService(dbClient, logger)
	userService := services.NewUserService(dbClient, logger)
	vectorService := services.NewVectorService(dbClient, logger)

	bot := &Bot{
		api:             api, // Use the real API here
		db:              dbClient,
		logger:          logger,
		expenseService:  expenseService,
		categoryService: categoryService,
		userService:     userService,
		vectorService:   vectorService,
		states:          make(map[int64]*models.UserState),
		stateTimeout:    30 * time.Minute,                                      // Default timeout of 30 minutes
		rateLimiter:     rate.NewLimiter(rate.Every(100*time.Millisecond), 10), // 10 requests per second
	}

	// Initialize metrics
	bot.metrics.lastUpdateTime = time.Now()

	// Start cleanup routine
	bot.startCleanupRoutine(ctx)

	return bot, nil
}

// HandleUpdate processes incoming updates from Telegram
func (b *Bot) HandleUpdate(ctx context.Context, update *tgbotapi.Update) error {
	// Create context with request ID
	reqCtx := context.WithValue(ctx, logger.RequestIDKey, fmt.Sprintf("update_%d", update.UpdateID))

	// Handle callback queries
	if update.CallbackQuery != nil {
		return b.handleCallbackQuery(reqCtx, update.CallbackQuery)
	}

	// Handle messages
	if update.Message == nil {
		return nil
	}

	// Get or create user state
	state := b.getState(update.Message.Chat.ID)
	if state == nil {
		state = models.NewUserState()
		b.setState(update.Message.Chat.ID, state)
	}

	// Handle commands
	if update.Message.IsCommand() {
		return b.handleCommand(reqCtx, update.Message)
	}

	// Handle state-based messages
	if state.Step == models.StepStart {
		return b.sendWelcome(reqCtx, update.Message)
	}
	return b.sendMessage(reqCtx, update.Message.Chat.ID, "Please use one of the available commands.")
}

// GetMetrics returns the current metrics
func (b *Bot) GetMetrics() map[string]any {
	b.metricsMutex.RLock()
	defer b.metricsMutex.RUnlock()

	return map[string]any{
		"message_count":    b.metrics.messageCount,
		"command_count":    b.metrics.commandCount,
		"error_count":      b.metrics.errorCount,
		"expense_count":    b.metrics.expenseCount,
		"active_users":     b.metrics.activeUsers,
		"last_update_time": b.metrics.lastUpdateTime,
	}
}

// Start starts the bot
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info(ctx, "Starting bot...")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			// Create context with request ID
			reqCtx := context.WithValue(ctx, logger.RequestIDKey, fmt.Sprintf("bot_%d", update.UpdateID))

			// Handle callback queries
			if update.CallbackQuery != nil {
				if err := b.handleCallbackQuery(reqCtx, update.CallbackQuery); err != nil {
					b.logger.Error(reqCtx, "Failed to handle callback query", logger.ErrorField(err))
				}
				continue
			}

			// Handle messages
			if update.Message == nil {
				continue
			}

			// Handle message
			if err := b.handleMessage(reqCtx, update.Message); err != nil {
				b.logger.Error(reqCtx, "Failed to handle message", logger.ErrorField(err))
			}
		}
	}
}

// Private methods
func (b *Bot) startCleanupRoutine(ctx context.Context) {
	b.cleanupTicker = time.NewTicker(5 * time.Minute)
	go func() {
		for range b.cleanupTicker.C {
			b.cleanupExpiredStates(ctx)
		}
	}()
}

func (b *Bot) cleanupExpiredStates(ctx context.Context) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()

	now := time.Now()
	for userID, state := range b.states {
		if now.Sub(state.LastActivity) > b.stateTimeout {
			delete(b.states, userID)
			b.logger.Info(ctx, "Cleaned up expired state",
				zap.Int64("user_id", userID))
		}
	}
}

func (b *Bot) getState(userID int64) *models.UserState {
	b.stateMutex.RLock()
	defer b.stateMutex.RUnlock()
	return b.states[userID]
}

func (b *Bot) setState(userID int64, state *models.UserState) {
	b.stateMutex.Lock()
	defer b.stateMutex.Unlock()
	state.LastActivity = time.Now()
	b.states[userID] = state
}

func (b *Bot) incrementMetric(metric *int64) {
	b.metricsMutex.Lock()
	defer b.metricsMutex.Unlock()
	*metric++
}

func (b *Bot) handleMessage(ctx context.Context, message *tgbotapi.Message) error {
	// Update metrics
	b.incrementMetric(&b.metrics.messageCount)
	b.metricsMutex.Lock()
	b.metrics.lastUpdateTime = time.Now()
	b.metricsMutex.Unlock()

	// Check rate limit
	if !b.rateLimiter.Allow() {
		b.incrementMetric(&b.metrics.errorCount)
		b.logger.Warn(ctx, "Rate limit exceeded", zap.Int64("user_id", message.From.ID))
		return b.sendMessage(ctx, message.Chat.ID, "Too many requests. Please try again later.")
	}

	// Get or create user state
	state := b.getState(message.From.ID)
	if state == nil {
		state = models.NewUserState()
		b.setState(message.From.ID, state)
		b.metricsMutex.Lock()
		b.metrics.activeUsers++
		b.metricsMutex.Unlock()
	}

	// Update activity
	state.UpdateActivity()

	// Special-case: if user is in StepNotes and sends /skip, treat as message, not command
	if state.Step == models.StepNotes && message.Text == "/skip" {
		return b.handleState(ctx, message, state)
	}

	// Handle command
	if message.IsCommand() {
		b.incrementMetric(&b.metrics.commandCount)
		return b.handleCommand(ctx, message)
	}

	// Handle state
	return b.handleState(ctx, message, state)
}

func (b *Bot) handleCommand(ctx context.Context, message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		return b.sendWelcome(ctx, message)
	case "help":
		return b.sendHelp(ctx, message)
	case "add":
		return b.handleAddCommand(ctx, message)
	case "list":
		return b.handleListCommand(ctx, message)
	case "edit":
		return b.handleEditCommand(ctx, message)
	case "delete":
		return b.handleDeleteCommand(ctx, message)
	case "report":
		return b.handleReportCommand(ctx, message)
	case "dashboard":
		return b.handleDashboardCommand(ctx, message)
	case "search":
		return b.handleSearchCommand(ctx, message)
	case "cancel":
		delete(b.states, message.Chat.ID)
		return b.sendMessage(ctx, message.Chat.ID, "Operation cancelled.")
	default:
		return b.sendMessage(ctx, message.Chat.ID, "Unknown command. Use /help to see available commands.")
	}
}

func (b *Bot) handleState(ctx context.Context, message *tgbotapi.Message, state *models.UserState) error {
	// Handle keyboard button text messages
	switch message.Text {
	case "📝 Add Expense":
		return b.handleAddCommand(ctx, message)
	case "📋 List Expenses":
		return b.handleListCommand(ctx, message)
	case "✏️ Edit Expense":
		return b.handleEditCommand(ctx, message)
	case "🗑️ Delete Expense":
		return b.handleDeleteCommand(ctx, message)
	case "📊 Reports":
		return b.handleReportCommand(ctx, message)
	case "📈 Dashboard":
		return b.handleDashboardCommand(ctx, message)
	}

	// Add semantic search state handling
	if state.Step == models.StepSearchExpense {
		return b.handleSearchQuery(ctx, message)
	}

	switch state.Step {
	case models.StepStart:
		return b.sendWelcome(ctx, message)
	case models.StepOdometer:
		// Parse odometer reading using helper
		odometer, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "odometer reading")
		if err != nil {
			return err
		}
		state.TempExpense.Odometer = odometer
		state.Step = models.StepPetrolPrice
		return b.sendMessage(ctx, message.Chat.ID, "⛽ Please enter the petrol price per liter:")
	case models.StepPetrolPrice:
		// Parse petrol price using helper
		price, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "petrol price")
		if err != nil {
			return err
		}
		state.TempExpense.PetrolPrice = price
		state.Step = models.StepTotalPrice
		return b.sendMessage(ctx, message.Chat.ID, "💰 Please enter the total price:")
	case models.StepTotalPrice:
		// Parse total price using helper
		total, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "total price")
		if err != nil {
			return err
		}
		state.TempExpense.TotalPrice = total
		state.Step = models.StepNotes
		return b.sendMessage(ctx, message.Chat.ID, "📝 Add any notes (or send /skip to skip):")
	case models.StepNotes:
		// Handle notes
		if message.Text != "/skip" {
			state.TempExpense.Notes = message.Text
		}

		// Set timestamp and user ID
		state.TempExpense.Timestamp = time.Now()
		state.TempExpense.UserID = message.From.ID

		// Create user if it doesn't exist
		_, err := b.userService.GetOrCreateUser(ctx, message.From.ID, message.From.UserName, message.From.FirstName, message.From.LastName)
		if err != nil {
			b.logger.Error(ctx, "Failed to create user", logger.ErrorField(err))
			return b.sendError(ctx, message.Chat.ID, err)
		}

		// Get category by name
		category, err := b.categoryService.GetCategoryByName(ctx, state.TempExpense.CategoryName)
		if err != nil {
			b.logger.Error(ctx, "Failed to get category", logger.ErrorField(err))
			return b.sendError(ctx, message.Chat.ID, err)
		}
		if category == nil {
			return b.sendMessage(ctx, message.Chat.ID, "Category not found. Please try again.")
		}

		// Create expense object
		expense := &models.Expense{
			CategoryName: state.TempExpense.CategoryName,
			TotalPrice:   state.TempExpense.TotalPrice,
			Odometer:     state.TempExpense.Odometer,
			PetrolPrice:  state.TempExpense.PetrolPrice,
			Notes:        state.TempExpense.Notes,
			Timestamp:    time.Now(),
		}

		// Save expense to database
		if err := b.expenseService.CreateExpense(ctx, expense, message.From.ID); err != nil {
			b.logger.Error(ctx, "Failed to create expense", logger.ErrorField(err))
			return b.sendError(ctx, message.Chat.ID, err)
		}

		// Reset state
		delete(b.states, message.Chat.ID)

		// Send confirmation
		return b.sendMessage(ctx, message.Chat.ID, "✅ Expense added successfully!")
	case models.StepEditOdometer:
		// Parse odometer reading for editing using helper
		odometer, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "odometer reading")
		if err != nil {
			return err
		}
		state.TempExpense.Odometer = odometer
		state.Step = models.StepEditExpense

		// Show updated expense and edit options
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Updated expense:\n%s - %s: ₹%.2f\nOdometer: %.1f km\n\nSelect what to edit:",
			state.TempExpense.Timestamp.Format("02 Jan 2006"),
			state.TempExpense.CategoryName,
			state.TempExpense.TotalPrice,
			state.TempExpense.Odometer))
		msg.ReplyMarkup = GetEditFieldKeyboard()
		_, err = b.api.Send(msg)
		return err
	case models.StepEditPetrolPrice:
		// Parse petrol price for editing using helper
		price, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "petrol price")
		if err != nil {
			return err
		}
		state.TempExpense.PetrolPrice = price
		state.Step = models.StepEditExpense

		// Show updated expense and edit options
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Updated expense:\n%s - %s: ₹%.2f\nPetrol Price: ₹%.2f/L\n\nSelect what to edit:",
			state.TempExpense.Timestamp.Format("02 Jan 2006"),
			state.TempExpense.CategoryName,
			state.TempExpense.TotalPrice,
			state.TempExpense.PetrolPrice))
		msg.ReplyMarkup = GetEditFieldKeyboard()
		_, err = b.api.Send(msg)
		return err
	case models.StepEditTotalPrice:
		// Parse total price for editing using helper
		total, err := b.parseFloatOrReply(ctx, message.Chat.ID, message.Text, "total price")
		if err != nil {
			return err
		}
		state.TempExpense.TotalPrice = total
		state.Step = models.StepEditExpense

		// Show updated expense and edit options
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Updated expense:\n%s - %s: ₹%.2f\n\nSelect what to edit:",
			state.TempExpense.Timestamp.Format("02 Jan 2006"),
			state.TempExpense.CategoryName,
			state.TempExpense.TotalPrice))
		msg.ReplyMarkup = GetEditFieldKeyboard()
		_, err = b.api.Send(msg)
		return err
	case models.StepEditNotes:
		// Handle notes for editing
		if message.Text != "/skip" {
			state.TempExpense.Notes = message.Text
		} else {
			state.TempExpense.Notes = ""
		}
		state.Step = models.StepEditExpense

		// Show updated expense and edit options
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Updated expense:\n%s - %s: ₹%.2f\nNotes: %s\n\nSelect what to edit:",
			state.TempExpense.Timestamp.Format("02 Jan 2006"),
			state.TempExpense.CategoryName,
			state.TempExpense.TotalPrice,
			state.TempExpense.Notes))
		msg.ReplyMarkup = GetEditFieldKeyboard()
		_, err := b.api.Send(msg)
		return err
	default:
		return b.sendMessage(ctx, message.Chat.ID, "Please use one of the available commands or buttons.")
	}
}

func (b *Bot) sendMessage(ctx context.Context, chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		b.logger.Error(ctx, "Failed to send message", logger.ErrorField(err))
		return err
	}
	return nil
}

func (b *Bot) sendError(ctx context.Context, chatID int64, err error) error {
	b.logger.Error(ctx, "Error occurred", logger.ErrorField(err))
	return b.sendMessage(ctx, chatID, fmt.Sprintf("An error occurred: %v", err))
}

func (b *Bot) sendWelcome(ctx context.Context, message *tgbotapi.Message) error {
	text := `Welcome to the Expense Tracker Bot! 🚗💰

I can help you track your vehicle expenses. Here's what you can do:

📝 Add Expense - Add a new expense
📋 List Expenses - View all expenses
✏️ Edit Expense - Modify existing expenses
🗑️ Delete Expense - Remove expenses
📊 Reports - Generate expense reports
📈 Dashboard - View expense dashboard

You can also use commands like /add, /list, /report, etc.

Let's get started! Use the buttons below or type /add to record your first expense.`

	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = GetMainMenuKeyboard()
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) sendHelp(ctx context.Context, message *tgbotapi.Message) error {
	text := `Here are the available commands:

/add - Add a new expense
/list - List your expenses
/edit - Edit an existing expense
/delete - Delete an expense
/search - Search expenses using natural language
/help - Show this help message
/cancel - Cancel current operation

To add an expense:
1. Use /add
2. Select a category
3. Enter vehicle type
4. Enter odometer reading
5. Enter petrol price
6. Enter total price
7. Add optional notes

To edit or delete an expense:
1. Use /edit or /delete
2. Select the expense from the list
3. Follow the prompts

To search expenses:
1. Use /search
2. Enter a natural language query
3. View matching expenses`
	return b.sendMessage(ctx, message.Chat.ID, text)
}

// handleCallbackQuery processes callback queries from inline keyboards
func (b *Bot) handleCallbackQuery(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	if callback == nil || callback.Message == nil {
		return errors.New("invalid callback query")
	}

	// Get or create user state
	state := b.getState(callback.Message.Chat.ID)
	if state == nil {
		state = models.NewUserState()
		b.setState(callback.Message.Chat.ID, state)
	}

	// Initialize TempExpense if nil
	if state.TempExpense == nil {
		state.TempExpense = &models.Expense{
			UserID: callback.Message.Chat.ID,
		}
	}

	// Process callback data
	data := callback.Data
	b.logger.Info(ctx, "Callback data received", logger.String("data", data))
	switch {
	case strings.HasPrefix(data, "group_"):
		// Handle category group selection
		groupName := strings.TrimPrefix(data, "group_")
		categories, err := b.categoryService.GetCategoriesByGroup(ctx, groupName)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			"Select a category:",
			GetCategoryGroupKeyboard(categories),
		)
		_, err = b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "category_"):
		// Handle category selection
		categoryName := strings.TrimPrefix(data, "category_")
		category, err := b.categoryService.GetCategoryByName(ctx, categoryName)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}
		if category == nil {
			return b.sendError(ctx, callback.Message.Chat.ID, fmt.Errorf("category not found: %s", categoryName))
		}

		// Store category in state
		state.TempExpense.CategoryName = category.Name

		if category.Group == "Vehicle" {
			// Vehicle categories: Vehicle Type → Odometer/Petrol Price → Total Price → Notes
			state.Step = models.StepVehicleType

			// Create vehicle type keyboard
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚗 Car", "vehicle_CAR"),
					tgbotapi.NewInlineKeyboardButtonData("🏍️ Bike", "vehicle_BIKE"),
				),
			)

			msg := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				fmt.Sprintf("Selected category: %s %s\nSelect vehicle type:", category.Emoji, category.Name),
				keyboard,
			)
			_, err = b.api.Send(msg)
			return err
		}
		// Non-vehicle categories: Directly to Total Price → Notes
		state.Step = models.StepTotalPrice

		msg := tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			fmt.Sprintf("Selected category: %s %s\n💰 Please enter the total price:", category.Emoji, category.Name),
		)
		_, err = b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "vehicle_"):
		// Handle vehicle type selection
		vehicleType := strings.TrimPrefix(data, "vehicle_")
		state.TempExpense.VehicleType = sql.NullString{String: vehicleType, Valid: true}

		// Determine next step based on category
		if state.TempExpense.CategoryName == "Petrol" {
			// Petrol category: Vehicle Type → Odometer → Petrol Price → Total Price → Notes
			state.Step = models.StepOdometer
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🔢 Please enter the odometer reading (in km):")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.api.Send(msg)
			return err
		}
		// Other vehicle categories: Vehicle Type → Total Price → Notes
		state.Step = models.StepTotalPrice
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💰 Please enter the total price:")
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		_, err := b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "edit_field_"):
		// Handle edit field selection
		field := strings.TrimPrefix(data, "edit_field_")
		switch field {
		case "category":
			state.Step = models.StepEditCategory
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				"Select new category:",
				GetCategoryKeyboard(),
			)
			_, err := b.api.Send(msg)
			return err
		case "vehicle":
			state.Step = models.StepEditVehicleType
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🚗 Car", "edit_vehicle_CAR"),
					tgbotapi.NewInlineKeyboardButtonData("🏍️ Bike", "edit_vehicle_BIKE"),
				),
			)
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				"Select new vehicle type:",
				keyboard,
			)
			_, err := b.api.Send(msg)
			return err
		case "odometer":
			state.Step = models.StepEditOdometer
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "🔢 Enter new odometer reading (in km):")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.api.Send(msg)
			return err
		case "petrol":
			state.Step = models.StepEditPetrolPrice
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "⛽ Enter new petrol price per liter:")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.api.Send(msg)
			return err
		case "total":
			state.Step = models.StepEditTotalPrice
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "💰 Enter new total price:")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.api.Send(msg)
			return err
		case "notes":
			state.Step = models.StepEditNotes
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "📝 Enter new notes (or send /skip to clear):")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			_, err := b.api.Send(msg)
			return err
		default:
			return b.sendMessage(ctx, callback.Message.Chat.ID, "Invalid edit field selection.")
		}

	case strings.HasPrefix(data, "edit_vehicle_"):
		// Handle vehicle type selection for editing
		vehicleType := strings.TrimPrefix(data, "edit_vehicle_")
		if state.TempExpense != nil {
			state.TempExpense.VehicleType = sql.NullString{String: vehicleType, Valid: true}
			// Show updated expense and edit options
			msg := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				fmt.Sprintf("Updated expense:\n%s - %s: ₹%.2f\nVehicle: %s\n\nSelect what to edit:",
					state.TempExpense.Timestamp.Format("02 Jan 2006"),
					state.TempExpense.CategoryName,
					state.TempExpense.TotalPrice,
					state.TempExpense.VehicleType.String),
				GetEditFieldKeyboard(),
			)
			_, err := b.api.Send(msg)
			return err
		}
		return b.sendError(ctx, callback.Message.Chat.ID, errors.New("no expense to edit"))

	case data == "edit_save":
		// Save the edited expense
		if state.TempExpense == nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, "No expense to save.")
		}

		// Update expense in database
		if err := b.expenseService.UpdateExpense(ctx, state.TempExpense, callback.Message.Chat.ID); err != nil {
			b.logger.Error(ctx, "Failed to update expense", logger.ErrorField(err))
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}

		// Reset state
		delete(b.states, callback.Message.Chat.ID)

		// Send confirmation
		msg := tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			"✅ Expense updated successfully!")
		_, err := b.api.Send(msg)
		return err

	case data == "edit_cancel":
		// Cancel editing
		delete(b.states, callback.Message.Chat.ID)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "❌ Edit cancelled.")
		msg.ReplyMarkup = GetMainMenuKeyboard()
		_, err := b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "edit_"):
		// Handle expense editing
		expenseIDStr := strings.TrimPrefix(data, "edit_")
		expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}

		// Get expense from database
		expenseToEdit, err := b.expenseService.GetExpenseByID(ctx, expenseID)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}
		if expenseToEdit == nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, "Expense not found.")
		}

		// Check if user owns this expense using helper
		if err := b.checkExpenseOwnership(ctx, callback.Message.Chat.ID, expenseToEdit); err != nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, err.Error())
		}

		// Store expense in state for editing
		state.TempExpense = expenseToEdit
		state.Step = models.StepEditExpense

		// Send edit options
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			fmt.Sprintf("Editing expense:\n%s - %s: ₹%.2f\n\nSelect what to edit:",
				expenseToEdit.Timestamp.Format("02 Jan 2006"),
				expenseToEdit.CategoryName,
				expenseToEdit.TotalPrice),
			GetEditFieldKeyboard(),
		)
		_, err = b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "delete_"):
		// Handle expense deletion
		expenseIDStr := strings.TrimPrefix(data, "delete_")
		expenseID, err := strconv.ParseInt(expenseIDStr, 10, 64)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}

		// Get expense from database
		expenseToDelete, err := b.expenseService.GetExpenseByID(ctx, expenseID)
		if err != nil {
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}
		if expenseToDelete == nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, "Expense not found.")
		}

		// Check if user owns this expense using helper
		if err := b.checkExpenseOwnership(ctx, callback.Message.Chat.ID, expenseToDelete); err != nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, err.Error())
		}

		// Store expense in state for deletion
		state.DeleteExpense = expenseToDelete
		state.Step = models.StepConfirmDelete

		// Send confirmation
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			fmt.Sprintf("Are you sure you want to delete this expense?\n\n%s - %s: ₹%.2f",
				expenseToDelete.Timestamp.Format("02 Jan 2006"),
				expenseToDelete.CategoryName,
				expenseToDelete.TotalPrice),
			GetConfirmationKeyboard(),
		)
		_, err = b.api.Send(msg)
		return err

	case data == "confirm_delete":
		// Handle delete confirmation
		if state.DeleteExpense == nil {
			return b.sendMessage(ctx, callback.Message.Chat.ID, "No expense selected for deletion.")
		}

		// Delete expense from database
		if err := b.expenseService.DeleteExpense(ctx, state.DeleteExpense.ID, callback.Message.Chat.ID); err != nil {
			b.logger.Error(ctx, "Failed to delete expense", logger.ErrorField(err))
			return b.sendError(ctx, callback.Message.Chat.ID, err)
		}

		// Reset state
		delete(b.states, callback.Message.Chat.ID)

		// Send confirmation
		msg := tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			"✅ Expense deleted successfully!")
		_, err := b.api.Send(msg)
		return err

	case data == "back_to_groups":
		// Handle back to groups
		msg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			"Select a category group:",
			GetCategoryKeyboard(),
		)
		_, err := b.api.Send(msg)
		return err

	case data == "back_to_main":
		// Handle back to main menu
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Main menu:")
		msg.ReplyMarkup = GetMainMenuKeyboard()
		_, err := b.api.Send(msg)
		return err

	case strings.HasPrefix(data, "confirm_"):
		// Handle confirmation
		confirmed := data == "confirm_yes"
		if confirmed {
			// Process the confirmed action
			if state.Step == models.StepConfirmDelete {
				if state.DeleteExpense != nil {
					// TODO: Replace with database storage
					b.logger.Info(ctx, "Expense deleted (placeholder)",
						logger.String("category", state.DeleteExpense.CategoryName),
						logger.Float64("total_price", state.DeleteExpense.TotalPrice),
						logger.Int("user_id", int(state.DeleteExpense.UserID)))

					msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "✅ Expense deleted successfully!")
					msg.ReplyMarkup = GetMainMenuKeyboard()
					_, err := b.api.Send(msg)
					delete(b.states, callback.Message.Chat.ID)
					return err
				}
			}
		}
		// Reset state and return to main menu
		delete(b.states, callback.Message.Chat.ID)
		msg := tgbotapi.NewMessage(callback.Message.Chat.ID, "Operation cancelled.")
		msg.ReplyMarkup = GetMainMenuKeyboard()
		_, err := b.api.Send(msg)
		return err

	default:
		// Unknown callback data
		return b.sendMessage(ctx, callback.Message.Chat.ID, "Invalid selection. Please try again.")
	}
}
