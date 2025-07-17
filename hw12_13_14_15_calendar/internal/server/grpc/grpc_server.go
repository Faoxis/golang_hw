package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/server"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	return nil
}

func (s *CalendarGRPCServer) Start(_ context.Context) error {
	lis, err := net.Listen("tcp", s.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpcServer = grpc.NewServer()
	api.RegisterCalendarServiceServer(s.grpcServer, s)

	s.logger.Info("gRPC server started on port " + s.port)
	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}

// CreateEvent - создание нового события
func (s *CalendarGRPCServer) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.EventResponse, error) {
	s.logger.Info("gRPC CreateEvent called")

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	id := uuid.New().String()
	startTime := req.StartTime.AsTime()
	duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
	notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())

	err := s.app.CreateEvent(
		ctx,
		id, req.Title, req.Description, req.UserId,
		startTime, duration, notifyBefore,
	)
	if err != nil {
		s.logger.Error("Failed to create event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to create event")
	}

	// Получаем созданное событие для ответа
	event, err := s.app.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get created event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to get created event")
	}

	return mapStorageEventToProtoEvent(event), nil
}

// UpdateEvent - обновление существующего события
func (s *CalendarGRPCServer) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.EventResponse, error) {
	s.logger.Info("gRPC UpdateEvent called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	startTime := req.StartTime.AsTime()
	duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
	notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())

	err := s.app.UpdateEvent(
		ctx,
		req.Id, req.Title, req.Description, req.UserId,
		startTime, duration, notifyBefore,
	)
	if err != nil {
		s.logger.Error("Failed to update event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to update event")
	}

	// Получаем обновлённое событие
	event, err := s.app.GetEventByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get updated event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to get updated event")
	}

	return mapStorageEventToProtoEvent(event), nil
}

// DeleteEvent - удаление события
func (s *CalendarGRPCServer) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	s.logger.Info("gRPC DeleteEvent called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to delete event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to delete event")
	}

	return &api.DeleteEventResponse{Success: true}, nil
}

// GetEvent - получение события по ID
func (s *CalendarGRPCServer) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.EventResponse, error) {
	s.logger.Info("gRPC GetEvent called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	event, err := s.app.GetEventByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get event: " + err.Error())
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return mapStorageEventToProtoEvent(event), nil
}

// ListEvents - получение списка событий за период
func (s *CalendarGRPCServer) ListEvents(ctx context.Context, req *api.ListEventsRequest) (*api.ListEventsResponse, error) {
	s.logger.Info("gRPC ListEvents called")

	if req.Date == nil {
		return nil, status.Error(codes.InvalidArgument, "date is required")
	}

	if req.Period == "" {
		return nil, status.Error(codes.InvalidArgument, "period is required")
	}

	date := req.Date.AsTime()
	var events []storage.Event
	var err error

	switch req.Period {
	case "day":
		events, err = s.app.ListEventsForDay(ctx, date)
	case "week":
		events, err = s.app.ListEventsForWeek(ctx, date)
	case "month":
		events, err = s.app.ListEventsForMonth(ctx, date)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid period: must be 'day', 'week', or 'month'")
	}

	if err != nil {
		s.logger.Error("Failed to list events: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to list events")
	}

	protoEvents := make([]*api.EventResponse, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, mapStorageEventToProtoEvent(event))
	}

	return &api.ListEventsResponse{Events: protoEvents}, nil
}

func mapStorageEventToProtoEvent(event storage.Event) *api.EventResponse {
	return &api.EventResponse{
		Id:           event.ID,
		Title:        event.Title,
		StartTime:    timestamppb.New(event.StartTime),
		Duration:     durationpb.New(time.Duration(event.Duration)),
		Description:  event.Description,
		UserId:       event.UserID,
		NotifyBefore: durationpb.New(time.Duration(event.NotifyBefore)),
	}
}
