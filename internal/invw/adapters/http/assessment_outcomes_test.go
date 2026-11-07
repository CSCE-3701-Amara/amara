package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	authdomain "github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	authports "github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
)

type outcomeSessions struct{}

func (outcomeSessions) ValidateSession(_ context.Context, token string) (authports.Identity, error) {
	switch token {
	case "candidate":
		return authports.Identity{AccountType: authdomain.AccountTypeCandidate, AccountID: "00000000-0000-0000-0000-000000000001"}, nil
	case "company":
		return authports.Identity{AccountType: authdomain.AccountTypeCompany, CompanyID: "00000000-0000-0000-0000-000000000002"}, nil
	default:
		return authports.Identity{}, authdomain.ErrInvalidSession
	}
}

func TestAssessmentOutcomePathsPostgres(t *testing.T) {
	dsn := os.Getenv("AMARA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL for the worker and HTTP outcome integration test")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	admin, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close(context.Background())
	schema := pgx.Identifier{fmt.Sprintf("outcomes_%d", time.Now().UnixNano())}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Error(err)
		}
	}()
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	migrations, err := filepath.Glob("../../../../supabase/migrations/*.sql")
	if err != nil || len(migrations) == 0 {
		t.Fatalf("locate migrations: %v", err)
	}
	for _, path := range migrations {
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		up, _, ok := strings.Cut(string(migration), "-- +goose Down")
		if !ok {
			t.Fatalf("missing Down section in %s", path)
		}
		exec(up)
	}
	exec(`
		INSERT INTO candidates (id, name, email, password_hash, salt, created_at)
		VALUES ('00000000-0000-0000-0000-000000000001', 'Candidate', 'candidate@example.test', '', '', now());
		INSERT INTO companies (id, email, password_hash, salt, legal_name)
		VALUES ('00000000-0000-0000-0000-000000000002', 'company@example.test', '', '', 'Company');
		INSERT INTO job_postings (id, company_id, title)
		VALUES ('00000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000002', 'Role');
		INSERT INTO applications (id, posting_id, cand_id, status, submitted_at, screened_out)
		VALUES ('00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000003',
			'00000000-0000-0000-0000-000000000001', 'admitted', now(), false);
		INSERT INTO pipelines (id, posting_id)
		VALUES ('00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000003');
		INSERT INTO pipeline_versions (id, pipeline_id, version_number, is_published)
		VALUES ('00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000005', 1, true);
		UPDATE pipelines SET published_version_id = '00000000-0000-0000-0000-000000000006';
		INSERT INTO pipeline_stages (id, version_id, name, module_type, stage_order, scheduled_start_at, scheduled_end_at)
		SELECT ('00000000-0000-0000-0000-00000000001' || n)::uuid, '00000000-0000-0000-0000-000000000006',
			'Stage ' || n, 'assessment', n, now() - interval '2 hours',
			now() + CASE WHEN n = 1 THEN interval '-1 hour' ELSE interval '1 hour' END
		FROM generate_series(1, 5) n;
		INSERT INTO modules (id, name, type, published)
		VALUES ('00000000-0000-0000-0000-000000000020', 'Module', 'behavioral', true);
		INSERT INTO test_configurations (id, stage_id, test_module_id, configuration)
		SELECT gen_random_uuid(), id, '00000000-0000-0000-0000-000000000020',
			jsonb_build_object('definition', CASE stage_order
				WHEN 3 THEN '{"type":"coding","minimum_tests_passed":1,"questions":[{"prompt":"Code","format":"code","minimum_cases_passed":1,"test_cases":[{"input":"secret"}]}]}'::jsonb
				WHEN 4 THEN '{"type":"recorded_video","behavioral":true,"questions":[{"prompt":"Video","format":"video"}]}'::jsonb
				WHEN 5 THEN '{"type":"recorded_video","questions":[{"prompt":"Video","format":"video"}]}'::jsonb
				ELSE '{"type":"behavioral","questions":[{"prompt":"Describe","format":"long_text"}]}'::jsonb END)
		FROM pipeline_stages;
		INSERT INTO pipeline_stages (id, version_id, name, module_type, stage_order,
			scheduled_start_at, scheduled_end_at, interview_configuration)
		VALUES ('00000000-0000-0000-0000-000000000016', '00000000-0000-0000-0000-000000000006',
			'Interview', 'live_interview', 6, now() - interval '2 hours', now() - interval '1 hour',
			'{"behavioral":true}'),
			('00000000-0000-0000-0000-000000000017', '00000000-0000-0000-0000-000000000006',
			'Scored interview', 'live_interview', 7, now() - interval '2 hours', now() - interval '1 hour', '{}');
	`)
	sweep := func() {
		t.Helper()
		if err := database.SpawnAssessmentInstances(ctx, pool); err != nil {
			t.Fatal(err)
		}
	}
	sweep()
	exec(`
		INSERT INTO stage_executions (id, pipeline_execution_id, stage_id, status, available_at, deadline)
		SELECT gen_random_uuid(), pe.id, s.id, 'available', now(), s.scheduled_end_at
		FROM pipeline_executions pe JOIN pipeline_stages s ON s.version_id = pe.pipeline_version_id
		WHERE s.stage_order > 1
	`)
	sweep()
	sweep()
	exec(`
		INSERT INTO interviews (id, stage_execution_id, scheduled_start_at, scheduled_end_at, status, rubric_snapshot)
		SELECT ('00000000-0000-0000-0000-00000000003' || (s.stage_order - 6))::uuid, se.id,
			now() - interval '90 minutes', now() - interval '60 minutes', 'confirmed', '[{"id":"criterion"}]'
		FROM stage_executions se JOIN pipeline_stages s ON s.id = se.stage_id WHERE s.stage_order IN (6, 7);
		INSERT INTO interviewers (id, company_id, name, email, password_hash, salt, created_at)
		VALUES ('00000000-0000-0000-0000-000000000040', '00000000-0000-0000-0000-000000000002',
			'Interviewer', 'interviewer@example.test', '', '', now());
		INSERT INTO interviewer_assignments (interview_id, interviewer_id, scorecard_status)
		SELECT id, '00000000-0000-0000-0000-000000000040', 'submitted' FROM interviews;
		INSERT INTO interview_feedback (id, interview_id, interviewer_id, ratings, recommendation, submitted_at)
		SELECT gen_random_uuid(), interview_id, interviewer_id, '{"criterion":4}', 'Recommend', now()
		FROM interviewer_assignments;
	`)

	repo := database.NewAssessmentRepository(pool)
	router := gin.New()
	NewAssessmentHandler(repo, outcomeSessions{}).RegisterRoutes(router.Group("/v1"))
	request := func(method, path, token, body string, status int) []byte {
		t.Helper()
		req := httptest.NewRequest(method, "/v1"+path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "amara_session", Value: token})
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
		if res.Code != status {
			t.Fatalf("%s %s: got %d want %d: %s", method, path, res.Code, status, res.Body)
		}
		return res.Body.Bytes()
	}
	load := func(stage int) (id, instance, status string, auto bool) {
		t.Helper()
		err := pool.QueryRow(ctx, `
			SELECT a.id::text, a.instance_id::text, a.status, a.auto_submitted
			FROM assessment_attempts a JOIN stage_executions se ON se.id = a.stage_execution_id
			JOIN pipeline_stages s ON s.id = se.stage_id WHERE s.stage_order = $1
		`, stage).Scan(&id, &instance, &status, &auto)
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	_, _, status, auto := load(1)
	if status != "missed" || auto {
		t.Fatalf("unspawned expiry: %s, auto=%v", status, auto)
	}
	start := func(stage int) (path, id, question string) {
		t.Helper()
		id, instance, _, _ := load(stage)
		raw := request(http.MethodPost, "/assessment-instances/"+instance+"/start", "candidate", "", http.StatusOK)
		var response struct {
			Questions []struct{ ID string }
		}
		if err := json.Unmarshal(raw, &response); err != nil || len(response.Questions) != 1 || strings.Contains(string(raw), "secret") {
			t.Fatalf("unsafe or invalid candidate content: %s (%v)", raw, err)
		}
		return "/assessment-attempts/" + id, id, response.Questions[0].ID
	}
	for _, stage := range []int{2, 4, 5} {
		path, _, question := start(stage)
		answer := `{"response":"example"}`
		if stage != 2 {
			answer = `{"response":{"videoUrl":"https://example.test/recording"}}`
		}
		request(http.MethodPut, path+"/submissions/"+question, "candidate", answer, http.StatusOK)
		first := request(http.MethodPost, path+"/submission", "candidate", "", http.StatusOK)
		_, _, submittedStatus, wasAutomatic := load(stage)
		wantStatus := "submitted"
		if stage == 5 {
			wantStatus = "pending_grading"
		}
		if submittedStatus != wantStatus || wasAutomatic {
			t.Fatalf("stage %d submitted with status=%s, automatic=%v", stage, submittedStatus, wasAutomatic)
		}
		retry := request(http.MethodPost, path+"/submission", "candidate", "", http.StatusOK)
		if string(first) != string(retry) {
			t.Fatalf("submission retry changed attempt: %s vs %s", first, retry)
		}
		evidence := fmt.Sprintf(`{"evidence":[{"questionId":%q,"notes":"reviewed"}]}`, question)
		if stage == 5 {
			evidence = fmt.Sprintf(`{"evidence":[{"questionId":%q,"score":80}]}`, question)
		}
		request(http.MethodPost, path+"/evaluation", "company", evidence, http.StatusNoContent)
		request(http.MethodPost, path+"/evaluation", "company", evidence, http.StatusNoContent)
		review := request(http.MethodGet, path+"/evaluation", "company", "", http.StatusOK)
		want := `"finalScore":null`
		if stage == 5 {
			want = `"finalScore":80`
		}
		if !strings.Contains(string(review), want) {
			t.Fatalf("unexpected review: %s", review)
		}
	}
	codingPath, id, question := start(3)
	exec("UPDATE assessment_attempts SET deadline = now() - interval '1 second' WHERE id = $1::uuid", id)
	sweep()
	_, _, status, auto = load(3)
	if status != "pending_grading" || !auto {
		t.Fatalf("started expiry: %s, auto=%v", status, auto)
	}
	passed, total := 2, 1
	evidence := []database.QuestionEvidence{{QuestionID: question, PassedCases: &passed, TotalCases: &total}}
	if err := repo.FinalizeEvidence(ctx, id, "", evidence); err != database.ErrAssessmentConflict {
		t.Fatalf("invalid evaluator evidence: %v", err)
	}
	_, _, status, _ = load(3)
	if status != "pending_grading" {
		t.Fatalf("invalid evidence changed status to %s", status)
	}
	passed = 1
	for range 2 {
		if err := repo.FinalizeEvidence(ctx, id, "", evidence); err != nil {
			t.Fatal(err)
		}
	}
	_, _, status, _ = load(3)
	if status != "scored" {
		t.Fatalf("expected finalized coding evidence, got %s", status)
	}
	review := request(http.MethodGet, codingPath+"/evaluation", "company", "", http.StatusOK)
	for _, want := range []string{`"finalScore":null`, `"passedCases":1`, `"totalCases":1`, question} {
		if !strings.Contains(string(review), want) {
			t.Fatalf("coding evidence missing %s: %s", want, review)
		}
	}
	for _, suffix := range []string{"0", "1"} {
		interviewID := "00000000-0000-0000-0000-00000000003" + suffix
		for range 2 {
			request(http.MethodPost, "/interviews/"+interviewID+"/outcome", "company", `{"outcome":"held"}`, http.StatusOK)
		}
		sweep()
		var outcome string
		var interviewScore *int
		if err := pool.QueryRow(ctx, "SELECT outcome, score FROM interviews WHERE id = $1::uuid", interviewID).
			Scan(&outcome, &interviewScore); err != nil {
			t.Fatal(err)
		}
		if outcome != "held" || (suffix == "0" && interviewScore != nil) ||
			(suffix == "1" && (interviewScore == nil || *interviewScore != 75)) {
			t.Fatalf("interview %s outcome=%s score=%v", suffix, outcome, interviewScore)
		}
	}
}
