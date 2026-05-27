package main

import (
	"context"
	"net/http"

	"github.com/gpbPiazza/internal/wes"
	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
)

func main() {
	logger := log.New()
	ctx := log.SetContext(context.Background(), logger)

	dbx, err := db.Init(ctx)
	if err != nil {
		logger.Error("db init", "err", err)
		return
	}
	defer func() { _ = dbx.Close() }()
	ctx = db.SetContext(ctx, dbx)

	server := http.NewServeMux()
	wes.RegisterHandlers(ctx, dbx, server)

	logger.Info("starting wes server")

	if err := http.ListenAndServe("localhost:8080", server); err != nil {
		logger.Error("listen and server", "err", err)
		return
	}
}
