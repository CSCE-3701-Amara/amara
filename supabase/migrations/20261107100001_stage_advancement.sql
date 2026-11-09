-- +goose Up
ALTER TABLE stage_executions DROP CONSTRAINT stage_executions_decision_type_check;
ALTER TABLE stage_executions ADD CONSTRAINT stage_executions_decision_type_check
    CHECK (decision_type IN ('manual', 'threshold', 'completion', 'missed', 'override'));
-- Historical overrides are readable. The decision endpoints never write one.

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
UPDATE stage_executions SET decision_type = NULL WHERE decision_type = 'missed';
ALTER TABLE stage_executions DROP CONSTRAINT stage_executions_decision_type_check;
ALTER TABLE stage_executions ADD CONSTRAINT stage_executions_decision_type_check
    CHECK (decision_type IN ('manual', 'threshold', 'completion', 'override'));
*/
