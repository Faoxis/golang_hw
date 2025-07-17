package internalgrpc

import (
	"context"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/api"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setupTestGRPCServer(t *testing.T) (*CalendarGRPCServer, *app.App) {
	logg := logger.New("debug")
	storage := memorystorage.New(logg)
	calendar := app.New(logg, storage)

	server := NewCalendarGRPCServer(logg, ":0", calendar)
	return server, calendar
}

func TestCreateEvent(t *testing.T) {
	server, _ := setupTestGRPCServer(t)

	req := &api.CreateEventRequest{
		Title:        "Test Event",
		Description:  "Test Description",
		UserId:       "user123",
		StartTime:    timestamppb.New(time.Now().Add(time.Hour)),
		Duration:     durationpb.New(time.Hour),
		NotifyBefore: durationpb.New(15 * time.Minute),
	}

	resp, err := server.CreateEvent(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, req.Description, resp.Description)
	assert.Equal(t, req.UserId, resp.UserId)
}

func TestGetEvent(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие
	eventID := "test-event-123"
	err := calendar.CreateEvent(
		context.Background(),
		eventID, "Test Event", "Test Description", "user123",
		time.Now().Add(time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Получаем событие
	req := &api.GetEventRequest{Id: eventID}
	resp, err := server.GetEvent(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, eventID, resp.Id)
	assert.Equal(t, "Test Event", resp.Title)
}

func TestUpdateEvent(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие
	eventID := "test-event-456"
	err := calendar.CreateEvent(
		context.Background(),
		eventID, "Original Title", "Original Description", "user123",
		time.Now().Add(time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Обновляем событие
	req := &api.UpdateEventRequest{
		Id:           eventID,
		Title:        "Updated Title",
		Description:  "Updated Description",
		UserId:       "user123",
		StartTime:    timestamppb.New(time.Now().Add(2 * time.Hour)),
		Duration:     durationpb.New(2 * time.Hour),
		NotifyBefore: durationpb.New(30 * time.Minute),
	}

	resp, err := server.UpdateEvent(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, "Updated Title", resp.Title)
	assert.Equal(t, "Updated Description", resp.Description)
}

func TestDeleteEvent(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие
	eventID := "test-event-789"
	err := calendar.CreateEvent(
		context.Background(),
		eventID, "Test Event", "Test Description", "user123",
		time.Now().Add(time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Удаляем событие
	req := &api.DeleteEventRequest{Id: eventID}
	resp, err := server.DeleteEvent(context.Background(), req)
	require.NoError(t, err)
	assert.True(t, resp.Success)

	// Проверяем что событие удалено
	getReq := &api.GetEventRequest{Id: eventID}
	_, err = server.GetEvent(context.Background(), getReq)
	assert.Error(t, err) // Должна быть ошибка, так как событие удалено
}

func TestListEventsForDay(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие на сегодня
	today := time.Now().Truncate(24 * time.Hour)
	err := calendar.CreateEvent(
		context.Background(),
		"event1", "Today Event", "Description", "user123",
		today.Add(10*time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Получаем события за день
	req := &api.ListEventsRequest{
		Date:   timestamppb.New(today),
		Period: "day",
	}
	resp, err := server.ListEvents(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "Today Event", resp.Events[0].Title)
}

func TestListEventsForWeek(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие на этой неделе
	today := time.Now().Truncate(24 * time.Hour)
	err := calendar.CreateEvent(
		context.Background(),
		"event1", "Week Event", "Description", "user123",
		today.Add(10*time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Получаем события за неделю
	req := &api.ListEventsRequest{
		Date:   timestamppb.New(today),
		Period: "week",
	}
	resp, err := server.ListEvents(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "Week Event", resp.Events[0].Title)
}

func TestListEventsForMonth(t *testing.T) {
	server, calendar := setupTestGRPCServer(t)

	// Создаем событие в этом месяце
	today := time.Now().Truncate(24 * time.Hour)
	err := calendar.CreateEvent(
		context.Background(),
		"event1", "Month Event", "Description", "user123",
		today.Add(10*time.Hour),
		calendar_types.CalendarDuration(time.Hour),
		calendar_types.CalendarDuration(15*time.Minute),
	)
	require.NoError(t, err)

	// Получаем события за месяц
	req := &api.ListEventsRequest{
		Date:   timestamppb.New(today),
		Period: "month",
	}
	resp, err := server.ListEvents(context.Background(), req)
	require.NoError(t, err)

	assert.Len(t, resp.Events, 1)
	assert.Equal(t, "Month Event", resp.Events[0].Title)
}

func TestInvalidEventID(t *testing.T) {
	server, _ := setupTestGRPCServer(t)

	// Тестируем несуществующий ID
	req := &api.GetEventRequest{Id: "non-existent-id"}
	_, err := server.GetEvent(context.Background(), req)
	assert.Error(t, err)
}
