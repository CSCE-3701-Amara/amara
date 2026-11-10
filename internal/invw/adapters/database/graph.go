package database

import (
	"context"
	"github.com/CSCE-3701-Amara/amara/internal/invw/domain"
	"github.com/jackc/pgx/v5"
)

func (r *PipelineRepository) Graph(ctx context.Context, company, version string) (domain.Graph, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadOnly})
	if err != nil {
		return domain.Graph{}, err
	}
	defer tx.Rollback(ctx)
	if _, err := pipelineVersion(ctx, tx, company, version, false); err != nil {
		return domain.Graph{}, err
	}
	g := domain.Graph{Stages: []domain.Stage{}, Links: []domain.Link{}}
	rows, err := tx.Query(ctx, `SELECT s.id::text,s.name,s.module_type,s.scheduled_start_at,s.scheduled_end_at,s.interview_configuration,s.advancement_rule,s.score_threshold,s.below_threshold_action,s.position_x::float8,s.position_y::float8,ARRAY(SELECT test_module_id::text FROM test_configurations WHERE stage_id=s.id ORDER BY (configuration->>'position')::int) FROM pipeline_stages s WHERE version_id=$1::uuid ORDER BY stage_order,id`, version)
	if err != nil {
		return g, err
	}
	for rows.Next() {
		var s domain.Stage
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.OpensAt, &s.Deadline, &s.Interview, &s.Rule, &s.ScoreThreshold, &s.Below, &s.X, &s.Y, &s.ModuleIDs); err != nil {
			rows.Close()
			return g, err
		}
		g.Stages = append(g.Stages, s)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return g, err
	}
	rows, err = tx.Query(ctx, `SELECT from_stage_id::text,to_stage_id::text FROM stage_transitions WHERE pipeline_version_id=$1::uuid ORDER BY id`, version)
	if err != nil {
		return g, err
	}
	for rows.Next() {
		var e domain.Link
		if err := rows.Scan(&e.From, &e.To); err != nil {
			rows.Close()
			return g, err
		}
		g.Links = append(g.Links, e)
	}
	err = rows.Err()
	rows.Close()
	if err != nil {
		return g, err
	}
	if len(g.Links) == 0 {
		// Linear drafts predate stored positions and all default to (0,0).
		// Project a readable layout until the user saves explicit graph positions.
		unpositioned := true
		for _, stage := range g.Stages {
			unpositioned = unpositioned && stage.X == 0 && stage.Y == 0
		}
		if unpositioned {
			for i := range g.Stages {
				g.Stages[i].X = float64(40 + (i%3)*260)
				g.Stages[i].Y = float64(40 + (i/3)*200)
			}
		}
		for i := 1; i < len(g.Stages); i++ {
			g.Links = append(g.Links, domain.Link{From: g.Stages[i-1].ID, To: g.Stages[i].ID})
		}
	}
	return g, tx.Commit(ctx)
}
func (r *PipelineRepository) SaveGraph(ctx context.Context, company, version string, g domain.Graph) (domain.Graph, error) {
	if err := g.Validate(); err != nil {
		return g, err
	}
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return g, err
	}
	defer tx.Rollback(ctx)
	published, err := pipelineVersion(ctx, tx, company, version, true)
	if err != nil {
		return g, err
	}
	if published {
		return g, domain.ErrPipelineLocked
	}
	defs := map[string]domain.Definition{}
	for _, s := range g.Stages {
		if s.Interview != nil {
			if err := validateInterviewers(ctx, tx, company, s.Interview.InterviewerIDs); err != nil {
				return g, err
			}
		}
		for _, id := range s.ModuleIDs {
			d, err := moduleSnapshot(ctx, tx, company, id)
			if err != nil {
				return g, err
			}
			defs[id] = d
		}
	}
	for _, q := range []string{`DELETE FROM stage_transitions WHERE pipeline_version_id=$1::uuid`, `DELETE FROM test_configurations WHERE stage_id IN (SELECT id FROM pipeline_stages WHERE version_id=$1::uuid)`, `DELETE FROM pipeline_stages WHERE version_id=$1::uuid`} {
		if _, err := tx.Exec(ctx, q, version); err != nil {
			return g, err
		}
	}
	ids := map[string]string{}
	for i := range g.Stages {
		key := g.Stages[i].ID
		if err := insertStage(ctx, tx, version, i+1, &g.Stages[i], defs); err != nil {
			return g, err
		}
		ids[key] = g.Stages[i].ID
		if _, err := tx.Exec(ctx, `UPDATE pipeline_stages SET position_x=$2,position_y=$3 WHERE id=$1::uuid`, g.Stages[i].ID, g.Stages[i].X, g.Stages[i].Y); err != nil {
			return g, err
		}
	}
	for i := range g.Links {
		e := &g.Links[i]
		e.From, e.To = ids[e.From], ids[e.To]
		if _, err := tx.Exec(ctx, `INSERT INTO stage_transitions(id,pipeline_version_id,from_stage_id,to_stage_id) VALUES(gen_random_uuid(),$1::uuid,$2::uuid,$3::uuid)`, version, e.From, e.To); err != nil {
			return g, err
		}
	}
	return g, tx.Commit(ctx)
}
