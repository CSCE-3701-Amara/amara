-- +goose Up
CREATE TABLE interviewers (
    id              UUID        PRIMARY KEY,
    company_id      UUID        REFERENCES companies(id),
    name            TEXT        NOT NULL,
    email           TEXT        NOT NULL,
    password_hash   BYTEA       NOT NULL,
    salt            BYTEA       NOT NULL,
    created_at      TIMESTAMP   NOT NULL,
    updated_at      TIMESTAMP,
    deleted_at      TIMESTAMP,
    deactivated_at  TIMESTAMP,
    last_login_at   TIMESTAMP
);

-- +goose Down
DROP TABLE interviewers;