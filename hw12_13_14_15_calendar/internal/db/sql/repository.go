package sql

import (
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
)

var _ baserepo.EventsRepo = (*Repo)(nil)

type Repo struct {
	base *sqlx.DB
}

func NewDBConnection(sourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", sourceName)
	if err != nil {
		log.Fatalln(err)
	}

	return db, err
}

func (r *Repo) Add(event baserepo.Event) (uuid.UUID, error) {
	id := uuid.New()

	_, err := r.base.Exec(`INSERT INTO events (id, title, description, start_at, finished_at, created_at, updated_at, user_id, notify_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, id.String(), event.Title, event.Description, event.StartedAt, event.FinishedAt, event.CreatedAt, event.UpdatedAt, event.UserID, event.NotifyAt)

	return id, err
}

func (r *Repo) Update(uuid uuid.UUID, event baserepo.Event) error {
	event, err := r.GetEventByID(uuid)
	if err != nil {
		return baserepo.ErrEventNotFound
	}
	_, err = r.base.NamedExec(`UPDATE events SET title=:title, description=:description, start_at=:start_at, finished_at=:finished_at,  user_id=:user_id, notify_at=:notify_at WHERE :uuid = :uuid`, event)
	return err
}

func (r *Repo) Delete(uuid uuid.UUID) error {
	_, err := r.GetEventByID(uuid)
	if err != nil {
		return baserepo.ErrEventNotFound
	}

	id := uuid.String()
	_, err = r.base.Exec("DELETE FROM events WHERE uuid = ?", id)

	return err
}

func (r *Repo) GetList() ([]baserepo.Event, error) {
	events := []baserepo.Event{}
	err := r.base.Select(&events, "SELECT * FROM events")

	return events, err
}

func (r *Repo) GetEventByID(uuid uuid.UUID) (baserepo.Event, error) {
	id := uuid.String()
	event := baserepo.Event{}
	err := r.base.Get(&event, "SELECT * FROM events WHERE uuid = ?", id)

	return event, err
}
