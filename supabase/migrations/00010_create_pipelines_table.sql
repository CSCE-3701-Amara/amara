-- +goose Up
CREATE TABLE pipelines (
    id          UUID    PRIMARY KEY,
    version     INTEGER,
    status      TEXT
);

-- +goose Down
DROP TABLE pipelines;