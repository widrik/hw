package grpc

import (
	"context"
	"net"

	"github.com/google/uuid"
	zap_grpc "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/widrik/hw/hw12_13_14_15_calendar/api/spec"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
	zap_logger "go.uber.org/zap"
	"google.golang.org/grpc"
)

var calenderApp *app.Calendar

type Server struct {
	server  *grpc.Server
	address string
}

func NewServer(calender *app.Calendar, listenAddress string) *Server {
	calenderApp = calender

	zapL := zap_logger.L()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(zap_grpc.UnaryServerInterceptor(zapL)))

	srv := &Server{}
	spec.RegisterCalendarServiceServer(grpcServer, srv)
	srv.server = grpcServer
	srv.address = listenAddress

	return srv
}

func (srv *Server) Start() error {
	lis, err := net.Listen("tcp", srv.address)
	if err != nil {
		return err
	}
	err = srv.server.Serve(lis)
	return err
}

func (srv *Server) Stop() {
	srv.server.GracefulStop()
}

func (srv *Server) Add(ctx context.Context, request *spec.AddRequest) (*spec.AddResponse, error) {
	event, err := grpcToEvent(request.Event)
	if err != nil {
		return nil, err
	}

	id, err := calenderApp.Add(event)

	if err != nil {
		return nil, err
	} else {
		return &spec.AddResponse{
			Uuid: id.String(),
		}, nil
	}
}

func (srv *Server) Update(ctx context.Context, request *spec.UpdateRequest) (*spec.UpdateResponse, error) {
	event, err := grpcToEvent(request.Event)
	if err != nil {
		return nil, err
	}

	return &spec.UpdateResponse{}, calenderApp.Update(event, event.ID)
}

func (srv *Server) Delete(ctx context.Context, request *spec.DeleteRequest) (*spec.DeleteResponse, error) {
	id, err := uuid.Parse(request.Uuid)

	if err != nil {
		return nil, err
	} else {
		return &spec.DeleteResponse{}, calenderApp.Delete(id)
	}
}

func (srv *Server) GetById(ctx context.Context, request *spec.GetByIdRequest) (*spec.GetByIdResponse, error) {
	id, err := uuid.Parse(request.Uuid)

	if err != nil {
		return nil, err
	} else {
		event, err := calenderApp.GetEventByID(id)

		if err != nil {
			return nil, err
		} else {
			grpcEvent, err := eventToGrpc(event)

			if err != nil {
				return nil, err
			}

			response := spec.GetByIdResponse{
				Event: grpcEvent,
			}

			return &response, nil
		}
	}
}

func (srv *Server) GetList(ctx context.Context, request *spec.GetListRequest) (*spec.GetListResponse, error) {
	events, err := calenderApp.GetList()

	if err != nil {
		return nil, err
	} else {
		responseEvents := make([]*spec.Event, len(events))

		response := spec.GetListResponse{
			Event: responseEvents,
		}
		for _, event := range events {
			grpcEvent, err := eventToGrpc(event)
			if err != nil {
				return nil, err
			}
			response.Event = append(response.Event, grpcEvent)
		}
		return &response, nil
	}
}