-- 002_recommendations_schema.up.sql
-- Tabla de recomendaciones inicial (sin triggers)

CREATE TABLE IF NOT EXISTS recommendations (
  id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticker              STRING NOT NULL,
  company_name        STRING NOT NULL,
  score               DECIMAL(5,4)  NOT NULL CHECK (score      BETWEEN 0 AND 1),
  confidence          DECIMAL(5,4)  NOT NULL CHECK (confidence BETWEEN 0 AND 1),
  recommendation_type rating_grade  NOT NULL,
  tier                STRING NOT NULL DEFAULT 'basic'
                     CHECK (tier IN ('basic','enriched','premium')),
  factors             JSONB,
  basic_factors       JSONB,
  external_data       JSONB,
  ai_insights         JSONB,
  explanation         STRING,
  created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at          TIMESTAMPTZ,
  updated_at          TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
  CONSTRAINT expire_future CHECK (expires_at IS NULL OR expires_at > created_at)
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_reco__ticker        ON recommendations (ticker);
CREATE INDEX IF NOT EXISTS idx_reco__score_desc    ON recommendations (score DESC);
CREATE INDEX IF NOT EXISTS idx_reco__created_desc  ON recommendations (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_reco__expires       ON recommendations (expires_at);
CREATE INDEX IF NOT EXISTS idx_reco__ticker_tier   ON recommendations (ticker,tier);
CREATE INDEX IF NOT EXISTS idx_reco__company_name  ON recommendations (company_name);
