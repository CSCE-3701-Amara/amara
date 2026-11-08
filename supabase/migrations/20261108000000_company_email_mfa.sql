-- +goose Up
CREATE TABLE auth_company_email_mfa (
    company_id UUID PRIMARY KEY REFERENCES companies(id) ON DELETE CASCADE,
    code_mac BYTEA NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    resend_at TIMESTAMPTZ NOT NULL,
    attempts_left INTEGER NOT NULL CHECK (attempts_left BETWEEN 0 AND 5)
);

-- +goose Down
DROP TABLE auth_company_email_mfa;
