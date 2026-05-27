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

	repo := wes.NewRepository(dbx)
	wesHandler := wes.NewHandler(repo)

	server.Handle("/wes", ctxMiddleware(wesHandler, ctx))

	logger.Info("starting wes server")

	if err := http.ListenAndServe("localhost:8080", server); err != nil {
		logger.Error("listen and server", "err", err)
		return
	}
}

func ctxMiddleware(next http.Handler, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
