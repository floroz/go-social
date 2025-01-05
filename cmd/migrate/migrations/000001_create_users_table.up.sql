CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL CHECK (LENGTH(TRIM(first_name)) > 0),
    last_name VARCHAR(50) NOT NULL CHECK (LENGTH(TRIM(last_name)) > 0),
    username VARCHAR(30) NOT NULL UNIQUE CHECK (LENGTH(username) > 0),
    email VARCHAR(255) NOT NULL CHECK (LENGTH(TRIM(email)) > 0),
    password VARCHAR(255) NOT NULL CHECK (LENGTH(password) >= 10),
    -- mechanism to track failed logins and lock account - possible improvement
    -- failed_login_attempts INT DEFAULT 0,
    profile_picture_url VARCHAR(255),
    bio TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT users_email_unique UNIQUE (LOWER(email))
);

-- Index for email lookups (case insensitive)
CREATE UNIQUE INDEX idx_users_email_lower ON users (LOWER(email));

-- Index for pagination and sorting
CREATE INDEX idx_users_created_at ON users (created_at DESC);

-- Trigger to update updated_at timestamp
CREATE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();