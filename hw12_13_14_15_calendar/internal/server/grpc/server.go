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
	Server  *grpc.Server
	Address string
}

func NewServer(calender *app.Calendar, listenAddress string) *Server {
	calenderApp = calender

	zapL := zap_logger.L()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(zap_grpc.UnaryServerInterceptor(zapL)))

	srv := &Server{}
	spec.RegisterCalendarServiceServer(grpcServer, srv)

	srv.Server = grpcServer
	srv.Address = listenAddress

	return srv
}

func (srv *Server) Start() error {
	lis, err := net.Listen("tcp", srv.Address)
	if err != nil {
		return err
	}
	err = srv.Server.Serve(lis)

	return err
}

func (srv *Server) Stop() {
	srv.Server.GracefulStop()
}

func (srv *Server) Add(ctx context.Context, request *spec.AddRequest) (*spec.AddResponse, error) {
	event, err := grpcToEvent(request.Event)
	if err != nil {
		return nil, err
	}

	id, err := calenderApp.Add(event)
	if err != nil {
		return nil, err
	}

	return &spec.AddResponse{
		Uuid: id.String(),
	}, nil
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
	}

	return &spec.DeleteResponse{}, calenderApp.Delete(id)
}

func (srv *Server) GetByID(ctx context.Context, request *spec.GetByIdRequest) (*spec.GetByIdResponse, error) {
	id, err := uuid.Parse(request.Uuid)
	if err != nil {
		return nil, err
	}

	event, err := calenderApp.GetEventByID(id)
	if err != nil {
		return nil, err
	}

	grpcEvent, err := EventToGrpc(event)
	if err != nil {
		return nil, err
	}

	response := spec.GetByIdResponse{
		Event: grpcEvent,
	}

	return &response, nil
}

func (srv *Server) GetList(ctx context.Context, request *spec.GetListRequest) (*spec.GetListResponse, error) {
	events, err := calenderApp.GetList()
	if err != nil {
		return nil, err
	}

	responseEvents := make([]*spec.Event, len(events))

	response := spec.GetListResponse{
		Event: responseEvents,
	}
	for _, event := range events {
		grpcEvent, err := EventToGrpc(event)
		if err != nil {
			return nil, err
		}
		response.Event = append(response.Event, grpcEvent)
	}
	
	return &response, nil
}
