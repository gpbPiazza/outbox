// Package migrator runs SQL migrations on startup using pressly/goose.
package migrator

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/gpbPiazza/pkg/envs"
	"github.com/gpbPiazza/pkg/log"
)

func Run(ctx context.Context, db *sql.DB) error {
	logger := log.FromContext(ctx)
	dir := envs.Envs().DB.GooseMigrationsDir

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	logger.Info("running migrations", "dir", dir)
	if err := goose.UpContext(ctx, db, dir); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}
	logger.Info("migrations complete")
	return nil
}
