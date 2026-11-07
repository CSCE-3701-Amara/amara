-- name: EmailExists :one
SELECT EXISTS (
    SELECT 1 FROM candidates WHERE lower(email) = lower($1)
    UNION ALL
    SELECT 1 FROM interviewers WHERE lower(email) = lower($1)
    UNION ALL
    SELECT 1 FROM companies WHERE lower(email) = lower($1)
    UNION ALL
    SELECT 1 FROM auth_company_registration_requests
    WHERE lower(email) = lower($1)
      AND status IN ('pending_review', 'approved')
);

-- name: FindApprovedCompanyBySecret :one
SELECT id::text
FROM companies
WHERE registration_secret_hash = $1
  AND approved = TRUE
  AND deleted_at IS NULL
  AND deactivated_at IS NULL;

-- name: CreateCandidateAccount :one
INSERT INTO candidates (id, name, email, password_hash, salt, created_at, account_status)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, 'pending_verification')
RETURNING id::text, email, account_status, created_at;

-- name: CreateInterviewerAccount :one
INSERT INTO interviewers (id, company_id, name, email, password_hash, salt, created_at, account_status)
VALUES (gen_random_uuid(), $1::uuid, $2, $3, $4, $5, $6, 'pending_verification')
RETURNING id::text, email, account_status, created_at;

-- name: SaveEmailVerificationChallenge :exec
WITH deleted AS (
    DELETE FROM auth_email_verification_challenges
    WHERE account_id = $1::uuid
      AND account_type = $2
),
challenge AS (
    INSERT INTO auth_email_verification_challenges
        (account_id, account_type, code_hash, expires_at)
    VALUES ($1::uuid, $2, $3, $4)
)
INSERT INTO auth_email_verification_deliveries
        (account_id, account_type, recipient, expires_at)
    VALUES ($1::uuid, $2, $5, $4)
    ON CONFLICT (account_id) DO UPDATE SET
        account_type = EXCLUDED.account_type,
        recipient = EXCLUDED.recipient,
        expires_at = EXCLUDED.expires_at;

-- name: FindPendingVerification :one
SELECT c.id::text, c.email, c.account_status, c.created_at,
       d.account_type, d.recipient, d.expires_at
FROM auth_email_verification_deliveries d
JOIN candidates c ON d.account_id = c.id
WHERE d.account_type = 'candidate'
  AND lower(c.email) = lower($1)
  AND c.account_status = 'pending_verification'
UNION ALL
SELECT i.id::text, i.email, i.account_status, i.created_at,
       d.account_type, d.recipient, d.expires_at
FROM auth_email_verification_deliveries d
JOIN interviewers i ON d.account_id = i.id
WHERE d.account_type = 'interviewer'
  AND lower(i.email) = lower($1)
  AND i.account_status = 'pending_verification'
LIMIT 1;
