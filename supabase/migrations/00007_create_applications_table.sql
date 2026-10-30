CREATE TABLE applications (
    id                  UUID    PRIMARY KEY,
    reference_number    INTEGER NOT NULL    UNIQUE,
    posting_id          UUID    NOT NULL    UNIQUE  REFERENCES job_postings(id),
    cand_id             UUID    NOT NULL    UNIQUE  REFERENCES candidates(id),
    status              TEXT,
    created_at          TIMESTAMP,
    updated_at          TIMESTAMP,
    submitted_at        TIMESTAMP,
    withdrawn_at        TIMESTAMP,
    final_decision      BOOLEAN,
    final_feedback      TEXT,
    screened_out        BOOLEAN,
    screening_completed_at  TIMESTAMP
);