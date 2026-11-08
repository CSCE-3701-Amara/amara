CREATE TABLE IF NOT EXISTS candidate_education (
    cand_id     UUID    REFERENCES candidates(id),
    degree_type TEXT,
    provider    TEXT,
    start_date  DATE,
    end_date    DATE,
    summary     TEXT
);
CREATE TABLE IF NOT EXISTS candidate_work_experience (
    cand_id         UUID    REFERENCES candidates(id),
    position        TEXT,
    company_name    TEXT,
    start_date      DATE,
    end_date        DATE,
    summary         TEXT
);
CREATE TABLE IF NOT EXISTS candidate_languages (
    cand_id     UUID    REFERENCES candidates(id),
    language    TEXT    NOT NULL,
    level       TEXT
);
