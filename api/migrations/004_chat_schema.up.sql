-- 004_chat_schema.up.sql
-- Chat sessions & messages (sin triggers)

CREATE TYPE IF NOT EXISTS chat_status AS ENUM ('active','closed');
CREATE TYPE IF NOT EXISTS chat_role   AS ENUM ('user','assistant');

CREATE TABLE IF NOT EXISTS chat_sessions (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title      STRING NOT NULL CHECK (char_length(title) BETWEEN 1 AND 200),
  status     chat_status NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE INDEX IF NOT EXISTS idx_chat_sessions__user_id  ON chat_sessions (user_id);
CREATE INDEX IF NOT EXISTS idx_chat_sessions__created  ON chat_sessions (created_at DESC);

CREATE TABLE IF NOT EXISTS chat_messages (
  id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
  role       chat_role NOT NULL,
  content    STRING NOT NULL CHECK (char_length(content) > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_chat_msgs__session_id  ON chat_messages (session_id);
CREATE INDEX IF NOT EXISTS idx_chat_msgs__created_asc ON chat_messages (created_at ASC);
