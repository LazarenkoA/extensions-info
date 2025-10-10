package repository

import (
	"context"
	onec "github.com/LazarenkoA/extensions-info/internal/1c"
	"github.com/LazarenkoA/extensions-info/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/pkg/errors"
)

func (p *PG) GetConfigurationInfo(ctx context.Context, dbID int32) (*models.ConfigurationInfo, error) {
	builder := sq.Select("c.id", "c.description", "c.version", "c.name").
		Columns(`case when ext.conf_id is not null then
							   jsonb_agg(
								jsonb_build_object(
										'ID', ext.id,
										'ConfID', ext.conf_id,
										'Description', ext.description,
										'Version', ext.version,
										'Name', ext.name,
										'Purpose', ext.purpose
								)
							   ) end AS extensions`).
		From("conf_info c").
		LeftJoin("extensions_info ext on c.id = ext.conf_id").
		Where(sq.Eq{"database_id": dbID}).
		GroupBy("c.id", "ext.conf_id").
		PlaceholderFormat(sq.Dollar)

	var result models.ConfigurationInfo

	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &result, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return nil, errors.Wrap(err, "query error")
	}

	return &result, nil
}

func (p *PG) StoreConfigurationInfo(ctx context.Context, dbID int32, confInfo *onec.ConfigurationInfo) (int32, error) {
	builder := sq.Insert("conf_info").
		SetMap(map[string]interface{}{
			"database_id": dbID,
			"name":        confInfo.Name,
			"description": confInfo.Synonym,
			"version":     confInfo.Version,
		}).
		Suffix("ON CONFLICT (database_id) DO UPDATE " +
			"SET name = EXCLUDED.name," +
			"description = EXCLUDED.description," +
			"version = EXCLUDED.version").
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	var id int32
	query, args, _ := builder.ToSql()
	err := pgxscan.Get(ctx, p.pool, &id, query, args...)
	if err != nil && !pgxscan.NotFound(err) {
		return 0, errors.Wrap(err, "query error")
	}

	return id, nil
}

func (p *PG) StoreExtensionsInfo(ctx context.Context, confID int32, confInfo []onec.ConfigurationInfo) error {
	names := make([]string, len(confInfo))
	descriptions := make([]string, len(confInfo))
	versions := make([]string, len(confInfo))
	purposes := make([]string, len(confInfo))

	for i, m := range confInfo {
		names[i], descriptions[i] = m.Name, m.Synonym
		versions[i], purposes[i] = m.Version, m.Purpose
	}

	query := `INSERT INTO extensions_info (conf_id, name, description, version, purpose)
				SELECT $1, name, description, version, purpose
					FROM unnest($2::text[], $3::text[], $4::text[], $5::text[]) AS input(name, description, version, purpose)
			  ON CONFLICT (conf_id, name) DO UPDATE 
			      SET name = EXCLUDED.name, 
			      description = EXCLUDED.description,
			      version = EXCLUDED.version,
			      purpose = EXCLUDED.purpose`
	_, err := p.pool.Exec(ctx, query, confID, names, descriptions, versions, purposes)
	if err != nil {
		return errors.Wrap(err, "query error")
	}

	return nil
}

func (p *PG) GetExtensionsInfo(ctx context.Context, confID int32) ([]onec.ConfigurationInfo, error) {
	var result []onec.ConfigurationInfo

	builder := sq.Select("id", "name", "description as synonym", "version", "purpose").
		From("extensions_info").
		Where(sq.Eq{"conf_id": confID}).
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()
	err := pgxscan.Select(ctx, p.pool, &result, query, args...)
	if err != nil {
		return result, errors.Wrap(err, "query error")
	}

	return result, err
}
