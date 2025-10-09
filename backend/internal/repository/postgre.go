package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"log"
	"time"
)

type PG struct {
	pool *pgxpool.Pool
}

func NewPG(ctx context.Context, connStr string) (*PG, error) {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse config")
	}

	// Настройка пула (по желанию)
	config.MaxConns = 10 // максимум соединений
	config.MinConns = 2  // минимум соединений
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 10 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create pool")
	}

	go func() {
		<-ctx.Done()
		pool.Close()
	}()

	// Проверка соединения
	err = pool.Ping(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to ping database")
	}

	log.Println("successfully connected with connection pool!")
	return &PG{
		pool: pool,
	}, nil
}
