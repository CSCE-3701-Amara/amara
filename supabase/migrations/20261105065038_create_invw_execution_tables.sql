CREATE TABLE pipeline_executions (
    id UUID PRIMARY KEY,
    application_id UUID NOT NULL UNIQUE REFERENCES applications(id),
    pipeline_version_id UUID NOT NULL REFERENCES pipeline_versions(id),
    status TEXT NOT NULL DEFAULT 'active'
        CHECK (status IN ('active', 'completed', 'rejected', 'withdrawn', 'cancelled')),
    anonymous_identifier UUID NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    started_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMPTZ,
    rejection_feedback JSONB,
    feedback_is_ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    feedback_approved_at TIMESTAMPTZ,
    feedback_released_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);

CREATE TABLE stage_executions (
    id UUID PRIMARY KEY,
    pipeline_execution_id UUID NOT NULL REFERENCES pipeline_executions(id),
    stage_id UUID NOT NULL REFERENCES pipeline_stages(id),
    status TEXT NOT NULL DEFAULT 'locked' CHECK (
        status IN ('locked', 'available', 'in_progress', 'pending_grading', 'held',
                   'passed', 'rejected', 'cancelled', 'missed')
    ),
    available_at TIMESTAMPTZ,
    deadline TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    score INTEGER CHECK (score BETWEEN 0 AND 100),
    decision_type TEXT CHECK (decision_type IN ('manual', 'threshold', 'completion', 'override')),
    decided_by UUID REFERENCES companies(id),
    decided_at TIMESTAMPTZ,
    decision_reason TEXT,
    score_at_decision INTEGER CHECK (score_at_decision BETWEEN 0 AND 100),
    decision_history JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(decision_history) = 'array'),
    feedback TEXT,
    feedback_is_ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    feedback_approved_at TIMESTAMPTZ,
    feedback_released_at TIMESTAMPTZ,
    UNIQUE (pipeline_execution_id, stage_id),
    CHECK (decision_type <> 'override' OR NULLIF(BTRIM(decision_reason), '') IS NOT NULL)
);

-- Keep the existing assessment instance and progress records, linking the new
-- execution model to them rather than replacing or deleting historical data.
ALTER TABLE assessment_instances
    ADD CONSTRAINT assessment_instances_application_stage_key UNIQUE (application_id, stage_id);

CREATE TABLE assessment_attempts (
    id UUID PRIMARY KEY,
    stage_execution_id UUID NOT NULL UNIQUE REFERENCES stage_executions(id),
    candidate_id UUID NOT NULL REFERENCES candidates(id),
    instance_id UUID NOT NULL UNIQUE REFERENCES assessment_instances(id),
    status TEXT NOT NULL DEFAULT 'available' CHECK (
        status IN ('available', 'in_progress', 'submitted', 'pending_grading', 'scored', 'missed')
    ),
    selected_content JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(selected_content) = 'array'),
    configuration_snapshot JSONB NOT NULL DEFAULT '{}',
    started_at TIMESTAMPTZ,
    deadline TIMESTAMPTZ NOT NULL,
    submitted_at TIMESTAMPTZ,
    auto_submitted BOOLEAN NOT NULL DEFAULT FALSE,
    extension_history JSONB NOT NULL DEFAULT '[]' CHECK (jsonb_typeof(extension_history) = 'array')
);

CREATE TABLE assessment_submissions (
    id UUID PRIMARY KEY,
    attempt_id UUID NOT NULL REFERENCES assessment_attempts(id),
    question_id UUID NOT NULL REFERENCES assessment_questions(id),
    response JSONB,
    max_points NUMERIC NOT NULL CHECK (max_points > 0),
    earned_points NUMERIC CHECK (earned_points >= 0 AND earned_points <= max_points),
    grading_status TEXT NOT NULL DEFAULT 'pending'
        CHECK (grading_status IN ('pending', 'evaluation_pending', 'graded')),
    grader_id UUID REFERENCES interviewers(id),
    grader_company_id UUID REFERENCES companies(id),
    graded_at TIMESTAMPTZ,
    evaluator_result JSONB,
    saved_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (attempt_id, question_id),
    CHECK (num_nonnulls(grader_id, grader_company_id) <= 1)
);

ALTER TABLE assessment_results
    ADD COLUMN attempt_id UUID UNIQUE REFERENCES assessment_attempts(id),
    ADD COLUMN earned_points NUMERIC CHECK (earned_points >= 0),
    ADD COLUMN max_points NUMERIC CHECK (max_points > 0),
    ADD COLUMN feedback_is_ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN feedback_approved_at TIMESTAMPTZ,
    ADD COLUMN disclosed_at TIMESTAMPTZ,
    ADD COLUMN scored_at TIMESTAMPTZ,
    ADD CONSTRAINT assessment_results_score_check CHECK (score BETWEEN 0 AND 100),
    ADD CONSTRAINT assessment_results_points_check CHECK (earned_points <= max_points);

ALTER TABLE candidate_stage_progress
    ADD COLUMN stage_execution_id UUID UNIQUE REFERENCES stage_executions(id);

CREATE INDEX pipeline_executions_version_idx ON pipeline_executions (pipeline_version_id);
CREATE INDEX stage_executions_stage_status_idx ON stage_executions (stage_id, status);
