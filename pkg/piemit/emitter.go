package piemit

import "context"

type Emitter struct {
	querier       Querier
	eventProvider Provider
}

func NewEmitter() *Emitter {
	return &Emitter{}
}

type Provider interface {
	Emit(ctx context.Context, topicName string, payload any) error
}

func (e *Emitter) Start(ctx context.Context) error {
	for {

		event := e.querier.NextPending(ctx)

		err := event.SetActive()

		err := e.querier.ActiveEvent(ctx, event)

	}

	return nil
}

func (e *Emitter) Stop() error // graceful shutdown
