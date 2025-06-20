package storage

import (
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"time"
)

type Event struct {
	ID           string                          // Уникальный идентификатор (например, UUID)
	Title        string                          // Название события
	StartTime    time.Time                       // Время начала события
	Duration     calendar_types.CalendarDuration // Длительность события
	Description  string                          // Подробное описание (опционально)
	UserID       string                          // Идентификатор пользователя
	NotifyBefore calendar_types.CalendarDuration // За сколько заранее отправить уведомление (опционально)
}
