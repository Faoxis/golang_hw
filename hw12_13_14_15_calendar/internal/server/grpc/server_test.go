package grpcserver

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	app2 "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	logger2 "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockLogger - мок для логгера
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string) { m.Called(msg) }
func (m *MockLogger) Info(msg string)  { m.Called(msg) }
func (m *MockLogger) Error(msg string) { m.Called(msg) }
func (m *MockLogger) Warn(msg string)  { m.Called(msg) }

// MockApplication - мок для приложения
type MockApplication struct {
	mock.Mock
}

func (m *MockApplication) CreateEvent(
	ctx context.Context,
	id, title, description, userID string,
	startTime time.Time,
	duration, notifyBefore calendar_types.CalendarDuration,
) error {
	args := m.Called(ctx, id, title, description, userID, startTime, duration, notifyBefore)
	return args.Error(0)
}

func (m *MockApplication) UpdateEvent(
	ctx context.Context,
	id, title, description, userID string,
	startTime time.Time,
	duration, notifyBefore calendar_types.CalendarDuration,
) error {
	args := m.Called(ctx, id, title, description, userID, startTime, duration, notifyBefore)
	return args.Error(0)
}

func (m *MockApplication) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockApplication) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func (m *MockApplication) ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	args := m.Called(ctx, date)
	return args.Get(0).([]storage.Event), args.Error(1)
}

func TestCreateEvent(t *testing.T) {
	mockApp := new(MockApplication)
	mockLogger := new(MockLogger)

	server := &Server{
		logger:      mockLogger,
		application: mockApp,
	}

	startTime := time.Now()
	req := &api.CreateEventRequest{
		Title:        "Test Event",
		StartTime:    timestamppb.New(startTime),
		Duration:     durationpb.New(time.Hour),
		Description:  "Test Description",
		UserId:       "user123",
		NotifyBefore: durationpb.New(30 * time.Minute),
	}

	// Мокаем вызовы логгера
	mockLogger.On("Info", "gRPC CreateEvent called").Return()
	mockLogger.On("Info", mock.Anything).Return() // для других логов

	mockApp.On("CreateEvent", mock.Anything, mock.Anything, req.Title, req.Description, req.UserId, mock.Anything, calendar_types.CalendarDuration(time.Hour), calendar_types.CalendarDuration(30*time.Minute)).Return(nil)

	expectedEvent := storage.Event{
		ID:           "test-id",
		Title:        req.Title,
		StartTime:    startTime,
		Duration:     calendar_types.CalendarDuration(time.Hour),
		Description:  req.Description,
		UserID:       req.UserId,
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	mockApp.On("GetEventByID", mock.Anything, mock.Anything).Return(expectedEvent, nil)

	response, err := server.CreateEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, req.Title, response.Event.Title)

	mockApp.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestCreateEventValidation(t *testing.T) {
	mockApp := new(MockApplication)
	mockLogger := new(MockLogger)

	server := &Server{
		logger:      mockLogger,
		application: mockApp,
	}

	// Мокаем вызовы логгера
	mockLogger.On("Info", "gRPC CreateEvent called").Return()

	// Тест без заголовка
	req := &api.CreateEventRequest{
		StartTime: timestamppb.New(time.Now()),
		Duration:  durationpb.New(time.Hour),
		UserId:    "user123",
	}

	response, err := server.CreateEvent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, response)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	mockLogger.AssertExpectations(t)
}

func TestGetEvent(t *testing.T) {
	mockApp := new(MockApplication)
	mockLogger := new(MockLogger)

	server := &Server{
		logger:      mockLogger,
		application: mockApp,
	}

	// Мокаем вызовы логгера
	mockLogger.On("Info", "gRPC GetEvent called").Return()

	expectedEvent := storage.Event{
		ID:           "test-id",
		Title:        "Test Event",
		StartTime:    time.Now(),
		Duration:     calendar_types.CalendarDuration(time.Hour),
		Description:  "Test Description",
		UserID:       "user123",
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	mockApp.On("GetEventByID", mock.Anything, "test-id").Return(expectedEvent, nil)

	req := &api.GetEventRequest{
		Id: "test-id",
	}

	response, err := server.GetEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedEvent.ID, response.Event.Id)

	mockApp.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestDeleteEvent(t *testing.T) {
	mockApp := new(MockApplication)
	mockLogger := new(MockLogger)

	server := &Server{
		logger:      mockLogger,
		application: mockApp,
	}

	// Мокаем вызовы логгера
	mockLogger.On("Info", "gRPC DeleteEvent called").Return()

	mockApp.On("DeleteEvent", mock.Anything, "test-id").Return(nil)

	req := &api.DeleteEventRequest{
		Id: "test-id",
	}

	response, err := server.DeleteEvent(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.True(t, response.Success)

	mockApp.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestListEventsForDay(t *testing.T) {
	mockApp := new(MockApplication)
	mockLogger := new(MockLogger)

	server := &Server{
		logger:      mockLogger,
		application: mockApp,
	}

	// Мокаем вызовы логгера
	mockLogger.On("Info", "gRPC ListEvents called").Return()

	expectedEvents := []storage.Event{
		{
			ID:           "event1",
			Title:        "Event 1",
			StartTime:    time.Now(),
			Duration:     calendar_types.CalendarDuration(time.Hour),
			Description:  "Description 1",
			UserID:       "user123",
			NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
		},
		{
			ID:           "event2",
			Title:        "Event 2",
			StartTime:    time.Now().Add(time.Hour),
			Duration:     calendar_types.CalendarDuration(2 * time.Hour),
			Description:  "Description 2",
			UserID:       "user123",
			NotifyBefore: calendar_types.CalendarDuration(1 * time.Hour),
		},
	}

	date := time.Now()
	mockApp.On("ListEventsForDay", mock.Anything, mock.Anything).Return(expectedEvents, nil)

	req := &api.ListEventsRequest{
		Date:   timestamppb.New(date),
		Period: "day",
	}

	response, err := server.ListEvents(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Len(t, response.Events, 2)

	mockApp.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
}

func TestGRPCCreateAndGetEvent_Integration(t *testing.T) {
	// 1. Настроить in-memory storage и gRPC сервер
	logger := logger2.New("debug")
	storage := memorystorage.New(logger)
	app := app2.New(logger, storage)
	grpcSrv := NewServer(logger, "localhost", 0, app)
	lis, _ := net.Listen("tcp", "localhost:0")
	go grpcSrv.server.Serve(lis)
	defer grpcSrv.server.Stop()

	// 2. Создать gRPC клиент
	conn, _ := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	client := api.NewEventServiceClient(conn)

	// 3. Вызвать CreateEvent
	startTime := time.Now().Truncate(time.Second)
	resp, err := client.CreateEvent(context.Background(), &api.CreateEventRequest{
		Title:        "Integration Event",
		StartTime:    timestamppb.New(startTime),
		Duration:     durationpb.New(time.Hour),
		Description:  "Integration Test",
		UserId:       "user42",
		NotifyBefore: durationpb.New(30 * time.Minute),
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp.Event)
	assert.Equal(t, "Integration Event", resp.Event.Title)

	// 4. Вызвать GetEvent и проверить результат
	getResp, err := client.GetEvent(context.Background(), &api.GetEventRequest{Id: resp.Event.Id})
	assert.NoError(t, err)
	assert.Equal(t, resp.Event.Id, getResp.Event.Id)
	assert.Equal(t, "Integration Event", getResp.Event.Title)
}
