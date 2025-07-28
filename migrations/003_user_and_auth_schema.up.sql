-- 003_user_and_auth_schema.up.sql
-- Gestión de usuarios, subscripciones y sesiones (sin triggers)

CREATE TYPE IF NOT EXISTS user_tier            AS ENUM ('guest','basic','premium');
CREATE TYPE IF NOT EXISTS subscription_status  AS ENUM ('active','cancelled','expired','pending');
CREATE TYPE IF NOT EXISTS subscription_plan    AS ENUM ('monthly','yearly');

-- ---------- USERS ----------
CREATE TABLE IF NOT EXISTS users (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email         STRING  NOT NULL UNIQUE
                 CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
  password_hash STRING  NOT NULL,
  first_name    STRING  NOT NULL CHECK (char_length(first_name) BETWEEN 1 AND 100),
  last_name     STRING  NOT NULL CHECK (char_length(last_name)  BETWEEN 1 AND 100),
  tier          user_tier NOT NULL DEFAULT 'basic',
  is_verified   BOOL NOT NULL DEFAULT FALSE,
  last_login    TIMESTAMPTZ,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE INDEX IF NOT EXISTS idx_users__tier         ON users (tier);
CREATE INDEX IF NOT EXISTS idx_users__created_desc ON users (created_at DESC);

-- ---------- SUBSCRIPTIONS ----------
CREATE TABLE IF NOT EXISTS subscriptions (
  id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  plan        subscription_plan NOT NULL,
  status      subscription_status NOT NULL DEFAULT 'pending',
  price       DECIMAL(10,2) NOT NULL CHECK (price > 0),
  currency    STRING NOT NULL DEFAULT 'USD' CHECK (currency IN ('USD','EUR','GBP')),
  start_date  TIMESTAMPTZ NOT NULL,
  end_date    TIMESTAMPTZ NOT NULL CHECK (end_date > start_date),
  payment_reference STRING,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE INDEX IF NOT EXISTS idx_subs__user_id     ON subscriptions (user_id);
CREATE INDEX IF NOT EXISTS idx_subs__status      ON subscriptions (status);
CREATE INDEX IF NOT EXISTS idx_subs__end_date    ON subscriptions (end_date);
CREATE INDEX IF NOT EXISTS idx_subs__payment_ref ON subscriptions (payment_reference);

-- ---------- SESSIONS ----------
CREATE TABLE IF NOT EXISTS sessions (
  id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  refresh_token STRING NOT NULL UNIQUE,
  user_agent    STRING,
  ip_address    INET,
  expires_at    TIMESTAMPTZ NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
  -- No CHECK con now() porque CockroachDB no permite funciones no deterministas en constraints
);

CREATE INDEX IF NOT EXISTS idx_sessions__user_id     ON sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_sessions__expires_at  ON sessions (expires_at);
CREATE INDEX IF NOT EXISTS idx_sessions__created_at  ON sessions (created_at DESC);
