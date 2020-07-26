package inmemory

import (
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"sync"
)

var _ baserepo.EventsRepo = (*Repo)(nil)

type Repo struct {
	mx   sync.Mutex
	id   int64
	list map[uuid.UUID]baserepo.Event
}

func (r *Repo) init() {
	r.list = make(map[uuid.UUID]baserepo.Event)
}

func (r *Repo) Add(event baserepo.Event) (uuid.UUID, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if len(r.list) == 0 {
		r.init()
	}

	var err error

	for _, e := range r.list {
		if e.StartedAt == event.StartedAt {
			err = baserepo.ErrDateBusy
		}
	}

	r.list[event.Id] = event

	return event.Id, err
}

func (r *Repo) Update(uuid uuid.UUID, event baserepo.Event) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, e := range r.list {
		if e.Id == uuid {
			r.list[e.Id] = event

			return nil
		}
	}
	return baserepo.ErrEventNotFound
}

func (r *Repo) Delete(uuid uuid.UUID) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	delete(r.list, uuid)

	return nil
}

func (r *Repo) GetList() ([]baserepo.Event, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	var eventList = make([]baserepo.Event, 0, len(r.list))

	for _, event := range r.list {
		eventList = append(eventList, event)
	}

	return eventList, nil
}
