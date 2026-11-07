package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/cand/domain"
)

func TestApplicationTrackingPostgres(t *testing.T) {
	dsn := os.Getenv("AMARA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL to run PostgreSQL tracking tests")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	admin, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = admin.Close(context.Background()) }()
	schema := pgx.Identifier{fmt.Sprintf("tracking_test_%d", time.Now().UnixNano())}.Sanitize()
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if _, err := admin.Exec(context.Background(), "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Errorf("remove test schema: %v", err)
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

	migrations, err := filepath.Glob("../../../../supabase/migrations/*.sql")
	if err != nil || len(migrations) == 0 {
		t.Fatalf("find migrations: %v", err)
	}
	for _, path := range migrations {
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		up, _, ok := strings.Cut(string(migration), "-- +goose Down")
		if !ok {
			t.Fatalf("migration %s has no Down section", path)
		}
		if _, err := pool.Exec(ctx, up); err != nil {
			t.Fatalf("apply %s: %v", path, err)
		}
	}

	const candidate = "00000000-0000-0000-0000-000000000001"
	const other = "00000000-0000-0000-0000-000000000002"
	const rejected = "00000000-0000-0000-0000-000000000031"
	_, err = pool.Exec(ctx, `
		INSERT INTO candidates (id, name, email, password_hash, salt, created_at)
		VALUES ('00000000-0000-0000-0000-000000000001', 'Owner', 'owner@example.test', '', '', now()),
		       ('00000000-0000-0000-0000-000000000002', 'Other', 'other@example.test', '', '', now());
		INSERT INTO companies (id, email, password_hash, salt, legal_name)
		VALUES ('00000000-0000-0000-0000-000000000010', 'company@example.test', '', '', 'Company');
		INSERT INTO job_postings (id, company_id, title)
		SELECT ('00000000-0000-0000-0000-00000000002' || n)::uuid,
		       '00000000-0000-0000-0000-000000000010', 'Role ' || n
		FROM generate_series(1, 3) AS n;
		INSERT INTO applications (id, posting_id, cand_id, status, submitted_at, screened_out, final_feedback)
		VALUES
		('00000000-0000-0000-0000-000000000031', '00000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000001', 'admitted', now(), false, 'Private until released'),
		('00000000-0000-0000-0000-000000000032', '00000000-0000-0000-0000-000000000022', '00000000-0000-0000-0000-000000000001', 'admitted', now() - interval '1 day', false, NULL),
		('00000000-0000-0000-0000-000000000033', '00000000-0000-0000-0000-000000000023', '00000000-0000-0000-0000-000000000001', 'admitted', now() - interval '2 days', true, NULL),
		('00000000-0000-0000-0000-000000000034', '00000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000002', 'admitted', now(), false, NULL);
		INSERT INTO pipelines (id, posting_id, anonymity_enabled, identity_reveal_point)
		VALUES ('00000000-0000-0000-0000-000000000040', '00000000-0000-0000-0000-000000000021', true, 'finalist');
		INSERT INTO pipeline_versions (id, version_number, pipeline_id)
		VALUES ('00000000-0000-0000-0000-000000000041', 1, '00000000-0000-0000-0000-000000000040');
		INSERT INTO pipeline_executions (id, application_id, pipeline_version_id, status, rejection_feedback, feedback_released_at)
		VALUES
		('00000000-0000-0000-0000-000000000051', '00000000-0000-0000-0000-000000000031', '00000000-0000-0000-0000-000000000041', 'rejected', '{"summary":"Try again"}', now() + interval '1 hour'),
		('00000000-0000-0000-0000-000000000052', '00000000-0000-0000-0000-000000000032', '00000000-0000-0000-0000-000000000041', 'active', NULL, NULL);
	`)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewApplicationTrackingRepository(pool)
	page, err := repo.List(ctx, candidate, domain.ApplicationFilter{Limit: 3})
	if err != nil {
		t.Fatal(err)
	}
	if page.Pagination.Total != 3 || len(page.Data) != 3 {
		t.Fatalf("unexpected history: %+v", page)
	}
	for i, want := range []struct {
		status string
		closed bool
	}{{"rejected", true}, {"in_progress", false}, {"screened_out", true}} {
		if row := page.Data[i]; row.Status != want.status || row.Closed != want.closed {
			t.Fatalf("row %d = %+v, want %+v", i, row, want)
		}
	}
	for _, offset := range []int{0, 1} {
		page, err := repo.List(ctx, candidate, domain.ApplicationFilter{Limit: 1, Offset: offset, Status: "rejected"})
		if err != nil || page.Pagination.Total != 1 || len(page.Data) != 1-offset ||
			(offset == 0 && page.Data[0].ID != rejected) {
			t.Fatalf("filtered history at offset %d: %+v, %v", offset, page, err)
		}
	}

	if _, err := repo.Get(ctx, other, rejected); !errors.Is(err, domain.ErrApplicationNotFound) {
		t.Fatalf("cross-candidate detail error = %v, want not found", err)
	}
	detail, err := repo.Get(ctx, candidate, rejected)
	if err != nil || !detail.IdentityHidden || detail.PipelineExecutionID == nil ||
		detail.FinalFeedback != nil || detail.RejectionFeedback != nil {
		t.Fatalf("detail before release: %+v, %v", detail, err)
	}

	_, err = pool.Exec(ctx, `
		UPDATE pipeline_executions SET feedback_released_at = now() - interval '1 minute' WHERE application_id = $1::uuid
	`, rejected)
	if err != nil {
		t.Fatal(err)
	}
	detail, err = repo.Get(ctx, candidate, rejected)
	if err != nil || detail.FinalFeedback == nil || *detail.FinalFeedback != "Private until released" ||
		!strings.Contains(string(detail.RejectionFeedback), "Try again") {
		t.Fatalf("detail after release: %+v, %v", detail, err)
	}
}
