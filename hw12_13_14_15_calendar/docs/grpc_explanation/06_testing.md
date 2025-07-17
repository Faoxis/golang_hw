# 6. Тестирование

## 🧪 Стратегия тестирования

```
┌─────────────────────────────────────────────────────────┐
│                    Testing Pyramid                      │
│  ┌─────────────────────────────────────────────────────┐ │
│  │              Integration Tests                      │ │
│  │  ┌─────────────────┐  ┌─────────────────┐          │ │
│  │  │   HTTP API      │  │   gRPC API      │          │ │
│  │  │   Tests         │  │   Tests         │          │ │
│  │  └─────────────────┘  └─────────────────┘          │ │
│  └─────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────┐ │
│  │              Unit Tests                             │ │
│  │  ┌─────────────────┐  ┌─────────────────┐          │ │
│  │  │   Application   │  │   Storage       │          │ │
│  │  │   Tests         │  │   Tests         │          │ │
│  │  └─────────────────┘  └─────────────────┘          │ │
│  └─────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────┐ │
│  │              Manual Tests                           │ │
│  │  ┌─────────────────┐  ┌─────────────────┐          │ │
│  │  │   curl          │  │   grpcurl       │          │ │
│  │  │   Commands      │  │   Commands      │          │ │
│  │  └─────────────────┘  └─────────────────┘          │ │
│  └─────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## 🔧 Unit тесты

### 1. **Application тесты**
```go
// internal/app/app_test.go
func TestApplication_CreateEvent(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    storage := &mockStorage{}
    app := New(logger, storage)
    
    ctx := context.Background()
    id := "test-id"
    title := "Test Event"
    description := "Test Description"
    userID := "user123"
    startTime := time.Now()
    duration := calendar_types.CalendarDuration(time.Hour)
    notifyBefore := calendar_types.CalendarDuration(15 * time.Minute)
    
    // Act
    err := app.CreateEvent(ctx, id, title, description, userID, startTime, duration, notifyBefore)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, 1, storage.addEventCalled)
}

func TestApplication_CreateEvent_ValidationError(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    storage := &mockStorage{}
    app := New(logger, storage)
    
    ctx := context.Background()
    
    // Act - пустой заголовок
    err := app.CreateEvent(ctx, "id", "", "desc", "user", time.Now(), time.Hour, 15*time.Minute)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "title is required")
}
```

### 2. **Storage тесты**
```go
// internal/storage/memory/storage_test.go
func TestMemoryStorage_AddEvent(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    storage := New(logger)
    
    event := storage.Event{
        ID:           "test-id",
        Title:        "Test Event",
        Description:  "Test Description",
        UserID:       "user123",
        StartTime:    time.Now(),
        Duration:     calendar_types.CalendarDuration(time.Hour),
        NotifyBefore: calendar_types.CalendarDuration(15 * time.Minute),
    }
    
    // Act
    err := storage.AddEvent(context.Background(), event)
    
    // Assert
    assert.NoError(t, err)
    
    // Проверяем, что событие сохранено
    retrievedEvent, err := storage.GetEventByID(context.Background(), "test-id")
    assert.NoError(t, err)
    assert.Equal(t, event.Title, retrievedEvent.Title)
}

func TestMemoryStorage_AddEvent_DuplicateID(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    storage := New(logger)
    
    event := storage.Event{ID: "test-id", Title: "Test Event"}
    storage.AddEvent(context.Background(), event)
    
    // Act - добавляем событие с тем же ID
    err := storage.AddEvent(context.Background(), event)
    
    // Assert
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "event already exists")
}
```

## 🔗 Integration тесты

### 1. **HTTP API тесты**
```go
// internal/server/http/server_test.go
func TestHTTPServer_CreateEvent(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    application := &mockApplication{}
    server := NewServer(logger, "localhost", 0, application)
    
    // Запускаем сервер на случайном порту
    listener, err := net.Listen("tcp", ":0")
    require.NoError(t, err)
    
    go func() {
        http.Serve(listener, server.Router())
    }()
    defer listener.Close()
    
    // Получаем порт сервера
    port := listener.Addr().(*net.TCPAddr).Port
    
    // Act
    requestBody := `{
        "title": "Test Event",
        "description": "Test Description",
        "user_id": "user123",
        "start_time": "2025-07-15T12:00:00Z",
        "duration": "1h0m0s",
        "notify_before": "15m0s"
    }`
    
    resp, err := http.Post(
        fmt.Sprintf("http://localhost:%d/events", port),
        "application/json",
        strings.NewReader(requestBody),
    )
    
    // Assert
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    // Проверяем ответ
    var response EventResponse
    err = json.NewDecoder(resp.Body).Decode(&response)
    require.NoError(t, err)
    
    assert.Equal(t, "Test Event", response.Title)
    assert.Equal(t, "user123", response.UserID)
}

