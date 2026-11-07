package database

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/CSCE-3701-Amara/amara/internal/auth/adapters/database/sqlc"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type fakeProfileQueries struct {
	candidate      sqlc.UpdateCandidateProfileParams
	company        sqlc.UpdateCompanyProfileParams
	interviewer    sqlc.UpdateInterviewerProfileParams
	candidateRow   sqlc.UpdateCandidateProfileRow
	companyRow     sqlc.UpdateCompanyProfileRow
	interviewerRow sqlc.UpdateInterviewerProfileRow
	candidateErr   error
	companyErr     error
	interviewerErr error
}

func (f *fakeProfileQueries) UpdateCandidateProfile(_ context.Context, params sqlc.UpdateCandidateProfileParams) (sqlc.UpdateCandidateProfileRow, error) {
	f.candidate = params
	return f.candidateRow, f.candidateErr
}
func (f *fakeProfileQueries) UpdateCompanyProfile(_ context.Context, params sqlc.UpdateCompanyProfileParams) (sqlc.UpdateCompanyProfileRow, error) {
	f.company = params
	return f.companyRow, f.companyErr
}
func (f *fakeProfileQueries) UpdateInterviewerProfile(_ context.Context, params sqlc.UpdateInterviewerProfileParams) (sqlc.UpdateInterviewerProfileRow, error) {
	f.interviewer = params
	return f.interviewerRow, f.interviewerErr
}

func TestNewProfileRepositoryRejectsNilPool(t *testing.T) {
	if _, err := NewProfileRepository(nil); err == nil {
		t.Fatal("expected nil pool error")
	}
	if _, err := NewProfileRepositoryWithQueries(nil); err == nil {
		t.Fatal("expected nil queries error")
	}
}

func TestProfileRepositoryUsesGeneratedQueries(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	id := "00000000-0000-0000-0000-000000000001"
	queries := &fakeProfileQueries{
		candidateRow:   sqlc.UpdateCandidateProfileRow{ID: id, Name: "Mona", Email: "mona@example.com", CreatedAt: pgtype.Timestamp{Time: now, Valid: true}},
		companyRow:     sqlc.UpdateCompanyProfileRow{ID: id, LegalName: "Example", Approved: true},
		interviewerRow: sqlc.UpdateInterviewerProfileRow{ID: id, CompanyID: "company-id", Name: "Mona", Email: "mona@example.com", AccountStatus: "active", CreatedAt: pgtype.Timestamp{Time: now, Valid: true}},
	}
	repository, err := NewProfileRepositoryWithQueries(queries)
	if err != nil {
		t.Fatal(err)
	}
	name := "Mona"
	if _, err := repository.UpdateCandidateProfile(context.Background(), id, domain.CandidateProfilePatch{Name: &name}); err != nil {
		t.Fatal(err)
	}
	if !queries.candidate.SetName || queries.candidate.Name != name {
		t.Fatalf("candidate params = %+v", queries.candidate)
	}

	if _, err := repository.UpdateCompanyProfile(context.Background(), id, domain.CompanyProfilePatch{ClearDisplayName: true}); err != nil {
		t.Fatal(err)
	}
	if !queries.company.SetDisplayName || queries.company.DisplayName.Valid {
		t.Fatalf("company params = %+v", queries.company)
	}
	if _, err := repository.UpdateInterviewerProfile(context.Background(), id, "Mona"); err != nil {
		t.Fatal(err)
	}
	if queries.interviewer.Name != "Mona" {
		t.Fatalf("interviewer params = %+v", queries.interviewer)
	}
}

