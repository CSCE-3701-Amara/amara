-- +goose Up
-- Attempt snapshots contain private evaluator and grading material.
ALTER TABLE assessment_attempts ENABLE ROW LEVEL SECURITY;

-- +goose Down
ALTER TABLE assessment_attempts DISABLE ROW LEVEL SECURITY;
