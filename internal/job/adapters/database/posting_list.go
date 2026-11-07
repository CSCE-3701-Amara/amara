package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/job/domain"
	"github.com/CSCE-3701-Amara/amara/internal/job/ports"
)

const createdAtSortKey = "COALESCE(p.created_at, '0001-01-01 00:00:00'::timestamp)"

var likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)

// CompanyPostingListRepository reads a company's postings in every state.
type CompanyPostingListRepository struct {
	pool *pgxpool.Pool
}

// NewCompanyPostingListRepository creates the company posting list repository.
func NewCompanyPostingListRepository(pool *pgxpool.Pool) *CompanyPostingListRepository {
	return &CompanyPostingListRepository{pool: pool}
}

func (r *CompanyPostingListRepository) ListByCompany(ctx context.Context, filter ports.CompanyPostingQuery) (ports.CompanyPostingResult, error) {
	if r == nil || r.pool == nil {
		return ports.CompanyPostingResult{}, fmt.Errorf("job posting database is not configured")
	}
	args := []any{filter.CompanyID}
	where := []string{"p.company_id = $1", "p.deleted_at IS NULL"}
	add := func(expression string, value any) {
		args = append(args, value)
		where = append(where, fmt.Sprintf(expression, len(args)))
	}
	if filter.Title != "" {
		add(`p.title ILIKE $%d ESCAPE '\'`, "%"+escapeLikePattern(filter.Title)+"%")
	}
	if filter.Status != "" {
		add("COALESCE(p.status, 'draft') = $%d", filter.Status)
	}
	if filter.AfterCreatedAt != nil {
		args = append(args, *filter.AfterCreatedAt, filter.AfterPostingID)
		where = append(where, fmt.Sprintf("("+createdAtSortKey+", p.id) < ($%d, $%d)", len(args)-1, len(args)))
	}
	args = append(args, filter.Limit+1)
	// A posting without a status is treated as Draft, the state every posting
	query := `SELECT p.id, p.title, COALESCE(p.status, 'draft'), p.application_deadline, p.published_at, p.archived_at, p.created_at, p.updated_at, (SELECT COUNT(*) FROM applications a WHERE a.posting_id = p.id) FROM job_postings p WHERE ` +
		strings.Join(where, " AND ") + fmt.Sprintf(" ORDER BY %s DESC, p.id DESC LIMIT $%d", createdAtSortKey, len(args))
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return ports.CompanyPostingResult{}, fmt.Errorf("query company postings: %w", err)
	}
	defer rows.Close()

	postings := make([]domain.CompanyPosting, 0, filter.Limit+1)
	for rows.Next() {
		var posting domain.CompanyPosting
		if err := rows.Scan(
			&posting.ID, &posting.Title, &posting.Status, &posting.ApplicationDeadline,
			&posting.PublishedAt, &posting.ArchivedAt, &posting.CreatedAt, &posting.UpdatedAt,
			&posting.ApplicationCount,
		); err != nil {
			return ports.CompanyPostingResult{}, fmt.Errorf("scan company posting row: %w", err)
		}
		postings = append(postings, posting)
	}
	if err := rows.Err(); err != nil {
		return ports.CompanyPostingResult{}, fmt.Errorf("read company posting rows: %w", err)
	}
	return trimPage(postings, filter.Limit), nil
}

// trimPage drops the lookahead row fetched to detect another page.
func trimPage(postings []domain.CompanyPosting, limit int) ports.CompanyPostingResult {
	if len(postings) > limit {
		return ports.CompanyPostingResult{Postings: postings[:limit], HasMore: true}
	}
	return ports.CompanyPostingResult{Postings: postings}
}

// escapeLikePattern makes user input match literally inside an ILIKE pattern.
func escapeLikePattern(value string) string {
	return likeEscaper.Replace(value)
}
