package application

import (
	"context"
	"errors"
	"net/mail"
	"net/url"
	"strings"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type ProfileService struct{ repository ports.ProfileRepository }

func NewProfileService(repository ports.ProfileRepository) (*ProfileService, error) {
	if repository == nil {
		return nil, errors.New("profile repository must not be nil")
	}
	return &ProfileService{repository: repository}, nil
}

func (s *ProfileService) UpdateCandidate(ctx context.Context, id string, patch domain.CandidateProfilePatch) (domain.CandidateProfile, error) {
	if id == "" || patch.Name != nil && !domain.ValidateProfileText(*patch.Name, 100) ||
		patch.PhoneNumber != nil && len([]rune(*patch.PhoneNumber)) > 50 ||
		patch.Gender != nil && len([]rune(*patch.Gender)) > 50 ||
		patch.Nationality != nil && len([]rune(*patch.Nationality)) > 100 ||
		patch.Address != nil && len([]rune(*patch.Address)) > 500 ||
		patch.Biography != nil && len([]rune(*patch.Biography)) > 5000 ||
		patch.YearsOfExperience != nil && *patch.YearsOfExperience < 0 ||
		!validURLs(patch.GithubProfile, patch.LinkedinProfile, patch.PortfolioLink) {
		return domain.CandidateProfile{}, domain.ErrInvalidProfile
	}
	return s.repository.UpdateCandidateProfile(ctx, id, patch)
}

func (s *ProfileService) UpdateCompany(ctx context.Context, id string, patch domain.CompanyProfilePatch) (domain.CompanyProfile, error) {
	if id == "" || !companyPatchHasFields(patch) || patch.LegalName != nil && !domain.ValidateProfileText(*patch.LegalName, 255) ||
		patch.DisplayName != nil && len([]rune(*patch.DisplayName)) > 255 ||
		patch.Description != nil && len([]rune(*patch.Description)) > 5000 ||
		patch.ContactPhone != nil && len([]rune(*patch.ContactPhone)) > 50 ||
		patch.ContactEmail != nil && !validEmail(*patch.ContactEmail) ||
		!validURLs(patch.WebsiteLink, patch.LogoURL, patch.LinkedinProfile) {
		return domain.CompanyProfile{}, domain.ErrInvalidProfile
	}

	return s.repository.UpdateCompanyProfile(ctx, id, patch)
}

func companyPatchHasFields(p domain.CompanyProfilePatch) bool {
	return p.LegalName != nil || p.DisplayName != nil || p.ClearDisplayName ||
		p.Description != nil || p.ClearDescription || p.WebsiteLink != nil || p.ClearWebsiteLink ||
		p.LogoURL != nil || p.ClearLogoURL || p.ContactPhone != nil || p.ClearContactPhone ||
		p.ContactEmail != nil || p.ClearContactEmail || p.LinkedinProfile != nil || p.ClearLinkedinProfile
}

func (s *ProfileService) UpdateInterviewer(ctx context.Context, id, name string) (domain.InterviewerProfile, error) {
	if id == "" || !domain.ValidateProfileText(name, 100) {
		return domain.InterviewerProfile{}, domain.ErrInvalidProfile
	}
	return s.repository.UpdateInterviewerProfile(ctx, id, strings.TrimSpace(name))
}

func validEmail(value string) bool {
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value && !strings.ContainsAny(value, " \r\n")
}

func validURLs(values ...*string) bool {
	for _, value := range values {
		if value == nil {
			continue
		}
		parsed, err := url.ParseRequestURI(*value)
		if err != nil || parsed.Scheme == "" || parsed.Host == "" {
			return false
		}
	}
	return true
}
