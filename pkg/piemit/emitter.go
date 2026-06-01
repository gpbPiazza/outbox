package piemit

import (
	"context"

	"github.com/gpbPiazza/pkg/db"
)

type Emitter struct {
	querier       Querier
	eventProvider Provider
	db            db.DBTX
}

func NewEmitter(db db.DBTX) *Emitter {
	return &Emitter{
		db: db,
	}
}

type Provider interface {
	Emit(ctx context.Context, topicName string, payload any) error
}

func (e *Emitter) Start(ctx context.Context) error {
	for {
		// event := (ctx)
		//
		// err := event.SetActive()
		//
		// err := e.querier.ActiveEvent(ctx, event)
	}
}

func (e *Emitter) Stop() error // graceful shutdown
