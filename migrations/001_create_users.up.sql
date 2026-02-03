CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    github_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    name VARCHAR(255),
    avatar_url TEXT,
    encrypted_token TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_username ON users(username);
