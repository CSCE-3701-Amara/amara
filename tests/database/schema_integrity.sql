-- Run after all Goose Up sections. Fixtures and changes are rolled back.
\set ON_ERROR_STOP on
BEGIN;

CREATE FUNCTION pg_temp.expect_schema_failure(statement TEXT, expected_state TEXT)
RETURNS VOID LANGUAGE plpgsql AS $$
DECLARE actual_state TEXT;
BEGIN
    BEGIN
        EXECUTE statement;
    EXCEPTION WHEN OTHERS THEN
        GET STACKED DIAGNOSTICS actual_state = RETURNED_SQLSTATE;
        IF actual_state = expected_state THEN RETURN; END IF;
        RAISE;
    END;
    RAISE EXCEPTION 'Expected SQLSTATE % for %', expected_state, statement;
END;
$$;

INSERT INTO companies (id, email, password_hash, salt, legal_name)
VALUES (md5('schema-test:company-a')::uuid, 'schema-company-a@example.test', '\x01', '\x02', 'Company A'),
       (md5('schema-test:company-b')::uuid, 'schema-company-b@example.test', '\x01', '\x02', 'Company B');

INSERT INTO roles (id, tenant_id, name)
VALUES (md5('schema-test:role-a')::uuid, md5('schema-test:company-a')::uuid, 'hiring_manager'),
       (md5('schema-test:role-b')::uuid, md5('schema-test:company-b')::uuid, 'hiring_manager'),
       (md5('schema-test:platform-role')::uuid, NULL, 'platform_admin');
INSERT INTO permissions (id, name) VALUES (md5('schema-test:permission')::uuid, 'schema-test:postings.manage');
INSERT INTO role_permissions (role_id, permission_id)
VALUES (md5('schema-test:role-a')::uuid, md5('schema-test:permission')::uuid);

INSERT INTO users (id, tenant_id, role_id, email, name)
VALUES (md5('schema-test:user-a')::uuid, md5('schema-test:company-a')::uuid, md5('schema-test:role-a')::uuid, 'schema-user-a@example.test', 'Manager A'),
       (md5('schema-test:user-b')::uuid, md5('schema-test:company-b')::uuid, md5('schema-test:role-b')::uuid, 'schema-user-b@example.test', 'Manager B'),
       (md5('schema-test:candidate-user')::uuid, NULL, NULL, 'schema-candidate@example.test', 'Candidate'),
       (md5('schema-test:admin-user')::uuid, NULL, md5('schema-test:platform-role')::uuid, 'schema-admin@example.test', 'Administrator');

SELECT pg_temp.expect_schema_failure($q$UPDATE users SET role_id = md5('schema-test:role-b')::uuid WHERE id = md5('schema-test:user-a')::uuid$q$, '23503');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO users (email, name) VALUES ('SCHEMA-USER-A@example.test', 'Duplicate')$q$, '23505');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO roles (name) VALUES ('platform_admin')$q$, '23505');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO role_permissions VALUES (md5('schema-test:role-a')::uuid, md5('schema-test:permission')::uuid, CURRENT_TIMESTAMP)$q$, '23505');

UPDATE companies SET hiring_manager_user_id = md5('schema-test:user-a')::uuid WHERE id = md5('schema-test:company-a')::uuid;
SELECT pg_temp.expect_schema_failure($q$UPDATE companies SET hiring_manager_user_id = md5('schema-test:user-b')::uuid WHERE id = md5('schema-test:company-a')::uuid$q$, '23503');

INSERT INTO auth_sessions (user_id, token_hash, expires_at)
VALUES (md5('schema-test:user-a')::uuid, '\x112233', CURRENT_TIMESTAMP + INTERVAL '1 day');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO auth_sessions (user_id, token_hash, expires_at) VALUES (md5('schema-test:user-a')::uuid, '\x445566', CURRENT_TIMESTAMP - INTERVAL '1 day')$q$, '23514');
INSERT INTO company_settings (tenant_id) VALUES (md5('schema-test:company-a')::uuid);
SELECT pg_temp.expect_schema_failure($q$INSERT INTO company_settings (tenant_id) VALUES (md5('schema-test:company-a')::uuid)$q$, '23505');
INSERT INTO company_verifications (tenant_id, reviewed_by, status, verification_method, reviewed_at)
VALUES (md5('schema-test:company-a')::uuid, md5('schema-test:user-a')::uuid, 'approved', 'manual', CURRENT_TIMESTAMP);
SELECT pg_temp.expect_schema_failure($q$INSERT INTO company_verifications (tenant_id, status, verification_method, reviewed_at) VALUES (md5('schema-test:company-a')::uuid, 'approved', 'manual', CURRENT_TIMESTAMP)$q$, '23514');
INSERT INTO system_administrators (user_id) VALUES (md5('schema-test:admin-user')::uuid);

