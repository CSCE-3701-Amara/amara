ALTER TABLE applications
    ADD CONSTRAINT applications_id_posting_unique UNIQUE (id, posting_id);

ALTER TABLE job_postings
    ADD CONSTRAINT job_postings_id_company_unique UNIQUE (id, company_id);

CREATE TABLE pipeline_versions (
    id              UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id      UUID        NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    version_number  INTEGER     NOT NULL CHECK (version_number > 0),
    status          TEXT        NOT NULL DEFAULT 'draft'
                                CHECK (status IN ('draft', 'published', 'retired')),
    published_at    TIMESTAMP,
    created_at      TIMESTAMP   NOT NULL DEFAULT now(),
    CHECK (status <> 'published' OR published_at IS NOT NULL),
    UNIQUE (posting_id, version_number),
    UNIQUE (id, posting_id)
);

CREATE UNIQUE INDEX pipeline_versions_one_published_per_posting
    ON pipeline_versions (posting_id) WHERE status = 'published';

CREATE TABLE pipeline_stages (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    pipeline_version_id UUID        NOT NULL REFERENCES pipeline_versions(id) ON DELETE CASCADE,
    name                TEXT        NOT NULL,
    stage_type          TEXT        NOT NULL,
    position            INTEGER     NOT NULL DEFAULT 0,
    UNIQUE (pipeline_version_id, id),
    UNIQUE (pipeline_version_id, name)
);

CREATE TABLE pipeline_links (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    pipeline_version_id UUID        NOT NULL REFERENCES pipeline_versions(id) ON DELETE CASCADE,
    from_stage_id       UUID        NOT NULL,
    to_stage_id         UUID        NOT NULL,
    CHECK (from_stage_id <> to_stage_id),
    UNIQUE (pipeline_version_id, from_stage_id, to_stage_id),
    FOREIGN KEY (pipeline_version_id, from_stage_id)
        REFERENCES pipeline_stages(pipeline_version_id, id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_version_id, to_stage_id)
        REFERENCES pipeline_stages(pipeline_version_id, id) ON DELETE CASCADE
);

CREATE TABLE application_pipeline_assignments (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id      UUID        NOT NULL UNIQUE,
    posting_id          UUID        NOT NULL,
    pipeline_version_id UUID        NOT NULL,
    assigned_at         TIMESTAMP   NOT NULL DEFAULT now(),
    UNIQUE (application_id, pipeline_version_id),
    FOREIGN KEY (application_id, posting_id)
        REFERENCES applications(id, posting_id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_version_id, posting_id)
        REFERENCES pipeline_versions(id, posting_id)
);

CREATE TABLE application_stage_progress (
    application_id      UUID        NOT NULL REFERENCES application_pipeline_assignments(application_id) ON DELETE CASCADE,
    pipeline_version_id UUID        NOT NULL,
    stage_id            UUID        NOT NULL,
    state               TEXT        NOT NULL CHECK (state IN ('active', 'completed')),
    passed              BOOLEAN,
    completed_at        TIMESTAMP,
    PRIMARY KEY (application_id, stage_id),
    FOREIGN KEY (application_id, pipeline_version_id)
        REFERENCES application_pipeline_assignments(application_id, pipeline_version_id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_version_id, stage_id)
        REFERENCES pipeline_stages(pipeline_version_id, id) ON DELETE CASCADE,
    CHECK (
        (state = 'active' AND passed IS NULL AND completed_at IS NULL)
        OR (state = 'completed' AND passed IS NOT NULL AND completed_at IS NOT NULL)
    )
);

CREATE INDEX application_stage_progress_stage_idx
    ON application_stage_progress (stage_id, application_id);

CREATE TABLE amara_awards (
    id                  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    application_id      UUID        NOT NULL REFERENCES application_pipeline_assignments(application_id) ON DELETE CASCADE,
    pipeline_version_id UUID        NOT NULL,
    stage_id            UUID        NOT NULL,
    posting_id          UUID        NOT NULL,
    company_id          UUID        NOT NULL REFERENCES companies(id),
    amara_name          TEXT        NOT NULL,
    crest               TEXT        NOT NULL,
    skill_domain        TEXT        NOT NULL,
    level               TEXT,
    level_description   TEXT,
    strengths           JSONB       NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(strengths) = 'array'),
    development_areas  JSONB       NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(development_areas) = 'array'),
    evidence            JSONB       NOT NULL DEFAULT '[]'::jsonb CHECK (jsonb_typeof(evidence) = 'array'),
    awarded_at          TIMESTAMP   NOT NULL DEFAULT now(),
    FOREIGN KEY (application_id, pipeline_version_id)
        REFERENCES application_pipeline_assignments(application_id, pipeline_version_id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_version_id, stage_id)
        REFERENCES pipeline_stages(pipeline_version_id, id) ON DELETE CASCADE,
    FOREIGN KEY (pipeline_version_id, posting_id)
        REFERENCES pipeline_versions(id, posting_id) ON DELETE CASCADE,
    FOREIGN KEY (posting_id, company_id)
        REFERENCES job_postings(id, company_id)
);

CREATE INDEX amara_awards_application_stage_idx
    ON amara_awards (application_id, stage_id, awarded_at);
