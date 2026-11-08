package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *RegistrationRepository) ReserveCompanyEmailCode(ctx context.Context, id string, mac []byte, expiresAt, resendAt time.Time) (bool, error) {
	result, err := r.pool.Exec(ctx, `INSERT INTO auth_company_email_mfa
 (company_id, code_mac, expires_at, resend_at, attempts_left) VALUES ($1::uuid, $2, $3, $4, 5)
 ON CONFLICT (company_id) DO UPDATE SET code_mac = EXCLUDED.code_mac,
 expires_at = EXCLUDED.expires_at, resend_at = EXCLUDED.resend_at, attempts_left = 5
 WHERE auth_company_email_mfa.resend_at <= now()`, id, mac, expiresAt, resendAt)
	return result.RowsAffected() == 1, err
}

func (r *RegistrationRepository) ConsumeCompanyEmailCode(ctx context.Context, id string, mac []byte) (bool, error) {
	var matched bool
	err := r.pool.QueryRow(ctx, `UPDATE auth_company_email_mfa
 SET attempts_left = CASE WHEN code_mac = $2 THEN 0 ELSE attempts_left - 1 END
 WHERE company_id = $1::uuid AND expires_at > now() AND attempts_left > 0
 RETURNING code_mac = $2`, id, mac).Scan(&matched)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	return matched, err
}

func (r *RegistrationRepository) ReleaseCompanyEmailCode(ctx context.Context, id string, mac []byte) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM auth_company_email_mfa
 WHERE company_id = $1::uuid AND code_mac = $2`, id, mac)
	return err
}
