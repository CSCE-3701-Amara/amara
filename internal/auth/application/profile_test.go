package application

import (
	"context"
	"errors"
	"testing"

	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type profileRepository struct {
	candidate       domain.CandidateProfilePatch
	company         domain.CompanyProfilePatch
	interviewerName string
	err             error
}

func (r *profileRepository) UpdateCandidateProfile(_ context.Context, _ string, patch domain.CandidateProfilePatch) (domain.CandidateProfile, error) {
	r.candidate = patch
	return domain.CandidateProfile{Name: "Mona"}, r.err
}
func (r *profileRepository) UpdateCompanyProfile(_ context.Context, _ string, patch domain.CompanyProfilePatch) (domain.CompanyProfile, error) {
	r.company = patch
	return domain.CompanyProfile{LegalName: "Example"}, r.err
}
func (r *profileRepository) UpdateInterviewerProfile(_ context.Context, _ string, name string) (domain.InterviewerProfile, error) {
	r.interviewerName = name
	return domain.InterviewerProfile{Name: name}, r.err
}

func TestNewProfileServiceRejectsNilRepository(t *testing.T) {
	if _, err := NewProfileService(nil); err == nil {
		t.Fatal("expected nil repository error")
	}
}

func TestProfileServiceUpdatesAllProfileTypes(t *testing.T) {
	repository := &profileRepository{}
	service, err := NewProfileService(repository)
	if err != nil {
		t.Fatal(err)
	}
	name := " Mona "
	url := "https://example.com/profile"
	if _, err := service.UpdateCandidate(context.Background(), "candidate-id", domain.CandidateProfilePatch{Name: &name, GithubProfile: &url}); err != nil {
		t.Fatal(err)
	}
	if repository.candidate.Name != &name || repository.candidate.GithubProfile != &url {
		t.Fatalf("candidate patch was not passed through: %+v", repository.candidate)
	}
	display := "Example"
	if _, err := service.UpdateCompany(context.Background(), "company-id", domain.CompanyProfilePatch{DisplayName: &display}); err != nil {
		t.Fatal(err)
	}
	if repository.company.DisplayName != &display {
		t.Fatalf("company patch was not passed through: %+v", repository.company)
	}
	if _, err := service.UpdateInterviewer(context.Background(), "interviewer-id", " Mona "); err != nil {
		t.Fatal(err)
	}
	if repository.interviewerName != "Mona" {
		t.Fatalf("interviewer name = %q", repository.interviewerName)
	}
}

func TestProfileServiceRejectsInvalidUpdates(t *testing.T) {
	repository := &profileRepository{}
	service, _ := NewProfileService(repository)
	long := string(make([]rune, 101))
	badURL := "not a URL"
	badEmail := "not an email"
	negative := int32(-1)
	tests := []struct {
		name string
		call func() error
	}{
		{"candidate id", func() error {
			_, err := service.UpdateCandidate(context.Background(), "", domain.CandidateProfilePatch{})
			return err
		}},
		{"candidate name", func() error {
			_, err := service.UpdateCandidate(context.Background(), "id", domain.CandidateProfilePatch{Name: &long})
			return err
		}},
		{"candidate URL", func() error {
			_, err := service.UpdateCandidate(context.Background(), "id", domain.CandidateProfilePatch{GithubProfile: &badURL})
			return err
		}},
		{"candidate experience", func() error {
			_, err := service.UpdateCandidate(context.Background(), "id", domain.CandidateProfilePatch{YearsOfExperience: &negative})
			return err
		}},
		{"company empty", func() error {
			_, err := service.UpdateCompany(context.Background(), "id", domain.CompanyProfilePatch{})
			return err
		}},
		{"company email", func() error {
			_, err := service.UpdateCompany(context.Background(), "id", domain.CompanyProfilePatch{ContactEmail: &badEmail})
			return err
		}},
		{"interviewer name", func() error { _, err := service.UpdateInterviewer(context.Background(), "id", " "); return err }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !errors.Is(test.call(), domain.ErrInvalidProfile) {
				t.Fatal("expected invalid profile")
			}
		})
	}
}

func TestProfileServicePropagatesRepositoryErrors(t *testing.T) {
	repository := &profileRepository{err: errors.New("storage down")}
	service, _ := NewProfileService(repository)
	name := "Mona"
	if _, err := service.UpdateCandidate(context.Background(), "id", domain.CandidateProfilePatch{Name: &name}); !errors.Is(err, repository.err) {
		t.Fatalf("candidate error = %v", err)
	}
	if _, err := service.UpdateCompany(context.Background(), "id", domain.CompanyProfilePatch{ClearDescription: true}); !errors.Is(err, repository.err) {
		t.Fatalf("company error = %v", err)
	}
	if _, err := service.UpdateInterviewer(context.Background(), "id", name); !errors.Is(err, repository.err) {
		t.Fatalf("interviewer error = %v", err)
	}
}
