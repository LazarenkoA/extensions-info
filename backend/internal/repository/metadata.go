package repository

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"
)

func (p *PG) SetMetadata(ctx context.Context, confID int32, value []byte) error {
	builder := sq.Insert("metadata").
		Columns("struct", "conf_id").
		Values(value, confID).
		Suffix("ON CONFLICT (conf_id) DO UPDATE " +
			"SET struct = EXCLUDED.struct").
		PlaceholderFormat(sq.Dollar)

	var id int32
	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &id, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return errors.Wrap(err, "query error")
	}

	return nil
}

func (p *PG) GetMetadata(ctx context.Context, confID int32) ([]byte, error) {
	var result []byte

	builder := sq.Select("struct").
		From("metadata").
		Where(sq.Eq{"conf_id": confID}).
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &result, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return result, errors.Wrap(err, "query error")
	}

	return result, nil
}
