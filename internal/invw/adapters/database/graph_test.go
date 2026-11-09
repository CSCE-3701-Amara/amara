package database

import (
	"errors"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
	"testing"
	"time"
)

func TestGraphPersistencePostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	const company = "00000000-0000-0000-0000-000000000010"
	const job = "00000000-0000-0000-0000-000000000020"
	_, err := pool.Exec(ctx, `INSERT INTO companies(id,email,password_hash,salt,legal_name) VALUES($1::uuid,'graph@example.test','','','Graph')`, company)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO job_postings(id,company_id,title) VALUES($1::uuid,$2::uuid,'Role')`, job, company); err != nil {
		t.Fatal(err)
	}
	modules := NewModuleRepository(pool)
	m, err := modules.Create(ctx, company, domain.Definition{Name: "Questions", Type: domain.Behavioral, Questions: []domain.Question{{Prompt: "Describe", Format: "long_text"}}})
	if err != nil {
		t.Fatal(err)
	}
	repo := NewPipelineRepository(pool)
	p, err := repo.Create(ctx, company, job)
	if err != nil {
		t.Fatal(err)
	}
	g := domain.Graph{Stages: []domain.Stage{}, Links: []domain.Link{{From: "a", To: "b"}, {From: "a", To: "c"}, {From: "b", To: "d"}, {From: "c", To: "d"}}}
	for i, id := range []string{"a", "b", "c", "d"} {
		g.Stages = append(g.Stages, domain.Stage{ID: id, Name: id, Type: domain.AssessmentStage, ModuleIDs: []string{m.ID}, OpensAt: time.Now(), Deadline: time.Now().Add(time.Hour), X: float64(i * 240), Y: 100})
	}
	if _, err := repo.ReplaceStages(ctx, company, p.VersionID, append([]domain.Stage{}, g.Stages...)); err != nil {
		t.Fatal(err)
	}
	legacy, err := repo.Graph(ctx, company, p.VersionID)
	if err != nil || len(legacy.Links) != 3 || legacy.Stages[0].X == legacy.Stages[1].X {
		t.Fatalf("legacy layout: %+v (%v)", legacy, err)
	}
	saved, err := repo.SaveGraph(ctx, company, p.VersionID, g)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := repo.Graph(ctx, company, p.VersionID)
	if err != nil || len(loaded.Links) != 4 || loaded.Stages[1].X != 240 {
		t.Fatalf("reload %+v: %v", loaded, err)
	}
	if _, err := repo.Graph(ctx, "00000000-0000-0000-0000-000000000099", p.VersionID); !errors.Is(err, domain.ErrPipelineNotFound) {
		t.Fatalf("foreign company: %v", err)
	}
	invalid := saved
	invalid.Links = append(invalid.Links, domain.Link{From: saved.Stages[3].ID, To: saved.Stages[0].ID})
	if _, err := repo.SaveGraph(ctx, company, p.VersionID, invalid); err == nil {
		t.Fatal("accepted cycle")
	}
	if err := repo.Publish(ctx, company, p.VersionID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SaveGraph(ctx, company, p.VersionID, saved); !errors.Is(err, domain.ErrPipelineLocked) {
		t.Fatalf("published graph edit: %v", err)
	}

	// Exercise published graph admission and real assessment/decision transitions.
	const candidate = "00000000-0000-0000-0000-000000000030"
	if _, err := pool.Exec(ctx, `INSERT INTO candidates(id,name,email,password_hash,salt,created_at) VALUES($1::uuid,'Candidate','dag@example.test','','',now())`, candidate); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `INSERT INTO applications(id,posting_id,cand_id,status,screened_out) VALUES(gen_random_uuid(),$1::uuid,$2::uuid,'admitted',false)`, job, candidate); err != nil {
		t.Fatal(err)
	}
	if err := SpawnAssessmentInstances(ctx, pool); err != nil {
		t.Fatal(err)
	}
	assessments := NewAssessmentRepository(pool)
	pass := func(index int) {
		t.Helper()
		var attempt, execution, instance string
		if err := pool.QueryRow(ctx, `SELECT a.id::text,se.id::text,a.instance_id::text FROM assessment_attempts a JOIN stage_executions se ON se.id=a.stage_execution_id WHERE se.stage_id=$1::uuid`, saved.Stages[index].ID).Scan(&attempt, &execution, &instance); err != nil {
			t.Fatal(err)
		}
		if _, err := assessments.CandidateAttempt(ctx, candidate, instance, true, false); err != nil {
			t.Fatal(err)
		}
		if _, err := assessments.CandidateAttempt(ctx, candidate, attempt, false, true); err != nil {
			t.Fatal(err)
		}
		if _, err := assessments.DecideStages(ctx, []string{execution}, company, saved.Stages[index].ID, "pass", "Reviewed", ""); err != nil {
			t.Fatal(err)
		}
		if err := SpawnAssessmentInstances(ctx, pool); err != nil {
			t.Fatal(err)
		}
	}
	available := func(index int) bool {
		t.Helper()
		var yes bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM stage_executions WHERE stage_id=$1::uuid AND status='available')`, saved.Stages[index].ID).Scan(&yes); err != nil {
			t.Fatal(err)
		}
		return yes
	}
	if !available(0) || available(1) || available(2) || available(3) {
		t.Fatal("incorrect graph roots")
	}
	pass(0)
	if !available(1) || !available(2) || available(3) {
		t.Fatal("branch did not open correctly")
	}
	pass(1)
	if available(3) {
		t.Fatal("join opened before both parents passed")
	}
	pass(2)
	if !available(3) {
		t.Fatal("join did not open after both parents passed")
	}
}
