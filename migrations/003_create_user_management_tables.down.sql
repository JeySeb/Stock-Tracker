-- Rollback User Management Tables Migration
-- Drops tables in reverse order to avoid foreign key constraints

DROP TRIGGER IF EXISTS update_subscriptions_updated_at ON subscriptions;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order of creation to handle foreign key dependencies
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS subscription_plan;
DROP TYPE IF EXISTS subscription_status;
DROP TYPE IF EXISTS user_tier; 