-- name: CompanyRegistrationEmailExists :one
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

-- name: CompanyRegistrationNumberExists :one
SELECT EXISTS (
    SELECT 1 FROM auth_company_registration_requests
    WHERE lower(registration_number) = lower($1)
      AND status IN ('pending_review', 'approved')
);

-- name: CreateCompanyRegistration :one
WITH request AS (
    INSERT INTO auth_company_registration_requests
        (company_name, registration_number, email, password_hash, salt, status, created_at, reviewed_at)
    VALUES ($1, $2, $3, $4, $5,
            CASE WHEN sqlc.arg(auto_approve)::boolean THEN 'approved' ELSE 'pending_review' END, sqlc.arg(created_at)::timestamp,
            CASE WHEN sqlc.arg(auto_approve)::boolean THEN sqlc.arg(created_at)::timestamp END)
    RETURNING id, company_name, email, password_hash, salt, status, created_at
), company AS (
    INSERT INTO companies (id, legal_name, email, password_hash, salt, approved, created_at)
    SELECT id, company_name, email, password_hash, salt, TRUE, created_at
    FROM request WHERE status = 'approved'
    RETURNING id
)
SELECT request.id::text AS id, request.email,
       (CASE WHEN company.id IS NOT NULL THEN 'active' ELSE request.status END)::text AS status,
       request.created_at
FROM request LEFT JOIN company ON company.id = request.id;
