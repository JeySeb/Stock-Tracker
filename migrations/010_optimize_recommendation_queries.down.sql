-- Rollback migration: Remove recommendation system optimizations
-- Purpose: Clean rollback of all recommendation system database optimizations

-- ============================================
-- 1. DROP FUNCTION
-- ============================================

DROP FUNCTION IF EXISTS refresh_recommendation_stats();

-- ============================================
-- 2. DROP VIEWS
-- ============================================

DROP VIEW IF EXISTS recommendation_query_stats;

-- ============================================
-- 3. DROP MATERIALIZED VIEWS
-- ============================================

DROP MATERIALIZED VIEW IF EXISTS ticker_stats;
DROP MATERIALIZED VIEW IF EXISTS brokerage_stats;

-- ============================================
-- 4. DROP INDEXES
-- ============================================

-- Drop recommendation-specific indexes
DROP INDEX IF EXISTS idx_stocks_event_time_desc;
DROP INDEX IF EXISTS idx_stocks_ticker_event_time;
DROP INDEX IF EXISTS idx_stocks_brokerage_event_time;
DROP INDEX IF EXISTS idx_stocks_recent_activity;
DROP INDEX IF EXISTS idx_stocks_targets;
DROP INDEX IF EXISTS idx_stocks_ratings;

-- Note: We don't drop indexes that might be used by other parts of the system 