INSERT INTO invitations (tenant_id, invited_by, role_id, email, token_hash, expires_at)
VALUES (md5('schema-test:company-a')::uuid, md5('schema-test:user-a')::uuid, md5('schema-test:role-a')::uuid, 'schema-invite@example.test', '\xaabbcc', CURRENT_TIMESTAMP + INTERVAL '7 days');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO invitations (tenant_id, invited_by, email, token_hash, expires_at) VALUES (md5('schema-test:company-a')::uuid, md5('schema-test:user-b')::uuid, 'schema-invalid-invite@example.test', '\xaabbdd', CURRENT_TIMESTAMP + INTERVAL '7 days')$q$, '23503');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO invitations (tenant_id, invited_by, email, token_hash, expires_at) VALUES (md5('schema-test:company-a')::uuid, md5('schema-test:user-a')::uuid, 'SCHEMA-INVITE@example.test', '\xaabbee', CURRENT_TIMESTAMP + INTERVAL '7 days')$q$, '23505');

INSERT INTO candidates (id, name, email, password_hash, salt, created_at, user_id)
VALUES (md5('schema-test:candidate-a')::uuid, 'Candidate A', 'schema-candidate-a@example.test', '\x01', '\x02', CURRENT_TIMESTAMP, md5('schema-test:candidate-user')::uuid),
       (md5('schema-test:candidate-b')::uuid, 'Candidate B', 'schema-candidate-b@example.test', '\x01', '\x02', CURRENT_TIMESTAMP, NULL);
INSERT INTO candidate_education (cand_id, degree_type) VALUES (md5('schema-test:candidate-a')::uuid, 'BSc');
INSERT INTO candidate_work_experience (cand_id, position) VALUES (md5('schema-test:candidate-a')::uuid, 'Engineer');
INSERT INTO candidate_languages (cand_id, language) VALUES (md5('schema-test:candidate-a')::uuid, 'English');
INSERT INTO candidate_certifications (cand_id, certification_name) VALUES (md5('schema-test:candidate-a')::uuid, 'Certificate');
INSERT INTO candidate_documents (cand_id, doc_url) VALUES (md5('schema-test:candidate-a')::uuid, 'https://example.test/resume.pdf');

INSERT INTO resumes (id, candidate_id, title, is_default)
VALUES (md5('schema-test:resume')::uuid, md5('schema-test:candidate-a')::uuid, 'Dynamic resume', TRUE);
SELECT pg_temp.expect_schema_failure($q$INSERT INTO resumes (candidate_id, title, is_default) VALUES (md5('schema-test:candidate-a')::uuid, 'Another default', TRUE)$q$, '23505');
INSERT INTO resume_versions (id, resume_id, version_number, content)
VALUES (md5('schema-test:resume-version')::uuid, md5('schema-test:resume')::uuid, 1, '{"experience_years":3}');
SELECT pg_temp.expect_schema_failure($q$UPDATE resume_versions SET content = '{}' WHERE id = md5('schema-test:resume-version')::uuid$q$, '23514');

INSERT INTO skills (id, name) VALUES (md5('schema-test:skill')::uuid, 'Schema test programming');
INSERT INTO candidate_amarat (id, cand_id, amara_title)
VALUES (md5('schema-test:award')::uuid, md5('schema-test:candidate-a')::uuid, 'Programming');
INSERT INTO candidate_skills (candidate_id, skill_id, amara_award_id, source, verified_at)
VALUES (md5('schema-test:candidate-a')::uuid, md5('schema-test:skill')::uuid, md5('schema-test:award')::uuid, 'assessed', CURRENT_TIMESTAMP);
SELECT pg_temp.expect_schema_failure($q$INSERT INTO candidate_skills (candidate_id, skill_id, amara_award_id, source, verified_at) VALUES (md5('schema-test:candidate-b')::uuid, md5('schema-test:skill')::uuid, md5('schema-test:award')::uuid, 'assessed', CURRENT_TIMESTAMP)$q$, '23503');

INSERT INTO job_postings (id, company_id, created_by, title)
VALUES (md5('schema-test:posting-a')::uuid, md5('schema-test:company-a')::uuid, md5('schema-test:user-a')::uuid, 'Engineer A'),
       (md5('schema-test:posting-b')::uuid, md5('schema-test:company-a')::uuid, md5('schema-test:user-a')::uuid, 'Engineer B');
