-- Stable namespace, insert-only reruns; published graphs are never rewritten.
SELECT pg_advisory_xact_lock(hashtext('amara-demo-parallel-v1'));
INSERT INTO job_postings(id,company_id,title,description,language,currency,seniority_level,employment_type,location,salary_min,salary_max,salary_visible,open_positions,application_opening,application_deadline,max_applications,finalist_target,status,published_at,created_at,updated_at)
SELECT md5('amara-demo-parallel-job-'||n)::uuid,md5('amara-demo-company-v1')::uuid,
 CASE n WHEN 1 THEN 'Full-stack Engineer · Parallel assessments' ELSE 'Platform Engineer · Parallel assessments' END,
 'Synthetic demo: candidates complete coding and collaboration assessments in parallel. Both branches must pass before the final review opens.',
 'English','EGP','Mid-level','remote full-time','Cairo',30000,50000,true,3,
 now()-interval '1 day',now()+interval '45 days',100,5,'published',now(),now(),now()
FROM generate_series(1,2) n ON CONFLICT(id) DO NOTHING;
INSERT INTO pipelines(id,posting_id,status)
SELECT md5('amara-demo-parallel-pipeline-'||n)::uuid,md5('amara-demo-parallel-job-'||n)::uuid,'draft'
FROM generate_series(1,2) n ON CONFLICT(id) DO NOTHING;
INSERT INTO pipeline_versions(id,pipeline_id,version_number,is_published)
SELECT md5('amara-demo-parallel-version-'||n)::uuid,md5('amara-demo-parallel-pipeline-'||n)::uuid,1,false
FROM generate_series(1,2) n ON CONFLICT(id) DO NOTHING;
INSERT INTO pipeline_stages(id,version_id,name,module_type,advancement_rule,stage_order,position_x,position_y,scheduled_start_at,scheduled_end_at)
SELECT md5('amara-demo-parallel-stage-'||n||'-'||s)::uuid,md5('amara-demo-parallel-version-'||n)::uuid,
 (ARRAY['Coding assessment','Collaboration assessment','Final review (both branches required)'])[s],
 'assessment',CASE WHEN s=3 THEN 'manual' ELSE 'completion' END,s,
 CASE WHEN s=3 THEN 400 ELSE 40 END,CASE WHEN s=2 THEN 260 ELSE 60 END,
 now()-interval '1 day',now()+interval '30 days'
FROM generate_series(1,2) n CROSS JOIN generate_series(1,3) s ON CONFLICT(id) DO NOTHING;
INSERT INTO stage_transitions(id,pipeline_version_id,from_stage_id,to_stage_id)
SELECT md5('amara-demo-parallel-edge-'||n||'-'||s)::uuid,md5('amara-demo-parallel-version-'||n)::uuid,
 md5('amara-demo-parallel-stage-'||n||'-'||s)::uuid,md5('amara-demo-parallel-stage-'||n||'-3')::uuid
FROM generate_series(1,2) n CROSS JOIN generate_series(1,2) s ON CONFLICT(id) DO NOTHING;
INSERT INTO test_configurations(id,stage_id,test_module_id,active_time_limit,configuration)
SELECT md5('amara-demo-parallel-config-'||n||'-'||s)::uuid,md5('amara-demo-parallel-stage-'||n||'-'||s)::uuid,
 m.id,interval '30 minutes',jsonb_build_object('position',0,'definition',m.authoring_definition)
FROM generate_series(1,2) n CROSS JOIN generate_series(1,3) s JOIN modules m
ON m.id=CASE WHEN s=1 THEN md5('amara-demo-module-code-v1')::uuid ELSE md5('amara-demo-module-behavior-v1')::uuid END
ON CONFLICT(id) DO NOTHING;
UPDATE pipeline_versions SET is_published=true,published_at=now(),published_by=md5('amara-demo-company-v1')::uuid
WHERE id IN (SELECT md5('amara-demo-parallel-version-'||n)::uuid FROM generate_series(1,2) n) AND is_published IS NOT TRUE;
UPDATE pipelines SET published_version_id=md5('amara-demo-parallel-version-'||n)::uuid,status='published'
FROM generate_series(1,2) n WHERE id=md5('amara-demo-parallel-pipeline-'||n)::uuid AND published_version_id IS NULL;
UPDATE job_postings SET pipeline_id=md5('amara-demo-parallel-pipeline-'||n)::uuid
FROM generate_series(1,2) n WHERE id=md5('amara-demo-parallel-job-'||n)::uuid AND pipeline_id IS NULL;
INSERT INTO applications(id,posting_id,cand_id,status,created_at,updated_at,submitted_at,screened_out,screening_completed_at)
SELECT md5('amara-demo-parallel-application-'||j||'-'||n)::uuid,md5('amara-demo-parallel-job-'||j)::uuid,
 md5('amara-demo-candidate-v1-'||n)::uuid,'admitted',now(),now(),now(),false,now()
FROM generate_series(1,2) j CROSS JOIN generate_series(1,6) n ON CONFLICT(id) DO NOTHING;
-- The worker creates the two entry-stage executions and assessment snapshots.
DO $$ BEGIN
 IF (SELECT count(*) FROM stage_transitions WHERE pipeline_version_id IN
 (SELECT md5('amara-demo-parallel-version-'||n)::uuid FROM generate_series(1,2) n)) <> 4 THEN
 RAISE EXCEPTION 'Expected four parallel-branch connections'; END IF;
END $$;
