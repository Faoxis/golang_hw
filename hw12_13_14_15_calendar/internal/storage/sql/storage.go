package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	db     *sql.DB
	logger app.Logger
}

func (storage *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	query := `
		INSERT INTO events (id, title, description, start_time, duration, user_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := storage.db.ExecContext(ctx, query, event.ID, event.Title, event.Description, event.StartTime, event.Duration, event.UserID, event.NotifyBefore)
	return err
}

func (storage *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `
		UPDATE events
		SET title = $2, description = $3, start_time = $4, duration = $5, user_id = $6, notify_before = $7
		WHERE id = $1
	`
	res, err := storage.db.ExecContext(
		ctx,
		query,
		event.ID,
		event.Title,
		event.Description,
		event.StartTime,
		event.Duration,
		event.UserID,
		event.NotifyBefore,
	)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (storage *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`
	res, err := storage.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (strg *Storage) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
	query := `
		SELECT id, title, description, start_time, duration, user_id, notify_before
		FROM events
		WHERE id = $1
	`
	var e storage.Event
	err := strg.db.QueryRowContext(ctx, query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.StartTime, &e.Duration, &e.UserID, &e.NotifyBefore,
	)
	if err != nil {
		return storage.Event{}, errors.New(fmt.Sprintf("Failed to get event from db by id %s: %s\n", id, err))
	}
	return e, nil
}

func (storage *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)
	return storage.listEvents(ctx, start, end)
}

func (storage *Storage) ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 7)
	return storage.listEvents(ctx, start, end)
}

func (storage *Storage) ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	return storage.listEvents(ctx, start, end)
}

func (strg *Storage) listEvents(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	query := `SELECT * FROM events WHERE start_time >= $1 AND start_time < $2;`
	rows, err := strg.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.Duration, &e.UserID, &e.NotifyBefore); err != nil {
			strg.logger.Error(fmt.Sprintf("Failed to scan event from db: %s\n", err))
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func New(dsn string, logger app.Logger) (app.Storage, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping db: %w", err)
	}
	return &Storage{
		db:     db,
		logger: logger,
	}, nil
}

func (storage *Storage) Close() error {
	return storage.db.Close()
}
