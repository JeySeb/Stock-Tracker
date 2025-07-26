-- Migration 005: Update recommendations table schema
-- Purpose: Add missing columns to support current Recommendation model structure

-- Add missing columns to recommendations table
ALTER TABLE recommendations ADD COLUMN IF NOT EXISTS company_name STRING;
ALTER TABLE recommendations ADD COLUMN IF NOT EXISTS tier STRING NOT NULL DEFAULT 'basic';
ALTER TABLE recommendations ADD COLUMN IF NOT EXISTS basic_factors JSONB;
ALTER TABLE recommendations ADD COLUMN IF NOT EXISTS external_data JSONB;
ALTER TABLE recommendations ADD COLUMN IF NOT EXISTS ai_insights JSONB;

-- Add constraints for tier enum values
ALTER TABLE recommendations ADD CONSTRAINT chk_recommendations_tier 
CHECK (tier IN ('basic', 'enriched', 'premium'));

-- Add constraint for recommendation_type enum values
ALTER TABLE recommendations DROP CONSTRAINT IF EXISTS chk_recommendations_type;
ALTER TABLE recommendations ADD CONSTRAINT chk_recommendations_type 
CHECK (recommendation_type IN ('Strong Buy', 'Buy', 'Hold', 'Sell', 'Strong Sell'));

-- Create additional indexes for new columns
CREATE INDEX IF NOT EXISTS idx_recommendations_tier ON recommendations(tier);
CREATE INDEX IF NOT EXISTS idx_recommendations_company_name ON recommendations(company_name);
CREATE INDEX IF NOT EXISTS idx_recommendations_ticker_tier ON recommendations(ticker, tier);

-- Update existing records to have default values
UPDATE recommendations 
SET 
    company_name = ticker,  -- Fallback: use ticker as company name
    basic_factors = '[]'::JSONB  -- Empty array for existing records
WHERE company_name IS NULL OR basic_factors IS NULL; 