func TestHTTPServer_CreateEvent_ValidationError(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    application := &mockApplication{}
    server := NewServer(logger, "localhost", 0, application)
    
    listener, err := net.Listen("tcp", ":0")
    require.NoError(t, err)
    
    go func() {
        http.Serve(listener, server.Router())
    }()
    defer listener.Close()
    
    port := listener.Addr().(*net.TCPAddr).Port
    
    // Act - запрос без заголовка
    requestBody := `{
        "description": "Test Description",
        "user_id": "user123",
        "start_time": "2025-07-15T12:00:00Z",
        "duration": "1h0m0s",
        "notify_before": "15m0s"
    }`
    
    resp, err := http.Post(
        fmt.Sprintf("http://localhost:%d/events", port),
        "application/json",
        strings.NewReader(requestBody),
    )
    
    // Assert
    require.NoError(t, err)
    defer resp.Body.Close()
    
    assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    
    body, err := io.ReadAll(resp.Body)
    require.NoError(t, err)
    assert.Contains(t, string(body), "title is required")
}
```

### 2. **gRPC API тесты**
```go
// internal/server/grpc/server_test.go
func TestGRPCServer_CreateEvent(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    application := &mockApplication{}
    server := NewServer(logger, "localhost", 0, application)
    
    // Запускаем gRPC сервер
    listener, err := net.Listen("tcp", ":0")
    require.NoError(t, err)
    
    grpcServer := grpc.NewServer()
    api.RegisterEventServiceServer(grpcServer, server)
    
    go func() {
        grpcServer.Serve(listener)
    }()
    defer grpcServer.Stop()
    
    // Подключаемся к серверу
    conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
    require.NoError(t, err)
    defer conn.Close()
    
    client := api.NewEventServiceClient(conn)
    
    // Act
    startTime := timestamppb.New(time.Date(2025, 7, 15, 12, 0, 0, 0, time.UTC))
    duration := durationpb.New(time.Hour)
    notifyBefore := durationpb.New(15 * time.Minute)
    
    req := &api.CreateEventRequest{
        Title:        "Test Event",
        Description:  "Test Description",
        UserId:       "user123",
        StartTime:    startTime,
        Duration:     duration,
        NotifyBefore: notifyBefore,
    }
    
    resp, err := client.CreateEvent(context.Background(), req)
    
    // Assert
    require.NoError(t, err)
    assert.NotNil(t, resp.Event)
    assert.Equal(t, "Test Event", resp.Event.Title)
    assert.Equal(t, "user123", resp.Event.UserId)
}

func TestGRPCServer_CreateEvent_ValidationError(t *testing.T) {
    // Arrange
    logger := &mockLogger{}
    application := &mockApplication{}
    server := NewServer(logger, "localhost", 0, application)
    
    listener, err := net.Listen("tcp", ":0")
    require.NoError(t, err)
    
    grpcServer := grpc.NewServer()
    api.RegisterEventServiceServer(grpcServer, server)
    
    go func() {
        grpcServer.Serve(listener)
    }()
    defer grpcServer.Stop()
    
    conn, err := grpc.Dial(listener.Addr().String(), grpc.WithInsecure())
    require.NoError(t, err)
    defer conn.Close()
    
    client := api.NewEventServiceClient(conn)
    
    // Act - запрос без заголовка
    req := &api.CreateEventRequest{
        Description:  "Test Description",
        UserId:       "user123",
        StartTime:    timestamppb.New(time.Now()),
        Duration:     durationpb.New(time.Hour),
        NotifyBefore: durationpb.New(15 * time.Minute),
    }
    
    resp, err := client.CreateEvent(context.Background(), req)
    
    // Assert
    assert.Error(t, err)
    assert.Nil(t, resp)
    
    // Проверяем gRPC код ошибки
    st, ok := status.FromError(err)
    assert.True(t, ok)
    assert.Equal(t, codes.InvalidArgument, st.Code())
    assert.Contains(t, st.Message(), "title is required")
}
```

## 🧪 Mock объекты

### 1. **Mock Logger**
```go
type mockLogger struct {
    debugCalled int
    infoCalled  int
    errorCalled int
    warnCalled  int
    messages    []string
}

func (m *mockLogger) Debug(msg string) {
    m.debugCalled++
    m.messages = append(m.messages, msg)
}

func (m *mockLogger) Info(msg string) {
    m.infoCalled++
    m.messages = append(m.messages, msg)
}

func (m *mockLogger) Error(msg string) {
    m.errorCalled++
    m.messages = append(m.messages, msg)
}

func (m *mockLogger) Warn(msg string) {
    m.warnCalled++
    m.messages = append(m.messages, msg)
}
```

### 2. **Mock Application**
```go
type mockApplication struct {
    createEventCalled int
    getEventCalled    int
    updateEventCalled int
    deleteEventCalled int
    listEventsCalled  int
    
    createEventError error
    getEventError    error
    getEventResult   storage.Event
    listEventsResult []storage.Event
}

func (m *mockApplication) CreateEvent(ctx context.Context, id, title, description, userID string, startTime time.Time, duration, notifyBefore calendar_types.CalendarDuration) error {
    m.createEventCalled++
    return m.createEventError
}

func (m *mockApplication) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
    m.getEventCalled++
    return m.getEventResult, m.getEventError
}

