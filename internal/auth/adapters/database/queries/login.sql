-- name: FindCredentials :one
SELECT id::text, email, account_type, password_hash, salt, account_status,
       deactivated, banned, COALESCE(company_id::text, '') AS company_id, company_approved
FROM (
    SELECT c.id, c.email, 'candidate'::text AS account_type, c.password_hash, c.salt,
           c.account_status, (c.deactivated_at IS NOT NULL OR c.deleted_at IS NOT NULL) AS deactivated,
           COALESCE(c.banned, FALSE) AS banned, NULL::uuid AS company_id,
           TRUE AS company_approved
    FROM candidates c
    WHERE lower(c.email) = lower($1)
    UNION ALL
    SELECT i.id, i.email, 'interviewer'::text AS account_type, i.password_hash, i.salt,
           i.account_status, (i.deactivated_at IS NOT NULL OR i.deleted_at IS NOT NULL) AS deactivated,
           FALSE AS banned, i.company_id,
           COALESCE(co.approved AND co.deactivated_at IS NULL AND co.deleted_at IS NULL AND COALESCE(co.banned, FALSE) = FALSE, FALSE) AS company_approved
    FROM interviewers i
    LEFT JOIN companies co ON co.id = i.company_id
    WHERE lower(i.email) = lower($1)
    UNION ALL
    SELECT co.id, co.email, 'company'::text AS account_type, co.password_hash, co.salt,
           CASE WHEN co.approved THEN 'active' ELSE 'pending_review' END AS account_status,
           (co.deactivated_at IS NOT NULL OR co.deleted_at IS NOT NULL) AS deactivated,
           COALESCE(co.banned, FALSE) AS banned, co.id AS company_id,
           co.approved AS company_approved
    FROM companies co
    WHERE lower(co.email) = lower($1)
) accounts
LIMIT 1;

-- name: CreateSession :exec
INSERT INTO auth_sessions (token, account_id, account_type, expires_at)
VALUES ($1, $2::uuid, $3, $4);

-- name: TerminateSession :exec
DELETE FROM auth_sessions
WHERE token = $1;
