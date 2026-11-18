package database

import (
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

// Listing loads requirements and Amarat in one batch; each posting must still
// receive exactly its own rows, and a posting without any gets empty slices.
func TestListByCompanyBatchesScreeningDataPostgres(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx := t.Context()
	if _, err := pool.Exec(ctx, `
 INSERT INTO companies(id,email,password_hash,salt,legal_name,approved)
 VALUES(md5('list-company')::uuid,'list@example.test','','','Company',true);
 INSERT INTO job_postings(id,company_id,title,created_at)
 SELECT md5('posting-'||n)::uuid, md5('list-company')::uuid, 'Posting '||n, now() - n * interval '1 minute'
 FROM generate_series(1,3) n;
 INSERT INTO job_posting_requirements(id,posting_id,requirement_text,requirement_type)
 VALUES (gen_random_uuid(), md5('posting-1')::uuid, 'Go', 'must_have'),
        (gen_random_uuid(), md5('posting-1')::uuid, 'SQL', 'nice_to_have'),
        (gen_random_uuid(), md5('posting-2')::uuid, 'React', 'must_have');
 INSERT INTO job_posting_amarat(id,posting_id,amara_title)
 VALUES (gen_random_uuid(), md5('posting-2')::uuid, 'Communication');`); err != nil {
		t.Fatal(err)
	}
	var company string
	if err := pool.QueryRow(ctx, "SELECT md5('list-company')::uuid::text").Scan(&company); err != nil {
		t.Fatal(err)
	}
	postings, err := NewPostingManagementRepository(pool).ListByCompany(ctx, company)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string][2]int{}
	for _, p := range postings {
		if p.Requirements == nil || p.Amarat == nil {
			t.Fatalf("%s: screening slices must be empty, not nil", p.Title)
		}
		got[p.Title] = [2]int{len(p.Requirements), len(p.Amarat)}
	}
	want := map[string][2]int{"Posting 1": {2, 0}, "Posting 2": {1, 1}, "Posting 3": {0, 0}}
	for title, counts := range want {
		if got[title] != counts {
			t.Errorf("%s: got %v requirements/Amarat, want %v", title, got[title], counts)
		}
	}
}
