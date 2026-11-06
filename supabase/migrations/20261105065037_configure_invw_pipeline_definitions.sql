ALTER TABLE modules
    ADD COLUMN source TEXT NOT NULL DEFAULT 'company'
        CHECK (source IN ('platform', 'company')),
    ADD COLUMN completion_window INTERVAL CHECK (completion_window > INTERVAL '0 seconds'),
    ADD COLUMN selection_mode TEXT NOT NULL DEFAULT 'specific'
        CHECK (selection_mode IN ('specific', 'random')),
    ADD COLUMN question_count INTEGER CHECK (question_count > 0),
    ADD COLUMN difficulty_mix JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN amara_definitions JSONB NOT NULL DEFAULT '[]'
        CHECK (jsonb_typeof(amara_definitions) = 'array'),
    ADD COLUMN rubric JSONB NOT NULL DEFAULT '[]'
        CHECK (jsonb_typeof(rubric) = 'array'),
    ADD COLUMN interview_duration_minutes INTEGER CHECK (interview_duration_minutes > 0),
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMPTZ;

-- Definitions include the Amara's name, crest, domain, version, and level bands.
-- Rubric items and questions name exactly one of the configured Amarat.
COMMENT ON COLUMN modules.amara_definitions IS
    'Versioned Amara definitions and level bands; publication validates names, bands and item mappings.';

ALTER TABLE pipelines
    ADD COLUMN anonymity_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN identity_reveal_point TEXT
        CHECK (identity_reveal_point IN ('stage', 'finalist')),
    ADD COLUMN reveal_stage_id UUID REFERENCES pipeline_stages(id),
    ADD COLUMN rejection_feedback_mode TEXT NOT NULL DEFAULT 'ai_with_approval'
        CHECK (rejection_feedback_mode IN ('human', 'ai', 'ai_with_approval')),
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN updated_at TIMESTAMPTZ,
    ADD CONSTRAINT pipelines_identity_reveal_check CHECK (
        NOT anonymity_enabled OR (
            identity_reveal_point IS NOT NULL AND (
                (identity_reveal_point = 'stage' AND reveal_stage_id IS NOT NULL)
                OR (identity_reveal_point = 'finalist' AND reveal_stage_id IS NULL)
            )
        )
    );

-- Hiring Managers currently authenticate through companies; interviewers have
-- their own accounts. These actor references use the existing account tables.
ALTER TABLE pipeline_versions
    ADD COLUMN previous_version_id UUID REFERENCES pipeline_versions(id),
    ADD COLUMN published_by UUID REFERENCES companies(id),
    ADD COLUMN migration_plan JSONB NOT NULL DEFAULT '{}',
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADD COLUMN published_at TIMESTAMPTZ,
    ADD CONSTRAINT pipeline_versions_number_key UNIQUE (pipeline_id, version_number);

CREATE UNIQUE INDEX pipeline_versions_one_published_idx
    ON pipeline_versions (pipeline_id) WHERE is_published IS TRUE;

ALTER TABLE pipeline_stages
    ADD COLUMN advancement_rule TEXT NOT NULL DEFAULT 'manual'
        CHECK (advancement_rule IN ('manual', 'threshold', 'completion')),
    ADD COLUMN score_threshold INTEGER CHECK (score_threshold BETWEEN 0 AND 100),
    ADD COLUMN below_threshold_action TEXT NOT NULL DEFAULT 'hold'
        CHECK (below_threshold_action IN ('hold', 'reject')),
    ADD COLUMN feedback_mode TEXT NOT NULL DEFAULT 'ai_with_approval'
        CHECK (feedback_mode IN ('human', 'ai', 'ai_with_approval')),
    ADD COLUMN scheduled_start_at TIMESTAMPTZ,
    ADD COLUMN scheduled_end_at TIMESTAMPTZ,
    ADD COLUMN position_x NUMERIC NOT NULL DEFAULT 0,
    ADD COLUMN position_y NUMERIC NOT NULL DEFAULT 0,
    ADD CONSTRAINT pipeline_stages_name_key UNIQUE (version_id, name),
    ADD CONSTRAINT pipeline_stages_version_id_key UNIQUE (version_id, id),
    ADD CONSTRAINT pipeline_stages_schedule_check CHECK (scheduled_end_at > scheduled_start_at),
    ADD CONSTRAINT pipeline_stages_threshold_check CHECK (
        advancement_rule <> 'threshold' OR score_threshold IS NOT NULL
    );

CREATE TABLE stage_transitions (
    id UUID PRIMARY KEY,
    pipeline_version_id UUID NOT NULL REFERENCES pipeline_versions(id),
    from_stage_id UUID NOT NULL,
    to_stage_id UUID NOT NULL,
    FOREIGN KEY (pipeline_version_id, from_stage_id) REFERENCES pipeline_stages(version_id, id),
    FOREIGN KEY (pipeline_version_id, to_stage_id) REFERENCES pipeline_stages(version_id, id),
    UNIQUE (pipeline_version_id, from_stage_id, to_stage_id),
    CHECK (from_stage_id <> to_stage_id)
);

CREATE TABLE assessment_questions (
    id UUID PRIMARY KEY,
    test_module_id UUID NOT NULL REFERENCES modules(id),
    prompt TEXT NOT NULL,
    response_format TEXT NOT NULL CHECK (
        response_format IN ('multiple_choice', 'short_text', 'long_text', 'code', 'video')
    ),
    scoring_method TEXT NOT NULL CHECK (
        scoring_method IN ('answer_key', 'code_evaluator', 'human')
    ),
    max_points NUMERIC NOT NULL CHECK (max_points > 0),
    answer_key JSONB,
    options JSONB,
    test_cases JSONB,
    difficulty TEXT,
    amara_name TEXT NOT NULL,
    grader_id UUID REFERENCES interviewers(id),
    grader_company_id UUID REFERENCES companies(id),
    ai_generated BOOLEAN NOT NULL DEFAULT FALSE,
    reviewed_by UUID REFERENCES companies(id),
    reviewed_at TIMESTAMPTZ,
    position INTEGER NOT NULL DEFAULT 0 CHECK (position >= 0),
    configuration JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ,
    CHECK (num_nonnulls(grader_id, grader_company_id) <= 1)
);

CREATE TABLE test_configurations (
    id UUID PRIMARY KEY,
    stage_id UUID NOT NULL UNIQUE REFERENCES pipeline_stages(id),
    test_module_id UUID NOT NULL REFERENCES modules(id),
    active_time_limit INTERVAL CHECK (active_time_limit > INTERVAL '0 seconds'),
    completion_window INTERVAL CHECK (completion_window > INTERVAL '0 seconds'),
    configuration JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ
);
