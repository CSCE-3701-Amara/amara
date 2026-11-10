package database

import (
	"context"
	"encoding/json"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
	"testing"
)

type weightedRunner struct{}

func (weightedRunner) Run(_ context.Context, _, _, input string) (string, bool, error) {
	if input == "heavy" {
		return "yes", true, nil
	}
	return "no", true, nil
}
func TestLegacyWeightedGradingAndRecoveryPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	_, err := pool.Exec(ctx, `INSERT INTO candidates(id,name,email,password_hash,salt,created_at) VALUES(md5('c')::uuid,'C','c@example.test','','',now());
 INSERT INTO companies(id,email,password_hash,salt,legal_name) VALUES(md5('co')::uuid,'co@example.test','','','Company');
 INSERT INTO job_postings(id,company_id,title) VALUES(md5('j')::uuid,md5('co')::uuid,'Role');
 INSERT INTO applications(id,posting_id,cand_id,status,screened_out) VALUES(md5('a')::uuid,md5('j')::uuid,md5('c')::uuid,'admitted',false);
 INSERT INTO pipelines(id,posting_id) VALUES(md5('p')::uuid,md5('j')::uuid);
 INSERT INTO pipeline_versions(id,pipeline_id,version_number) VALUES(md5('v')::uuid,md5('p')::uuid,1);
 INSERT INTO pipeline_stages(id,version_id,name,module_type,stage_order,advancement_rule) VALUES(md5('s')::uuid,md5('v')::uuid,'Code','assessment',1,'completion');
 INSERT INTO pipeline_executions(id,application_id,pipeline_version_id) VALUES(md5('pe')::uuid,md5('a')::uuid,md5('v')::uuid);
 INSERT INTO stage_executions(id,pipeline_execution_id,stage_id,status) VALUES(md5('se')::uuid,md5('pe')::uuid,md5('s')::uuid,'pending_grading');
 INSERT INTO assessment_instances(id,application_id,stage_id,status,deadline) VALUES(md5('i')::uuid,md5('a')::uuid,md5('s')::uuid,'pending_grading',now()+interval '1 hour');
 INSERT INTO assessment_attempts(id,stage_execution_id,candidate_id,instance_id,status,started_at,deadline,submitted_at,selected_content,configuration_snapshot,saved_responses)
 VALUES(md5('t')::uuid,md5('se')::uuid,md5('c')::uuid,md5('i')::uuid,'pending_grading',now(),now()+interval '1 hour',now(),
 '[{"module_id":"m","questions":[{"prompt":"Code","format":"code","test_cases":[{"input":"heavy","expected_output":"yes","weight":9},{"input":"light","expected_output":"yes","weight":1}]}]}]',
 '{"modules":[{"id":"m","configuration":{"definition":{"type":"coding","threshold":80}}}]}',jsonb_build_object(md5(md5('t')::uuid::text||':1:1')::uuid::text,jsonb_build_object('language','python','code','print()')));`)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewAssessmentRepository(pool)
	for n := 0; n < 2; n++ {
		if n == 1 {
			_, err = pool.Exec(ctx, `UPDATE assessment_attempts SET final_score=NULL,final_evidence=jsonb_build_array(jsonb_build_object('questionId',md5(id::text||':1:1')::uuid::text,'passedCases',1,'totalCases',2)),code_grading_retry_at=NULL`)
			if err != nil {
				t.Fatal(err)
			}
		}
		claimed, err := repo.GradeNextCodeAttempt(ctx, weightedRunner{})
		if err != nil || !claimed {
			t.Fatalf("grade: %v %v", claimed, err)
		}
		var score int
		var evidence json.RawMessage
		if err := pool.QueryRow(ctx, `SELECT final_score,final_evidence FROM assessment_attempts`).Scan(&score, &evidence); err != nil || score != 90 {
			t.Fatalf("score=%d evidence=%s error=%v", score, evidence, err)
		}
		if claimed, err := repo.GradeNextCodeAttempt(ctx, weightedRunner{}); claimed || err != nil {
			t.Fatalf("repeat grading: %v %v", claimed, err)
		}
	}
	if err := ReconcileAdvancement(ctx, pool); err != nil {
		t.Fatal(err)
	}
	var status string
	if err := pool.QueryRow(ctx, `SELECT status FROM stage_executions`).Scan(&status); err != nil || status != "passed" {
		t.Fatalf("advancement=%s error=%v", status, err)
	}
}
