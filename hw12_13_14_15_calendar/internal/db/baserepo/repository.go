package baserepo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	Title       string
	Description string
	StartedAt   time.Time
	FinishedAt  time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
	NotifyAt    time.Time
	UserID      int64
}

type EventsRepo interface {
	Add(event Event) (uuid.UUID, error)
	Update(uuid uuid.UUID, event Event) error
	Delete(uuid uuid.UUID) error
	GetList() ([]Event, error)
}

type Repo struct{}

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("time is busy")
)
