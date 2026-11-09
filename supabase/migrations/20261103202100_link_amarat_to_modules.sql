-- +goose Up
ALTER TABLE candidate_amarat ADD COLUMN module_id UUID REFERENCES modules(id);

-- +goose Down
-- Manual rollback reference only; Supabase executes the entire migration file.
/*
ALTER TABLE candidate_amarat DROP COLUMN module_id;
*/
