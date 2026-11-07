package ports

import (
	"context"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type ProfileRepository interface {
	GetCompanyProfile(context.Context, string, bool) (domain.CompanyProfile, error)
	ListCompanyInterviewers(context.Context, string, int32, int32) ([]domain.InterviewerProfile, int64, error)
	GetInterviewerProfile(context.Context, string) (domain.InterviewerProfile, error)
	GetPublicCompanyProfile(context.Context, string) (domain.CompanyProfile, error)
	UpdateCandidateProfile(context.Context, string, domain.CandidateProfilePatch) (domain.CandidateProfile, error)
	UpdateCompanyProfile(context.Context, string, domain.CompanyProfilePatch) (domain.CompanyProfile, error)
	UpdateInterviewerProfile(context.Context, string, string) (domain.InterviewerProfile, error)
}
