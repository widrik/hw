package grpc

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
)

func grpcToEvent(gerpcEvent *spec.Event) (baserepo.Event, error) {
	var event baserepo.Event
	idString, err := uuid.Parse(gerpcEvent.Uuid)
	if err != nil {
		return event, err
	}

	event = baserepo.Event{
		ID:          idString,
		Title:       gerpcEvent.Title,
		Description: gerpcEvent.Description,
		StartedAt:   time.Unix(gerpcEvent.Start.GetSeconds(), 0),
		FinishedAt:  time.Unix(gerpcEvent.Finish.GetSeconds(), 0),
		NotifyAt:    time.Unix(gerpcEvent.NotifyTime.GetSeconds(), 0),
		UserID:      gerpcEvent.UserId,
	}

	return event, nil
}

func eventToGrpc(event baserepo.Event) (*spec.Event, error) {
	var grpcEvent *spec.Event

	startedAt, err := ptypes.TimestampProto(event.StartedAt)
	if err != nil {
		return nil, err
	}

	finishedAt, err := ptypes.TimestampProto(event.FinishedAt)
	if err != nil {
		return nil, err
	}

	notifiedAt, err := ptypes.TimestampProto(event.NotifyAt)
	if err != nil {
		return nil, err
	}

	grpcEvent = &spec.Event{
		Uuid:        event.ID.String(),
		Title:       event.Title,
		Description: event.Description,
		Start:       startedAt,
		Finish:      finishedAt,
		NotifyTime:  notifiedAt,
		UserId:      event.UserID,
	}

	return grpcEvent, nil
}
