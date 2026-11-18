package http

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/invw/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	jobdb "github.com/CSCE-3701-Amara/amara/internal/job/adapters/database"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

func TestPipelineGraphPostgres(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 45*time.Second)
	defer cancel()
	pool := testutil.Postgres(t)
	exec := func(sql string, args ...any) {
		t.Helper()
		if _, err := pool.Exec(ctx, sql, args...); err != nil {
			t.Fatal(err)
		}
	}
	exec(`
 INSERT INTO companies(id,email,password_hash,salt,legal_name) VALUES ('00000000-0000-0000-0000-000000000002','graph@example.test','','','Graph');
 INSERT INTO job_postings(id,company_id,title) VALUES ('00000000-0000-0000-0000-000000000003','00000000-0000-0000-0000-000000000002','Role');
 INSERT INTO modules(id,company_id,name,type,published,authoring_definition) VALUES ('00000000-0000-0000-0000-000000000020','00000000-0000-0000-0000-000000000002','Questions','behavioral',true,'{"name":"Questions","type":"behavioral","activeMinutes":10,"questions":[]}');
 `)
	const company = "00000000-0000-0000-0000-000000000002"
	const posting = "00000000-0000-0000-0000-000000000003"
	repo := database.NewPipelineRepository(pool)
	d, err := repo.Create(ctx, company, posting)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Publish(ctx, company, d.VersionID); err == nil {
		t.Fatal("published empty draft")
	}
	g := domain.Graph{Links: []domain.Link{{From: "a", To: "c"}, {From: "b", To: "c"}}}
	for i, id := range []string{"a", "b", "c"} {
		g.Stages = append(g.Stages, domain.Stage{ID: id, Name: id, Type: "assessment", OpensAt: time.Now().Add(-time.Hour), Deadline: time.Now().Add(time.Hour), ModuleIDs: []string{"00000000-0000-0000-0000-000000000020"}, X: float64(i*300 + 100), Y: 120})
	}
	saved, err := repo.SaveGraph(ctx, company, d.VersionID, g)
	if err != nil {
		t.Fatal(err)
	}
	read, err := repo.Graph(ctx, company, d.VersionID)
	if err != nil {
		t.Fatal(err)
	}
	if len(read.Stages) != 3 || len(read.Links) != 2 || read.Stages[1].X != 400 {
		t.Fatalf("roundtrip lost graph: %+v", read)
	}
	if saved.Links[0].From != saved.Stages[0].ID || saved.Links[0].To != saved.Stages[2].ID {
		t.Fatal("connections did not map to saved stage IDs")
	}
	if _, err := repo.Graph(ctx, "00000000-0000-0000-0000-000000000099", d.VersionID); !errors.Is(err, domain.ErrPipelineNotFound) {
		t.Fatalf("ownership: %v", err)
	}
	if _, err := repo.ReplaceStages(ctx, company, d.VersionID, read.Stages); !errors.Is(err, domain.ErrPipelineLocked) {
		t.Fatal("linear API destroyed graph")
	}

	invalid := domain.Graph{Stages: read.Stages, Links: []domain.Link{{From: read.Stages[0].ID, To: read.Stages[0].ID}}}
	if _, err := repo.SaveGraph(ctx, company, d.VersionID, invalid); err == nil {
		t.Fatal("saved cyclic graph")
	}
	unchanged, err := repo.Graph(ctx, company, d.VersionID)
	if err != nil || len(unchanged.Links) != 2 {
		t.Fatal("failed save destroyed previous graph")
	}
	if err := repo.Publish(ctx, company, d.VersionID); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.SaveGraph(ctx, company, d.VersionID, read); !errors.Is(err, domain.ErrPipelineLocked) {
		t.Fatal("published graph was editable")
	}
	exec(`
 INSERT INTO candidates(id,name,email,password_hash,salt,created_at) VALUES ('00000000-0000-0000-0000-000000000001','Candidate','graph-candidate@example.test','','',now());
 INSERT INTO applications(id,posting_id,cand_id,status,submitted_at,screened_out) VALUES ('00000000-0000-0000-0000-000000000004','00000000-0000-0000-0000-000000000003','00000000-0000-0000-0000-000000000001','admitted',now(),false);
 `)
	if err := database.SpawnAssessmentInstances(ctx, pool); err != nil {
		t.Fatal(err)
	}
	var entries int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM stage_executions").Scan(&entries); err != nil {
		t.Fatal(err)
	}
	if entries != 2 {
		t.Fatalf("expected two parallel entry stages, got %d", entries)
	}
	copy, err := repo.NewDraft(ctx, company, posting)
	if err != nil {
		t.Fatal(err)
	}
	cloned, err := repo.Graph(ctx, company, copy.VersionID)
	if err != nil {
		t.Fatal(err)
	}
	if copy.VersionNumber != 2 || len(cloned.Links) != 2 || cloned.Stages[0].ID == read.Stages[0].ID {
		t.Fatal("copy lost graph or reused stage IDs")
	}

	exec("UPDATE test_configurations SET completion_window=interval '2 hours',active_time_limit=interval '7 minutes' WHERE stage_id=$1::text::uuid", cloned.Stages[0].ID)
	cloned.Stages[0].Name = "Renamed stage"
	cloned, err = repo.SaveGraph(ctx, company, copy.VersionID, cloned)
	if err != nil {
		t.Fatal(err)
	}
	var preserved bool
	if err = pool.QueryRow(ctx, "SELECT completion_window=interval '2 hours' AND active_time_limit=interval '7 minutes' FROM test_configurations WHERE stage_id=$1::text::uuid", cloned.Stages[0].ID).Scan(&preserved); err != nil || !preserved {
		t.Fatalf("saving graph lost assessment timing: %v", err)
	}
	if err := repo.Publish(ctx, company, copy.VersionID); err != nil {
		t.Fatal(err)
	}
	var pinned string
	if err := pool.QueryRow(ctx, "SELECT pipeline_version_id::text FROM pipeline_executions").Scan(&pinned); err != nil {
		t.Fatal(err)
	}
	if pinned != d.VersionID {
		t.Fatal("publication moved an existing application")
	}
	if _, err := repo.SaveGraph(ctx, company, d.VersionID, read); !errors.Is(err, domain.ErrPipelineLocked) {
		t.Fatal("superseded version became editable")
	}
	current, err := repo.Get(ctx, company, posting)
	if err != nil || !current.Published || current.VersionID != copy.VersionID {
		t.Fatalf("publication state: %+v %v", current, err)
	}

	// Deleting the posting ends its pipeline: active runs are cancelled and the
	// pipeline can no longer be read or edited through the posting.
	if err := jobdb.NewPostingManagementRepository(pool).SoftDelete(ctx, company, posting, false); err != nil {
		t.Fatal(err)
	}
	var active int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM pipeline_executions WHERE status='active'").Scan(&active); err != nil || active != 0 {
		t.Fatalf("deleted posting kept %d active executions: %v", active, err)
	}
	if _, err := repo.Get(ctx, company, posting); !errors.Is(err, domain.ErrPipelineNotFound) {
		t.Fatalf("deleted posting pipeline still readable: %v", err)
	}
	if _, err := repo.NewDraft(ctx, company, posting); !errors.Is(err, domain.ErrPipelineNotFound) {
		t.Fatalf("deleted posting pipeline still editable: %v", err)
	}
}
