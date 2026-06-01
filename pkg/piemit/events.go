package piemit

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	// Pending : event is ready to be processed and will be picked up by a free worker.
	Pending Status = "pending"

	// Active : event is being processed by a worker (i.e. handler is invoked with the event).
	Active Status = "active"

	// Retry : event failed to process the event and the task is waiting to be retried in the future.
	Retry Status = "retry"

	// Archived : event reached its max retry and stored in an archive for manual inspection.
	Archived Status = "archived"

	// Completed: event was successfully processed and retained until retention TTL expires (Only applies to tasks with Retention option).
	Completed Status = "completed"
)

// Register -> é cliente da tabela events registrando no SQL
// Emitter -> é pacote events publicando o evento apartida do Register

// event is the event type registed in the DB
type event struct {
	ID uuid.UUID

	// Event status
	Status Status
	// Event payload
	Payload any
	// Event topic name
	Topic string
	// Attempets is the quantity of times tried to execute a event
	Attempts int
	// Errors is the attempt number + the error occured
	// TODO: implement attempt ID or number + error feature relation
	LastErr   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
