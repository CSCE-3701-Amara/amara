-- +goose Up
CREATE TABLE auth_sessions (
    token        TEXT        PRIMARY KEY,
    account_id   UUID        NOT NULL,
    account_type TEXT        NOT NULL CHECK (account_type IN ('candidate', 'interviewer', 'company')),
    expires_at   TIMESTAMP   NOT NULL,
    created_at   TIMESTAMP   NOT NULL DEFAULT now()
);

CREATE INDEX auth_sessions_account_idx ON auth_sessions (account_id, account_type);
CREATE INDEX auth_sessions_expiry_idx ON auth_sessions (expires_at);

-- +goose Down
DROP TABLE auth_sessions;
