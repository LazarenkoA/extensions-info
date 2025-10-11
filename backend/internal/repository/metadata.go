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

func (p *PG) SetCode(ctx context.Context, extID int32, key, code string) error {
	builder := sq.Insert("code").
		Columns("ext_id", "key", "code").
		Values(extID, key, code).
		Suffix("ON CONFLICT (ext_id, key) DO UPDATE " +
			"SET code = EXCLUDED.code").
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	_, err := p.pool.Exec(ctx, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return errors.Wrap(err, "query error")
	}

	return nil
}

func (p *PG) GetCode(ctx context.Context, extID int32, key string) (string, error) {
	var result string

	builder := sq.Select("code").
		From("code").
		Where(sq.Eq{"ext_id": extID}).
		Where(sq.Eq{"key": key}).
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &result, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return "", errors.Wrap(err, "query error")
	}

	return result, nil
}
