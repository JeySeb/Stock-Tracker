-- Drop tables in reverse order to avoid foreign key constraints
DROP TABLE IF EXISTS ingestion_logs;
DROP TABLE IF EXISTS recommendations;
DROP TABLE IF EXISTS stocks;
DROP TABLE IF EXISTS brokers; 