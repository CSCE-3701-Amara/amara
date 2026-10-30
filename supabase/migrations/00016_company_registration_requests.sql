-- +goose Up
CREATE TABLE auth_company_registration_requests (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_name        TEXT NOT NULL,
    registration_number TEXT NOT NULL,
    email               TEXT NOT NULL,
    password_hash       BYTEA NOT NULL,
    salt                BYTEA NOT NULL,
    status              TEXT NOT NULL DEFAULT 'pending_review',
    created_at          TIMESTAMP NOT NULL DEFAULT now(),
    reviewed_at         TIMESTAMP,
    review_reason       TEXT
);

CREATE UNIQUE INDEX auth_company_registration_email_idx
    ON auth_company_registration_requests (lower(email))
    WHERE status IN ('pending_review', 'approved');
CREATE UNIQUE INDEX auth_company_registration_number_idx
    ON auth_company_registration_requests (lower(registration_number))
    WHERE status IN ('pending_review', 'approved');

-- +goose Down
DROP INDEX auth_company_registration_number_idx;
DROP INDEX auth_company_registration_email_idx;
DROP TABLE auth_company_registration_requests;