SELECT pg_temp.expect_schema_failure($q$UPDATE job_postings SET created_by = md5('schema-test:user-b')::uuid WHERE id = md5('schema-test:posting-a')::uuid$q$, '23503');
INSERT INTO job_requirements (job_id, criterion_type, description, field_name, comparison_operator, expected_value)
VALUES (md5('schema-test:posting-a')::uuid, 'must_have', 'Two years experience', 'years_of_experience', 'gte', '2');
SELECT pg_temp.expect_schema_failure($q$INSERT INTO job_requirements (job_id, criterion_type, description) VALUES (md5('schema-test:posting-a')::uuid, 'must_have', 'Unstructured criterion')$q$, '23514');
INSERT INTO job_skills (job_id, skill_id) VALUES (md5('schema-test:posting-a')::uuid, md5('schema-test:skill')::uuid);

-- Both many-candidate and many-posting cases must succeed.
INSERT INTO applications (id, posting_id, cand_id, reference_number)
VALUES (md5('schema-test:application-a')::uuid, md5('schema-test:posting-a')::uuid, md5('schema-test:candidate-a')::uuid, -1301001),
       (md5('schema-test:application-b')::uuid, md5('schema-test:posting-a')::uuid, md5('schema-test:candidate-b')::uuid, -1301002),
       (md5('schema-test:application-c')::uuid, md5('schema-test:posting-b')::uuid, md5('schema-test:candidate-a')::uuid, -1301003);
SELECT pg_temp.expect_schema_failure($q$INSERT INTO applications (id, posting_id, cand_id, reference_number) VALUES (md5('schema-test:duplicate-app')::uuid, md5('schema-test:posting-a')::uuid, md5('schema-test:candidate-a')::uuid, -1301004)$q$, '23505');
SELECT pg_temp.expect_schema_failure($q$UPDATE applications SET job_id = md5('schema-test:posting-b')::uuid WHERE id = md5('schema-test:application-a')::uuid$q$, '428C9');

INSERT INTO application_snapshots (application_id, resume_version_id, resume_snapshot, answers)
VALUES (md5('schema-test:application-a')::uuid, md5('schema-test:resume-version')::uuid, '{"experience_years":3}', '{"available":true}');
SELECT pg_temp.expect_schema_failure($q$UPDATE application_snapshots SET answers = '{}' WHERE application_id = md5('schema-test:application-a')::uuid$q$, '23514');
INSERT INTO screening_decisions (application_id, outcome)
VALUES (md5('schema-test:application-a')::uuid, 'admitted');
INSERT INTO screening_decisions (application_id, outcome, failed_criteria)
VALUES (md5('schema-test:application-b')::uuid, 'screened_out', '[{"field":"years_of_experience","required":2}]');
SELECT pg_temp.expect_schema_failure($q$UPDATE screening_decisions SET failed_criteria = '[]' WHERE application_id = md5('schema-test:application-b')::uuid$q$, '23514');
INSERT INTO candidate_history (candidate_id, application_id, event_type)
VALUES (md5('schema-test:candidate-a')::uuid, md5('schema-test:application-a')::uuid, 'admitted');

INSERT INTO pipelines (id, posting_id, tenant_id)
VALUES (md5('schema-test:pipeline')::uuid, md5('schema-test:posting-a')::uuid, md5('schema-test:company-a')::uuid);
SELECT pg_temp.expect_schema_failure($q$UPDATE pipelines SET tenant_id = md5('schema-test:company-b')::uuid WHERE id = md5('schema-test:pipeline')::uuid$q$, '23503');

DO $$
BEGIN
    IF EXISTS (SELECT FROM applications WHERE job_id IS DISTINCT FROM posting_id OR candidate_id IS DISTINCT FROM cand_id)
        OR EXISTS (SELECT FROM candidate_documents WHERE candidate_id IS DISTINCT FROM cand_id)
        OR EXISTS (SELECT FROM candidate_amarat WHERE candidate_id IS DISTINCT FROM cand_id)
        OR EXISTS (SELECT FROM job_postings WHERE tenant_id IS DISTINCT FROM company_id)
        OR EXISTS (SELECT FROM pipelines WHERE job_id IS DISTINCT FROM posting_id) THEN
        RAISE EXCEPTION 'Generated compatibility aliases are inconsistent';
    END IF;
END;
$$;

ROLLBACK;
\echo 'PASS: schema integrity checks and workflow fixtures; all test changes rolled back.'
