CREATE TABLE auth_email_verification_deliveries (
    account_id UUID PRIMARY KEY,
    account_type TEXT NOT NULL CHECK (account_type IN ('candidate', 'interviewer')),
    recipient TEXT NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
