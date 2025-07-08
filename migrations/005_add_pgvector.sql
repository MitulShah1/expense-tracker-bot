-- Migration: 005_add_pgvector.sql
-- Description: Add pgvector extension and vector columns for semantic search
-- Created: 2024-01-01

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Add vector columns to expenses table for semantic search
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS notes_embedding vector(1536);
ALTER TABLE expenses ADD COLUMN IF NOT EXISTS category_embedding vector(1536);

-- Create indexes for vector similarity search
CREATE INDEX IF NOT EXISTS idx_expenses_notes_embedding ON expenses USING ivfflat (notes_embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX IF NOT EXISTS idx_expenses_category_embedding ON expenses USING ivfflat (category_embedding vector_cosine_ops) WITH (lists = 100);

-- Create a function to update embeddings when notes or category changes
CREATE OR REPLACE FUNCTION update_expense_embeddings()
RETURNS TRIGGER AS $$
BEGIN
    -- This function will be called by the application to update embeddings
    -- The actual embedding generation will be done in the application layer
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to update embeddings when expense is modified
CREATE TRIGGER update_expense_embeddings_trigger
    BEFORE INSERT OR UPDATE ON expenses
    FOR EACH ROW
    EXECUTE FUNCTION update_expense_embeddings();

-- Create a function for semantic similarity search
CREATE OR REPLACE FUNCTION search_expenses_by_similarity(
    user_id_param INTEGER,
    query_embedding vector(1536),
    similarity_threshold FLOAT DEFAULT 0.7,
    limit_count INTEGER DEFAULT 10
)
RETURNS TABLE (
    id INTEGER,
    total_price FLOAT,
    notes TEXT,
    timestamp TIMESTAMPTZ,
    category_name TEXT,
    similarity FLOAT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        e.id,
        e.total_price,
        e.notes,
        e.timestamp,
        c.name as category_name,
        1 - (e.notes_embedding <=> query_embedding) as similarity
    FROM expenses e
    JOIN categories c ON e.category_id = c.id
    WHERE e.user_id = user_id_param 
        AND e.deleted_at IS NULL
        AND e.notes_embedding IS NOT NULL
        AND 1 - (e.notes_embedding <=> query_embedding) > similarity_threshold
    ORDER BY e.notes_embedding <=> query_embedding
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;

-- Create a function for category-based similarity search
CREATE OR REPLACE FUNCTION find_similar_expenses(
    expense_id_param INTEGER,
    similarity_threshold FLOAT DEFAULT 0.8,
    limit_count INTEGER DEFAULT 5
)
RETURNS TABLE (
    id INTEGER,
    total_price FLOAT,
    notes TEXT,
    timestamp TIMESTAMPTZ,
    category_name TEXT,
    similarity FLOAT
) AS $$
DECLARE
    target_embedding vector(1536);
    target_user_id INTEGER;
BEGIN
    -- Get the embedding and user_id of the target expense
    SELECT notes_embedding, user_id INTO target_embedding, target_user_id
    FROM expenses 
    WHERE id = expense_id_param AND deleted_at IS NULL;
    
    IF target_embedding IS NULL THEN
        RETURN;
    END IF;
    
    RETURN QUERY
    SELECT 
        e.id,
        e.total_price,
        e.notes,
        e.timestamp,
        c.name as category_name,
        1 - (e.notes_embedding <=> target_embedding) as similarity
    FROM expenses e
    JOIN categories c ON e.category_id = c.id
    WHERE e.user_id = target_user_id 
        AND e.id != expense_id_param
        AND e.deleted_at IS NULL
        AND e.notes_embedding IS NOT NULL
        AND 1 - (e.notes_embedding <=> target_embedding) > similarity_threshold
    ORDER BY e.notes_embedding <=> target_embedding
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql; 