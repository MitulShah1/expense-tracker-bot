# pgvector Integration for Expense Tracker Bot

This document describes the integration of pgvector (PostgreSQL vector extension) into the expense tracker bot to enable semantic search capabilities.

## Overview

pgvector allows us to store and query vector embeddings in PostgreSQL, enabling semantic search functionality. This integration enhances the expense tracker bot with:

- **Semantic Search**: Find expenses using natural language queries
- **Similarity Matching**: Find similar expenses based on descriptions
- **Smart Categorization**: Automatically suggest categories based on expense descriptions
- **Pattern Recognition**: Identify spending patterns across users

## Architecture

### Database Schema

The integration adds vector columns to the existing `expenses` table:

```sql
-- Vector columns for semantic search
ALTER TABLE expenses ADD COLUMN notes_embedding vector(1536);
ALTER TABLE expenses ADD COLUMN category_embedding vector(1536);

-- Indexes for efficient similarity search
CREATE INDEX idx_expenses_notes_embedding ON expenses 
USING ivfflat (notes_embedding vector_cosine_ops) WITH (lists = 100);
```

### Key Components

1. **Migration (005_add_pgvector.sql)**: Sets up pgvector extension and database functions
2. **VectorSearchStorage Interface**: Defines vector search operations
3. **VectorService**: Business logic for embedding generation and search
4. **Bot Integration**: New `/search` command for semantic search

## Features

### 1. Semantic Search

Users can search for expenses using natural language:

```sh
/search
"Find all fuel expenses from last month"
"Show me expensive car repairs"
"Find expenses related to maintenance"
```

### 2. Similarity Matching

Find expenses similar to a given expense:

```go
similarExpenses, err := vectorService.FindSimilarExpenses(ctx, expenseID, 0.8, 5)
```

### 3. Automatic Embedding Generation

When expenses are created or updated, embeddings are automatically generated:

```go
err := vectorService.UpdateExpenseEmbeddings(ctx, expenseID)
```

## Implementation Details

### Embedding Generation

The current implementation uses a placeholder embedding generator. In production, you should:

1. **Use a real embedding service**:
   - OpenAI's `text-embedding-ada-002`
   - Cohere's `embed-multilingual-v3.0`
   - Hugging Face's sentence-transformers

2. **Implement caching** to avoid repeated API calls
3. **Handle rate limiting** and errors gracefully

### Database Functions

The migration creates several PostgreSQL functions:

```sql
-- Semantic similarity search
CREATE OR REPLACE FUNCTION search_expenses_by_similarity(
    user_id_param INTEGER,
    query_embedding vector(1536),
    similarity_threshold FLOAT DEFAULT 0.7,
    limit_count INTEGER DEFAULT 10
)

-- Find similar expenses
CREATE OR REPLACE FUNCTION find_similar_expenses(
    expense_id_param INTEGER,
    similarity_threshold FLOAT DEFAULT 0.8,
    limit_count INTEGER DEFAULT 5
)
```

### Vector Dimensions

The implementation uses 1536-dimensional vectors, which is compatible with:

- OpenAI's text-embedding-ada-002 model
- Most modern embedding models

## Usage Examples

### Basic Search

```go
// Search for expenses using natural language
expenses, err := vectorService.SearchExpensesByQuery(
    ctx, 
    telegramID, 
    "fuel expenses from last month", 
    0.7, 
    10
)
```

### Batch Processing

```go
// Update embeddings for all user expenses
err := vectorService.BatchUpdateEmbeddings(ctx, telegramID)
```

### Similarity Search

```go
// Find expenses similar to a specific expense
similarExpenses, err := vectorService.FindSimilarExpenses(
    ctx, 
    expenseID, 
    0.8, 
    5
)
```

## Configuration

### Docker Setup

The `docker/docker-compose.yml` has been updated to use the pgvector-enabled PostgreSQL image:

