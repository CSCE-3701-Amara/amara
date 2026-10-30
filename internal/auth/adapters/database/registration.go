package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/auth/adapters/database/sqlc"
	"github.com/CSCE-3701-Amara/amara/internal/auth/domain"
	"github.com/CSCE-3701-Amara/amara/internal/auth/ports"
)

type RegistrationRepository struct {
	queries *sqlc.Queries
}

func NewRegistrationRepository(pool *pgxpool.Pool) (*RegistrationRepository, error) {
	if pool == nil {
		return nil, errors.New("database pool must not be nil")
	}
	return &RegistrationRepository{queries: sqlc.New(pool)}, nil
}

func NewRegistrationRepositoryWithQueries(queries *sqlc.Queries) (*RegistrationRepository, error) {
	if queries == nil {
		return nil, errors.New("database queries must not be nil")
	}
	return &RegistrationRepository{queries: queries}, nil
}

func (r *RegistrationRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.queries.EmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("query registration email: %w", err)
	}
	return exists, nil
}

func (r *RegistrationRepository) FindApprovedCompanyBySecret(ctx context.Context, secret string) (string, error) {
	secretHash := sha256.Sum256([]byte(secret))
	companyID, err := r.queries.FindApprovedCompanyBySecret(ctx, secretHash[:])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrCompanyNotFound
		}
		return "", fmt.Errorf("query company association: %w", err)
	}
	return companyID, nil
}

func (r *RegistrationRepository) CreateCandidate(ctx context.Context, name, email string, passwordHash, salt []byte, createdAt time.Time) (domain.Account, error) {
	row, err := r.queries.CreateCandidateAccount(ctx, sqlc.CreateCandidateAccountParams{
		Name: name, Email: email, PasswordHash: passwordHash, Salt: salt, CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
	})
	if err != nil {
		return domain.Account{}, mapCreateError(err)
	}
	return domain.Account{
		ID: row.ID, Type: domain.AccountTypeCandidate, Email: row.Email,
		Status: row.AccountStatus, CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
	}, nil
}

func (r *RegistrationRepository) CreateInterviewer(ctx context.Context, name, email string, passwordHash, salt []byte, companyID string, createdAt time.Time) (domain.Account, error) {
	companyUUID, err := parseUUID(companyID)
	if err != nil {
		return domain.Account{}, fmt.Errorf("parse company id: %w", err)
	}
	row, err := r.queries.CreateInterviewerAccount(ctx, sqlc.CreateInterviewerAccountParams{
		Column1: companyUUID, Name: name, Email: email, PasswordHash: passwordHash, Salt: salt, CreatedAt: pgtype.Timestamp{Time: createdAt, Valid: true},
	})
	if err != nil {
		return domain.Account{}, mapCreateError(err)
	}
	return domain.Account{
		ID: row.ID, Type: domain.AccountTypeInterviewer, Email: row.Email,
		Status: row.AccountStatus, CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
	}, nil
}

func (r *RegistrationRepository) SaveVerificationChallenge(ctx context.Context, challenge ports.VerificationChallenge) error {
	accountUUID, err := parseUUID(challenge.AccountID)
	if err != nil {
		return fmt.Errorf("parse account id: %w", err)
	}
	if err := r.queries.SaveEmailVerificationChallenge(ctx, sqlc.SaveEmailVerificationChallengeParams{
		Column1: accountUUID, AccountType: string(challenge.AccountType),
		CodeHash: challenge.CodeHash, ExpiresAt: pgtype.Timestamp{Time: challenge.ExpiresAt, Valid: true},
	}); err != nil {
		return fmt.Errorf("insert verification challenge: %w", err)
	}
	return nil
}

func (r *RegistrationRepository) RegistrationNumberExists(ctx context.Context, number string) (bool, error) {
	exists, err := r.queries.CompanyRegistrationNumberExists(ctx, number)
	if err != nil {
		return false, fmt.Errorf("query registration number: %w", err)
	}
	return exists, nil
}

func (r *RegistrationRepository) CreateCompanyRegistration(ctx context.Context, record ports.CompanyRegistrationRecord) (domain.Account, error) {
	row, err := r.queries.CreateCompanyRegistration(ctx, sqlc.CreateCompanyRegistrationParams{
		CompanyName: record.CompanyName, RegistrationNumber: record.RegistrationNumber,
		Email: record.Email, PasswordHash: record.PasswordHash, Salt: record.Salt,
		CreatedAt: pgtype.Timestamp{Time: record.CreatedAt, Valid: true},
	})
	if err != nil {
		return domain.Account{}, mapCreateError(err)
	}
	return domain.Account{
		ID: row.ID, Type: "company", Email: row.Email, Status: row.Status,
		CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
	}, nil
}

func parseUUID(value string) (pgtype.UUID, error) {
	compact := strings.ReplaceAll(value, "-", "")
	if len(compact) != 32 {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID %q", value)
	}
	bytes, err := hex.DecodeString(compact)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID %q: %w", value, err)
	}
	var result pgtype.UUID
	copy(result.Bytes[:], bytes)
	result.Valid = true
	return result, nil
}

func mapCreateError(err error) error {
	var constraintError *pgconn.PgError
	if errors.As(err, &constraintError) && constraintError.Code == "23505" {
		return domain.ErrEmailAlreadyExists
	}
	return fmt.Errorf("insert account: %w", err)
}
