-- +goose Up
CREATE TABLE job_postings (
    id                      UUID    PRIMARY KEY,
    company_id              UUID    NOT NULL    REFERENCES companies(id),
    title                   TEXT    NOT NULL,
    description             TEXT,
    language                TEXT,
    currency                TEXT,
    seniority_level         TEXT,
    employment_type         TEXT,
    location                TEXT,
    salary_min              INTEGER,
    salary_max              INTEGER,
    salary_visible          BOOLEAN,
    open_positions          INTEGER,
    application_opening     TIMESTAMP,
    application_deadline    TIMESTAMP,
    max_applications        INTEGER,
    status                  TEXT,
    last_autosaved_at       TIMESTAMP,
    published_at            TIMESTAMP,
    archived_at             TIMESTAMP,
    created_at              TIMESTAMP,
    updated_at              TIMESTAMP,
    deleted_at              TIMESTAMP
);

-- +goose Down
DROP TABLE job_postings;