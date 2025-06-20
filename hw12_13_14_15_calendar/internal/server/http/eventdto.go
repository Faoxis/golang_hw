package internalhttp

import (
	"encoding/json"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"io"
	"net/http"
	"time"
)

type EventRequest struct {
	Title        string                          `json:"title"`
	StartTime    time.Time                       `json:"start_time"`
	Duration     calendar_types.CalendarDuration `json:"duration"`
	Description  string                          `json:"description"`
	UserID       string                          `json:"user_id"`
	NotifyBefore calendar_types.CalendarDuration `json:"notify_before"`
}

type EventResponse struct {
	ID           string                          // Уникальный идентификатор (например, UUID)
	Title        string                          // Название события
	StartTime    time.Time                       // Время начала события
	Duration     calendar_types.CalendarDuration // Длительность события
	Description  string                          // Подробное описание (опционально)
	UserID       string                          // Идентификатор пользователя
	NotifyBefore calendar_types.CalendarDuration // За сколько заранее отправить уведомление (опционально)
}

func mapStorageEventToEventResponse(storageEvent storage.Event) EventResponse {
	return EventResponse{
		ID:           storageEvent.ID,
		Title:        storageEvent.Title,
		StartTime:    storageEvent.StartTime,
		Duration:     storageEvent.Duration,
		Description:  storageEvent.Description,
		UserID:       storageEvent.UserID,
		NotifyBefore: storageEvent.NotifyBefore,
	}
}

type EventListResponse struct {
	events []EventResponse
}

func fromJson[T any](reader io.ReadCloser) (T, error) {
	decoder := json.NewDecoder(reader)
	var result T
	err := decoder.Decode(&result)
	if err != nil {
		var zero T
		return zero, err
	}
	return result, nil
}

func sendInResponse(writer http.ResponseWriter, data interface{}, status int) error {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(data); err != nil {
		return err
	}
	return nil
}
