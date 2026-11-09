-- +goose Up
ALTER TABLE assessment_attempts ADD COLUMN non_code_evidence jsonb;

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
ALTER TABLE assessment_attempts DROP COLUMN non_code_evidence;
*/
