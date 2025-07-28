-- 
-- Tablas analíticas y vista de monitorización

-- ---------- BROKERAGE STATS ----------
CREATE TABLE IF NOT EXISTS brokerage_stats (
  brokerage            STRING PRIMARY KEY,
  total_reports        INT NOT NULL,
  unique_tickers       INT NOT NULL,
  positive_ratio       DECIMAL(5,3) NOT NULL,
  reports_last_30d     INT NOT NULL,
  reports_last_90d     INT NOT NULL,
  last_activity        TIMESTAMPTZ NOT NULL,
  first_activity       TIMESTAMPTZ NOT NULL,
  avg_target_change_pct DECIMAL(8,4),
  stats_updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_brokerage_stats__activity
  ON brokerage_stats (total_reports DESC, last_activity DESC);

-- ---------- TICKER STATS ----------
CREATE TABLE IF NOT EXISTS ticker_stats (
  ticker               STRING PRIMARY KEY,
  company              STRING,
  total_events         INT NOT NULL,
  positive_events      INT NOT NULL,
  negative_events      INT NOT NULL,
  avg_target_change    DECIMAL(8,4),
  last_event_time      TIMESTAMPTZ NOT NULL,
  latest_target_price  DECIMAL(10,2),
  unique_brokers       INT NOT NULL,
  events_last_30d      INT NOT NULL,
  events_last_90d      INT NOT NULL,
  stats_updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ticker_stats__activity
  ON ticker_stats (total_events DESC, last_event_time DESC);

CREATE INDEX IF NOT EXISTS idx_ticker_stats__recent
  ON ticker_stats (events_last_90d DESC, last_event_time DESC);

-- ---------- MONITOR VIEW ----------
CREATE VIEW IF NOT EXISTS recommendation_query_stats AS
WITH base AS (
  SELECT COUNT(*) AS total, COUNT(DISTINCT ticker) AS unique_tickers FROM stocks
)
SELECT 'stocks_total'          AS metric, total::STRING           AS value, 'Total stock events'                         AS description FROM base
UNION ALL
SELECT 'unique_tickers'        AS metric, unique_tickers::STRING  AS value, 'Unique tickers'                              AS description FROM base
UNION ALL
SELECT 'stocks_recent_90d'     AS metric, COUNT(*)::STRING        AS value, 'Events last 90 days'                         AS description FROM stocks WHERE event_time > now() - INTERVAL '90 days'
UNION ALL
SELECT 'unique_brokers'        AS metric, COUNT(DISTINCT brokerage)::STRING AS value, 'Unique brokerages'                  AS description FROM stocks WHERE brokerage IS NOT NULL
UNION ALL
SELECT 'avg_events_per_ticker' AS metric, ROUND(AVG(cnt),2)::STRING AS value, 'Avg events per ticker' FROM (SELECT COUNT(*) cnt FROM stocks GROUP BY ticker);
