-- +goose Up
-- +goose StatementBegin
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
-- +goose StatementEnd

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP TABLE candidate_education, candidate_work_experience, candidate_languages;
*/
