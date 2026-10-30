-- name: GetMFASecret :one
SELECT account_id::text, account_type, secret, enabled, created_at
FROM auth_mfa_credentials
WHERE account_id = $1::uuid AND account_type = $2;

-- name: SaveMFASecret :exec
INSERT INTO auth_mfa_credentials (account_id, account_type, secret, enabled, created_at)
VALUES ($1::uuid, $2, $3, FALSE, $4)
ON CONFLICT (account_id, account_type)
DO UPDATE SET secret = EXCLUDED.secret, enabled = FALSE, created_at = EXCLUDED.created_at, enabled_at = NULL;

-- name: EnableMFA :exec
UPDATE auth_mfa_credentials
SET enabled = TRUE, enabled_at = $3
WHERE account_id = $1::uuid AND account_type = $2;
