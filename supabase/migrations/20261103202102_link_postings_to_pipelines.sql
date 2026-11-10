-- +goose Up
ALTER TABLE job_postings ADD COLUMN pipeline_id UUID REFERENCES pipelines(id);

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
ALTER TABLE job_postings DROP COLUMN pipeline_id;
*/
