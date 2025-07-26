-- Migration 005 DOWN: Reverse recommendations table schema changes
-- Purpose: Remove added columns and constraints

-- Drop indexes
DROP INDEX IF EXISTS idx_recommendations_tier;
DROP INDEX IF EXISTS idx_recommendations_company_name;
DROP INDEX IF EXISTS idx_recommendations_ticker_tier;

-- Drop constraints
ALTER TABLE recommendations DROP CONSTRAINT IF EXISTS chk_recommendations_tier;
ALTER TABLE recommendations DROP CONSTRAINT IF EXISTS chk_recommendations_type;

-- Remove added columns
ALTER TABLE recommendations DROP COLUMN IF EXISTS company_name;
ALTER TABLE recommendations DROP COLUMN IF EXISTS tier;
ALTER TABLE recommendations DROP COLUMN IF EXISTS basic_factors;
ALTER TABLE recommendations DROP COLUMN IF EXISTS external_data;
ALTER TABLE recommendations DROP COLUMN IF EXISTS ai_insights; 