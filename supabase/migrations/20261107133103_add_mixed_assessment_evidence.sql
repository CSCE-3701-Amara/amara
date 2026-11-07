-- +goose Up
ALTER TABLE assessment_attempts ADD COLUMN non_code_evidence jsonb;

-- +goose Down
ALTER TABLE assessment_attempts DROP COLUMN non_code_evidence;
