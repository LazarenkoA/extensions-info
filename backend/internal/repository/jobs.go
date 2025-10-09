package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"
	"your-app/internal/models"
)

func (p *PG) GetCronSettings(ctx context.Context) (*models.CRONInfo, error) {
	builder := sq.Select("database_id", "next_check", "cron as schedule", "status").
		From("jobs")

	var result models.CRONInfo

	query, args, _ := builder.ToSql()
	err := pgxscan.Select(ctx, p.pool, &result, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return nil, errors.Wrap(err, "query error")
	}

	return &result, err
}

func (p *PG) SetSchedule(ctx context.Context, bdID int32, schedule string) error {
	builder := sq.Insert("jobs").
		Columns("cron", "database_id").
		Values(schedule, bdID).
		Suffix("ON CONFLICT (database_id) DO UPDATE SET cron = EXCLUDED.cron").
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}

func (p *PG) DeleteSchedule(ctx context.Context, bdID int32) error {
	builder := sq.Delete("jobs").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"database_id": bdID})

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}

func (p *PG) SetJobState(ctx context.Context, bdID int32, newState string) error {
	builder := sq.Update("jobs").
		Set("status", newState).
		Where(sq.Eq{"database_id": bdID}).
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}
