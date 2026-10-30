-- +goose Up
-- +goose StatementBegin
CREATE TABLE pipeline_versions (
    id              UUID        PRIMARY KEY,
    version_number  INTEGER     NOT NULL,
    is_published    BOOLEAN
);

CREATE TABLE pipelines (
    id                      UUID    PRIMARY KEY,
    posting_id              UUID    NOT NULL    REFERENCES job_postings(id),
    published_version_id    UUID    REFERENCES pipeline_versions(id),
    status                  TEXT
);

ALTER TABLE pipeline_versions ADD COLUMN pipeline_id UUID REFERENCES pipelines(id);

CREATE TABLE pipeline_stages (
    id              UUID        PRIMARY KEY,
    version_id      UUID        NOT NULL    REFERENCES pipeline_versions(id),
    name            TEXT        NOT NULL,
    module_type     TEXT        NOT NULL
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE pipeline_versions, pipelines, pipeline_stages;