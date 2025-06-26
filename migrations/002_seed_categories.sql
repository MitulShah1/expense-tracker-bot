-- Migration: 002_seed_categories.sql
-- Description: Seed default categories for expense tracking
-- Created: 2024-01-01

-- Vehicle Expenses
INSERT INTO categories (name, emoji, "group") VALUES
('Petrol', '⛽', 'Vehicle'),
('Service', '🔧', 'Vehicle'),
('Repairs', '🛠️', 'Vehicle'),
('Insurance', '🚗', 'Vehicle'),
('Parking', '🅿️', 'Vehicle'),
('Toll', '🛣️', 'Vehicle'),
('Car Loan EMI', '🚘', 'Vehicle');

-- Home Expenses
INSERT INTO categories (name, emoji, "group") VALUES
('Home Loan EMI', '🏠', 'Home'),
('Electricity', '💡', 'Home'),
('Water', '💧', 'Home'),
('Gas', '🔥', 'Home'),
('Internet', '📶', 'Home'),
('Mobile', '📱', 'Home'),
('Cable/DTH', '📺', 'Home');

-- Daily Living
INSERT INTO categories (name, emoji, "group") VALUES
('Grocery', '🛒', 'Daily Living'),
('Dining', '🍽️', 'Daily Living'),
('Coffee/Tea', '☕', 'Daily Living'),
('Transportation', '🚕', 'Daily Living'),
('Shopping', '👕', 'Daily Living'),
('Personal Care', '💇', 'Daily Living');

-- Entertainment
INSERT INTO categories (name, emoji, "group") VALUES
('Movies', '🎬', 'Entertainment'),
('Gaming', '🎮', 'Entertainment'),
('Music', '🎵', 'Entertainment'),
('Books', '📚', 'Entertainment'),
('Hobbies', '🎨', 'Entertainment'),
('Fitness', '🏋️', 'Entertainment');

-- Health & Medical
INSERT INTO categories (name, emoji, "group") VALUES
('Medicines', '💊', 'Health'),
('Doctor', '👨‍⚕️', 'Health'),
('Hospital', '🏥', 'Health'),
('Wellness', '🧘', 'Health');

-- Education
INSERT INTO categories (name, emoji, "group") VALUES
('Courses', '📖', 'Education'),
('Stationery', '📝', 'Education'),
('Online Learning', '💻', 'Education');

-- Travel
INSERT INTO categories (name, emoji, "group") VALUES
('Flights', '✈️', 'Travel'),
('Hotels', '🏨', 'Travel'),
('Trains', '🚂', 'Travel'),
('Buses', '🚌', 'Travel');

-- Investments & Savings
INSERT INTO categories (name, emoji, "group") VALUES
('Savings', '💰', 'Investments'),
('Investments', '📈', 'Investments'),
('Bank Charges', '🏦', 'Investments');

-- Gifts & Donations
INSERT INTO categories (name, emoji, "group") VALUES
('Gifts', '🎁', 'Gifts'),
('Donations', '🤝', 'Gifts'),
('Celebrations', '🎉', 'Gifts');

-- Other
INSERT INTO categories (name, emoji, "group") VALUES
('Other', '📌', 'Other'); 