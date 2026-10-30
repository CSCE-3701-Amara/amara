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
DROP TABLE assessment_results;