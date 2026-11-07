package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
)

type ModuleRepository struct {
	pool *pgxpool.Pool
}

func NewModuleRepository(pool *pgxpool.Pool) *ModuleRepository {
	return &ModuleRepository{pool: pool}
}

const moduleColumns = "id::text, company_id::text, copied_from_module_id::text, published, authoring_definition"

func scanModule(row pgx.Row) (domain.Module, error) {
	var module domain.Module
	err := row.Scan(&module.ID, &module.CompanyID, &module.CopiedFrom, &module.Published, &module.Definition)
	if errors.Is(err, pgx.ErrNoRows) {
		err = domain.ErrModuleNotFound
	}
	return module, err
}

func (r *ModuleRepository) Create(ctx context.Context, company string, d domain.Definition) (domain.Module, error) {
	return scanModule(r.pool.QueryRow(ctx, `
		INSERT INTO modules (id, company_id, name, description, type, active_time_limit, published, authoring_definition)
		VALUES (gen_random_uuid(), $1::text::uuid, $2, $3, $4, make_interval(mins => $5), false, $6)
		RETURNING `+moduleColumns, company, d.Name, d.Description, d.Type, d.ActiveMinutes, d))
}

func (r *ModuleRepository) Get(ctx context.Context, company, id string) (domain.Module, error) {
	return scanModule(r.pool.QueryRow(ctx, `
		SELECT `+moduleColumns+` FROM modules
		WHERE id = $1::text::uuid AND company_id = $2::text::uuid AND authoring_definition IS NOT NULL`, id, company))
}

func (r *ModuleRepository) Update(ctx context.Context, company, id string, d domain.Definition) (domain.Module, error) {
	module, err := scanModule(r.pool.QueryRow(ctx, `
		UPDATE modules SET name = $3, description = $4, type = $5,
			active_time_limit = make_interval(mins => $6), authoring_definition = $7, updated_at = now()
		WHERE id = $1::text::uuid AND company_id = $2::text::uuid
			AND authoring_definition IS NOT NULL AND published = false AND copied_from_module_id IS NULL
		RETURNING `+moduleColumns, id, company, d.Name, d.Description, d.Type, d.ActiveMinutes, d))
	if errors.Is(err, domain.ErrModuleNotFound) {
		if _, getErr := r.Get(ctx, company, id); getErr != nil {
			return domain.Module{}, getErr
		}
		return domain.Module{}, domain.ErrModuleLocked
	}
	return module, err
}

type ModuleSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Published   bool   `json:"published"`
}

func (r *ModuleRepository) List(ctx context.Context, company string, library bool) ([]ModuleSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, name, COALESCE(description, ''), type, published FROM modules
		WHERE authoring_definition IS NOT NULL AND
			(($2 AND published) OR (NOT $2 AND company_id = $1::text::uuid))
		ORDER BY name, id`, company, library)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]ModuleSummary, 0)
	for rows.Next() {
		var item ModuleSummary
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Type, &item.Published); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ModuleRepository) Copy(ctx context.Context, company, id string) (domain.Module, error) {
	return scanModule(r.pool.QueryRow(ctx, `
		INSERT INTO modules (id, company_id, name, description, type, active_time_limit,
			published, authoring_definition, copied_from_module_id)
		SELECT gen_random_uuid(), $2::text::uuid, name, description, type, active_time_limit,
			false, authoring_definition, id
		FROM modules WHERE id = $1::text::uuid AND published AND authoring_definition IS NOT NULL
		RETURNING `+moduleColumns, id, company))
}

type problemQuerier interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
}

func (r *ModuleRepository) Problems(ctx context.Context) ([]domain.LibraryProblem, error) {
	return listProblems(ctx, r.pool)
}

func listProblems(ctx context.Context, db problemQuerier) ([]domain.LibraryProblem, error) {
	rows, err := db.Query(ctx, `
		SELECT m.id::text, q.position - 1, q.item->>'prompt', q.item->>'difficulty'
		FROM modules m
		CROSS JOIN LATERAL jsonb_array_elements(m.authoring_definition->'questions')
			WITH ORDINALITY AS q(item, position)
		WHERE m.published AND q.item->>'format' = 'code'
		ORDER BY m.id, q.position`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	problems := make([]domain.LibraryProblem, 0)
	for rows.Next() {
		var problem domain.LibraryProblem
		if err := rows.Scan(&problem.ModuleID, &problem.Position, &problem.Prompt, &problem.Difficulty); err != nil {
			return nil, err
		}
		problems = append(problems, problem)
	}
	return problems, rows.Err()
}

func (r *ModuleRepository) Publish(ctx context.Context, company, id string) (domain.Module, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return domain.Module{}, err
	}
	defer tx.Rollback(ctx)
	module, err := scanModule(tx.QueryRow(ctx, `
		SELECT `+moduleColumns+` FROM modules
		WHERE id = $1::text::uuid AND company_id = $2::text::uuid AND authoring_definition IS NOT NULL
		FOR UPDATE`, id, company))
	if err != nil {
		return domain.Module{}, err
	}
	if module.CopiedFrom != nil {
		return domain.Module{}, domain.ErrModuleLocked
	}
	if module.Published {
		return module, nil
	}

	if err := prepareModuleDefinition(ctx, tx, &module.Definition); err != nil {
		return domain.Module{}, err
	}

	_, err = tx.Exec(ctx, `UPDATE modules SET published = true, updated_at = now(),
		authoring_definition = $2 WHERE id = $1::text::uuid`, id, module.Definition)
	if err != nil {
		return domain.Module{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return domain.Module{}, err
	}
	module.Published = true
	return module, nil
}

func prepareModuleDefinition(ctx context.Context, tx pgx.Tx, definition *domain.Definition) error {
	var err error
	var problems []domain.LibraryProblem
	if definition.Library != nil {
		problems, err = listProblems(ctx, tx)
		if err != nil {
			return err
		}
	}
	if err := definition.ValidatePublication(problems); err != nil {
		return &domain.InvalidModuleError{Reason: err.Error()}
	}
	if definition.Library != nil && len(definition.Library.DifficultyMix) > 0 {
		for _, problem := range problems {
			if definition.Library.DifficultyMix[problem.Difficulty] > 0 {
				definition.Library.Pool = append(definition.Library.Pool, problem.ProblemRef)
			}
		}
	}

	return nil
}
