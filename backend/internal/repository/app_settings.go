package repository

import (
	"context"
	"github.com/LazarenkoA/extensions-info/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"
)

func (p *PG) GetAppSettings(ctx context.Context) (*models.AppSettings, error) {
	builder := sq.Select("platform_path").From("app_settings").Limit(1)

	var result models.AppSettings

	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &result, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return nil, errors.Wrap(err, "query error")
	}

	return &result, nil
}

func (p *PG) SetAppSettings(ctx context.Context, id int32, settings models.AppSettings) error {
	builder := sq.Insert("app_settings").
		Columns("platform_path", "id").
		Values(settings.PlatformPath, id).
		Suffix("ON CONFLICT (id) DO UPDATE SET platform_path = EXCLUDED.platform_path").
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return err
}
