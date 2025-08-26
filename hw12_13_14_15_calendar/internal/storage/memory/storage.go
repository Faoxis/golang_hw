package memorystorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     *sync.RWMutex //nolint:unused
	logger app.Logger
}

func (strg *Storage) AddEvent(ctx context.Context, e storage.Event) error {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	if _, ok := strg.events[e.ID]; ok {
		return errors.New("event already exists")
	}
	strg.events[e.ID] = e
	return nil
}

func (strg *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	strg.mu.Lock()
	defer strg.mu.Unlock()

	if _, ok := strg.events[e.ID]; !ok {
		return errors.New("event does not exist")
	}
	strg.events[e.ID] = e
	return nil
}

func (strg *Storage) DeleteEvent(ctx context.Context, id string) error {
	strg.mu.Lock()
	defer strg.mu.Unlock()
	delete(strg.events, id)
	return nil
}

func (strg *Storage) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
	strg.mu.RLock()
	defer strg.mu.RUnlock()
	event, ok := strg.events[id]
	if !ok {
		return storage.Event{}, errors.New("event does not exist")
	}
	return event, nil
}

func sameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (strg *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	strg.mu.RLock()
	defer strg.mu.RUnlock()
	foundEvents := []storage.Event{}
	for _, e := range strg.events {
		if sameDay(e.StartTime, date) {
			foundEvents = append(foundEvents, e)
		}
	}
	return foundEvents, nil
}

func sameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}
func (strg *Storage) ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error) {
	strg.mu.RLock()
	defer strg.mu.RUnlock()
	foundEvents := []storage.Event{}
	for _, e := range strg.events {
		if sameWeek(e.StartTime, date) {
			foundEvents = append(foundEvents, e)
		}
	}
	return foundEvents, nil
}

func sameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}
func (strg Storage) ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error) {
	strg.mu.RLock()
	defer strg.mu.RUnlock()
	foundEvents := []storage.Event{}
	for _, e := range strg.events {
		if sameMonth(e.StartTime, date) {
			foundEvents = append(foundEvents, e)
		}
	}
	return foundEvents, nil
}

func New(logger app.Logger) app.Storage {
	return &Storage{
		events: map[string]storage.Event{},
		mu:     &sync.RWMutex{},
		logger: logger,
	}
}

func (storage *Storage) Close() error {
	return nil
}
