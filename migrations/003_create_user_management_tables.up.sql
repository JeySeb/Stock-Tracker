-- User Management Tables Migration
-- Creates tables for users, subscriptions, chat sessions, and authentication

-- Create ENUM types for user tiers and subscription status
CREATE TYPE user_tier AS ENUM ('guest', 'basic', 'premium');
CREATE TYPE subscription_status AS ENUM ('active', 'cancelled', 'expired', 'pending');
CREATE TYPE subscription_plan AS ENUM ('monthly', 'yearly');

-- Users table - core user information and tier management
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING NOT NULL UNIQUE,
    password_hash STRING NOT NULL,
    first_name STRING NOT NULL,
    last_name STRING NOT NULL,
    tier user_tier NOT NULL DEFAULT 'basic',
    is_verified BOOL NOT NULL DEFAULT false,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    
    -- Constraints
    CONSTRAINT email_format CHECK (email ~ '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}$'),
    CONSTRAINT name_length CHECK (char_length(first_name) >= 1 AND char_length(first_name) <= 100),
    CONSTRAINT last_name_length CHECK (char_length(last_name) >= 1 AND char_length(last_name) <= 100),
    
    -- Indexes for performance
    INDEX idx_users_email (email),
    INDEX idx_users_tier (tier),
    INDEX idx_users_created_at (created_at DESC)
);

-- Subscriptions table - premium subscription management
CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    plan subscription_plan NOT NULL,
    status subscription_status NOT NULL DEFAULT 'pending',
    price DECIMAL(10,2) NOT NULL,
    currency STRING NOT NULL DEFAULT 'USD',
    start_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ NOT NULL,
    payment_reference STRING,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    
    -- Constraints
    CONSTRAINT positive_price CHECK (price > 0),
    CONSTRAINT valid_currency CHECK (currency IN ('USD', 'EUR', 'GBP')),
    CONSTRAINT valid_date_range CHECK (end_date > start_date),
    
    -- Indexes for performance
    INDEX idx_subscriptions_user_id (user_id),
    INDEX idx_subscriptions_status (status),
    INDEX idx_subscriptions_end_date (end_date),
    INDEX idx_subscriptions_payment_ref (payment_reference)
);

-- Sessions table - authentication session management
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token STRING NOT NULL UNIQUE,
    user_agent STRING,
    ip_address INET,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    
    -- Constraints
    CONSTRAINT future_expiry CHECK (expires_at > created_at),
    
    -- Indexes for performance
    INDEX idx_sessions_user_id (user_id),
    INDEX idx_sessions_refresh_token (refresh_token),
    INDEX idx_sessions_expires_at (expires_at),
    INDEX idx_sessions_created_at (created_at DESC)
);

-- Add trigger to update updated_at timestamp for users
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_subscriptions_updated_at BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

