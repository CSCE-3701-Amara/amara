-- +goose Up
ALTER TABLE application_assessments
    ADD COLUMN evaluation_status TEXT NOT NULL DEFAULT 'pending'
                                CHECK (evaluation_status IN ('pending', 'evaluated')),
    ADD COLUMN passed BOOLEAN,
    ADD COLUMN evaluated_at TIMESTAMP,
    ADD COLUMN result_disclosed_at TIMESTAMP,
    ADD COLUMN feedback_text TEXT,
    ADD COLUMN feedback_strengths JSONB NOT NULL DEFAULT '[]'::jsonb
                                      CHECK (jsonb_typeof(feedback_strengths) = 'array'),
    ADD COLUMN areas_for_improvement JSONB NOT NULL DEFAULT '[]'::jsonb
                                            CHECK (jsonb_typeof(areas_for_improvement) = 'array'),
    ADD COLUMN feedback_is_ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN feedback_disclosed_at TIMESTAMP,
    ADD CONSTRAINT application_assessments_evaluation_state_check CHECK (
        (evaluation_status = 'pending' AND passed IS NULL AND evaluated_at IS NULL)
        OR (evaluation_status = 'evaluated' AND passed IS NOT NULL AND evaluated_at IS NOT NULL)
    ),
    ADD CONSTRAINT application_assessments_result_disclosure_check CHECK (
        result_disclosed_at IS NULL OR evaluation_status = 'evaluated'
    );

CREATE INDEX application_assessments_result_lookup_idx
    ON application_assessments (application_id, created_at, id)
    WHERE status IN ('submitted', 'completed');

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP INDEX IF EXISTS application_assessments_result_lookup_idx;
ALTER TABLE application_assessments
    DROP CONSTRAINT application_assessments_result_disclosure_check,
    DROP CONSTRAINT application_assessments_evaluation_state_check,
    DROP COLUMN feedback_disclosed_at,
    DROP COLUMN feedback_is_ai_generated,
    DROP COLUMN areas_for_improvement,
    DROP COLUMN feedback_strengths,
    DROP COLUMN feedback_text,
    DROP COLUMN result_disclosed_at,
    DROP COLUMN evaluated_at,
    DROP COLUMN passed,
    DROP COLUMN evaluation_status;
*/
