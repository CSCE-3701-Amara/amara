ALTER TABLE modules
    ADD COLUMN authoring_definition JSONB,
    ADD COLUMN copied_from_module_id UUID REFERENCES modules(id),
    ADD CONSTRAINT authored_module_owner_check CHECK (
        authoring_definition IS NULL OR (company_id IS NOT NULL AND published IS NOT NULL)
    );

-- Module definitions contain private grading material. The backend uses
-- the database owner connection and enforces company authorization itself.
ALTER TABLE modules ENABLE ROW LEVEL SECURITY;
