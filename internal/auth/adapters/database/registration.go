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
	pool    *pgxpool.Pool
}

func (r *RegistrationRepository) FindCredentials(ctx context.Context, email string) (ports.Credentials, error) {
	row, err := r.queries.FindCredentials(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.Credentials{}, domain.ErrInvalidCredentials
		}
		return ports.Credentials{}, fmt.Errorf("query login credentials: %w", err)
	}
	if !row.Deactivated.Valid {
		return ports.Credentials{}, fmt.Errorf("query login credentials: invalid deactivated value")
	}
	deactivated := row.Deactivated.Bool
	companyID, ok := row.CompanyID.(string)
	if !ok {
		return ports.Credentials{}, fmt.Errorf("query login credentials: invalid company id")
	}
	return ports.Credentials{
		ID: row.ID, Email: row.Email, AccountType: domain.AccountType(row.AccountType),
		PasswordHash: row.PasswordHash, Salt: row.Salt, Status: row.AccountStatus,
		Deactivated: deactivated, Banned: row.Banned, CompanyID: companyID,
		CompanyApproved: row.CompanyApproved,
	}, nil
}

func (r *RegistrationRepository) CreateSession(ctx context.Context, session ports.Session) error {
	if r == nil || r.pool == nil {
		return errors.New("session creation database is not configured")
	}
	accountID, err := parseUUID(session.AccountID)
	if err != nil {
		return fmt.Errorf("parse session account id: %w", err)
	}
	_, sessionType := accountActionTable(session.AccountType)
	if sessionType == "" {
		return domain.ErrInvalidSession
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin session creation: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := r.queries.WithTx(tx)
	var lockErr error
	switch session.AccountType {
	case domain.AccountTypeCandidate:
		_, lockErr = queries.LockCandidateForSession(ctx, accountID)
	case domain.AccountTypeInterviewer:
		_, lockErr = queries.LockInterviewerForSession(ctx, accountID)
	case domain.AccountTypeCompany:
		_, lockErr = queries.LockCompanyForSession(ctx, accountID)
	}
	if lockErr != nil {
		if errors.Is(lockErr, pgx.ErrNoRows) {
			return domain.ErrAccountInactive
		}
		return fmt.Errorf("lock session account: %w", lockErr)
	}
	if err := queries.CreateSession(ctx, sqlc.CreateSessionParams{
		Token: session.Token, Column2: accountID, AccountType: sessionType,
		ExpiresAt: pgtype.Timestamp{Time: session.ExpiresAt, Valid: true},
	}); err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit session creation: %w", err)
	}
	return nil
}

func (r *RegistrationRepository) TerminateSession(ctx context.Context, token string) error {
	if strings.TrimSpace(token) == "" {
		return domain.ErrInvalidSession
	}
	if r.queries == nil {
		return errors.New("session termination database is not configured")
	}
	return r.queries.TerminateSession(ctx, token)
}

func (r *RegistrationRepository) PerformAccountAction(ctx context.Context, accountID, sessionToken string, accountType domain.AccountType, action domain.AccountAction) error {
	if r == nil || r.pool == nil {
		return errors.New("account action database is not configured")
	}
	accountID = strings.TrimSpace(accountID)
	sessionToken = strings.TrimSpace(sessionToken)
	if accountID == "" || sessionToken == "" {
		return domain.ErrInvalidSession
	}
	if err := action.Validate(); err != nil {
		return err
	}
	table, sessionType := accountActionTable(accountType)
	if table == "" {
		return domain.ErrInvalidSession
	}
	accountUUID, err := parseUUID(accountID)
	if err != nil {
		return domain.ErrInvalidSession
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin account action: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	queries := r.queries.WithTx(tx)
	var locked int
	if err := tx.QueryRow(ctx, "SELECT 1 FROM "+table+" WHERE id = $1::uuid FOR UPDATE", accountID).Scan(&locked); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInvalidSession
		}
		return fmt.Errorf("lock account: %w", err)
	}
	if _, err := queries.LockActiveSessionForAccount(ctx, sqlc.LockActiveSessionForAccountParams{
		Token: sessionToken, Column2: accountUUID, AccountType: sessionType,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInvalidSession
		}
		return fmt.Errorf("lock account session: %w", err)
	}

	column := "deactivated_at"
	if action == domain.AccountActionDelete {
		column = "deleted_at"
	}
	result, err := tx.Exec(ctx, "UPDATE "+table+" SET "+column+" = now(), updated_at = now() WHERE id = $1::uuid AND deleted_at IS NULL AND deactivated_at IS NULL", accountID)
	if err != nil {
		return fmt.Errorf("update %s account: %w", string(accountType), err)
	}
	if result.RowsAffected() == 0 {
		return domain.ErrInvalidSession
	}
	if _, err := tx.Exec(ctx, `DELETE FROM auth_sessions WHERE account_id = $1::uuid AND account_type = $2`, accountID, sessionType); err != nil {
		return fmt.Errorf("terminate account session: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit account action: %w", err)
	}
	return nil
}

