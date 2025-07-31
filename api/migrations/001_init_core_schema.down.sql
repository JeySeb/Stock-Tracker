-- Core: ingestion, stocks, brokers y enums
DROP TABLE IF EXISTS ingestion_logs CASCADE;
DROP TABLE IF EXISTS stocks         CASCADE;
DROP TABLE IF EXISTS brokers        CASCADE;

DROP TYPE  IF EXISTS ingest_status;
DROP TYPE  IF EXISTS rating_grade;
