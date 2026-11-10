-- +goose Up
-- Attempt snapshots contain private evaluator and grading material.
ALTER TABLE assessment_attempts ENABLE ROW LEVEL SECURITY;

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
ALTER TABLE assessment_attempts DISABLE ROW LEVEL SECURITY;
*/
