CREATE TABLE application_assessments (
    id                      UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id          UUID        NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    assessment_id           UUID        NOT NULL,
    status                  TEXT        NOT NULL DEFAULT 'not_started'
                                        CHECK (status IN ('not_started', 'in_progress', 'submitted', 'completed')),
    supports_interruption   BOOLEAN     NOT NULL DEFAULT FALSE,
    saved_progress          JSONB,
    created_at              TIMESTAMP   NOT NULL DEFAULT now(),
    updated_at              TIMESTAMP   NOT NULL DEFAULT now(),
    UNIQUE (application_id, assessment_id)
);

CREATE INDEX application_assessments_application_created_idx
    ON application_assessments (application_id, created_at, id);
