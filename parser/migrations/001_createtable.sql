-- +goose Up
CREATE TABLE logs (
    timestamp        timestamp,
    hash             TEXT,
    url              TEXT,
    method           TEXT,
    host             TEXT,
    user_agent       TEXT,
    accept_encoding  TEXT,
    accept           TEXT,
    cookie           TEXT,
    ip               INET NULL,
    protocol         TEXT,
    headers          JSONB,
    body             BYTEA,
    PRIMARY KEY(timestamp, hash)
);

-- +goose Down
DROP TABLE logs;