-- +goose Up
ALTER TABLE candidates
    ADD COLUMN account_status TEXT NOT NULL DEFAULT 'pending_verification';

ALTER TABLE interviewers
    ADD COLUMN account_status TEXT NOT NULL DEFAULT 'pending_verification';

ALTER TABLE companies
    ADD COLUMN approved BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN registration_secret_hash BYTEA;

CREATE UNIQUE INDEX candidates_email_lower_idx ON candidates (lower(email));
CREATE UNIQUE INDEX interviewers_email_lower_idx ON interviewers (lower(email));
CREATE UNIQUE INDEX companies_email_lower_idx ON companies (lower(email));

CREATE TABLE auth_email_verification_challenges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL,
    account_type TEXT NOT NULL CHECK (account_type IN ('candidate', 'interviewer')),
    code_hash BYTEA NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    consumed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX auth_email_verification_challenges_account_idx
    ON auth_email_verification_challenges (account_id, account_type);

-- +goose Down
DROP TABLE auth_email_verification_challenges;
DROP INDEX companies_email_lower_idx;
DROP INDEX interviewers_email_lower_idx;
DROP INDEX candidates_email_lower_idx;
ALTER TABLE companies
    DROP COLUMN registration_secret_hash,
    DROP COLUMN approved;
ALTER TABLE interviewers DROP COLUMN account_status;
ALTER TABLE candidates DROP COLUMN account_status;
