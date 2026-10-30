package domain

import "testing"

func TestMFARequired(t *testing.T) {
	if !MFARequired(AccountTypeInterviewer) || !MFARequired(AccountTypeCompany) {
		t.Fatal("interviewers and companies must require MFA")
	}
	if MFARequired(AccountTypeCandidate) {
		t.Fatal("candidate MFA should remain optional")
	}
}
