-- +goose Up
ALTER TABLE assessment_attempts ADD COLUMN code_grading_retry_at timestamptz;

CREATE INDEX assessment_attempts_code_grading_due
    ON assessment_attempts (code_grading_retry_at NULLS FIRST, submitted_at, id)
    WHERE status = 'pending_grading' AND submitted_at IS NOT NULL;

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP INDEX IF EXISTS assessment_attempts_code_grading_due;
ALTER TABLE assessment_attempts DROP COLUMN code_grading_retry_at;
*/
