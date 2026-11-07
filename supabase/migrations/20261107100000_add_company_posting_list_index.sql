-- +goose Up
CREATE INDEX job_postings_company_list_idx
    ON job_postings (
        company_id,
        (COALESCE(created_at, '0001-01-01 00:00:00'::timestamp)) DESC,
        id DESC
    )
    WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS job_postings_company_list_idx;
