package memorystorage

import (
	"context"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/logger"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

func TestAddAndGetEvent(t *testing.T) {
	strg := New(logger.New("debug"))
	ctx := context.Background()

	event := storage.Event{
		ID:        uuid.New().String(),
		Title:     "Meeting",
		StartTime: time.Now(),
		Duration:  calendar_types.CalendarDuration(time.Hour),
		UserID:    "user1",
	}

	err := strg.AddEvent(ctx, event)
	if err != nil {
		t.Fatalf("failed to add event: %v", err)
	}

	got, err := strg.GetEventByID(ctx, event.ID)
	if err != nil {
		t.Fatalf("failed to get event: %v", err)
	}

	if got.Title != event.Title {
		t.Errorf("expected title %q, got %q", event.Title, got.Title)
	}
}

func TestUpdateEvent(t *testing.T) {
	s := New(logger.New("debug"))
	ctx := context.Background()

	event := storage.Event{
		ID:        uuid.New().String(),
		Title:     "Original",
		StartTime: time.Now(),
		Duration:  calendar_types.CalendarDuration(time.Hour),
		UserID:    "user2",
	}
	_ = s.AddEvent(ctx, event)

	event.Title = "Updated"
	err := s.UpdateEvent(ctx, event)
	if err != nil {
		t.Fatalf("failed to update event: %v", err)
	}

	got, _ := s.GetEventByID(ctx, event.ID)
	if got.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %q", got.Title)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := New(logger.New("debug"))
	ctx := context.Background()

	event := storage.Event{
		ID:        uuid.New().String(),
		Title:     "To Delete",
		StartTime: time.Now(),
		Duration:  calendar_types.CalendarDuration(time.Hour),
		UserID:    "user3",
	}
	_ = s.AddEvent(ctx, event)

	err := s.DeleteEvent(ctx, event.ID)
	if err != nil {
		t.Fatalf("failed to delete event: %v", err)
	}

	_, err = s.GetEventByID(ctx, event.ID)
	if err == nil {
		t.Error("expected error when getting deleted event, got nil")
	}
}

func TestListEventsForDay(t *testing.T) {
	s := New(logger.New("debug"))
	ctx := context.Background()

	now := time.Now()
	eventToday := storage.Event{
		ID:        uuid.New().String(),
		Title:     "Today Event",
		StartTime: now,
		Duration:  calendar_types.CalendarDuration(time.Hour),
		UserID:    "user4",
	}
	eventTomorrow := storage.Event{
		ID:        uuid.New().String(),
		Title:     "Tomorrow Event",
		StartTime: now.Add(24 * time.Hour),
		Duration:  calendar_types.CalendarDuration(time.Hour),
		UserID:    "user4",
	}

	_ = s.AddEvent(ctx, eventToday)
	_ = s.AddEvent(ctx, eventTomorrow)

	list, err := s.ListEventsForDay(ctx, now)
	if err != nil {
		t.Fatalf("error listing events: %v", err)
	}

	if len(list) != 1 {
		t.Errorf("expected 1 event today, got %d", len(list))
	}
}
