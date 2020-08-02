package inmemory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
)

var _ baserepo.EventsRepo = (*Repo)(nil)

type Repo struct {
	mx   sync.Mutex
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

	r.list[event.ID] = event

	return event.ID, err
}

func (r *Repo) Update(uuid uuid.UUID, event baserepo.Event) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, e := range r.list {
		if e.ID == uuid {
			r.list[e.ID] = event

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

func (r *Repo) GetEventByID(uuid uuid.UUID) (baserepo.Event, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	for _, e := range r.list {
		if e.ID == uuid {
			return r.list[e.ID], nil
		}
	}

	emptyEvent := baserepo.Event{}

	return emptyEvent, baserepo.ErrEventNotFound
}

func (r *Repo) GetList() ([]baserepo.Event, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	eventList := make([]baserepo.Event, 0, len(r.list))

	for _, event := range r.list {
		eventList = append(eventList, event)
	}

	return eventList, nil
}
