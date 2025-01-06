-- COMMENTS
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent_comment_id INT REFERENCES comments (id),
    content TEXT NOT NULL CHECK (LENGTH(TRIM(content)) > 0),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN NOT NULL DEFAULT false
);

-- Index for nested comments
CREATE INDEX idx_comments_post_id_parent_comment_id ON comments (post_id, parent_comment_id, created_at DESC);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();