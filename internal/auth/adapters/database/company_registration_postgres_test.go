package database

import (
	"testing"
	"time"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
	"github.com/CSCE-3701-Amara/amara/internal/testutil"
)

func TestCompanyRegistrationApprovalPolicyPostgres(t *testing.T) {
	for _, autoApprove := range []bool{false, true} {
		name, wantAccountStatus, wantRequestStatus := "manual", domain.PendingReview, domain.PendingReview
		if autoApprove {
			name, wantAccountStatus, wantRequestStatus = "automatic", "active", "approved"
		}
		t.Run(name, func(t *testing.T) {
			pool := testutil.Postgres(t)
			repo, err := NewRegistrationRepository(pool)
			if err != nil {
				t.Fatal(err)
			}
			record := ports.CompanyRegistrationRecord{
				CompanyName: "Example Company", RegistrationNumber: "REG-APPROVAL",
				Email: "company@example.test", PasswordHash: []byte("hash"), Salt: []byte("salt"),
				CreatedAt: time.Now().UTC(), AutoApprove: autoApprove,
			}
			account, err := repo.CreateCompanyRegistration(t.Context(), record)
			if err != nil {
				t.Fatal(err)
			}
			if account.Status != wantAccountStatus || account.Email != record.Email {
				t.Fatalf("registration returned %+v; want status %s and email %s", account, wantAccountStatus, record.Email)
			}
			var status string
			var reviewed, activeCompany bool
			err = pool.QueryRow(t.Context(), `SELECT status, reviewed_at IS NOT NULL,
 EXISTS(SELECT 1 FROM companies WHERE id = $1::uuid AND approved)
 FROM auth_company_registration_requests WHERE id = $1::uuid`, account.ID).Scan(&status, &reviewed, &activeCompany)
			if err != nil {
				t.Fatal(err)
			}
			if status != wantRequestStatus || reviewed != autoApprove || activeCompany != autoApprove {
				t.Fatalf("persisted status=%s reviewed=%t activeCompany=%t; want %s, %t, %t",
					status, reviewed, activeCompany, wantRequestStatus, autoApprove, autoApprove)
			}
		})
	}
}
