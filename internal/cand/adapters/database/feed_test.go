package database

import (
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/cand/ports"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

func TestFeedSalaryFiltersAndCompanyEligibilityPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	_, err := pool.Exec(ctx, `
 INSERT INTO companies (id,email,password_hash,salt,legal_name,approved)
 VALUES (md5('feed-company')::uuid,'feed@example.test','','','Feed company',true);
 INSERT INTO job_postings (id,company_id,title,status,published_at,salary_min,salary_max,salary_visible)
 VALUES (md5('visible')::uuid,md5('feed-company')::uuid,'Visible','published',now(),40000,70000,true),
 (md5('hidden')::uuid,md5('feed-company')::uuid,'Hidden','published',now(),40000,70000,false),
 (md5('unset')::uuid,md5('feed-company')::uuid,'Unset','published',now(),40000,70000,NULL);
 `)
	if err != nil {
		t.Fatal(err)
	}
	repo := NewJobFeedRepository(pool)
	for _, filter := range []ports.JobFeedQuery{{Limit: 10, SalaryMin: new(int32(45000))}, {Limit: 10, SalaryMax: new(int32(65000))}} {
		result, err := repo.ListPublished(ctx, filter)
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Postings) != 1 || result.Postings[0].Title != "Visible" {
			t.Errorf("salary filter revealed hidden ranges: %+v", result.Postings)
		}
	}
	for _, field := range []string{"banned", "deactivated_at", "deleted_at"} {
		t.Run(field, func(t *testing.T) {
			value, clear := "now()", "NULL"
			if field == "banned" {
				value, clear = "true", "false"
			}
			if _, err := pool.Exec(ctx, "UPDATE companies SET "+field+"="+value); err != nil {
				t.Fatal(err)
			}
			result, err := repo.ListPublished(ctx, ports.JobFeedQuery{Limit: 10})
			if err != nil {
				t.Fatal(err)
			}
			if len(result.Postings) != 0 {
				t.Errorf("%s company still in feed: %d postings", field, len(result.Postings))
			}
			if _, err := pool.Exec(ctx, "UPDATE companies SET "+field+"="+clear); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// "%" and "_" typed into the search box are literal characters, not wildcards.
func TestFeedTitleSearchMatchesWildcardsLiterallyPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	if _, err := pool.Exec(ctx, `
 INSERT INTO companies (id,email,password_hash,salt,legal_name,approved)
 VALUES (md5('search-company')::uuid,'search@example.test','','','Search company',true);
 INSERT INTO job_postings (id,company_id,title,status,published_at)
 VALUES (md5('percent')::uuid,md5('search-company')::uuid,'100% remote','published',now()),
 (md5('plain')::uuid,md5('search-company')::uuid,'1000 offices','published',now()),
 (md5('under')::uuid,md5('search-company')::uuid,'snake_case role','published',now()),
 (md5('other')::uuid,md5('search-company')::uuid,'snakeXcase role','published',now());`); err != nil {
		t.Fatal(err)
	}
	repo := NewJobFeedRepository(pool)
	for query, want := range map[string]string{"100%": "100% remote", "snake_case": "snake_case role"} {
		result, err := repo.ListPublished(ctx, ports.JobFeedQuery{Limit: 10, Title: query})
		if err != nil {
			t.Fatal(err)
		}
		if len(result.Postings) != 1 || result.Postings[0].Title != want {
			t.Errorf("search %q returned %+v, want only %q", query, result.Postings, want)
		}
	}
}