func TestProfileRepositoryMapsDatabaseErrorsAndInvalidIDs(t *testing.T) {
	queries := &fakeProfileQueries{
		candidateErr:   errors.New("candidate down"),
		companyErr:     errors.New("company down"),
		interviewerErr: errors.New("interviewer down"),
	}
	repository, _ := NewProfileRepositoryWithQueries(queries)
	if _, err := repository.UpdateCandidateProfile(context.Background(), "bad", domain.CandidateProfilePatch{}); !errors.Is(err, domain.ErrProfileNotFound) {
		t.Fatalf("invalid candidate id error = %v", err)
	}
	if _, err := repository.UpdateCompanyProfile(context.Background(), "bad", domain.CompanyProfilePatch{}); !errors.Is(err, domain.ErrProfileNotFound) {
		t.Fatalf("invalid company id error = %v", err)
	}
	if _, err := repository.UpdateInterviewerProfile(context.Background(), "bad", "Mona"); !errors.Is(err, domain.ErrProfileNotFound) {
		t.Fatalf("invalid interviewer id error = %v", err)
	}
	if _, err := repository.UpdateCandidateProfile(context.Background(), "00000000-0000-0000-0000-000000000001", domain.CandidateProfilePatch{}); err == nil {
		t.Fatal("expected candidate database error")
	}
	if _, err := repository.UpdateCompanyProfile(context.Background(), "00000000-0000-0000-0000-000000000001", domain.CompanyProfilePatch{ClearDescription: true}); err == nil {
		t.Fatal("expected company database error")
	}
	if _, err := repository.UpdateInterviewerProfile(context.Background(), "00000000-0000-0000-0000-000000000001", "Mona"); err == nil {
		t.Fatal("expected interviewer database error")
	}
}

func TestProfileMappingHelpers(t *testing.T) {
	now := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	row := sqlc.UpdateCandidateProfileRow{
		ID: "candidate-id", Name: "Mona", Email: "mona@example.com",
		CreatedAt:         pgtype.Timestamp{Time: now, Valid: true},
		UpdatedAt:         pgtype.Timestamp{Time: now, Valid: true},
		PhoneNumber:       pgtype.Text{String: "123", Valid: true},
		GithubProfile:     pgtype.Text{String: "https://github.com/mona", Valid: true},
		LinkedinProfile:   pgtype.Text{String: "https://linkedin.com/in/mona", Valid: true},
		PortfolioLink:     pgtype.Text{String: "https://mona.dev", Valid: true},
		Gender:            pgtype.Text{String: "x", Valid: true},
		BirthDate:         pgtype.Date{Time: now, Valid: true},
		Nationality:       pgtype.Text{String: "EG", Valid: true},
		Address:           pgtype.Text{String: "Cairo", Valid: true},
		Biography:         pgtype.Text{String: "Bio", Valid: true},
		YearsOfExperience: pgtype.Int4{Int32: 5, Valid: true},
		Available:         pgtype.Bool{Bool: true, Valid: true},
	}
	profile := mapCandidateProfile(row)
	if profile.ID != row.ID || profile.UpdatedAt == nil || profile.BirthDate == nil ||
		profile.YearsOfExperience == nil || profile.Available == nil || *profile.YearsOfExperience != 5 || !*profile.Available {
		t.Fatalf("candidate mapping = %+v", profile)
	}

	companyRow := sqlc.UpdateCompanyProfileRow{
		ID: "company-id", LegalName: "Example", DisplayName: pgtype.Text{String: "Example", Valid: true}, Approved: true,
	}
	company := mapCompanyProfile(companyRow)
	if company.ID != companyRow.ID || company.DisplayName == nil || !company.Approved {
		t.Fatalf("company mapping = %+v", company)
	}
}

func TestProfileConversionHelpers(t *testing.T) {
	value := " Mona "
	if candidateString(&value) != "Mona" || candidateString(nil) != "" {
		t.Fatal("candidateString conversion failed")
	}
	if nullableString(&value, false).String != value || !nullableString(&value, false).Valid {
		t.Fatal("nullable string conversion failed")
	}
	if nullableString(nil, true).Valid {
		t.Fatal("cleared nullable string should be invalid")
	}
	date := time.Now()
	if !candidateDate(&date).Valid || candidateDate(nil).Valid {
		t.Fatal("date conversion failed")
	}
	experience := 4
	available := true
	if candidateInt(&experience) != 4 || candidateInt(nil) != 0 || !candidateBool(&available) || candidateBool(nil) {
		t.Fatal("scalar conversion failed")
	}
	if timestampPtr(pgtype.Timestamp{}) != nil || datePtr(pgtype.Date{}) != nil ||
		intPtr(pgtype.Int4{}) != nil || boolPtr(pgtype.Bool{}) != nil || textPtr(pgtype.Text{}) != nil {
		t.Fatal("invalid nullable values should map to nil")
	}
}
