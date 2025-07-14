package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	api.UnimplementedEventServiceServer
	logger      Logger
	application Application
	server      *grpc.Server
	port        int
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Warn(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		id, title, description, userID string,
		startTime time.Time,
		duration, notifyBefore calendar_types.CalendarDuration,
	) error

	UpdateEvent(
		ctx context.Context,
		id, title, description, userID string,
		startTime time.Time,
		duration, notifyBefore calendar_types.CalendarDuration,
	) error

	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (storage.Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

func NewServer(logger Logger, host string, port int, app Application) *Server {
	grpcServer := grpc.NewServer()
	server := &Server{
		logger:      logger,
		application: app,
		server:      grpcServer,
		port:        port,
	}
	api.RegisterEventServiceServer(grpcServer, server)
	return server
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.logger.Info(fmt.Sprintf("gRPC server is running on port %d...", s.port))

	go func() {
		if err := s.server.Serve(lis); err != nil {
			s.logger.Error("gRPC server error: " + err.Error())
		}
	}()

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	return nil
}

// CreateEvent реализует gRPC метод создания события
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
	s.logger.Info("gRPC CreateEvent called")

	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	startTime := req.StartTime.AsTime()
	duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
	notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())

	// Генерируем UUID для события
	id := generateUUID()

	err := s.application.CreateEvent(
		ctx,
		id,
		req.Title,
		req.Description,
		req.UserId,
		startTime,
		duration,
		notifyBefore,
	)
	if err != nil {
		s.logger.Error("Failed to create event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to create event")
	}

	// Получаем созданное событие
	event, err := s.application.GetEventByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get created event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to get created event")
	}

	return &api.CreateEventResponse{
		Event: mapStorageEventToProtoEvent(event),
	}, nil
}

// UpdateEvent реализует gRPC метод обновления события
func (s *Server) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
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

	err := s.application.UpdateEvent(
		ctx,
		req.Id,
		req.Title,
		req.Description,
		req.UserId,
		startTime,
		duration,
		notifyBefore,
	)
	if err != nil {
		s.logger.Error("Failed to update event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to update event")
	}

	// Получаем обновленное событие
	event, err := s.application.GetEventByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get updated event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to get updated event")
	}

	return &api.UpdateEventResponse{
		Event: mapStorageEventToProtoEvent(event),
	}, nil
}

// DeleteEvent реализует gRPC метод удаления события
func (s *Server) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
	s.logger.Info("gRPC DeleteEvent called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	err := s.application.DeleteEvent(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to delete event: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to delete event")
	}

	return &api.DeleteEventResponse{
		Success: true,
	}, nil
}

// GetEvent реализует gRPC метод получения события по ID
func (s *Server) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
	s.logger.Info("gRPC GetEvent called")

	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	event, err := s.application.GetEventByID(ctx, req.Id)
	if err != nil {
		s.logger.Error("Failed to get event: " + err.Error())
		return nil, status.Error(codes.NotFound, "event not found")
	}

	return &api.GetEventResponse{
		Event: mapStorageEventToProtoEvent(event),
	}, nil
}

// ListEvents реализует gRPC метод получения списка событий
func (s *Server) ListEvents(ctx context.Context, req *api.ListEventsRequest) (*api.ListEventsResponse, error) {
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
		events, err = s.application.ListEventsForDay(ctx, date)
	case "week":
		events, err = s.application.ListEventsForWeek(ctx, date)
	case "month":
		events, err = s.application.ListEventsForMonth(ctx, date)
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid period: must be 'day', 'week', or 'month'")
	}

	if err != nil {
		s.logger.Error("Failed to list events: " + err.Error())
		return nil, status.Error(codes.Internal, "failed to list events")
	}

	protoEvents := make([]*api.Event, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, mapStorageEventToProtoEvent(event))
	}

	return &api.ListEventsResponse{
		Events: protoEvents,
	}, nil
}

// Вспомогательные функции для маппинга
func mapStorageEventToProtoEvent(event storage.Event) *api.Event {
	return &api.Event{
		Id:           event.ID,
		Title:        event.Title,
		StartTime:    timestamppb.New(event.StartTime),
		Duration:     durationpb.New(time.Duration(event.Duration)),
		Description:  event.Description,
		UserId:       event.UserID,
		NotifyBefore: durationpb.New(time.Duration(event.NotifyBefore)),
	}
}

func generateUUID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "fallback-uuid"
	}
	return newUUID.String()
}
