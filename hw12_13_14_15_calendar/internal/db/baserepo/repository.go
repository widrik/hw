package baserepo

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Event struct {
	Id          uuid.UUID
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
	ErrEventNotFound = errors.New("Event not found")
	ErrDateBusy = errors.New("Time is busy")
)
