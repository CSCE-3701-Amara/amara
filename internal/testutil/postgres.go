// Package testutil provides disposable PostgreSQL fixtures for integration tests.
package testutil

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres applies the exact deployment SQL in an isolated schema. It never
// loads application credentials or changes an existing application's tables.
func Postgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("AMARA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set AMARA_TEST_DATABASE_URL for PostgreSQL integration tests")
	}
	ctx := t.Context()
	admin, err := pgx.Connect(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	schema := pgx.Identifier{fmt.Sprintf("audit_test_%d", time.Now().UnixNano())}.Sanitize()
	// Generated schema name is quoted with pgx.Identifier.Sanitize above.
	// nosemgrep: go.lang.security.audit.sqli.pgx-sqli.pgx-sqli
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		_ = admin.Close(context.Background())
		t.Fatal(err)
	}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Same generated, quoted identifier; no request data enters this SQL.
		// nosemgrep: go.lang.security.audit.sqli.pgx-sqli.pgx-sqli
		if _, err := admin.Exec(cleanup, "DROP SCHEMA "+schema+" CASCADE"); err != nil {
			t.Error(err)
		}
		_ = admin.Close(cleanup)
	})
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	_, file, _, _ := runtime.Caller(0)
	migrationRoot, err := os.OpenRoot(filepath.Join(filepath.Dir(file), "../../supabase/migrations"))
	if err != nil {
		t.Fatal(err)
	}
	defer migrationRoot.Close()
	migrations, err := fs.Glob(migrationRoot.FS(), "*.sql")
	if err != nil || len(migrations) == 0 {
		t.Fatalf("locate migrations: %v", err)
	}
	for _, path := range migrations {
		sql, err := migrationRoot.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			t.Fatalf("apply %s: %v", filepath.Base(path), err)
		}
	}
	return pool
}
