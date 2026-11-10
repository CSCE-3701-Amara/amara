-- Restore grading fields on databases that executed the original rollback SQL.
ALTER TABLE assessment_attempts
    ADD COLUMN IF NOT EXISTS code_grading_retry_at timestamptz,
    ADD COLUMN IF NOT EXISTS non_code_evidence jsonb;

CREATE INDEX IF NOT EXISTS assessment_attempts_code_grading_due
    ON assessment_attempts (code_grading_retry_at NULLS FIRST, submitted_at, id)
    WHERE status = 'pending_grading' AND submitted_at IS NOT NULL;
