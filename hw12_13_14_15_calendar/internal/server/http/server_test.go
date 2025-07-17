package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestServer(t *testing.T) (*httptest.Server, *app.App) {
	logg := logger.New("debug")
	storage := memorystorage.New(logg)
	calendar := app.New(logg, storage)

	// Используем порт 0 для автоматического выбора свободного порта
	server := NewServer(logg, "localhost", 0, calendar)

	ts := httptest.NewServer(server.(*HttpServer).server.Handler)
	return ts, calendar
}

func TestCreateEvent(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	event := EventRequest{
		Title:        "Test Event",
		Description:  "Test Description",
		UserID:       "user123",
		StartTime:    time.Now().Add(time.Hour),
		Duration:     calendar_types.CalendarDuration(time.Hour),
		NotifyBefore: calendar_types.CalendarDuration(15 * time.Minute),
	}

	body, _ := json.Marshal(event)
	resp, err := http.Post(ts.URL+"/events/", "application/json", bytes.NewBuffer(body))

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response EventResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, event.Title, response.Title)
	assert.Equal(t, event.Description, response.Description)
	assert.Equal(t, event.UserID, response.UserID)
}

func TestGetEvent(t *testing.T) {
	ts, calendar := setupTestServer(t)
	defer ts.Close()

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
	resp, err := http.Get(fmt.Sprintf("%s/events/%s", ts.URL, eventID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response EventResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, eventID, response.ID)
	assert.Equal(t, "Test Event", response.Title)
}

func TestUpdateEvent(t *testing.T) {
	ts, calendar := setupTestServer(t)
	defer ts.Close()

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
	updateEvent := EventRequest{
		Title:        "Updated Title",
		Description:  "Updated Description",
		UserID:       "user123",
		StartTime:    time.Now().Add(2 * time.Hour),
		Duration:     calendar_types.CalendarDuration(2 * time.Hour),
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	body, _ := json.Marshal(updateEvent)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/events/%s", ts.URL, eventID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response EventResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "Updated Title", response.Title)
	assert.Equal(t, "Updated Description", response.Description)
}

func TestDeleteEvent(t *testing.T) {
	ts, calendar := setupTestServer(t)
	defer ts.Close()

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
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/events/%s", ts.URL, eventID), nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем что событие удалено
	resp, err = http.Get(fmt.Sprintf("%s/events/%s", ts.URL, eventID))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestListEventsForDay(t *testing.T) {
	ts, calendar := setupTestServer(t)
	defer ts.Close()

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
	dateStr := today.Format("2006-01-02")
	resp, err := http.Get(fmt.Sprintf("%s/events/?day=%s", ts.URL, dateStr))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []EventResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Len(t, response, 1)
	assert.Equal(t, "Today Event", response[0].Title)
}

func TestInvalidDateFormat(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	// Тестируем неверный формат даты
	resp, err := http.Get(ts.URL + "/events/?day=invalid-date")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestInvalidEventID(t *testing.T) {
	ts, _ := setupTestServer(t)
	defer ts.Close()

	// Тестируем несуществующий ID
	resp, err := http.Get(ts.URL + "/events/non-existent-id")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
