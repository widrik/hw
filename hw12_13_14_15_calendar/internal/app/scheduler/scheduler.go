package scheduler

import (
	"context"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/entities"
	"go.uber.org/zap"
)

type Scheduler struct {
	Publisher  entities.PublisherInterface
	GrpcClient spec.CalendarServiceClient
}

func InitScheduler(publisher entities.PublisherInterface, grpcClient spec.CalendarServiceClient) entities.SchedulerInterface {
	return &Scheduler{
		Publisher:  publisher,
		GrpcClient: grpcClient,
	}
}

func (scheduler *Scheduler) Process() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventsResponse, err := scheduler.GrpcClient.GetList(ctx, &spec.GetListRequest{})
	if err != nil {
		zap.L().Error("getting events error", zap.Error(err))
	}

	timeNow := time.Now()
	date := timeNow.AddDate(-1, 0, 0)

	for _, event := range eventsResponse.Event {
		evStartAt, err := ptypes.Timestamp(event.Start)
		if err != nil {
			zap.L().Error("start time error", zap.Error(err))
		}

		evFinishedAt, err := ptypes.Timestamp(event.Finish)
		if err != nil {
			zap.L().Error("finish time error", zap.Error(err))
		}

		evNotifyAt, err := ptypes.Timestamp(event.NotifyTime)
		if err != nil {
			zap.L().Error("notify time error", zap.Error(err))
		}

		if evFinishedAt.Before(date) {
			id, err := uuid.Parse(event.Uuid)
			if err != nil {
				zap.L().Error("parse id error", zap.Error(err))
			}

			_, err = scheduler.GrpcClient.Delete(ctx, &spec.DeleteRequest{
				Uuid: id.String(),
			})
			if err != nil {
				zap.L().Error("delete error", zap.Error(err))
			}
		}

		if (time.Now().Before(evStartAt)) && timeNow.Equal(evNotifyAt) {
			convertedEvent := convertGrpcEvent(event)
			err := scheduler.Publish(&convertedEvent)
			if err != nil {
				zap.L().Error("publish error", zap.Error(err))
			}
		}
	}
}

func convertGrpcEvent(event *spec.Event) entities.Event {
	return entities.Event{
		ID:          event.Uuid,
		Title:       event.Title,
		Description: event.Description,
		StartedAt:   event.Start.String(),
		FinishedAt:  event.Finish.String(),
		NotifyAt:    event.Finish.String(),
		UserID:      strconv.FormatInt(event.UserId, 10),
	}
}

func (scheduler *Scheduler) Publish(event *entities.Event) error {
	err := scheduler.Publisher.Send(event)
	if err != nil {
		return err
	}

	return nil
}

func (scheduler *Scheduler) Stop() error {
	err := scheduler.Publisher.Close()
	if err != nil {
		return err
	}

	return nil
}
