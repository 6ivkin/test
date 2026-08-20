package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/6ivkin/test.git/internal/config"
	"github.com/6ivkin/test.git/internal/database"
	"github.com/6ivkin/test.git/internal/reader"
	"github.com/6ivkin/test.git/internal/repository/postgres"
	httptransport "github.com/6ivkin/test.git/internal/transport/http"
	"github.com/6ivkin/test.git/internal/transport/http/handler"
	"github.com/joho/godotenv"
)

const shutdownTimeout = 10 * time.Second

func Run() error {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	appCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	pool, err := database.NewPostgres(appCtx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}
	defer pool.Close()

	readerRepository :=
		postgres.NewReaderRepository(pool)

	readerService :=
		reader.NewService(readerRepository)

	readerHandler :=
		handler.NewReaderHandler(readerService)

	router :=
		httptransport.NewRouter(readerHandler)

	server := httptransport.NewServer(cfg.HTTPAddr, router)

	serverErrCh := make(chan error, 1)

	go func() {
		serverErrCh <- server.Run()
	}()

	select {
	case <-appCtx.Done():
	case err := <-serverErrCh:
		if err != nil {
			return fmt.Errorf("http server: %w", err)
		}

		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := <-serverErrCh; err != nil {
		return fmt.Errorf("http server after shutdown: %w", err)
	}

	return nil
}
