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
