CREATE TABLE job_posting_requirements (
    id                  UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id          UUID    NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    requirement_text    TEXT    NOT NULL,
    requirement_type    TEXT    NOT NULL CHECK (requirement_type IN ('must_have', 'nice_to_have')),
    created_at          TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX job_posting_requirements_posting_idx ON job_posting_requirements(posting_id);

CREATE TABLE job_posting_amarat (
    id                  UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    posting_id          UUID    NOT NULL REFERENCES job_postings(id) ON DELETE CASCADE,
    amara_title         TEXT    NOT NULL,
    amara_description   TEXT,
    created_at          TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX job_posting_amarat_posting_idx ON job_posting_amarat(posting_id);
CREATE UNIQUE INDEX job_posting_amarat_title_idx ON job_posting_amarat(posting_id, lower(amara_title));
