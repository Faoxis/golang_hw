package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/calendar_types"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/queue"
	"github.com/Faoxis/golang_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/assert"
)

type MockNotificationStorage struct {
	events []storage.Event
	err    error
}

func (m *MockNotificationStorage) GetEventsForNotification(ctx context.Context, now time.Time) ([]storage.Event, error) {
	return m.events, m.err
}

func (m *MockNotificationStorage) MarkEventNotified(ctx context.Context, eventID string) error {
	return nil
}

func (m *MockNotificationStorage) CleanOldEvents(ctx context.Context) error {
	return nil
}

func (m *MockNotificationStorage) Close() error {
	return nil
}

type MockQueue struct {
	messages []storage.Notification
	err      error
}

func (m *MockQueue) Put(queue string, exchange string, message queue.MessageQueue[storage.Notification]) error {
	m.messages = append(m.messages, message.Body)
	return m.err
}

func (m *MockQueue) Get(ctx context.Context, queueName string) (<-chan queue.MessageQueue[storage.Notification], <-chan error) {
	msgChan := make(chan queue.MessageQueue[storage.Notification])
	errChan := make(chan error)

	go func() {
		defer close(msgChan)
		defer close(errChan)

		if m.err != nil {
			errChan <- m.err
			return
		}

		for _, msg := range m.messages {
			select {
			case <-ctx.Done():
				return
			case msgChan <- queue.MessageQueue[storage.Notification]{ID: "test", Body: msg}:
			}
		}
	}()

	return msgChan, errChan
}

func (m *MockQueue) Close() error {
	return nil
}

type MockLogger struct {
	messages []string
}

func (m *MockLogger) Debug(msg string) {
	m.messages = append(m.messages, "DEBUG: "+msg)
}

func (m *MockLogger) Info(msg string) {
	m.messages = append(m.messages, "INFO: "+msg)
}

func (m *MockLogger) Error(msg string) {
	m.messages = append(m.messages, "ERROR: "+msg)
}

func (m *MockLogger) Warn(msg string) {
	m.messages = append(m.messages, "WARN: "+msg)
}

func TestScheduler_CheckAndSendNotifications(t *testing.T) {
	now := time.Now()
	event := storage.Event{
		ID:           "test-event-2",
		Title:        "Test Event Now",
		StartTime:    now.Add(30 * time.Minute),
		Duration:     calendar_types.CalendarDuration(time.Hour),
		UserID:       "user1",
		NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
	}

	mockStorage := &MockNotificationStorage{
		events: []storage.Event{event},
		err:    nil,
	}

	mockQueue := &MockQueue{
		messages: []storage.Notification{},
		err:      nil,
	}

	mockLogger := &MockLogger{
		messages: []string{},
	}

	scheduler := NewScheduler(mockLogger, mockStorage, mockQueue, "test-queue", "test-exchange", "1m")

	scheduler.checkAndSendNotifications(context.Background())

	if len(mockQueue.messages) == 0 {
		t.Error("Expected notification to be sent")
	}

	notification := mockQueue.messages[0]
	if notification.EventID != event.ID {
		t.Errorf("Expected EventID %s, got %s", event.ID, notification.EventID)
	}
	if notification.Title != event.Title {
		t.Errorf("Expected Title %s, got %s", event.Title, notification.Title)
	}
	if notification.UserID != event.UserID {
		t.Errorf("Expected UserID %s, got %s", event.UserID, notification.UserID)
	}
}

func TestScheduler_ErrorHandling(t *testing.T) {
	mockStorage := &MockNotificationStorage{
		events: nil,
		err:    assert.AnError,
	}

	mockQueue := &MockQueue{
		messages: []storage.Notification{},
		err:      nil,
	}

	mockLogger := &MockLogger{
		messages: []string{},
	}

	scheduler := NewScheduler(mockLogger, mockStorage, mockQueue, "test-queue", "test-exchange", "1m")

	scheduler.checkAndSendNotifications(context.Background())

	foundErrorLog := false
	for _, msg := range mockLogger.messages {
		if msg == "ERROR: Failed to get events for notification: assert.AnError general error for testing" {
			foundErrorLog = true
			break
		}
	}
	if !foundErrorLog {
		t.Error("Expected error log message")
	}

	if len(mockQueue.messages) > 0 {
		t.Error("Expected no messages to be sent when there's an error")
	}
}

func TestScheduler_WithRealisticScenario(t *testing.T) {
	now := time.Now()
	events := []storage.Event{
		{
			ID:           "event-1",
			Title:        "Event in 30 minutes",
			StartTime:    now.Add(30 * time.Minute),
			Duration:     calendar_types.CalendarDuration(time.Hour),
			UserID:       "user1",
			NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
		},
		{
			ID:           "event-2",
			Title:        "Event in 2 hours",
			StartTime:    now.Add(2 * time.Hour),
			Duration:     calendar_types.CalendarDuration(time.Hour),
			UserID:       "user2",
			NotifyBefore: calendar_types.CalendarDuration(1 * time.Hour),
		},
		{
			ID:           "event-3",
			Title:        "Event already started",
			StartTime:    now.Add(-30 * time.Minute),
			Duration:     calendar_types.CalendarDuration(time.Hour),
			UserID:       "user3",
			NotifyBefore: calendar_types.CalendarDuration(30 * time.Minute),
		},
	}

	mockStorage := &MockNotificationStorage{
		events: events,
		err:    nil,
	}

	mockQueue := &MockQueue{
		messages: []storage.Notification{},
		err:      nil,
	}

	mockLogger := &MockLogger{
		messages: []string{},
	}

	scheduler := NewScheduler(mockLogger, mockStorage, mockQueue, "test-queue", "test-exchange", "1m")

	scheduler.checkAndSendNotifications(context.Background())

	if len(mockQueue.messages) == 0 {
		t.Error("Expected notifications to be sent")
	}

	expectedEventIDs := map[string]bool{"event-1": true, "event-2": true, "event-3": true}
	for _, notification := range mockQueue.messages {
		if !expectedEventIDs[notification.EventID] {
			t.Errorf("Unexpected event ID: %s", notification.EventID)
		}
	}
}
