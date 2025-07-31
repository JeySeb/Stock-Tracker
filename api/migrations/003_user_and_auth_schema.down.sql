-- Gestión de usuarios, subscripciones y sesiones
DROP TABLE IF EXISTS sessions       CASCADE;
DROP TABLE IF EXISTS subscriptions  CASCADE;
DROP TABLE IF EXISTS users          CASCADE;

DROP TYPE  IF EXISTS subscription_plan;
DROP TYPE  IF EXISTS subscription_status;
DROP TYPE  IF EXISTS user_tier;
