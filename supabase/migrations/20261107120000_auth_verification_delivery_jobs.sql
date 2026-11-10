-- +goose Up
CREATE TABLE auth_email_verification_deliveries (
    account_id UUID PRIMARY KEY,
    account_type TEXT NOT NULL CHECK (account_type IN ('candidate', 'interviewer')),
    recipient TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP TABLE auth_email_verification_deliveries;
*/
