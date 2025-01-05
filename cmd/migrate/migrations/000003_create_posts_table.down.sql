DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;

DROP TRIGGER IF EXISTS update_posts_updated_at ON posts;

DROP FUNCTION IF EXISTS update_updated_at_column;

DROP TABLE IF EXISTS comments;

DROP TABLE IF EXISTS posts;