func accountActionTable(accountType domain.AccountType) (string, string) {
	switch accountType {
	case domain.AccountTypeCandidate:
		return "candidates", string(domain.AccountTypeCandidate)
	case domain.AccountTypeInterviewer:
		return "interviewers", string(domain.AccountTypeInterviewer)
	case domain.AccountTypeCompany:
		return "companies", string(domain.AccountTypeCompany)
	default:
		return "", ""
	}
}

func NewRegistrationRepository(pool *pgxpool.Pool) (*RegistrationRepository, error) {
	if pool == nil {
		return nil, errors.New("database pool must not be nil")
	}
	return &RegistrationRepository{queries: sqlc.New(pool), pool: pool}, nil
}

func NewRegistrationRepositoryWithQueries(queries *sqlc.Queries) (*RegistrationRepository, error) {
	if queries == nil {
		return nil, errors.New("database queries must not be nil")
	}
	return &RegistrationRepository{queries: queries}, nil
}

func (r *RegistrationRepository) ValidateSession(ctx context.Context, token string) (ports.Identity, error) {
	if r.pool == nil {
		return ports.Identity{}, errors.New("session validation database is not configured")
	}
	row, err := r.queries.ValidateSession(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.Identity{}, domain.ErrInvalidSession
		}
		return ports.Identity{}, fmt.Errorf("validate session: %w", err)
	}
	return ports.Identity{
		AccountID: row.SAccountID, AccountType: domain.AccountType(row.AccountType),
		CompanyID: row.CompanyID, ExpiresAt: row.ExpiresAt.Time,
	}, nil
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
		Recipient: challenge.Recipient,
	}); err != nil {
		return fmt.Errorf("insert verification challenge: %w", err)
	}

	return nil
}

func (r *RegistrationRepository) FindPendingVerification(ctx context.Context, email string) (ports.PendingVerification, error) {
	row, err := r.queries.FindPendingVerification(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.PendingVerification{}, domain.ErrPendingVerificationNotFound
		}
		return ports.PendingVerification{}, fmt.Errorf("query pending verification: %w", err)
	}
	accountType := domain.AccountType(row.AccountType)
	return ports.PendingVerification{
		Account: domain.Account{
			ID: row.ID, Type: accountType, Email: row.Email,
			Status: row.AccountStatus, CreatedAt: row.CreatedAt.Time.UTC().Format(time.RFC3339),
		},
		Recipient: row.Recipient,
		ExpiresAt: row.ExpiresAt.Time.UTC(),
	}, nil
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

func (r *RegistrationRepository) GetMFASecret(ctx context.Context, accountID string, accountType domain.AccountType) (ports.MFARecord, error) {
	accountUUID, err := parseUUID(accountID)
	if err != nil {
		return ports.MFARecord{}, err
	}
	row, err := r.queries.GetMFASecret(ctx, sqlc.GetMFASecretParams{
		Column1: accountUUID, AccountType: string(accountType),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ports.MFARecord{}, ports.ErrMFANotConfigured
		}
		return ports.MFARecord{}, fmt.Errorf("query MFA credentials: %w", err)
	}
	return ports.MFARecord{
		AccountID: row.AccountID, AccountType: domain.AccountType(row.AccountType),
		Secret: row.Secret, Enabled: row.Enabled, CreatedAt: row.CreatedAt.Time,
	}, nil
}

func (r *RegistrationRepository) SaveMFASecret(ctx context.Context, record ports.MFARecord) error {
	accountUUID, err := parseUUID(record.AccountID)
	if err != nil {
		return err
	}
	return r.queries.SaveMFASecret(ctx, sqlc.SaveMFASecretParams{
		Column1: accountUUID, AccountType: string(record.AccountType), Secret: record.Secret,
		CreatedAt: pgtype.Timestamp{Time: record.CreatedAt, Valid: true},
	})
}

func (r *RegistrationRepository) EnableMFA(ctx context.Context, accountID string, accountType domain.AccountType) error {
	accountUUID, err := parseUUID(accountID)
	if err != nil {
		return err
	}
	return r.queries.EnableMFA(ctx, sqlc.EnableMFAParams{
		Column1: accountUUID, AccountType: string(accountType),
		EnabledAt: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
	})
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
