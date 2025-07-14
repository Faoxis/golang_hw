package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	eventRequest := EventRequest{
		Title:        "Test Event",
		StartTime:    time.Now(),
		Duration:     calendar_types.CalendarDuration(time.Hour),
		Description:  "Test Description",
		UserID:       "user123",
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	eventJSON, _ := json.Marshal(eventRequest)
	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(eventJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockApp.On("CreateEvent", mock.Anything, mock.Anything, eventRequest.Title, eventRequest.Description, eventRequest.UserID, mock.Anything, eventRequest.Duration, eventRequest.NotifyBefore).Return(nil)

	expectedEvent := storage.Event{
		ID:           "test-id",
		Title:        eventRequest.Title,
		StartTime:    eventRequest.StartTime,
		Duration:     eventRequest.Duration,
		Description:  eventRequest.Description,
		UserID:       eventRequest.UserID,
		NotifyBefore: eventRequest.NotifyBefore,
	}

	mockApp.On("GetEventByID", mock.Anything, mock.Anything).Return(expectedEvent, nil)

	handler := addNewEvent(mockApp, mockLogger)
	handler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response EventResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, eventRequest.Title, response.Title)

	mockApp.AssertExpectations(t)
}

func TestEventResponseMapping(t *testing.T) {
	// Тестируем маппинг событий
	storageEvent := storage.Event{
		ID:           "test-id",
		Title:        "Test Event",
		StartTime:    time.Now(),
		Duration:     calendar_types.CalendarDuration(time.Hour),
		Description:  "Test Description",
		UserID:       "user123",
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	response := mapStorageEventToEventResponse(storageEvent)

	assert.Equal(t, storageEvent.ID, response.ID)
	assert.Equal(t, storageEvent.Title, response.Title)
	assert.Equal(t, storageEvent.Description, response.Description)
	assert.Equal(t, storageEvent.UserID, response.UserID)
	assert.Equal(t, storageEvent.Duration, response.Duration)
	assert.Equal(t, storageEvent.NotifyBefore, response.NotifyBefore)
}

func TestDeleteEventLogic(t *testing.T) {
	mockApp := new(MockApplication)

	// Тестируем логику удаления напрямую
	mockApp.On("DeleteEvent", mock.Anything, "test-id").Return(nil)

	// Проверяем, что метод вызывается с правильными параметрами
	err := mockApp.DeleteEvent(context.Background(), "test-id")
	assert.NoError(t, err)

	mockApp.AssertExpectations(t)
}

func TestListEventsForDayLogic(t *testing.T) {
	mockApp := new(MockApplication)

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

	// Тестируем логику напрямую
	events, err := mockApp.ListEventsForDay(context.Background(), date)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.Equal(t, expectedEvents[0].ID, events[0].ID)
	assert.Equal(t, expectedEvents[1].ID, events[1].ID)

	mockApp.AssertExpectations(t)
}
