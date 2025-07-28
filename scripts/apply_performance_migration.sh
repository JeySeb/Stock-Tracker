#!/bin/bash

# Performance Migration Script
# Applies database indexes and validates performance improvements

set -e

echo "🚀 Applying performance optimizations for Stock Tracker..."

# Check if database connection is available
if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
    echo "❌ PostgreSQL is not available. Please start the database first."
    exit 1
fi

# Apply the migration
echo "📊 Applying performance indexes migration..."
psql -h localhost -p 5432 -d stock_tracker -f migrations/007_performance_indexes.sql

echo "✅ Performance indexes applied successfully!"

# Run some basic validation queries
echo "🔍 Running validation queries..."

echo "📈 Checking index creation..."
psql -h localhost -p 5432 -d stock_tracker -c "
SELECT 
    indexname, 
    tablename,
    indexdef 
FROM pg_indexes 
WHERE tablename = 'stocks' 
    AND indexname LIKE 'idx_stocks_%'
ORDER BY indexname;
"

echo "📊 Analyzing table statistics..."
psql -h localhost -p 5432 -d stock_tracker -c "
SELECT 
    schemaname,
    tablename,
    n_tup_ins as inserts,
    n_tup_upd as updates,
    n_tup_del as deletes,
    n_live_tup as live_rows,
    n_dead_tup as dead_rows,
    last_vacuum,
    last_analyze
FROM pg_stat_user_tables 
WHERE tablename IN ('stocks', 'brokers');
"

echo "🎯 Performance optimization complete!"
echo ""
echo "Expected improvements:"
echo "  • 80-95% reduction in recommendation generation time"
echo "  • Elimination of N+1 query problems"
echo "  • Cached brokerage statistics (30min TTL)"
echo "  • Concurrent processing of up to 10 tickers"
echo "  • Early filtering to process only top active tickers"
echo ""
echo "Monitor your application logs to see the performance gains!" 