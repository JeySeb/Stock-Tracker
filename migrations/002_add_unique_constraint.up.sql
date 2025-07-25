-- Add unique constraint for ON CONFLICT to work (CockroachDB compatible)
-- Using CREATE UNIQUE INDEX IF NOT EXISTS which is more compatible than ALTER TABLE ADD CONSTRAINT
CREATE UNIQUE INDEX IF NOT EXISTS unique_ticker_event_time ON stocks (ticker, event_time); 