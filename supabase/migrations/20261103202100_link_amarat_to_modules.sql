-- +goose Up
ALTER TABLE candidate_amarat ADD COLUMN module_id UUID REFERENCES modules(id);

-- +goose Down
ALTER TABLE candidate_amarat DROP COLUMN module_id;
