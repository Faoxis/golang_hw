package scheduler

import (
	"context"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type NotificationStorage interface {
	GetEventsForNotification(ctx context.Context, now time.Time) ([]storage.Event, error)
	MarkEventNotified(ctx context.Context, eventID string) error
	CleanOldEvents(ctx context.Context) error
	Close() error
}
