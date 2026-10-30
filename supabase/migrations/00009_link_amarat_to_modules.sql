-- +goose Up
ALTER TABLE amarat ADD COLUMN module_id UUID REFERENCES companies(id);

-- +goose Down
ALTER TABLE amarat DROP COLUMN module_id;