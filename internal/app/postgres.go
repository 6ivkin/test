package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const postgresPingTimeout = 5 * time.Second

func newPostgres(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, errors.New("db url is empty")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}
