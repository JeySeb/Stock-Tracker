-- ==========================================
-- COMPLETE DATABASE RESET SCRIPT (CockroachDB Compatible)
-- ==========================================
-- This script will drop ALL tables, types, functions, and migration state
-- Use this when migration state is corrupted and you need a fresh start

-- Step 1: Drop all tables in reverse dependency order
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS ingestion_logs;
DROP TABLE IF EXISTS recommendations;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS brokers;

-- Step 2: Drop migration tracking table
DROP TABLE IF EXISTS schema_migrations;

-- Step 3: Drop custom types (CockroachDB compatible - no CASCADE)
DROP TYPE IF EXISTS subscription_plan;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS user_tier;

-- Step 4: Drop functions (CockroachDB compatible - no CASCADE)
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Step 5: Drop any remaining indexes
DROP INDEX IF EXISTS unique_ticker_event_time;

-- Verify cleanup
SELECT 'Database reset complete. All tables, types, and functions dropped.' AS status; 