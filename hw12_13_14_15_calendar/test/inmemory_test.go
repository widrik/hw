package main

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/inmemory"
	"testing"
	"time"
)

func TestMemRepo(t *testing.T) {
	t.Run("Create event with no errors", func(t *testing.T) {
		var r inmemory.Repo

		event := baserepo.Event{
			ID:          uuid.New(),
			Title:       "Testing",
			Description: "testing of what is here",
			StartedAt:   time.Now(),
			FinishedAt:  time.Now().Add(30 * time.Minute),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			NotifyAt:    time.Now(),
			UserID:      10,
		}

		_, err := r.Add(event)
		require.NoError(t, err)
	})

	t.Run("Create events on same date", func(t *testing.T) {
		var r inmemory.Repo

		timeForEvent := time.Now()

		event1 := baserepo.Event{
			ID:          uuid.New(),
			Title:       "Testing 1",
			Description: "testing of what is here",
			StartedAt:   timeForEvent,
			FinishedAt:  time.Now().Add(30 * time.Minute),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			NotifyAt:    time.Now(),
			UserID:      10,
		}

		event2 := baserepo.Event{
			ID:          uuid.New(),
			Title:       "Testing 2",
			Description: "testing of what is here",
			StartedAt:   timeForEvent,
			FinishedAt:  time.Now().Add(30 * time.Minute),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			NotifyAt:    time.Now(),
			UserID:      10,
		}

		_, err1 := r.Add(event1)
		require.NoError(t, err1)

		_, err2 := r.Add(event2)
		require.Equal(t,baserepo.ErrDateBusy,err2)
	})

	t.Run("Delete event with no errors", func(t *testing.T) {
		var r inmemory.Repo

		event := baserepo.Event{
			ID:          uuid.New(),
			Title:       "Testing",
			Description: "testing of what is here",
			StartedAt:   time.Now(),
			FinishedAt:  time.Now().Add(30 * time.Minute),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			NotifyAt:    time.Now(),
			UserID:      10,
		}

		id, err1 := r.Add(event)
		require.NoError(t, err1)

		err2 := r.Delete(id)
		require.Equal(t,nil,err2)
	})

}