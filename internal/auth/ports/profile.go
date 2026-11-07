package ports

import (
	"context"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type ProfileRepository interface {
	UpdateCandidateProfile(context.Context, string, domain.CandidateProfilePatch) (domain.CandidateProfile, error)
	UpdateCompanyProfile(context.Context, string, domain.CompanyProfilePatch) (domain.CompanyProfile, error)
	UpdateInterviewerProfile(context.Context, string, string) (domain.InterviewerProfile, error)
}
