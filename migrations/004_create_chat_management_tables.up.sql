-- Chat Sessions table - AI chat session management
CREATE TABLE IF NOT EXISTS chat_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title STRING NOT NULL,
    status STRING NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    
    -- Constraints
    CONSTRAINT title_length CHECK (char_length(title) >= 1 AND char_length(title) <= 200),
    CONSTRAINT valid_status CHECK (status IN ('active', 'closed')),
    
    -- Indexes for performance
    INDEX idx_chat_sessions_user_id (user_id),
    INDEX idx_chat_sessions_created_at (created_at DESC)
);

-- Chat Messages table - individual messages within sessions
CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES chat_sessions(id) ON DELETE CASCADE,
    role STRING NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    
    -- Constraints
    CONSTRAINT valid_role CHECK (role IN ('user', 'assistant')), -- Removed 'system' as it's not in MessageRole
    CONSTRAINT content_not_empty CHECK (char_length(content) > 0),
    
    -- Indexes for performance
    INDEX idx_chat_messages_session_id (session_id),
    INDEX idx_chat_messages_created_at (created_at ASC)
);

-- Add trigger to update updated_at timestamp for chat sessions
-- Note: The update_updated_at_column() function already exists from migration 003
CREATE TRIGGER update_chat_sessions_updated_at BEFORE UPDATE ON chat_sessions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();