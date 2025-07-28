-- 006_seed_stats_from_stocks.up.sql
-- Poblado inicial de tablas de estadísticas
-- ¡Ejecuta esto después de haber cargado al menos históricos de stocks!

-- ---------- BROKERAGE ----------
INSERT INTO brokerage_stats (
  brokerage, total_reports, unique_tickers, positive_ratio,
  reports_last_30d, reports_last_90d, last_activity, first_activity,
  avg_target_change_pct
)
SELECT
  brokerage,
  COUNT(*)                              AS total_reports,
  COUNT(DISTINCT ticker)                AS unique_tickers,
  ROUND(AVG(
    CASE WHEN LOWER(action) ~ '(upgrade|raised|initiated|outperform|buy)'
          OR target_to > target_from THEN 1 ELSE 0 END
  )::DECIMAL,3)                         AS positive_ratio,
  COUNT(*) FILTER (WHERE event_time > now() - INTERVAL '30 days') AS reports_last_30d,
  COUNT(*) FILTER (WHERE event_time > now() - INTERVAL '90 days') AS reports_last_90d,
  MAX(event_time)                       AS last_activity,
  MIN(event_time)                       AS first_activity,
  ROUND(AVG(
    CASE WHEN target_from > 0 AND target_to > 0
         THEN (target_to - target_from)/target_from END
  )::DECIMAL,4)                         AS avg_target_change_pct
FROM stocks
WHERE brokerage IS NOT NULL
  AND event_time > now() - INTERVAL '2 years'
GROUP BY brokerage
HAVING COUNT(*) >= 5;

-- ---------- TICKER ----------
INSERT INTO ticker_stats (
  ticker, company, total_events, positive_events, negative_events,
  avg_target_change, last_event_time, latest_target_price,
  unique_brokers, events_last_30d, events_last_90d
)
SELECT
  ticker,
  MAX(company)                    AS company,
  COUNT(*)                              AS total_events,
  SUM(CASE WHEN LOWER(action) ~ '(upgrade|raised|initiated)'
            OR target_to > target_from THEN 1 ELSE 0 END)          AS positive_events,
  SUM(CASE WHEN LOWER(action) ~ '(downgrade|lowered|reduced)'
            OR target_to < target_from THEN 1 ELSE 0 END)          AS negative_events,
  AVG(
    CASE WHEN target_from > 0 AND target_to > 0
         THEN (target_to - target_from)/target_from END
  )                               AS avg_target_change,
  MAX(event_time)                 AS last_event_time,
  (SELECT target_to FROM stocks s2 WHERE s2.ticker = stocks.ticker AND target_to > 0
   ORDER BY event_time DESC LIMIT 1)                           AS latest_target_price,
  COUNT(DISTINCT brokerage)       AS unique_brokers,
  COUNT(*) FILTER (WHERE event_time > now() - INTERVAL '30 days') AS events_last_30d,
  COUNT(*) FILTER (WHERE event_time > now() - INTERVAL '90 days') AS events_last_90d
FROM stocks
WHERE event_time > now() - INTERVAL '1 year'
GROUP BY ticker
HAVING COUNT(*) >= 2;
