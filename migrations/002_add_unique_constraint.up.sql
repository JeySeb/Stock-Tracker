-- Add unique constraint for ON CONFLICT to work
ALTER TABLE stocks ADD CONSTRAINT unique_ticker_event_time UNIQUE (ticker, event_time); 