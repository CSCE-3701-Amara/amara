CREATE TABLE candidate_documents (
    cand_id     UUID        REFERENCES candidates(id),
    docu_name   TEXT,
    doc_type    TEXT,
    doc_size    INTEGER,
    doc_url     TEXT        NOT NULL,
    uploaded_at TIMESTAMP
);

CREATE TABLE candidate_certifications (
    cand_id             UUID    REFERENCES candidates(id),
    certification_name  TEXT    NOT NULL,
    certification_type  TEXT,
    verification_link   TEXT
);

CREATE TABLE candidate_amarat (
    cand_id             UUID    REFERENCES candidates(id),
    amara_title         TEXT    NOT NULL,
    amara_description   TEXT
);
