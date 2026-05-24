-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE wes (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    height     TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL    DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL    DEFAULT NOW()
);

-- +goose Down
DROP TABLE wes;
