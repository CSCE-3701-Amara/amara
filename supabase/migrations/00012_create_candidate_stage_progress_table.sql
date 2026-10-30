-- +goose Up
CREATE TABLE candidate_stage_progress (
    id              UUID        PRIMARY KEY,
    application_id  UUID        NOT NULL    REFERENCES applications(id),
    stage_id        UUID        NOT NULL    REFERENCES pipeline_stages(id),
    status          TEXT        NOT NULL,
    amara_awarded   TEXT
);

-- +goose Down
DROP TABLE candidate_stage_progress;