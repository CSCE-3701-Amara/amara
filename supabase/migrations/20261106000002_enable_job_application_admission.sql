-- +goose Up
-- Applications are many-to-many; original UNIQUE constraints allowed one per posting and one per candidate globally.
ALTER TABLE applications DROP CONSTRAINT IF EXISTS applications_posting_id_key;
ALTER TABLE applications DROP CONSTRAINT IF EXISTS applications_cand_id_key;
ALTER TABLE applications ADD CONSTRAINT applications_posting_candidate_key UNIQUE(posting_id,cand_id);
CREATE SEQUENCE IF NOT EXISTS application_reference_number_seq;
SELECT setval('application_reference_number_seq',COALESCE(MAX(reference_number),1),MAX(reference_number) IS NOT NULL) FROM applications;
ALTER TABLE applications ALTER COLUMN reference_number SET DEFAULT nextval('application_reference_number_seq');
ALTER TABLE applications ADD COLUMN IF NOT EXISTS failed_amarat TEXT[] NOT NULL DEFAULT '{}';
CREATE INDEX IF NOT EXISTS applications_posting_queue_idx ON applications(posting_id,submitted_at,id);
-- +goose Down
DROP INDEX IF EXISTS applications_posting_queue_idx;
ALTER TABLE applications DROP COLUMN IF EXISTS failed_amarat;
ALTER TABLE applications ALTER COLUMN reference_number DROP DEFAULT;
DROP SEQUENCE IF EXISTS application_reference_number_seq;
ALTER TABLE applications DROP CONSTRAINT IF EXISTS applications_posting_candidate_key;
ALTER TABLE applications ADD CONSTRAINT applications_posting_id_key UNIQUE(posting_id);
ALTER TABLE applications ADD CONSTRAINT applications_cand_id_key UNIQUE(cand_id);
