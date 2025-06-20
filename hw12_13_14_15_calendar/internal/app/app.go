package app

import (
	"context"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Warn(msg string)
}

type Storage interface {
	AddEvent(ctx context.Context, e storage.Event) error
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (storage.Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
	Close() error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	id, title, description, userID string,
	startTime time.Time,
	duration, notifyBefore calendar_types.CalendarDuration,
) error {
	event := storage.Event{
		ID:           id,
		Title:        title,
		Description:  description,
		StartTime:    startTime,
		Duration:     duration,
		UserID:       userID,
		NotifyBefore: notifyBefore,
	}
	return a.storage.AddEvent(ctx, event)
}

func (a *App) UpdateEvent(
	ctx context.Context,
	id, title, description, userID string,
	startTime time.Time,
	duration, notifyBefore calendar_types.CalendarDuration,
) error {
	event := storage.Event{
		ID:           id,
		Title:        title,
		Description:  description,
		StartTime:    startTime,
		Duration:     duration,
		UserID:       userID,
		NotifyBefore: notifyBefore,
	}
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func (a *App) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
	return a.storage.GetEventByID(ctx, id)
}

func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForDay(ctx, date)
}

func (a *App) ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForWeek(ctx, date)
}

func (a *App) ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	return a.storage.ListEventsForMonth(ctx, date)
}
