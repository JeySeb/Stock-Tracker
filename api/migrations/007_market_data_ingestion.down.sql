-- 007_market_data_ingestion.down.sql
-- Rollback market data ingestion schema

-- Drop views first (dependencies)
DROP VIEW IF EXISTS market_data_monitor;
DROP VIEW IF EXISTS market_data_stats;

-- Drop tables
DROP TABLE IF EXISTS market_data_ingestion_logs;
DROP TABLE IF EXISTS market_data;

-- Drop enums
DROP TYPE IF EXISTS data_quality;
DROP TYPE IF EXISTS data_source; 