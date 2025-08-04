-- 001_init_core_schema.up.sql  (CockroachDB compatible, sin triggers)

-- ---------- ENUMS ----------
CREATE TYPE IF NOT EXISTS rating_grade AS ENUM
  ('Strong Buy','Buy','Hold','Sell','Strong Sell');

CREATE TYPE IF NOT EXISTS ingest_status AS ENUM
  ('running','completed','failed');

-- ---------- BROKERS ----------
CREATE TABLE IF NOT EXISTS brokers (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name              STRING NOT NULL UNIQUE,
  credibility_score DECIMAL(3,2) NOT NULL DEFAULT 0.60
    CHECK (credibility_score >= 0 AND credibility_score <= 1),
  created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at        TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now()
);

-- ---------- STOCKS ----------
CREATE TABLE IF NOT EXISTS stocks (
  id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  ticker       STRING NOT NULL,
  company      STRING NOT NULL,
  broker_id    UUID REFERENCES brokers(id) ON DELETE SET NULL,
  action       STRING NOT NULL,
  rating_from  STRING,
  rating_to    STRING,
  target_from  DECIMAL(10,2),
  target_to    DECIMAL(10,2),
  event_time   TIMESTAMPTZ NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now(),
  CONSTRAINT   uc_ticker_event UNIQUE (ticker,event_time)
);

CREATE INDEX IF NOT EXISTS idx_stocks__ticker            ON stocks (ticker);
CREATE INDEX IF NOT EXISTS idx_stocks__event_time_desc   ON stocks (event_time DESC);
CREATE INDEX IF NOT EXISTS idx_stocks__ticker_time_desc  ON stocks (ticker,event_time DESC);
CREATE INDEX IF NOT EXISTS idx_stocks__broker_id_time    ON stocks (broker_id,event_time DESC)
  WHERE broker_id IS NOT NULL;

-- ---------- INGESTION LOGS ----------
CREATE TABLE IF NOT EXISTS ingestion_logs (
  id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  batch_id           STRING NOT NULL,
  total_records      INT  NOT NULL CHECK (total_records      >= 0),
  successful_records INT  NOT NULL CHECK (successful_records >= 0),
  failed_records     INT  NOT NULL CHECK (failed_records     >= 0),
  status             ingest_status NOT NULL,
  error_details      JSONB,
  started_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  completed_at       TIMESTAMPTZ,
  created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at         TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE INDEX IF NOT EXISTS idx_ingestion_logs__batch_id   ON ingestion_logs (batch_id);
CREATE INDEX IF NOT EXISTS idx_ingestion_logs__status     ON ingestion_logs (status);
CREATE INDEX IF NOT EXISTS idx_ingestion_logs__started_at ON ingestion_logs (started_at DESC);
