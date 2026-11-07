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
INSERT INTO auth_company_registration_requests
    (company_name, registration_number, email, password_hash, salt, status, created_at)
VALUES ($1, $2, $3, $4, $5, 'pending_review', $6)
RETURNING id::text, email, status, created_at;
