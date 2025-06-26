-- Migration: 004_add_views.sql
-- Description: Add useful views for reporting and analytics
-- Created: 2024-01-01

-- View for expense summary by user and category
CREATE VIEW expense_summary AS
SELECT 
    u.telegram_id,
    u.username,
    c."group" as category_group,
    c.name as category_name,
    c.emoji as category_emoji,
    COUNT(e.id) as expense_count,
    SUM(e.total_price) as total_amount,
    AVG(e.total_price) as avg_amount,
    MIN(e.timestamp) as first_expense,
    MAX(e.timestamp) as last_expense
FROM users u
LEFT JOIN expenses e ON u.id = e.user_id AND e.deleted_at IS NULL
LEFT JOIN categories c ON e.category_id = c.id
GROUP BY u.telegram_id, u.username, c."group", c.name, c.emoji;

-- View for monthly expense summary
CREATE VIEW monthly_expense_summary AS
SELECT 
    u.telegram_id,
    u.username,
    DATE_TRUNC('month', e.timestamp) as month,
    COUNT(e.id) as expense_count,
    SUM(e.total_price) as total_amount,
    AVG(e.total_price) as avg_amount
FROM users u
LEFT JOIN expenses e ON u.id = e.user_id AND e.deleted_at IS NULL
GROUP BY u.telegram_id, u.username, DATE_TRUNC('month', e.timestamp);

-- View for category-wise expense breakdown
CREATE VIEW category_expense_breakdown AS
SELECT 
    c."group" as category_group,
    c.name as category_name,
    c.emoji as category_emoji,
    COUNT(e.id) as expense_count,
    SUM(e.total_price) as total_amount,
    AVG(e.total_price) as avg_amount,
    MIN(e.total_price) as min_amount,
    MAX(e.total_price) as max_amount
FROM categories c
LEFT JOIN expenses e ON c.id = e.category_id AND e.deleted_at IS NULL
GROUP BY c."group", c.name, c.emoji;

-- View for user expense statistics
CREATE VIEW user_expense_stats AS
SELECT 
    u.telegram_id,
    u.username,
    COUNT(e.id) as total_expenses,
    SUM(e.total_price) as total_spent,
    AVG(e.total_price) as avg_expense,
    MIN(e.total_price) as min_expense,
    MAX(e.total_price) as max_expense,
    MIN(e.timestamp) as first_expense_date,
    MAX(e.timestamp) as last_expense_date,
    COUNT(DISTINCT c."group") as category_groups_used,
    COUNT(DISTINCT c.id) as categories_used
FROM users u
LEFT JOIN expenses e ON u.id = e.user_id AND e.deleted_at IS NULL
LEFT JOIN categories c ON e.category_id = c.id
GROUP BY u.telegram_id, u.username;

-- View for recent expenses (last 30 days)
CREATE VIEW recent_expenses AS
SELECT 
    u.telegram_id,
    u.username,
    c.name as category_name,
    c.emoji as category_emoji,
    e.total_price,
    e.notes,
    e.timestamp,
    e.created_at
FROM users u
JOIN expenses e ON u.id = e.user_id AND e.deleted_at IS NULL
JOIN categories c ON e.category_id = c.id
WHERE e.timestamp >= NOW() - INTERVAL '30 days'
ORDER BY e.timestamp DESC; 