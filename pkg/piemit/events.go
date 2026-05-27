package piemit

import (
	"encoding/json"
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

type Outbox struct {
	MaxAttempts int
}

// Register -> é cliente da tabela events registrando no SQL
// Emitter -> é pacote events publicando o evento apartida do Register

type BoxParam struct {
	// Event payload
	Payload json.RawMessage
	// Event topic name
	Topic string
}

// Event is the event type registed in the DB
type Event struct {
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
	Errors    map[int]string
	CreatedAt time.Time
	UpdatedAt time.Time
}
