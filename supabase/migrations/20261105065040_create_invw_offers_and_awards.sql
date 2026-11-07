CREATE TABLE identity_reveal_records (
    id UUID PRIMARY KEY,
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    revealed_by UUID REFERENCES companies(id),
    application_id UUID NOT NULL UNIQUE REFERENCES applications(id),
    stage_id UUID REFERENCES pipeline_stages(id),
    revealed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Preserve repeated awards and the definition/level scale used at award time.
ALTER TABLE candidate_amarat
    ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ADD COLUMN company_id UUID REFERENCES companies(id),
    ADD COLUMN application_id UUID REFERENCES applications(id),
    ADD COLUMN source_stage_execution_id UUID REFERENCES stage_executions(id),
    ADD COLUMN definition_version INTEGER CHECK (definition_version > 0),
    ADD COLUMN definition_snapshot JSONB,
    ADD COLUMN score INTEGER CHECK (score BETWEEN 0 AND 100),
    ADD COLUMN level TEXT,
    ADD COLUMN strengths TEXT,
    ADD COLUMN areas_for_development TEXT,
    ADD COLUMN pass_decision JSONB,
    ADD COLUMN awarded_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP;

CREATE TABLE offers (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL REFERENCES applications(id),
    job_title TEXT NOT NULL,
    compensation_summary TEXT NOT NULL,
    proposed_start_date DATE NOT NULL,
    response_deadline TIMESTAMPTZ NOT NULL,
    message TEXT,
    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted', 'declined', 'expired', 'withdrawn')),
    response TEXT CHECK (response IN ('accepted', 'declined')),
    responded_at TIMESTAMPTZ,
    withdrawn_at TIMESTAMPTZ,
    deadline_history JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(deadline_history) = 'array'),
    capacity_confirmed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CHECK (response_deadline > created_at),
    CHECK (
        (response IS NULL AND responded_at IS NULL AND status NOT IN ('accepted', 'declined'))
        OR (response IS NOT NULL AND responded_at IS NOT NULL
            AND status = response AND responded_at < response_deadline)
    )
);

CREATE UNIQUE INDEX offers_one_pending_application_idx
    ON offers (application_id) WHERE status = 'pending';
CREATE INDEX offers_pending_deadline_idx ON offers (response_deadline) WHERE status = 'pending';
CREATE INDEX candidate_amarat_candidate_awarded_idx ON candidate_amarat (cand_id, awarded_at);
