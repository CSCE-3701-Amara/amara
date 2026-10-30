-- name: EmailExists :one
SELECT EXISTS (
    SELECT 1 FROM candidates WHERE lower(email) = lower($1)
    UNION ALL
    SELECT 1 FROM interviewers WHERE lower(email) = lower($1)
    UNION ALL
    SELECT 1 FROM companies WHERE lower(email) = lower($1)
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
INSERT INTO auth_email_verification_challenges
    (account_id, account_type, code_hash, expires_at)
VALUES ($1::uuid, $2, $3, $4);
