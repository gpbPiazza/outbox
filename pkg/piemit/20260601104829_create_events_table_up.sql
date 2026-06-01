CREATE TYPE pi_event_status IF NOT EXISTS ('pending', 'active', 'retry', 'archived', 'completed');

CREATE TABLE IF NOT EXISTS pi_events (
    id         UUID        PRIMARY KEY,
    status     pi_event_status        NOT NULL,
    topic      VARCHAR(255),
    attempts   INT32,
    last_error VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL    DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL    DEFAULT NOW()
);
