-- +goose Up
CREATE INDEX candidate_amarat_application_stage_idx
    ON candidate_amarat (application_id, source_stage_execution_id)
    WHERE application_id IS NOT NULL AND source_stage_execution_id IS NOT NULL;

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP INDEX IF EXISTS candidate_amarat_application_stage_idx;
*/
