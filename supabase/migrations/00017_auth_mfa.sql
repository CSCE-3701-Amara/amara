-- +goose Up
CREATE TABLE auth_mfa_credentials (
    account_id   UUID NOT NULL,
    account_type TEXT NOT NULL CHECK (account_type IN ('candidate', 'interviewer', 'company')),
    secret       BYTEA NOT NULL,
    enabled      BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP NOT NULL DEFAULT now(),
    enabled_at   TIMESTAMP,
    PRIMARY KEY (account_id, account_type)
);

-- +goose Down
DROP TABLE auth_mfa_credentials;
