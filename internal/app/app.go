package app

import (
	"context"
	"fmt"
	"time"

	"github.com/6ivkin/test.git/internal/config"
	"github.com/6ivkin/test.git/internal/reader"
	"github.com/6ivkin/test.git/internal/repository/postgres"
	httptransport "github.com/6ivkin/test.git/internal/transport/http"
	"github.com/6ivkin/test.git/internal/transport/http/handler"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type App struct {
	server *httptransport.Server
	pool   *pgxpool.Pool
}

const shutdownTimeout = 10 * time.Second

func (a *App) Run(ctx context.Context) error {
	serverErrCh := make(chan error, 1)

	go func() {
		serverErrCh <- a.server.Run()
	}()

	select {
	case <-ctx.Done():

	case err := <-serverErrCh:
		if err != nil {
			return fmt.Errorf(
				"http server: %w",
				err,
			)
		}

		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer cancel()

	if err := a.server.Shutdown(
		shutdownCtx,
	); err != nil {
		return fmt.Errorf(
			"shutdown http server: %w",
			err,
		)
	}

	if err := <-serverErrCh; err != nil {
		return fmt.Errorf(
			"http server after shutdown: %w",
			err,
		)
	}

	return nil
}

func New(
	ctx context.Context,
) (*App, error) {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf(
			"load config: %w",
			err,
		)
	}

	pool, err := newPostgres(
		ctx,
		cfg.DatabaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"connect postgres: %w",
			err,
		)
	}

	readerRepository :=
		postgres.NewReaderRepository(pool)

	readerService :=
		reader.NewService(readerRepository)

	readerHandler :=
		handler.NewReaderHandler(readerService)

	logger := newLogger()

	router :=
		httptransport.NewRouter(
			readerHandler,
			logger,
		)

	server :=
		httptransport.NewServer(
			cfg.HTTPAddr,
			router,
		)

	return &App{
		server: server,
		pool:   pool,
	}, nil
}

func (a *App) Close() {
	a.pool.Close()
}
