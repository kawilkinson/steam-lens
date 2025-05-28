-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    username TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL,
    steam_id TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
