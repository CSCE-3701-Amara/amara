-- +goose Up
CREATE TABLE modules (
    id                  UUID    PRIMARY KEY,
    name                TEXT    NOT NULL,
    type                TEXT    NOT NULL,
    description         TEXT,
    company_id          UUID    REFERENCES companies(id),
    active_time_limit   INTERVAL,
    published           BOOLEAN
);

-- +goose Down
DROP TABLE modules;