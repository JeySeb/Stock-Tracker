-- Migration: Optimize database for recommendation system (CockroachDB Compatible)
-- Purpose: Add indexes and statistics tables to improve recommendation query performance

-- ============================================
-- 1. INDEXES FOR RECOMMENDATION QUERIES
-- ============================================

-- Index for filtering stocks by event time (most common query pattern)
CREATE INDEX IF NOT EXISTS idx_stocks_event_time_desc 
ON stocks(event_time DESC);

-- Composite index for ticker + event_time (used in GetByTicker with date filtering)
CREATE INDEX IF NOT EXISTS idx_stocks_ticker_event_time 
ON stocks(ticker, event_time DESC);

-- Index for brokerage analysis (used in broker frequency scoring)
CREATE INDEX IF NOT EXISTS idx_stocks_brokerage_event_time 
ON stocks(brokerage, event_time DESC);

-- Index for recent stocks query optimization
CREATE INDEX IF NOT EXISTS idx_stocks_recent_activity 
ON stocks(event_time DESC, ticker) 
WHERE event_time > NOW() - INTERVAL '90 days';

-- Index for target price analysis
CREATE INDEX IF NOT EXISTS idx_stocks_targets 
ON stocks(ticker, target_from, target_to, event_time DESC) 
WHERE target_from > 0 AND target_to > 0;

-- Index for rating analysis
CREATE INDEX IF NOT EXISTS idx_stocks_ratings 
ON stocks(ticker, rating_from, rating_to, event_time DESC) 
WHERE rating_from != '' AND rating_to != '';

-- ============================================
-- 2. BROKERAGE STATISTICS TABLE (replaces materialized view)
-- ============================================

