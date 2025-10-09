package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"
	"time"
	"your-app/internal/models"
)

func (p *PG) GetDataBaseSettings(ctx context.Context) ([]models.DatabaseSettings, error) {
	builder := sq.Select("db.id", "db.connection_string", "db.status", "db.name", "db.last_check").
		Columns(`CASE WHEN j.database_id IS NOT NULL THEN jsonb_build_object(
								   'Schedule', j.cron,
								   'NextCheck', j.next_check ) END AS Cron`).
		From("database_info db").
		LeftJoin("jobs j on j.database_id = db.id").
		GroupBy("db.id", "j.database_id").
		PlaceholderFormat(sq.Dollar)

	var result []models.DatabaseSettings

	query, args, _ := builder.ToSql()
	err := pgxscan.Select(ctx, p.pool, &result, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	return result, err
}

func (p *PG) GetDataBaseByID(ctx context.Context, id int32) (*models.DatabaseSettings, error) {
	builder := sq.Select("db.id", "db.connection_string", "db.status", "db.name", "db.last_check").
		From("database_info db").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	var result models.DatabaseSettings

	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &result, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}

	return &result, err
}

func (p *PG) AddDataBase(ctx context.Context, data models.DatabaseSettings) error {
	builder := sq.Insert("database_info").
		SetMap(map[string]interface{}{
			"connection_string": data.ConnectionString,
			"name":              data.Name,
			"username":          data.Username,
			"password":          data.Password,
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}

func (p *PG) DeleteDataBase(ctx context.Context, id int32) error {
	builder := sq.Delete("database_info").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}

func (p *PG) SetBDState(ctx context.Context, bdID int32, newState string, lastCheck time.Time) error {
	builder := sq.Update("database_info").
		Set("status", newState).
		Where(sq.Eq{"id": bdID}).
		PlaceholderFormat(sq.Dollar)

	if !lastCheck.IsZero() {
		builder = builder.Set("last_check", lastCheck)
	}

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}
