package internalgrpc

import (
	"context"
	"fmt"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type CalendarGRPCServer struct {
	api.UnimplementedCalendarServiceServer
	logger     server.Logger
	port       string
	grpcServer *grpc.Server
	app        server.Application
}

func NewCalendarGRPCServer(logger server.Logger, port string, app server.Application) *CalendarGRPCServer {
	return &CalendarGRPCServer{
		port:   port,
		logger: logger,
		app:    app,
	}
}

func (s *CalendarGRPCServer) Stop(_ context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}

func (s *CalendarGRPCServer) Start(_ context.Context) error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	api.RegisterCalendarServiceServer(s.grpcServer, s)

	s.logger.Info("grpc server started on port " + s.port)
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

func (s *CalendarGRPCServer) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.EventResponse, error) {
	storageEvent, err := s.app.GetEventByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	eventResponse := convertToEventResponse(storageEvent)
	return &eventResponse, nil
}

func (s *CalendarGRPCServer) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.EventResponse, error) {
	s.logger.Info("CreateEvent called")
	id := uuid.New().String()
	duration, err := parseDuration(req.Duration)
	if err != nil {
		return nil, err
	}

	notifyBefore, err := parseDuration(req.NotifyBefore)
	if err != nil {
		return nil, err
	}

	err = s.app.CreateEvent(
		ctx,
		id, req.Title, req.Description, req.UserId,
		req.StartTime.AsTime(),
		*duration,
		*notifyBefore,
	)
	if err != nil {
		return &api.EventResponse{}, err
	}
	return &api.EventResponse{
		Id:           id,
		Title:        req.Title,
		StartTime:    req.StartTime,
		Duration:     req.Duration,
		Description:  req.Description,
		UserId:       req.UserId,
		NotifyBefore: req.NotifyBefore,
	}, nil
}

func parseDuration(s string) (*calendar_types.CalendarDuration, error) {
	dur, err := time.ParseDuration(s)
	if err != nil {
		return nil, fmt.Errorf("invalid duration: %w", err)
	}
	duration := calendar_types.CalendarDuration(dur)
	return &duration, nil
}