-- Create brokerage statistics table
CREATE TABLE IF NOT EXISTS brokerage_stats (
    brokerage STRING PRIMARY KEY,
    total_reports INT NOT NULL,
    unique_tickers INT NOT NULL,
    positive_ratio DECIMAL(5,3) NOT NULL,
    reports_last_30d INT NOT NULL,
    reports_last_90d INT NOT NULL,
    last_activity TIMESTAMPTZ NOT NULL,
    first_activity TIMESTAMPTZ NOT NULL,
    avg_target_change_pct DECIMAL(8,4),
    stats_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_brokerage_stats_activity 
ON brokerage_stats(total_reports DESC, last_activity DESC);

-- ============================================
-- 3. TICKER STATISTICS TABLE (replaces materialized view)
-- ============================================

-- Create ticker statistics table
CREATE TABLE IF NOT EXISTS ticker_stats (
    ticker STRING PRIMARY KEY,
    company STRING,
    total_events INT NOT NULL,
    positive_events INT NOT NULL,
    negative_events INT NOT NULL,
    avg_target_change DECIMAL(8,4),
    last_event_time TIMESTAMPTZ NOT NULL,
    latest_target_price DECIMAL(10,2),
    unique_brokers INT NOT NULL,
    events_last_30d INT NOT NULL,
    events_last_90d INT NOT NULL,
    stats_updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_ticker_stats_activity 
ON ticker_stats(total_events DESC, last_event_time DESC);

CREATE INDEX IF NOT EXISTS idx_ticker_stats_recent 
ON ticker_stats(events_last_90d DESC, last_event_time DESC);

-- ============================================
-- 4. PERFORMANCE MONITORING VIEW
-- ============================================

-- Create a view to monitor query performance
CREATE OR REPLACE VIEW recommendation_query_stats AS
SELECT 
    'stocks_total' as metric,
    COUNT(*)::STRING as value,
    'Total stock events in database' as description
FROM stocks

UNION ALL

SELECT 
    'stocks_recent_90d' as metric,
    COUNT(*)::STRING as value,
    'Stock events in last 90 days' as description
FROM stocks 
WHERE event_time > NOW() - INTERVAL '90 days'

UNION ALL

SELECT 
    'unique_tickers' as metric,
    COUNT(DISTINCT ticker)::STRING as value,
    'Unique tickers with events' as description
FROM stocks

UNION ALL

SELECT 
    'unique_brokers' as metric,
    COUNT(DISTINCT brokerage)::STRING as value,
    'Unique brokerages' as description
FROM stocks

UNION ALL

SELECT 
    'avg_events_per_ticker' as metric,
    ROUND(AVG(event_count::DECIMAL), 2)::STRING as value,
    'Average events per ticker' as description
FROM (
    SELECT COUNT(*) as event_count 
    FROM stocks 
    GROUP BY ticker
) ticker_counts;

-- ============================================
-- 5. FUNCTION TO REFRESH STATISTICS TABLES
-- ============================================

-- Create function to refresh all recommendation-related statistics
CREATE OR REPLACE FUNCTION refresh_recommendation_stats()
RETURNS STRING
LANGUAGE SQL
AS $$
    -- Clear existing stats
    DELETE FROM brokerage_stats;
    DELETE FROM ticker_stats;
    
    -- Populate brokerage stats
    INSERT INTO brokerage_stats (
        brokerage, total_reports, unique_tickers, positive_ratio,
        reports_last_30d, reports_last_90d, last_activity, first_activity,
        avg_target_change_pct, stats_updated_at
    )
    SELECT 
        brokerage,
        COUNT(*) as total_reports,
        COUNT(DISTINCT ticker) as unique_tickers,
        
        -- Calculate positive action ratio
        ROUND(
            AVG(CASE 
                WHEN LOWER(action) LIKE '%upgrade%' OR 
                     LOWER(action) LIKE '%raised%' OR 
                     LOWER(action) LIKE '%initiated%' OR
                     LOWER(action) LIKE '%outperform%' OR
                     LOWER(action) LIKE '%buy%' OR
                     target_to > target_from OR
                     (rating_to = 'Buy' OR rating_to = 'Strong Buy' OR rating_to = 'Outperform')
                THEN 1.0 
                ELSE 0.0 
            END)::DECIMAL, 3
        ) as positive_ratio,
        
        -- Recent activity metrics
        COUNT(CASE WHEN event_time > NOW() - INTERVAL '30 days' THEN 1 END) as reports_last_30d,
        COUNT(CASE WHEN event_time > NOW() - INTERVAL '90 days' THEN 1 END) as reports_last_90d,
        
        -- Time metrics
        MAX(event_time) as last_activity,
        MIN(event_time) as first_activity,
        
        -- Average target change percentage
        ROUND(
            AVG(
                CASE 
                    WHEN target_from > 0 AND target_to > 0 
                    THEN ((target_to - target_from) / target_from) 
                    ELSE NULL 
                END
            )::DECIMAL, 4
        ) as avg_target_change_pct,
        
        NOW() as stats_updated_at
        
    FROM stocks 
    WHERE event_time > NOW() - INTERVAL '2 years'  -- Only last 2 years for relevance
    GROUP BY brokerage
    HAVING COUNT(*) >= 5;  -- Only brokers with at least 5 reports

    -- Populate ticker stats
    INSERT INTO ticker_stats (
        ticker, company, total_events, positive_events, negative_events,
        avg_target_change, last_event_time, latest_target_price,
        unique_brokers, events_last_30d, events_last_90d, stats_updated_at
    )
    SELECT 
        ticker,
        company,
        COUNT(*) as total_events,
        
        -- Event type analysis
        COUNT(CASE 
            WHEN LOWER(action) LIKE '%upgrade%' OR 
                 LOWER(action) LIKE '%raised%' OR 
                 LOWER(action) LIKE '%initiated%' OR
                 target_to > target_from OR
                 (rating_to IN ('Buy', 'Strong Buy', 'Outperform') AND 
                  rating_from NOT IN ('Buy', 'Strong Buy', 'Outperform'))
            THEN 1 
        END) as positive_events,
        
        COUNT(CASE 
            WHEN LOWER(action) LIKE '%downgrade%' OR 
                 LOWER(action) LIKE '%lowered%' OR 
                 LOWER(action) LIKE '%reduced%' OR
                 target_to < target_from OR
                 (rating_to IN ('Sell', 'Strong Sell', 'Underperform') AND 
                  rating_from NOT IN ('Sell', 'Strong Sell', 'Underperform'))
            THEN 1 
        END) as negative_events,
        
        -- Price target metrics
        AVG(
            CASE 
                WHEN target_from > 0 AND target_to > 0 
                THEN ((target_to - target_from) / target_from) 
                ELSE NULL 
            END
        ) as avg_target_change,
        
        -- Latest metrics
        MAX(event_time) as last_event_time,
        
        -- Get latest target price (using subquery approach for CockroachDB)
        (SELECT target_to 
         FROM stocks s2 
         WHERE s2.ticker = stocks.ticker 
           AND s2.target_to > 0 
         ORDER BY s2.event_time DESC 
         LIMIT 1) as latest_target_price,
        
        -- Broker diversity
        COUNT(DISTINCT brokerage) as unique_brokers,
        
        -- Recent activity
        COUNT(CASE WHEN event_time > NOW() - INTERVAL '30 days' THEN 1 END) as events_last_30d,
        COUNT(CASE WHEN event_time > NOW() - INTERVAL '90 days' THEN 1 END) as events_last_90d,
        
        NOW() as stats_updated_at
        
    FROM stocks
    WHERE event_time > NOW() - INTERVAL '1 year'  -- Only last year for current relevance
    GROUP BY ticker, company
    HAVING COUNT(*) >= 2;  -- Only tickers with at least 2 events

    SELECT 'Recommendation statistics refreshed successfully';
$$;

-- ============================================
-- 6. COMMENTS FOR DOCUMENTATION
-- ============================================

COMMENT ON TABLE brokerage_stats IS 
'Aggregated statistics for brokerages used in recommendation scoring. Refreshed periodically.';

COMMENT ON TABLE ticker_stats IS 
'Aggregated statistics for tickers used in recommendation scoring. Refreshed periodically.';

COMMENT ON FUNCTION refresh_recommendation_stats() IS 
'Refreshes all statistics tables used by the recommendation system. Should be called periodically.';

-- ============================================
-- 7. INITIAL REFRESH
-- ============================================

-- Refresh the statistics tables after creation
SELECT refresh_recommendation_stats(); 