```yaml
postgres:
  image: pgvector/pgvector:pg15
  # ... other configuration
```

### Environment Variables

No additional environment variables are required for basic functionality. For production embedding services, you'll need:

```bash
# OpenAI (recommended)
OPENAI_API_KEY=your_openai_api_key

# Cohere (alternative)
COHERE_API_KEY=your_cohere_api_key
```

## Performance Considerations

### Indexing

The implementation uses IVFFlat indexes for efficient similarity search:

```sql
CREATE INDEX idx_expenses_notes_embedding ON expenses 
USING ivfflat (notes_embedding vector_cosine_ops) WITH (lists = 100);
```

### Batch Operations

For large datasets, use batch operations:

```go
// Update embeddings in batches
err := vectorService.BatchUpdateEmbeddings(ctx, telegramID)
```

### Caching

Consider implementing embedding caching to avoid repeated API calls:

```go
// Cache embeddings in Redis or similar
cacheKey := fmt.Sprintf("embedding:%s", textHash)
cachedEmbedding := cache.Get(cacheKey)
```

## Security Considerations

1. **API Key Management**: Store embedding service API keys securely
2. **Rate Limiting**: Implement rate limiting for embedding generation
3. **Data Privacy**: Ensure expense data is handled according to privacy policies
4. **Input Validation**: Validate and sanitize search queries

## Monitoring and Metrics

### Key Metrics to Track

1. **Search Performance**: Query response times
2. **Embedding Generation**: Success/failure rates
3. **Cache Hit Rates**: If implementing caching
4. **API Usage**: Embedding service API calls

### Logging

The implementation includes comprehensive logging:

```go
s.logger.Info(ctx, "Expense embeddings updated successfully",
    logger.Int("expense_id", int(expenseID)))
```

## Future Enhancements

### 1. Advanced Search Features

- **Date Range Filtering**: Combine semantic search with date filters
- **Category Filtering**: Search within specific categories
- **Amount Range Filtering**: Search by expense amount ranges

### 2. Machine Learning Integration

- **Automatic Categorization**: Suggest categories based on descriptions
- **Spending Pattern Detection**: Identify unusual spending patterns
- **Budget Recommendations**: Suggest budget adjustments based on patterns

### 3. Multi-language Support

- **Multilingual Embeddings**: Support for multiple languages
- **Translation**: Automatic translation of search queries

### 4. Real-time Updates

- **Webhook Integration**: Real-time embedding updates
- **Background Processing**: Async embedding generation

## Troubleshooting

### Common Issues

1. **pgvector Extension Not Available**

   ```bash
   # Ensure you're using the correct PostgreSQL image
   image: pgvector/pgvector:pg15
   ```

2. **Embedding Generation Fails**
   - Check API key configuration
   - Verify network connectivity to embedding service
   - Check rate limits

3. **Slow Search Performance**
   - Verify indexes are created correctly
   - Consider adjusting similarity thresholds
   - Monitor database performance

### Debug Commands

```sql
-- Check if pgvector is installed
SELECT * FROM pg_extension WHERE extname = 'vector';

-- Check vector columns
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'expenses' 
AND data_type LIKE '%vector%';

-- Test similarity search
SELECT id, notes, 1 - (notes_embedding <=> '[0.1,0.2,...]'::vector) as similarity
FROM expenses 
WHERE notes_embedding IS NOT NULL
ORDER BY notes_embedding <=> '[0.1,0.2,...]'::vector
LIMIT 5;
```

## Conclusion

The pgvector integration significantly enhances the expense tracker bot's search capabilities, making it more user-friendly and intelligent. The implementation is designed to be scalable and can be easily extended with additional features as needed.

For production deployment, ensure you:

1. Use a proper embedding service (OpenAI, Cohere, etc.)
2. Implement proper error handling and retry logic
3. Set up monitoring and alerting
4. Consider implementing caching for better performance
5. Follow security best practices for API key management
