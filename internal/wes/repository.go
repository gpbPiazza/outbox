package wes

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gpbPiazza/pkg/db"
	"github.com/gpbPiazza/pkg/log"
)

type Repository interface {
	Querier
	RunInTx(ctx context.Context, fn func(q Querier) error) error
}

type Querier interface {
	CreateMove(ctx context.Context, wesID uuid.UUID, status MoveStatus, description string) (Move, error)
	CreateWes(ctx context.Context, height string) (Wes, error)
	GetWesByID(ctx context.Context, id uuid.UUID) (Wes, error)
}

// queries is the bro that implements all the domain queries
type queries struct {
	db db.DBTX
}

// repository os the type wrapper to hold domain queries and access generic DB type
type repository struct {
	db *sql.DB
	*queries
}

var _ Querier = (*queries)(nil)

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
		queries: &queries{
			db: db,
		},
	}
}

func (q *queries) WithTx(tx *sql.Tx) *queries {
	return &queries{db: tx}
}

func (q *repository) RunInTx(ctx context.Context, fn func(qw Querier) error) error {
	tx, err := q.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	err = fn(q.WithTx(tx))
	if err != nil {
		log.FromContext(ctx).Error("erro during tx", "err", err)
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			log.FromContext(ctx).Error("rollback failed", "err", rbErr)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.FromContext(ctx).Error("erro during tx", "err", err)
	}

	return err
}

const createMove = `-- name: CreateMove :one
INSERT INTO moves (wes_id, status, description)
VALUES ($1, $2, $3)
RETURNING id, wes_id, status, description, created_at;`

func (q *queries) CreateMove(ctx context.Context, wesID uuid.UUID, status MoveStatus, description string) (Move, error) {
	row := q.db.QueryRowContext(ctx, createMove, wesID, status, description)
	var i Move
	if err := row.Scan(&i.ID, &i.WesID, &i.Status, &i.Description, &i.CreatedAt); err != nil {
		return i, fmt.Errorf("scan created move row: %w", err)
	}
	return i, nil
}

const createWes = `-- name: CreateWes :one
INSERT INTO wes (height)
VALUES ($1)
RETURNING id, height, created_at, updated_at;
`

func (q *queries) CreateWes(ctx context.Context, height string) (Wes, error) {
	row := q.db.QueryRowContext(ctx, createWes, height)
	var i Wes
	if err := row.Scan(
		&i.ID,
		&i.Height,
		&i.CreatedAt,
		&i.UpdatedAt,
	); err != nil {
		return i, fmt.Errorf("scan created wes row: %w", err)
	}
	return i, nil
}

const getWesByID = `-- name: GetWesByID :one
SELECT 
	w.id, w.height, w.created_at, w.updated_at
FROM wes w
WHERE w.id = $1
`

const getMovesByWesID = `-- name: GetMovesByWesID :many
SELECT
	id, wes_id, status, description, created_at
FROM moves m
WHERE m.wes_id = $1
`

func (q *queries) GetWesByID(ctx context.Context, id uuid.UUID) (Wes, error) {
	row := q.db.QueryRowContext(ctx, getWesByID, id)

	var i Wes
	err := row.Scan(
		&i.ID,
		&i.Height,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	if err != nil {
		return i, fmt.Errorf("scan wes row: %w", err)
	}

	rows, err := q.db.QueryContext(ctx, getMovesByWesID, i.ID)
	if err != nil {
		return i, fmt.Errorf("query moves by wes id: %w", err)
	}
	defer rows.Close()

	ms := []Move{}
	for rows.Next() {
		var m Move
		if err := rows.Scan(&m.ID, &m.WesID, &m.Status, &m.Description, &m.CreatedAt); err != nil {
			return i, fmt.Errorf("scan move row: %w", err)
		}

		ms = append(ms, m)
	}
	if err := rows.Err(); err != nil {
		return i, fmt.Errorf("iterate move rows: %w", err)
	}

	i.Moves = ms

	return i, nil
}
