package db

import (
	"context"
	"database/sql"

	"github.com/gpbPiazza/wes/log"
)

func Init(ctx context.Context) (*sql.DB, error) {
	log, _ := log.FromContext(ctx)

	db, err := sql.Open("", "wes-db")
	if err != nil {
		log.Error("err and server", "err", err)
		return nil, err
	}

	return db, nil
}
