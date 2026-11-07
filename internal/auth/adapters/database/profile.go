package database

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/auth/adapters/database/sqlc"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
)

type profileQueries interface {
	UpdateCandidateProfile(context.Context, sqlc.UpdateCandidateProfileParams) (sqlc.UpdateCandidateProfileRow, error)
	UpdateCompanyProfile(context.Context, sqlc.UpdateCompanyProfileParams) (sqlc.UpdateCompanyProfileRow, error)
	UpdateInterviewerProfile(context.Context, sqlc.UpdateInterviewerProfileParams) (sqlc.UpdateInterviewerProfileRow, error)
}

type ProfileRepository struct{ queries profileQueries }

func NewProfileRepository(pool *pgxpool.Pool) (*ProfileRepository, error) {
	if pool == nil {
		return nil, errors.New("database pool must not be nil")
	}
	return &ProfileRepository{queries: sqlc.New(pool)}, nil
}

func NewProfileRepositoryWithQueries(queries profileQueries) (*ProfileRepository, error) {
	if queries == nil {
		return nil, errors.New("database queries must not be nil")
	}
	return &ProfileRepository{queries: queries}, nil
}

func (r *ProfileRepository) UpdateCandidateProfile(ctx context.Context, id string, p domain.CandidateProfilePatch) (domain.CandidateProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.CandidateProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.UpdateCandidateProfile(ctx, sqlc.UpdateCandidateProfileParams{
		SetName: p.Name != nil, Name: candidateString(p.Name),
		SetPhoneNumber: p.PhoneNumber != nil, PhoneNumber: candidateString(p.PhoneNumber),
		SetGithubProfile: p.GithubProfile != nil, GithubProfile: candidateString(p.GithubProfile),
		SetLinkedinProfile: p.LinkedinProfile != nil, LinkedinProfile: candidateString(p.LinkedinProfile),
		SetPortfolioLink: p.PortfolioLink != nil, PortfolioLink: candidateString(p.PortfolioLink),
		SetGender: p.Gender != nil, Gender: candidateString(p.Gender),
		SetBirthDate: p.BirthDate != nil, BirthDate: candidateDate(p.BirthDate),
		SetNationality: p.Nationality != nil, Nationality: candidateString(p.Nationality),
		SetAddress: p.Address != nil, Address: candidateString(p.Address),
		SetBiography: p.Biography != nil, Biography: candidateString(p.Biography),
		SetYearsOfExperience: p.YearsOfExperience != nil, YearsOfExperience: candidateInt(p.YearsOfExperience),
		SetAvailable: p.Available != nil, Available: candidateBool(p.Available),
		ID: accountID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CandidateProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.CandidateProfile{}, fmt.Errorf("update candidate profile: %w", err)
	}
	return mapCandidateProfile(row), nil
}

func (r *ProfileRepository) UpdateCompanyProfile(ctx context.Context, id string, p domain.CompanyProfilePatch) (domain.CompanyProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.UpdateCompanyProfile(ctx, sqlc.UpdateCompanyProfileParams{
		SetLegalName: p.LegalName != nil, LegalName: candidateString(p.LegalName),
		SetDisplayName: p.DisplayName != nil || p.ClearDisplayName, DisplayName: nullableString(p.DisplayName, p.ClearDisplayName),
		SetDescription: p.Description != nil || p.ClearDescription, Description: nullableString(p.Description, p.ClearDescription),
		SetWebsiteLink: p.WebsiteLink != nil || p.ClearWebsiteLink, WebsiteLink: nullableString(p.WebsiteLink, p.ClearWebsiteLink),
		SetLogoUrl: p.LogoURL != nil || p.ClearLogoURL, LogoUrl: nullableString(p.LogoURL, p.ClearLogoURL),
		SetContactPhone: p.ContactPhone != nil || p.ClearContactPhone, ContactPhone: nullableString(p.ContactPhone, p.ClearContactPhone),
		SetContactEmail: p.ContactEmail != nil || p.ClearContactEmail, ContactEmail: nullableString(p.ContactEmail, p.ClearContactEmail),
		SetLinkedinProfile: p.LinkedinProfile != nil || p.ClearLinkedinProfile, LinkedinProfile: nullableString(p.LinkedinProfile, p.ClearLinkedinProfile),
		ID: accountID,
	})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.CompanyProfile{}, fmt.Errorf("update company profile: %w", err)
	}
	return mapCompanyProfile(row), nil
}

func (r *ProfileRepository) UpdateInterviewerProfile(ctx context.Context, id, name string) (domain.InterviewerProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.InterviewerProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.UpdateInterviewerProfile(ctx, sqlc.UpdateInterviewerProfileParams{Name: name, ID: accountID})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.InterviewerProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.InterviewerProfile{}, fmt.Errorf("update interviewer profile: %w", err)
	}
	return domain.InterviewerProfile{
		ID: row.ID, CompanyID: row.CompanyID, Name: row.Name, Email: row.Email,
		AccountStatus: row.AccountStatus, CreatedAt: row.CreatedAt.Time,
		UpdatedAt: timestampPtr(row.UpdatedAt),
	}, nil
}

func mapCandidateProfile(row sqlc.UpdateCandidateProfileRow) domain.CandidateProfile {
	return domain.CandidateProfile{
		ID: row.ID, Name: row.Name, Email: row.Email, CreatedAt: row.CreatedAt.Time,
		UpdatedAt: timestampPtr(row.UpdatedAt), PhoneNumber: row.PhoneNumber.String,
		GithubProfile: row.GithubProfile.String, LinkedinProfile: row.LinkedinProfile.String,
		PortfolioLink: row.PortfolioLink.String, Gender: row.Gender.String,
		BirthDate: datePtr(row.BirthDate), Nationality: row.Nationality.String,
		Address: row.Address.String, Biography: row.Biography.String,
		YearsOfExperience: intPtr(row.YearsOfExperience), Available: boolPtr(row.Available),
	}
}

func mapCompanyProfile(row sqlc.UpdateCompanyProfileRow) domain.CompanyProfile {
	return domain.CompanyProfile{
		ID: row.ID, LegalName: row.LegalName, DisplayName: textPtr(row.DisplayName),
		Description: textPtr(row.Description), WebsiteLink: textPtr(row.WebsiteLink),
		LogoURL: textPtr(row.LogoUrl), ContactPhone: textPtr(row.ContactPhone),
		ContactEmail: textPtr(row.ContactEmail), LinkedinProfile: textPtr(row.LinkedinProfile),
		Approved: row.Approved,
	}
}

func candidateString(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(*value)
}

func nullableString(value *string, clear bool) pgtype.Text {
	if value == nil && clear {
		return pgtype.Text{}
	}
	if value == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *value, Valid: true}
}

func candidateDate(value *time.Time) pgtype.Date {
	if value == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *value, Valid: true}
}

func candidateInt(value *int32) int32 {
	if value == nil {
		return 0
	}
	return int32(*value)
}

func candidateBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}

func timestampPtr(value pgtype.Timestamp) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func datePtr(value pgtype.Date) *time.Time {
	if !value.Valid {
		return nil
	}
	return &value.Time
}

func intPtr(value pgtype.Int4) *int {
	if !value.Valid {
		return nil
	}
	result := int(value.Int32)
	return &result
}

func boolPtr(value pgtype.Bool) *bool {
	if !value.Valid {
		return nil
	}
	result := value.Bool
	return &result
}

func textPtr(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
