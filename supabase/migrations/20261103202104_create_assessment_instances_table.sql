-- +goose Up
CREATE TABLE assessment_instances (
    id              UUID        PRIMARY KEY,
    application_id  UUID        NOT NULL    REFERENCES applications(id),
    stage_id        UUID        NOT NULL    REFERENCES pipeline_stages(id),
    status          TEXT        NOT NULL,
    started_at      TIMESTAMP,
    deadline        TIMESTAMP
);

-- +goose Down
DROP TABLE assessment_instances;
