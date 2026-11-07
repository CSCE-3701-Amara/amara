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
	GetCompanyProfile(context.Context, sqlc.GetCompanyProfileParams) (sqlc.GetCompanyProfileRow, error)
	GetPublicCompanyProfile(context.Context, pgtype.UUID) (sqlc.GetPublicCompanyProfileRow, error)
	GetInterviewerProfile(context.Context, pgtype.UUID) (sqlc.GetInterviewerProfileRow, error)
	ListCompanyInterviewers(context.Context, sqlc.ListCompanyInterviewersParams) ([]sqlc.ListCompanyInterviewersRow, error)
	CountCompanyInterviewers(context.Context, pgtype.UUID) (int64, error)
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

func (r *ProfileRepository) GetCompanyProfile(ctx context.Context, id string, includePrivate bool) (domain.CompanyProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.GetCompanyProfile(ctx, sqlc.GetCompanyProfileParams{ID: accountID, IncludePrivate: includePrivate})
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.CompanyProfile{}, fmt.Errorf("get company profile: %w", err)
	}
	return mapCompanyProfileRow(row), nil
}

func (r *ProfileRepository) GetPublicCompanyProfile(ctx context.Context, id string) (domain.CompanyProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.GetPublicCompanyProfile(ctx, accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.CompanyProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.CompanyProfile{}, fmt.Errorf("get public company profile: %w", err)
	}
	return mapPublicCompanyProfileRow(row), nil
}

func (r *ProfileRepository) GetInterviewerProfile(ctx context.Context, id string) (domain.InterviewerProfile, error) {
	accountID, err := parseUUID(id)
	if err != nil {
		return domain.InterviewerProfile{}, domain.ErrProfileNotFound
	}
	row, err := r.queries.GetInterviewerProfile(ctx, accountID)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.InterviewerProfile{}, domain.ErrProfileNotFound
	}
	if err != nil {
		return domain.InterviewerProfile{}, fmt.Errorf("get interviewer profile: %w", err)
	}
	return mapInterviewerProfileValues(row.IID, row.ICompanyID, row.Name, row.Email, row.AccountStatus, row.CreatedAt, row.UpdatedAt), nil
}

func (r *ProfileRepository) ListCompanyInterviewers(ctx context.Context, companyID string, limit, offset int32) ([]domain.InterviewerProfile, int64, error) {
	id, err := parseUUID(companyID)
	if err != nil {
		return nil, 0, domain.ErrProfileNotFound
	}
	rows, err := r.queries.ListCompanyInterviewers(ctx, sqlc.ListCompanyInterviewersParams{
		CompanyID: id, PageLimit: limit, PageOffset: offset,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("list company interviewers: %w", err)
	}
	total, err := r.queries.CountCompanyInterviewers(ctx, id)
	if err != nil {
		return nil, 0, fmt.Errorf("count company interviewers: %w", err)
	}
	profiles := make([]domain.InterviewerProfile, 0, len(rows))
	for _, row := range rows {
		profiles = append(profiles, mapInterviewerProfileValues(row.IID, row.ICompanyID, row.Name, row.Email, row.AccountStatus, row.CreatedAt, row.UpdatedAt))
	}
	return profiles, total, nil
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
	return mapCompanyProfileValues(row.ID, row.LegalName, row.DisplayName, row.Description, row.WebsiteLink, row.LogoUrl, row.ContactPhone, row.ContactEmail, row.LinkedinProfile, row.Approved)
}

func mapCompanyProfileRow(row sqlc.GetCompanyProfileRow) domain.CompanyProfile {
	return mapCompanyProfileValues(row.ID, row.LegalName, row.DisplayName, row.Description, row.WebsiteLink, row.LogoUrl, row.ContactPhone, row.ContactEmail, row.LinkedinProfile, row.Approved)
}

func mapPublicCompanyProfileRow(row sqlc.GetPublicCompanyProfileRow) domain.CompanyProfile {
	return mapCompanyProfileValues(row.ID, row.LegalName, row.DisplayName, row.Description, row.WebsiteLink, row.LogoUrl, row.ContactPhone, row.ContactEmail, row.LinkedinProfile, row.Approved)
}

func mapCompanyProfileValues(id, legalName string, displayName, description, websiteLink, logoURL, contactPhone, contactEmail, linkedinProfile pgtype.Text, approved bool) domain.CompanyProfile {
	return domain.CompanyProfile{
		ID: id, LegalName: legalName, DisplayName: textPtr(displayName),
		Description: textPtr(description), WebsiteLink: textPtr(websiteLink),
		LogoURL: textPtr(logoURL), ContactPhone: textPtr(contactPhone),
		ContactEmail: textPtr(contactEmail), LinkedinProfile: textPtr(linkedinProfile),
		Approved: approved,
	}
}

func mapInterviewerProfileValues(id, companyID, name, email, status string, createdAt, updatedAt pgtype.Timestamp) domain.InterviewerProfile {
	return domain.InterviewerProfile{
		ID: id, CompanyID: companyID, Name: name, Email: email,
		AccountStatus: status, CreatedAt: createdAt.Time, UpdatedAt: timestampPtr(updatedAt),
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
