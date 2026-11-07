-- +goose Up
ALTER TABLE pipelines ADD CONSTRAINT pipelines_posting_key UNIQUE (posting_id);

ALTER TABLE pipeline_stages
    ADD COLUMN stage_order INTEGER CHECK (stage_order > 0),
    ADD COLUMN interview_configuration JSONB,
    ADD CONSTRAINT pipeline_stages_order_key UNIQUE (version_id, stage_order);

ALTER TABLE test_configurations DROP CONSTRAINT test_configurations_stage_id_key;
ALTER TABLE test_configurations ADD CONSTRAINT test_configurations_stage_module_key
    UNIQUE (stage_id, test_module_id);

-- Stage rubrics and module snapshots contain private grading material.
ALTER TABLE pipeline_stages ENABLE ROW LEVEL SECURITY;
ALTER TABLE test_configurations ENABLE ROW LEVEL SECURITY;

-- +goose Down
ALTER TABLE test_configurations DISABLE ROW LEVEL SECURITY;
ALTER TABLE pipeline_stages DISABLE ROW LEVEL SECURITY;

ALTER TABLE test_configurations DROP CONSTRAINT test_configurations_stage_module_key;
ALTER TABLE test_configurations ADD CONSTRAINT test_configurations_stage_id_key UNIQUE (stage_id);

ALTER TABLE pipeline_stages
    DROP CONSTRAINT pipeline_stages_order_key,
    DROP COLUMN interview_configuration,
    DROP COLUMN stage_order;

ALTER TABLE pipelines DROP CONSTRAINT pipelines_posting_key;
