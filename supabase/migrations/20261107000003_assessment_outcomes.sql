-- +goose Up
ALTER TABLE assessment_attempts
    ADD COLUMN saved_responses JSONB NOT NULL DEFAULT '{}'
        CHECK (jsonb_typeof(saved_responses) = 'object'),
    ADD COLUMN final_evidence JSONB
        CHECK (jsonb_typeof(final_evidence) = 'array'),
    ADD COLUMN final_score INTEGER CHECK (final_score BETWEEN 0 AND 100),
    ADD COLUMN evaluated_at TIMESTAMPTZ,
    ADD CONSTRAINT assessment_evidence_atomic CHECK (
        (final_evidence IS NULL) = (evaluated_at IS NULL)
        AND (evaluated_at IS NULL OR (status = 'scored' AND submitted_at IS NOT NULL))
        AND (final_score IS NULL OR evaluated_at IS NOT NULL)
    );

-- +goose Down
ALTER TABLE assessment_attempts
    DROP CONSTRAINT assessment_evidence_atomic,
    DROP COLUMN evaluated_at,
    DROP COLUMN final_score,
    DROP COLUMN final_evidence,
    DROP COLUMN saved_responses;
