package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gpbPiazza/internal/wes"
	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	status := run(ctx)

	cancel()
	os.Exit(status)
}

func run(ctx context.Context) int {
	logger := log.New()
	ctx = log.SetContext(ctx, logger)

	dbx, err := db.Init(ctx)
	if err != nil {
		logger.Error("db init", "err", err)
		return 1
	}
	ctx = db.SetContext(ctx, dbx)

	mux := http.NewServeMux()
	wes.RegisterHandlers(ctx, dbx, mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	logger.Info("starting wes server")

	var shutDownErrs []error
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("listen and server", "err", err)
			shutDownErrs = append(shutDownErrs, err)
		}
	}()

	<-ctx.Done()

	logger.Info("gracefully shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("http server shutdown fail", "err", err)
		shutDownErrs = append(shutDownErrs, err)
	}

	if err := dbx.Close(); err != nil {
		logger.Error("close database con fail", "err", err)
		shutDownErrs = append(shutDownErrs, err)
	}

	if len(shutDownErrs) != 0 {
		return 1
	}

	return 0
}
