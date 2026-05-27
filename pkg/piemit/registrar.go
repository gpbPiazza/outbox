package piemit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Registrar interface {
	Register(ctx context.Context, topicName string, payload any) error
}

type registrar struct {
	querirer Querier
}

func (r *registrar) Register(ctx context.Context, topicName string, payload any) error {
	e := newPendingEvent(topicName, payload)
	return r.querirer.Create(ctx, e)
}

func newPendingEvent(topicName string, payload any) Event {
	return Event{
		ID:        uuid.New(),
		Status:    Pending,
		Payload:   payload,
		Topic:     topicName,
		Attempts:  0,
		Errors:    nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
