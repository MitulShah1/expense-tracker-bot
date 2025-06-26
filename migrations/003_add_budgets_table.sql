-- Migration: 003_add_budgets_table.sql
-- Description: Add budgets table for budget tracking
-- Created: 2024-01-01

-- Create budgets table
CREATE TABLE budgets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id),
    amount FLOAT NOT NULL CHECK (amount > 0),
    period TEXT NOT NULL CHECK (period IN ('daily', 'weekly', 'monthly', 'yearly')),
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create budget_limits table for category-specific limits
CREATE TABLE budget_limits (
    id SERIAL PRIMARY KEY,
    budget_id INTEGER REFERENCES budgets(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id),
    limit_amount FLOAT NOT NULL CHECK (limit_amount > 0),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create indexes for budgets
CREATE INDEX idx_budgets_user_id ON budgets(user_id);
CREATE INDEX idx_budgets_period ON budgets(period);
CREATE INDEX idx_budgets_is_active ON budgets(is_active);
CREATE INDEX idx_budget_limits_budget_id ON budget_limits(budget_id);
CREATE INDEX idx_budget_limits_category_id ON budget_limits(category_id);

-- Create triggers for updated_at
CREATE TRIGGER update_budgets_updated_at BEFORE UPDATE ON budgets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_budget_limits_updated_at BEFORE UPDATE ON budget_limits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 