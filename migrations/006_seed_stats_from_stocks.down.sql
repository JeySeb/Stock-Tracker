-- Revierte la carga inicial de estadísticas
TRUNCATE TABLE brokerage_stats  IF EXISTS;
TRUNCATE TABLE ticker_stats     IF EXISTS;
