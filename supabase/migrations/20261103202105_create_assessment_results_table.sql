-- +goose Up
CREATE TABLE assessment_results (
    id                      UUID        PRIMARY KEY,
    instance_id             UUID        NOT NULL    UNIQUE REFERENCES assessment_instances(id),
    passed                  BOOLEAN     NOT NULL,
    score                   INTEGER,
    feedback                TEXT,
    is_disclosed            BOOLEAN
);

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP TABLE assessment_results;
*/
