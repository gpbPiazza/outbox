package main

import (
	"context"
	"os"

	"github.com/gpbPiazza/internal/migrator"
	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
)

func main() {
	logger := log.New()
	ctx := log.SetContext(context.Background(), logger)

	dbx, err := db.Init(ctx)
	if err != nil {
		logger.Error("db init", "err", err)
		os.Exit(1)
	}
	defer func() { _ = dbx.Close() }()

	if err := migrator.Run(ctx, dbx); err != nil {
		logger.Error("migrate", "err", err)
		os.Exit(1)
	}
}
