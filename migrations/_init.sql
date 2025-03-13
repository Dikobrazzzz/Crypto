-- +goose Up
-- +goose StatementBegin

CREATE TABLE main (
    id SERIAL PRIMARY KEY,
    wallet_address TEXT NOT NULL,
    chain_name TEXT NOT NULL,
    crypto_name TEXT NOT NULL,
    tag TEXT,
    balance NUMERIC DEFAULT 0
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS main;
-- +goose StatementEnd
