package scheduler

import (
	"context"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type Scheduler struct {
	storage       NotificationStorage
	logger        app.Logger
	queue         queue.Queue[storage.Notification]
	queueName     string
	exchangeName  string
	checkInterval string
}

func NewScheduler(logger app.Logger, storage NotificationStorage, queue queue.Queue[storage.Notification], queueName, exchangeName, checkInterval string) *Scheduler {
	return &Scheduler{
		storage:       storage,
		logger:        logger,
		queue:         queue,
		queueName:     queueName,
		exchangeName:  exchangeName,
		checkInterval: checkInterval,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	interval, err := time.ParseDuration(s.checkInterval)
	if err != nil {
		s.logger.Error("Failed to parse check interval, using default 1m: " + err.Error())
		interval = 1 * time.Minute
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Очистка старых событий каждые 24 часа
	cleanupTicker := time.NewTicker(24 * time.Hour)
	defer cleanupTicker.Stop()

	s.logger.Info("Scheduler started with interval: " + interval.String())

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Scheduler stopped")
			return
		case <-ticker.C:
			s.checkAndSendNotifications(ctx)
		case <-cleanupTicker.C:
			s.cleanOldEvents(ctx)
		}
	}
}

func (s *Scheduler) checkAndSendNotifications(ctx context.Context) {
	events, err := s.storage.GetEventsForNotification(ctx, time.Now())
	if err != nil {
		s.logger.Error("Failed to get events for notification: " + err.Error())
		return
	}

	for _, event := range events {
		notification := storage.Notification{
			EventID:   event.ID,
			Title:     event.Title,
			EventTime: event.StartTime,
			UserID:    event.UserID,
		}

		msg := queue.MessageQueue[storage.Notification]{
			ID:   event.ID,
			Body: notification,
		}

		err := s.queue.Put(s.queueName, s.exchangeName, msg)
		if err != nil {
			s.logger.Error("Failed to send notification to queue: " + err.Error())
			continue
		}

		if err := s.storage.MarkEventNotified(ctx, event.ID); err != nil {
			s.logger.Error("Failed to mark event as notified: " + err.Error())
		}

		s.logger.Info("Notification sent to queue for event: " + event.ID)
	}
}

func (s *Scheduler) cleanOldEvents(ctx context.Context) {
	if err := s.storage.CleanOldEvents(ctx); err != nil {
		s.logger.Error("Failed to clean old events: " + err.Error())
	}
}