func (m *mockApplication) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
    m.listEventsCalled++
    return m.listEventsResult, nil
}

// ... остальные методы
```

### 3. **Mock Storage**
```go
type mockStorage struct {
    addEventCalled    int
    updateEventCalled int
    deleteEventCalled int
    getEventCalled    int
    listEventsCalled  int
    
    addEventError    error
    updateEventError error
    deleteEventError error
    getEventError    error
    getEventResult   storage.Event
    listEventsResult []storage.Event
}

func (m *mockStorage) AddEvent(ctx context.Context, e storage.Event) error {
    m.addEventCalled++
    return m.addEventError
}

func (m *mockStorage) GetEventByID(ctx context.Context, id string) (storage.Event, error) {
    m.getEventCalled++
    return m.getEventResult, m.getEventError
}

// ... остальные методы
```

## 🚀 Запуск тестов

### 1. **Все тесты**
```bash
go test ./...
```

### 2. **Конкретный пакет**
```bash
go test ./internal/app/
go test ./internal/server/http/
go test ./internal/server/grpc/
```

### 3. **С покрытием**
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 4. **Бенчмарки**
```bash
go test -bench=. ./internal/app/
```

## 🔧 Manual тестирование

### 1. **HTTP API с curl**
```bash
# Создание события
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "user_id": "user123",
    "start_time": "2025-07-15T12:00:00Z",
    "duration": "1h0m0s",
    "notify_before": "15m0s"
  }'

# Получение события
curl "http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000"

# Обновление события
curl -X PUT http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Meeting",
    "description": "Updated description",
    "user_id": "user123",
    "start_time": "2025-07-15T13:00:00Z",
    "duration": "2h0m0s",
    "notify_before": "30m0s"
  }'

# Удаление события
curl -X DELETE http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000

# Получение событий за день
curl "http://localhost:8080/events?day=2025-07-15"
```

### 2. **gRPC API с grpcurl**
```bash
# Создание события
grpcurl -plaintext -proto api/EventService.proto \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "user_id": "user123",
    "start_time": "2025-07-15T12:00:00Z",
    "duration": "3600s",
    "notify_before": "900s"
  }' \
  localhost:50051 calendar.EventService/CreateEvent

# Получение события
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 calendar.EventService/GetEvent

# Обновление события
grpcurl -plaintext -proto api/EventService.proto \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Meeting",
    "description": "Updated description",
    "user_id": "user123",
    "start_time": "2025-07-15T13:00:00Z",
    "duration": "7200s",
    "notify_before": "1800s"
  }' \
  localhost:50051 calendar.EventService/UpdateEvent

# Удаление события
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 calendar.EventService/DeleteEvent

# Получение событий за день
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"date": "2025-07-15T00:00:00Z"}' \
  localhost:50051 calendar.EventService/ListEventsForDay
```

## 📊 Покрытие тестами

### 1. **Генерация отчета**
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

### 2. **HTML отчет**
```bash
go tool cover -html=coverage.out -o coverage.html
```

### 3. **Целевое покрытие**
```bash
# Проверка покрытия
go test -cover -coverpkg=./internal/app/ ./internal/app/
go test -cover -coverpkg=./internal/server/http/ ./internal/server/http/
go test -cover -coverpkg=./internal/server/grpc/ ./internal/server/grpc/
```

## 🎯 Лучшие практики

### 1. **Структура тестов**
```go
func TestFunctionName_Scenario_ExpectedResult(t *testing.T) {
    // Arrange - подготовка данных
    // Act - выполнение действия
    // Assert - проверка результата
}
```

### 2. **Именование тестов**
```go
// Хорошо
func TestCreateEvent_Success(t *testing.T)
func TestCreateEvent_EmptyTitle_ReturnsError(t *testing.T)
func TestCreateEvent_DuplicateID_ReturnsError(t *testing.T)

// Плохо
func TestCreateEvent(t *testing.T)
func TestCreateEvent2(t *testing.T)
```

### 3. **Использование таблиц**
```go
func TestParseDuration(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected time.Duration
        hasError bool
    }{
        {"valid duration", "1h30m", 90 * time.Minute, false},
        {"invalid duration", "invalid", 0, true},
        {"empty string", "", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := calendar_types.ParseDuration(tt.input)
            
            if tt.hasError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result.Duration())
            }
        })
    }
}
```

## 🎓 Заключение

Тестирование в нашем проекте обеспечивает:

### **Unit тесты:**
- ✅ Проверка бизнес-логики
- ✅ Изолированное тестирование компонентов
- ✅ Быстрое выполнение
- ✅ Легкая отладка

### **Integration тесты:**
- ✅ Проверка взаимодействия компонентов
- ✅ Тестирование HTTP и gRPC API
- ✅ Проверка реальных сценариев
- ✅ Валидация контрактов

### **Manual тесты:**
- ✅ Проверка в реальных условиях
- ✅ Отладка проблем
- ✅ Демонстрация функциональности
- ✅ Валидация производительности

Это обеспечивает надежность и качество кода, а также упрощает рефакторинг и добавление новых функций. 