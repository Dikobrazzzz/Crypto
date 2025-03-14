-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS main (
    id SERIAL PRIMARY KEY,
    wallet_address TEXT NOT NULL,
    chain_name TEXT NOT NULL,
    crypto_name TEXT NOT NULL,
    tag TEXT,
    balance NUMERIC DEFAULT 0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS main;
-- +goose StatementEnd
