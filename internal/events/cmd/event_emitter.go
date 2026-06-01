package main

import (
	"context"

	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
	"github.com/gpbPiazza/pkg/piemit"
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

	emmiter := piemit.NewEmitter(dbx)

	emmiter.Start(ctx)
}
