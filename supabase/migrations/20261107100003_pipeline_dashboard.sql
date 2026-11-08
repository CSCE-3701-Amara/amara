ALTER TABLE job_postings ADD COLUMN finalist_target INTEGER CHECK (finalist_target > 0);
ALTER TABLE applications ADD COLUMN nice_to_have_score INTEGER CHECK (nice_to_have_score BETWEEN 0 AND 100);

CREATE TABLE pipeline_threshold_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stage_id UUID NOT NULL REFERENCES pipeline_stages(id),
    changed_by UUID NOT NULL REFERENCES companies(id),
    previous_threshold INTEGER NOT NULL CHECK (previous_threshold BETWEEN 0 AND 100),
    new_threshold INTEGER NOT NULL CHECK (new_threshold BETWEEN 0 AND 100),
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
