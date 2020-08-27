package calendar

import (
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
)

type Calendar struct {
	Repository baserepo.EventsRepo
}

func NewCalendar(r baserepo.EventsRepo) Calendar {
	return Calendar{
		Repository: r,
	}
}

func (c *Calendar) Add(event baserepo.Event) (uuid.UUID, error) {
	newEventUUID, err := c.Repository.Add(event)

	return newEventUUID, err
}

func (c *Calendar) Update(event baserepo.Event, uuid uuid.UUID) error {
	return c.Repository.Update(uuid, event)
}

func (c *Calendar) Delete(uuid uuid.UUID) error {
	return c.Repository.Delete(uuid)
}

func (c *Calendar) GetEventByID(uuid uuid.UUID) (baserepo.Event, error) {
	return c.Repository.GetEventByID(uuid)
}

func (c *Calendar) GetList() ([]baserepo.Event, error) {
	return c.Repository.GetList()
}
