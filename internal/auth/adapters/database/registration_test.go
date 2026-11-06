package database

import (
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

func TestAccountActionTable(t *testing.T) {
	tests := []struct {
		accountType domain.AccountType
		table       string
		sessionType string
	}{
		{domain.AccountTypeCandidate, "candidates", "candidate"},
		{domain.AccountTypeInterviewer, "interviewers", "interviewer"},
		{domain.AccountTypeCompany, "companies", "company"},
	}
	for _, test := range tests {
		table, sessionType := accountActionTable(test.accountType)
		if table != test.table || sessionType != test.sessionType {
			t.Fatalf("accountActionTable(%q) = %q, %q; want %q, %q", test.accountType, table, sessionType, test.table, test.sessionType)
		}
	}

	if table, sessionType := accountActionTable("admin"); table != "" || sessionType != "" {
		t.Fatalf("accountActionTable(admin) = %q, %q; want empty values", table, sessionType)
	}
}
