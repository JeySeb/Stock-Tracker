-- Remove unique constraint (CockroachDB syntax)
DROP INDEX IF EXISTS stocks@unique_ticker_event_time CASCADE; 