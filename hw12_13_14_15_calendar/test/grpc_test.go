package main

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/inmemory"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	grpcSrv "github.com/widrik/hw/hw12_13_14_15_calendar/internal/server/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const buffer = 1024 * 1024

func TestGrpcEndpoints(t *testing.T) {
	event, _ := grpcSrv.EventToGrpc(getEvent())

	configuration, _ := config.Init("testdata/config/correctData.json")
	_ = logging.Init(configuration.Logging.Level, configuration.Logging.File)

	repo := new(inmemory.Repo)
	calenderApp := app.Calendar{
		Repository: repo,
	}

	listener := bufconn.Listen(buffer)
	address := net.JoinHostPort(configuration.GRPCServer.Host, configuration.GRPCServer.Port)
	srv := grpcSrv.NewServer(&calenderApp, address)

	go func() {
		err := srv.Server.Serve(listener)
		require.Nil(t, err)
	}()

	ctx := context.Background()

	connection, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(ctx context.Context, s string) (conn net.Conn, err error) {
		return listener.Dial()
	}), grpc.WithInsecure())
	require.Nil(t, err)
	defer connection.Close()

	client := spec.NewCalendarServiceClient(connection)

	t.Run("add and delete success", func(t *testing.T) {
		request := &spec.AddRequest{Event: event}
		response, err := client.Add(ctx, request)
		require.Nil(t, err)

		id, err := uuid.Parse(response.Uuid)
		require.Nil(t, err)

		request2 := &spec.DeleteRequest{Uuid: id.String()}
		_, err = client.Delete(ctx, request2)
		require.Nil(t, err)
	})

	t.Run("delete error", func(t *testing.T) {
		request := &spec.DeleteRequest{Uuid: "blabla"}
		_, err = client.Delete(ctx, request)
		require.NotNil(t, err)
	})

	t.Run("get list success", func(t *testing.T) {
		request := &spec.GetListRequest{}
		_, err = client.GetList(ctx, request)
		require.Nil(t, err)
	})

	t.Run("get by id not valid id", func(t *testing.T) {
		request := &spec.GetByIdRequest{Uuid: uuid.New().String()}
		_, err = client.GetByID(ctx, request)
		require.NotNil(t, err)
	})

	t.Run("get by id valid id", func(t *testing.T) {
		request := &spec.AddRequest{Event: event}
		response, err := client.Add(ctx, request)
		require.Nil(t, err)

		id, err := uuid.Parse(response.Uuid)
		require.Nil(t, err)

		request2 := &spec.GetByIdRequest{Uuid: id.String()}
		_, err = client.GetByID(ctx, request2)
		require.Nil(t, err)
	})
}

func getEvent() baserepo.Event {
	startAt := time.Date(2021, 1, 1, 10, 10, 0, 0, time.UTC)
	finishedAt := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)
	notifyAt := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)

	baseRepoEvent := baserepo.Event{
		Title:       "Встреча очень важная",
		Description: "Тестовая очень важная встреча",
		StartedAt:   startAt,
		FinishedAt:  finishedAt,
		NotifyAt:    notifyAt,
		UserID:      1,
	}

	return baseRepoEvent
}