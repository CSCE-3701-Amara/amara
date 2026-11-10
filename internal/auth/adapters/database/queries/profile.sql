-- name: GetCompanyProfile :one
SELECT id::text, legal_name, display_name, description, website_link, logo_url,
       contact_phone, contact_email, linkedin_profile, approved
FROM companies
WHERE id = @id::uuid
  AND deleted_at IS NULL
  AND (@include_private::bool OR (approved = TRUE AND deactivated_at IS NULL AND COALESCE(banned, FALSE) = FALSE));

-- name: GetCandidateProfile :one
SELECT id::text, name, email, created_at, updated_at, phone_number, github_profile,
       linkedin_profile, portfolio_link, gender, birth_date, nationality, address,
       biography, years_of_experience, available,
       ARRAY(SELECT DISTINCT a.amara_title FROM candidate_amarat a
             WHERE a.cand_id = candidates.id ORDER BY a.amara_title)::text[] AS amarat
FROM candidates
WHERE id = @id::uuid
  AND deleted_at IS NULL;

-- name: GetPublicCompanyProfile :one
SELECT id::text, legal_name, display_name, description, website_link, logo_url,
       contact_phone, contact_email, linkedin_profile, approved
FROM companies
WHERE id = @id::uuid
  AND approved = TRUE
  AND deleted_at IS NULL
  AND deactivated_at IS NULL
  AND COALESCE(banned, FALSE) = FALSE;

-- name: ListCompanyInterviewers :many
SELECT i.id::text, i.company_id::text, i.name, i.email, i.account_status,
       i.created_at, i.updated_at
FROM interviewers i
JOIN companies c ON c.id = i.company_id
WHERE i.company_id = @company_id::uuid
  AND i.deleted_at IS NULL
  AND c.deleted_at IS NULL
  AND c.deactivated_at IS NULL
  AND c.approved = TRUE
ORDER BY i.created_at DESC, i.id
LIMIT @page_limit OFFSET @page_offset;

-- name: CountCompanyInterviewers :one
SELECT COUNT(*)
FROM interviewers i
JOIN companies c ON c.id = i.company_id
WHERE i.company_id = @company_id::uuid
  AND i.deleted_at IS NULL
  AND c.deleted_at IS NULL
  AND c.deactivated_at IS NULL
  AND c.approved = TRUE;

-- name: GetInterviewerProfile :one
SELECT i.id::text, i.company_id::text, i.name, i.email, i.account_status,
       i.created_at, i.updated_at
FROM interviewers i
JOIN companies c ON c.id = i.company_id
WHERE i.id = @id::uuid
  AND i.deleted_at IS NULL
  AND c.deleted_at IS NULL
  AND c.deactivated_at IS NULL
  AND c.approved = TRUE;

-- name: UpdateCandidateProfile :one
UPDATE candidates
SET name = CASE WHEN @set_name::bool THEN @name::text ELSE name END,
    phone_number = CASE WHEN @set_phone_number::bool THEN @phone_number::text ELSE phone_number END,
    github_profile = CASE WHEN @set_github_profile::bool THEN @github_profile::text ELSE github_profile END,
    linkedin_profile = CASE WHEN @set_linkedin_profile::bool THEN @linkedin_profile::text ELSE linkedin_profile END,
    portfolio_link = CASE WHEN @set_portfolio_link::bool THEN @portfolio_link::text ELSE portfolio_link END,
    gender = CASE WHEN @set_gender::bool THEN @gender::text ELSE gender END,
    birth_date = CASE WHEN @set_birth_date::bool THEN @birth_date::date ELSE birth_date END,
    nationality = CASE WHEN @set_nationality::bool THEN @nationality::text ELSE nationality END,
    address = CASE WHEN @set_address::bool THEN @address::text ELSE address END,
    biography = CASE WHEN @set_biography::bool THEN @biography::text ELSE biography END,
    years_of_experience = CASE WHEN @set_years_of_experience::bool THEN @years_of_experience::int ELSE years_of_experience END,
    available = CASE WHEN @set_available::bool THEN @available::bool ELSE available END,
    updated_at = now()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id::text, name, email, created_at, updated_at, phone_number, github_profile,
          linkedin_profile, portfolio_link, gender, birth_date, nationality, address,
          biography, years_of_experience, available;

-- name: UpdateCompanyProfile :one
UPDATE companies
SET legal_name = CASE WHEN @set_legal_name::bool THEN @legal_name::text ELSE legal_name END,
    display_name = CASE WHEN @set_display_name::bool THEN sqlc.narg('display_name') ELSE display_name END,
    description = CASE WHEN @set_description::bool THEN sqlc.narg('description') ELSE description END,
    website_link = CASE WHEN @set_website_link::bool THEN sqlc.narg('website_link') ELSE website_link END,
    logo_url = CASE WHEN @set_logo_url::bool THEN sqlc.narg('logo_url') ELSE logo_url END,
    contact_phone = CASE WHEN @set_contact_phone::bool THEN sqlc.narg('contact_phone') ELSE contact_phone END,
    contact_email = CASE WHEN @set_contact_email::bool THEN sqlc.narg('contact_email') ELSE contact_email END,
    linkedin_profile = CASE WHEN @set_linkedin_profile::bool THEN sqlc.narg('linkedin_profile') ELSE linkedin_profile END,
    updated_at = now()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id::text, legal_name, display_name, description, website_link, logo_url,
          contact_phone, contact_email, linkedin_profile, approved;

-- name: UpdateInterviewerProfile :one
UPDATE interviewers
SET name = @name::text, updated_at = now()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id::text, company_id::text, name, email, account_status, created_at, updated_at;
