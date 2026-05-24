-- +goose Up
CREATE TYPE move_status AS ENUM (
    'Wes sendo Bagre',
    'Wes fez a boa',
    'Wes com o livro de baixo do braço',
    'Wes esquecendo a historia enquanto conta'
);

CREATE TABLE moves (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    wes_id      UUID        NOT NULL    REFERENCES wes(id) ON DELETE CASCADE,
    status      move_status NOT NULL,
    description TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL    DEFAULT NOW()
);

CREATE INDEX moves_wes_id_idx ON moves (wes_id);

-- +goose Down
DROP TABLE moves;
DROP TYPE move_status;
