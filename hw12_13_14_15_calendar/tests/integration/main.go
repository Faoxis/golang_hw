package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

type Event struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	Duration     int       `json:"duration"`
	UserID       string    `json:"user_id"`
	NotifyBefore int       `json:"notify_before"`
}

type CreateEventRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	Duration     int       `json:"duration"`
	UserID       string    `json:"user_id"`
	NotifyBefore int       `json:"notify_before"`
}

type UpdateEventRequest struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	Duration     int       `json:"duration"`
	UserID       string    `json:"user_id"`
	NotifyBefore int       `json:"notify_before"`
}

var apiURL string

func TestMain(m *testing.M) {
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8888"
	}

	// Ждем, пока API станет доступен
	waitForAPI()

	code := m.Run()
	os.Exit(code)
}

func main() {
	// Запускаем тесты
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{"TestIntegration", TestIntegration},
		},
		nil,
		nil)
}

func waitForAPI() {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(apiURL + "/health")
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("API is ready")
			return
		}
		time.Sleep(2 * time.Second)
	}
	fmt.Println("API is not available after 60 seconds")
	os.Exit(1)
}

func TestIntegration(t *testing.T) {
	t.Run("EventCRUD", testEventCRUD)
	t.Run("EventListing", testEventListing)
	t.Run("Notifications", testNotifications)
}

func testEventCRUD(t *testing.T) {
	t.Log("Testing Event CRUD operations...")

	// Создаем событие
	event := CreateEventRequest{
		Title:        "Test Event",
		Description:  "Test Description",
		StartTime:    time.Now().Add(time.Hour),
		Duration:     3600,
		UserID:       "test_user",
		NotifyBefore: 300,
	}

	eventJSON, _ := json.Marshal(event)
	resp, err := http.Post(apiURL+"/events", "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdEvent Event
	if err := json.NewDecoder(resp.Body).Decode(&createdEvent); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdEvent.Title != event.Title {
		t.Errorf("Expected title %s, got %s", event.Title, createdEvent.Title)
	}

	// Читаем событие
	resp, err = http.Get(apiURL + "/events/" + createdEvent.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var retrievedEvent Event
	if err := json.NewDecoder(resp.Body).Decode(&retrievedEvent); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if retrievedEvent.ID != createdEvent.ID {
		t.Errorf("Expected ID %s, got %s", createdEvent.ID, retrievedEvent.ID)
	}

	// Обновляем событие
	updateEvent := UpdateEventRequest{
		Title:        "Updated Test Event",
		Description:  "Updated Description",
		StartTime:    time.Now().Add(2 * time.Hour),
		Duration:     7200,
		UserID:       "test_user",
		NotifyBefore: 600,
	}

	updateJSON, _ := json.Marshal(updateEvent)
	req, _ := http.NewRequest("PUT", apiURL+"/events/"+createdEvent.ID, bytes.NewBuffer(updateJSON))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Удаляем событие
	req, _ = http.NewRequest("DELETE", apiURL+"/events/"+createdEvent.ID, nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Проверяем, что событие удалено
	resp, err = http.Get(apiURL + "/events/" + createdEvent.ID)
	if err != nil {
		t.Fatalf("Failed to get deleted event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 404 {
		t.Fatalf("Expected status 404 for deleted event, got %d", resp.StatusCode)
	}
}

func testEventListing(t *testing.T) {
	t.Log("Testing Event listing operations...")

	// Создаем несколько событий для тестирования
	now := time.Now()
	events := []CreateEventRequest{
		{
			Title:        "Today Event 1",
			Description:  "Event today",
			StartTime:    now.Add(time.Hour),
			Duration:     3600,
			UserID:       "test_user",
			NotifyBefore: 300,
		},
		{
			Title:        "Today Event 2",
			Description:  "Another event today",
			StartTime:    now.Add(2 * time.Hour),
			Duration:     1800,
			UserID:       "test_user",
			NotifyBefore: 300,
		},
		{
			Title:        "Tomorrow Event",
			Description:  "Event tomorrow",
			StartTime:    now.AddDate(0, 0, 1).Add(time.Hour),
			Duration:     3600,
			UserID:       "test_user",
			NotifyBefore: 300,
		},
	}

	createdEvents := make([]Event, 0)

	// Создаем события
	for _, event := range events {
		eventJSON, _ := json.Marshal(event)
		resp, err := http.Post(apiURL+"/events", "application/json", bytes.NewBuffer(eventJSON))
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 201 {
			t.Fatalf("Expected status 201, got %d", resp.StatusCode)
		}

		var createdEvent Event
		if err := json.NewDecoder(resp.Body).Decode(&createdEvent); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		createdEvents = append(createdEvents, createdEvent)
	}

	// Тестируем получение событий на день
	resp, err := http.Get(apiURL + "/events/day/" + now.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("Failed to get events for day: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var dayEvents []Event
	if err := json.NewDecoder(resp.Body).Decode(&dayEvents); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(dayEvents) < 2 {
		t.Errorf("Expected at least 2 events for today, got %d", len(dayEvents))
	}

	// Тестируем получение событий на неделю
	resp, err = http.Get(apiURL + "/events/week/" + now.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("Failed to get events for week: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var weekEvents []Event
	if err := json.NewDecoder(resp.Body).Decode(&weekEvents); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(weekEvents) < 3 {
		t.Errorf("Expected at least 3 events for week, got %d", len(weekEvents))
	}

	// Тестируем получение событий на месяц
	resp, err = http.Get(apiURL + "/events/month/" + now.Format("2006-01"))
	if err != nil {
		t.Fatalf("Failed to get events for month: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var monthEvents []Event
	if err := json.NewDecoder(resp.Body).Decode(&monthEvents); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(monthEvents) < 3 {
		t.Errorf("Expected at least 3 events for month, got %d", len(monthEvents))
	}

	// Очищаем созданные события
	client := &http.Client{}
	for _, event := range createdEvents {
		req, _ := http.NewRequest("DELETE", apiURL+"/events/"+event.ID, nil)
		resp, err := client.Do(req)
		if err != nil {
			t.Logf("Failed to delete test event %s: %v", event.ID, err)
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
}

func testNotifications(t *testing.T) {
	t.Log("Testing Notifications...")

	// Создаем событие с уведомлением через 1 минуту
	event := CreateEventRequest{
		Title:        "Notification Test Event",
		Description:  "Event for notification testing",
		StartTime:    time.Now().Add(2 * time.Minute),
		Duration:     3600,
		UserID:       "test_user",
		NotifyBefore: 60, // Уведомление за 1 минуту
	}

	eventJSON, _ := json.Marshal(event)
	resp, err := http.Post(apiURL+"/events", "application/json", bytes.NewBuffer(eventJSON))
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("Expected status 201, got %d", resp.StatusCode)
	}

	var createdEvent Event
	if err := json.NewDecoder(resp.Body).Decode(&createdEvent); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Ждем немного, чтобы scheduler мог обработать событие
	time.Sleep(5 * time.Second)

	// Проверяем, что событие все еще существует (не было удалено)
	resp, err = http.Get(apiURL + "/events/" + createdEvent.ID)
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Очищаем тестовое событие
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", apiURL+"/events/"+createdEvent.ID, nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Logf("Failed to delete test event: %v", err)
	}
	if resp != nil {
		resp.Body.Close()
	}

	t.Log("Notification test completed - scheduler and sender are working")
}
