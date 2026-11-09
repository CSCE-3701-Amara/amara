-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION sync_assessment_stage_execution() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    UPDATE stage_executions se
    SET status = CASE
            WHEN NEW.status IN ('submitted', 'scored') THEN 'held'
            ELSE NEW.status
        END,
        started_at = NEW.started_at,
        deadline = NEW.deadline,
        completed_at = CASE WHEN NEW.status = 'missed' THEN NEW.deadline ELSE NEW.submitted_at END,
        score = NEW.final_score
    FROM pipeline_executions pe
    WHERE se.id = NEW.stage_execution_id
        AND pe.id = se.pipeline_execution_id AND pe.status = 'active'
        AND se.status IN ('available', 'in_progress', 'pending_grading', 'held', 'missed')
        AND se.decided_at IS NULL;
    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- A held assessment is ready for advancement to apply the configured decision.
CREATE TRIGGER assessment_attempt_sync_stage
AFTER INSERT OR UPDATE OF status, deadline ON assessment_attempts
FOR EACH ROW EXECUTE FUNCTION sync_assessment_stage_execution();

-- Replay the projection for attempts recorded before the trigger existed.
UPDATE assessment_attempts SET status = status;

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
DROP TRIGGER assessment_attempt_sync_stage ON assessment_attempts;
DROP FUNCTION sync_assessment_stage_execution();
*/
