package server

import (
	"context"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"time"
)

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

type CalculatorServer interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
}
