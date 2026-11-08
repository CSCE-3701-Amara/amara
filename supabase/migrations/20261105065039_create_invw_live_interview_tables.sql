CREATE TABLE interviews (
    id UUID PRIMARY KEY,
    stage_execution_id UUID NOT NULL UNIQUE REFERENCES stage_executions(id),
    scheduled_start_at TIMESTAMPTZ NOT NULL,
    scheduled_end_at TIMESTAMPTZ NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending_confirmation'
        CHECK (status IN ('pending_confirmation', 'confirmed', 'completed', 'cancelled')),
    outcome TEXT CHECK (outcome IN ('held', 'candidate_absent', 'interviewer_absent')),
    provider TEXT,
    provider_meeting_id TEXT,
    join_url TEXT,
    reschedule_count SMALLINT NOT NULL DEFAULT 0 CHECK (reschedule_count BETWEEN 0 AND 2),
    reschedule_history JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(reschedule_history) = 'array'),
    rubric_snapshot JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(rubric_snapshot) = 'array'),
    score INTEGER CHECK (score BETWEEN 0 AND 100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CHECK (scheduled_end_at > scheduled_start_at)
);

CREATE TABLE interviewer_assignments (
    interview_id UUID NOT NULL REFERENCES interviews(id),
    interviewer_id UUID NOT NULL REFERENCES interviewers(id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    available_slots JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(available_slots) = 'array'),
    scorecard_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (scorecard_status IN ('pending', 'submitted', 'waived')),
    waived_by UUID REFERENCES companies(id),
    waived_at TIMESTAMPTZ,
    waiver_reason TEXT,
    PRIMARY KEY (interview_id, interviewer_id),
    CHECK (scorecard_status <> 'waived' OR (
        waived_by IS NOT NULL AND waived_at IS NOT NULL
        AND NULLIF(BTRIM(waiver_reason), '') IS NOT NULL
    ))
);

CREATE TABLE interview_feedback (
    id UUID PRIMARY KEY,
    interview_id UUID NOT NULL,
    interviewer_id UUID NOT NULL,
    notes TEXT,
    notes_updated_at TIMESTAMPTZ,
    ratings JSONB NOT NULL DEFAULT '{}' CHECK (jsonb_typeof(ratings) = 'object'),
    recommendation TEXT,
    submitted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    UNIQUE (interview_id, interviewer_id),
    FOREIGN KEY (interview_id, interviewer_id)
        REFERENCES interviewer_assignments(interview_id, interviewer_id)
);

COMMENT ON COLUMN interview_feedback.notes IS
    'Private hiring-team notes; never included in Candidate feedback.';
