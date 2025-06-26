# Test Data for Expense Tracker Bot

## Sample Expenses

### Vehicle Expenses

1. Petrol Expense
   - Category Group: Vehicle
   - Category: Petrol
   - Vehicle Type: CAR
   - Odometer: 15000
   - Petrol Price: 96.72
   - Total Price: 1000
   - Notes: Full tank fill

2. Service Expense
   - Category Group: Vehicle
   - Category: Service
   - Vehicle Type: CAR
   - Odometer: 15200
   - Total Price: 2500
   - Notes: Regular service

### Food Expenses

1. Groceries
   - Category Group: Food
   - Category: Groceries
   - Total Price: 1500
   - Notes: Monthly groceries

2. Restaurant
   - Category Group: Food
   - Category: Restaurant
   - Total Price: 800
   - Notes: Dinner with friends

### Utility Expenses

1. Electricity
   - Category Group: Utility
   - Category: Electricity
   - Total Price: 1200
   - Notes: Monthly bill

2. Water
   - Category Group: Utility
   - Category: Water
   - Total Price: 300
   - Notes: Monthly bill

## Test Scenarios

### 1. Add Expense Flow

1. Start with /add command
2. Follow the prompts to add each expense
3. Verify the expense is added correctly

### 2. Edit Expense Flow

1. Use /edit command
2. Select an expense to edit
3. Modify one or more fields
4. Verify the changes are saved

### 3. Delete Expense Flow

1. Use /delete command
2. Select an expense to delete
3. Confirm deletion
4. Verify the expense is removed

### 4. Report Generation

1. Use /report command
2. Verify all expenses are included
3. Check category-wise totals
4. Verify monthly totals

### 5. Dashboard Display

1. Use /dashboard command
2. Verify overall metrics
3. Check recent expenses
4. Verify efficiency calculations

### 6. Error Cases

1. Invalid inputs
   - Non-numeric values for prices
   - Invalid vehicle types
   - Missing required fields
2. Rate limiting
   - Multiple rapid commands
3. Network issues
   - Disconnect during operation
   - Reconnect and verify state
