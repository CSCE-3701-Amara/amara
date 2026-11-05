ALTER TABLE job_postings ADD COLUMN pipeline_id UUID REFERENCES pipelines(id);
