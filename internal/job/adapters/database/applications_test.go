package database

import (
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

func TestAdmissionCompanyEligibilityAndConcurrentCapPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	_, err := pool.Exec(ctx, `
 INSERT INTO companies(id,email,password_hash,salt,legal_name,approved)
 VALUES(md5('admission-company')::uuid,'admission@example.test','','','Company',true);
 INSERT INTO candidates(id,name,email,password_hash,salt,created_at,account_status)
 SELECT md5('candidate-'||n)::uuid,'Candidate',n||'@example.test','','',now(),'active' FROM generate_series(1,3)n;
 INSERT INTO job_postings(id,company_id,title,status,published_at,max_applications)
 VALUES(md5('admission-job')::uuid,md5('admission-company')::uuid,'Role','published',now(),1);
 `)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewApplicationRepository(pool)
	var posting, candidate string
	if err := pool.QueryRow(ctx, "SELECT md5('admission-job')::uuid::text,md5('candidate-1')::uuid::text").Scan(&posting, &candidate); err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"banned", "deactivated_at", "deleted_at"} {
		t.Run("reject "+field+" employer", func(t *testing.T) {
			value, clear := "now()", "NULL"
			if field == "banned" {
				value, clear = "true", "false"
			}
			if _, err := pool.Exec(ctx, "UPDATE companies SET "+field+"="+value); err != nil {
				t.Fatal(err)
			}
			if _, err := repo.Apply(ctx, posting, candidate); !errors.Is(err, domain.ErrApplicationPostingUnavailable) {
				t.Errorf("ineligible employer admitted application: %v", err)
			}
			if _, err := pool.Exec(ctx, "DELETE FROM applications"); err != nil {
				t.Fatal(err)
			}
			if _, err := pool.Exec(ctx, "UPDATE companies SET "+field+"="+clear); err != nil {
				t.Fatal(err)
			}
		})
	}
	t.Run("concurrent applications cannot exceed cap", func(t *testing.T) {
		var wg sync.WaitGroup
		errs := make(chan error, 3)
		for n := 1; n <= 3; n++ {
			var id string
			if err := pool.QueryRow(ctx, "SELECT md5('candidate-'||$1::text)::uuid::text", strconv.Itoa(n)).Scan(&id); err != nil {
				t.Fatal(err)
			}
			wg.Go(func() { _, err := repo.Apply(ctx, posting, id); errs <- err })
		}
		wg.Wait()
		close(errs)
		success, full := 0, 0
		for err := range errs {
			if err == nil {
				success++
			} else if errors.Is(err, domain.ErrApplicationLimitReached) {
				full++
			} else {
				t.Errorf("unexpected admission error: %v", err)
			}
		}
		if success != 1 || full != 2 {
			t.Errorf("admissions=%d, cap rejections=%d", success, full)
		}
	})
}
