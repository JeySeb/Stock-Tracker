-- 007_market_data_ingestion.up.sql
-- Market data ingestion from external APIs (Yahoo Finance, Alpha Vantage, etc.)
-- Hourly data collection for each ticker

-- ---------- MARKET DATA ENUMS ----------
CREATE TYPE IF NOT EXISTS data_source AS ENUM
  ('yahoo_finance', 'alpha_vantage', 'manual');

CREATE TYPE IF NOT EXISTS data_quality AS ENUM
  ('excellent', 'good', 'fair', 'poor');

-- ---------- MARKET DATA TABLE ----------
CREATE TABLE IF NOT EXISTS market_data (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticker                STRING NOT NULL,
  data_source           data_source NOT NULL,
  data_quality          data_quality NOT NULL DEFAULT 'good',
  
  -- Price data
  current_price         DECIMAL(10,4) NOT NULL CHECK (current_price > 0),
  day_change            DECIMAL(10,4) NOT NULL,
  day_change_percent    DECIMAL(8,4) NOT NULL,
  
  -- Volume and market metrics
  volume                BIGINT NOT NULL CHECK (volume >= 0),
  market_cap            BIGINT,
  
  -- Fundamental ratios (nullable)
  pe_ratio              DECIMAL(10,4),
  dividend_yield        DECIMAL(8,4),
  
  -- Technical levels
  week_52_high          DECIMAL(10,4),
  week_52_low           DECIMAL(10,4),
  avg_volume            BIGINT,
  
  -- Metadata
  collected_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  data_timestamp        TIMESTAMPTZ NOT NULL, -- When the data was actually collected from source
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
  
  -- Constraints
  CONSTRAINT uc_ticker_source_timestamp UNIQUE (ticker, data_source, data_timestamp),
  CONSTRAINT chk_day_change_percent CHECK (day_change_percent BETWEEN -100 AND 1000),
  CONSTRAINT chk_pe_ratio CHECK (pe_ratio IS NULL OR pe_ratio > 0),
  CONSTRAINT chk_dividend_yield CHECK (dividend_yield IS NULL OR dividend_yield >= 0),
  CONSTRAINT chk_week_52_high CHECK (week_52_high IS NULL OR week_52_high > 0),
  CONSTRAINT chk_week_52_low CHECK (week_52_low IS NULL OR week_52_low > 0),
  CONSTRAINT chk_avg_volume CHECK (avg_volume IS NULL OR avg_volume >= 0)
);

-- ---------- INDEXES FOR PERFORMANCE ----------
CREATE INDEX IF NOT EXISTS idx_market_data__ticker_timestamp 
  ON market_data (ticker, data_timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_market_data__source_timestamp 
  ON market_data (data_source, data_timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_market_data__collected_at 
  ON market_data (collected_at DESC);

CREATE INDEX IF NOT EXISTS idx_market_data__quality_timestamp 
  ON market_data (data_quality, data_timestamp DESC);

-- ---------- MARKET DATA INGESTION LOGS ----------
CREATE TABLE IF NOT EXISTS market_data_ingestion_logs (
  id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  batch_id              STRING NOT NULL,
  data_source           data_source NOT NULL,
  total_tickers         INT NOT NULL CHECK (total_tickers >= 0),
  successful_tickers    INT NOT NULL CHECK (successful_tickers >= 0),
  failed_tickers        INT NOT NULL CHECK (failed_tickers >= 0),
  skipped_tickers       INT NOT NULL CHECK (skipped_tickers >= 0),
  status                ingest_status NOT NULL,
  error_details         JSONB,
  started_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at          TIMESTAMPTZ,
  created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at            TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
  
  CONSTRAINT chk_successful_failed CHECK (successful_tickers + failed_tickers <= total_tickers)
);

CREATE INDEX IF NOT EXISTS idx_market_ingestion__batch_id 
  ON market_data_ingestion_logs (batch_id);

CREATE INDEX IF NOT EXISTS idx_market_ingestion__source_status 
  ON market_data_ingestion_logs (data_source, status);

CREATE INDEX IF NOT EXISTS idx_market_ingestion__started_at 
  ON market_data_ingestion_logs (started_at DESC);

-- ---------- MARKET DATA STATS VIEW ----------
CREATE VIEW IF NOT EXISTS market_data_stats AS
SELECT 
  ticker,
  COUNT(*) AS total_records,
  COUNT(*) FILTER (WHERE data_source = 'yahoo_finance') AS yahoo_records,
  COUNT(*) FILTER (WHERE data_source = 'alpha_vantage') AS alpha_records,
  MAX(data_timestamp) AS latest_data_timestamp,
  MAX(collected_at) AS latest_collection,
  AVG(current_price) AS avg_price,
  AVG(ABS(day_change_percent)) AS avg_volatility,
  COUNT(*) FILTER (WHERE data_quality = 'excellent') AS excellent_quality,
  COUNT(*) FILTER (WHERE data_quality = 'good') AS good_quality,
  COUNT(*) FILTER (WHERE data_quality = 'fair') AS fair_quality,
  COUNT(*) FILTER (WHERE data_quality = 'poor') AS poor_quality
FROM market_data
WHERE data_timestamp > now() - INTERVAL '30 days'
GROUP BY ticker
HAVING COUNT(*) >= 1;

-- ---------- MARKET DATA MONITORING VIEW ----------
CREATE VIEW IF NOT EXISTS market_data_monitor AS
SELECT 
  'total_market_records' AS metric,
  COUNT(*)::STRING AS value,
  'Total market data records' AS description
FROM market_data
WHERE data_timestamp > now() - INTERVAL '24 hours'

UNION ALL

SELECT 
  'unique_tickers_today' AS metric,
  COUNT(DISTINCT ticker)::STRING AS value,
  'Unique tickers with data today' AS description
FROM market_data
WHERE data_timestamp > now() - INTERVAL '24 hours'

UNION ALL

SELECT 
  'yahoo_finance_records' AS metric,
  COUNT(*)::STRING AS value,
  'Yahoo Finance records today' AS description
FROM market_data
WHERE data_source = 'yahoo_finance' 
  AND data_timestamp > now() - INTERVAL '24 hours'

UNION ALL

SELECT 
  'alpha_vantage_records' AS metric,
  COUNT(*)::STRING AS value,
  'Alpha Vantage records today' AS description
FROM market_data
WHERE data_source = 'alpha_vantage' 
  AND data_timestamp > now() - INTERVAL '24 hours'

UNION ALL

SELECT 
  'avg_data_quality' AS metric,
  ROUND(AVG(
    CASE data_quality
      WHEN 'excellent' THEN 1.0
      WHEN 'good' THEN 0.8
      WHEN 'fair' THEN 0.6
      WHEN 'poor' THEN 0.3
    END
  ), 3)::STRING AS value,
  'Average data quality score' AS description
FROM market_data
WHERE data_timestamp > now() - INTERVAL '24 hours'; 