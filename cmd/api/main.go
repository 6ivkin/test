package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/6ivkin/test.git/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	if err := application.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
