-- Drop triggers first
DROP TRIGGER IF EXISTS update_chat_sessions_updated_at ON chat_sessions;

-- Note: Not dropping update_updated_at_column() function as it's shared with migration 003 tables

DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;