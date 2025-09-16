package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type SQLNotificationStorage struct {
	db     *sql.DB
	logger app.Logger
}

func NewSQLNotificationStorage(dsn string, logger app.Logger) (NotificationStorage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}
	return &SQLNotificationStorage{
		db:     db,
		logger: logger,
	}, nil
}

func (ns *SQLNotificationStorage) GetEventsForNotification(ctx context.Context, now time.Time) ([]storage.Event, error) {
	// Ищем события, для которых время уведомления попадает в интервал ±1 минута от текущего времени
	oneMinuteAgo := now.Add(-time.Minute)
	oneMinuteLater := now.Add(time.Minute)

	query := `
		SELECT id, title, description, start_time, duration, user_id, notify_before
		FROM events 
		WHERE (start_time - INTERVAL '1 second' * notify_before) >= $1 
		AND (start_time - INTERVAL '1 second' * notify_before) <= $2 
		AND notify_before > 0
		ORDER BY start_time
	`

	rows, err := ns.db.QueryContext(ctx, query, oneMinuteAgo, oneMinuteLater)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(
			&event.ID, &event.Title, &event.Description,
			&event.StartTime, &event.Duration, &event.UserID, &event.NotifyBefore,
		); err != nil {
			ns.logger.Error(fmt.Sprintf("Failed to scan event: %s", err))
			continue
		}
		events = append(events, event)
	}

	return events, nil
}

func (ns *SQLNotificationStorage) MarkEventNotified(ctx context.Context, eventID string) error {
	// Устанавливаем notify_before в 0, чтобы событие больше не попадало в выборку
	query := `UPDATE events SET notify_before = 0 WHERE id = $1`

	result, err := ns.db.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to mark event as notified: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		ns.logger.Error("Failed to get rows affected: " + err.Error())
	} else if rowsAffected == 0 {
		ns.logger.Error("No event found to mark as notified: " + eventID)
	}

	return nil
}

func (ns *SQLNotificationStorage) CleanOldEvents(ctx context.Context) error {
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	query := `DELETE FROM events WHERE (start_time + INTERVAL '1 second' * duration) < $1`

	result, err := ns.db.ExecContext(ctx, query, oneYearAgo)
	if err != nil {
		return fmt.Errorf("failed to clean old events: %w", err)
	}

	deletedCount, err := result.RowsAffected()
	if err != nil {
		ns.logger.Error("Failed to get deleted count: " + err.Error())
	} else {
		ns.logger.Info(fmt.Sprintf("Cleaned %d old events", deletedCount))
	}

	return nil
}

func (ns *SQLNotificationStorage) Close() error {
	return ns.db.Close()
}
