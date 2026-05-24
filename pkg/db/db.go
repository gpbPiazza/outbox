package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gpbPiazza/pkg/envs"
	"github.com/gpbPiazza/pkg/log"
)

func Init(ctx context.Context) (*sql.DB, error) {
	logger := log.FromContext(ctx)
	cfg := envs.Envs().DB

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		logger.Error("open db", "err", err)
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MinConns)
	db.SetConnMaxLifetime(cfg.MaxConnsLifeTime)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		logger.Error("ping db", "err", err)
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	logger.Info("db connected")

	return db, nil
}

type dbCtxKey struct{}

func SetContext(ctx context.Context, db *sql.DB) context.Context {
	return context.WithValue(ctx, dbCtxKey{}, db)
}

func FromContext(ctx context.Context) *sql.DB {
	d := ctx.Value(dbCtxKey{})
	db, ok := d.(*sql.DB)
	if !ok {
		panic("database not in the context")
	}
	return db
}
