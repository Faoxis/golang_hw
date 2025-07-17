# 5. HTTP vs gRPC

## 🔄 Сравнение протоколов

### Общая архитектура
```
┌─────────────────────────────────────────────────────────┐
│                    Calendar Service                     │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   HTTP API      │  │   gRPC API      │              │
│  │   (port 8080)   │  │   (port 50051)  │              │
│  │                 │  │                 │              │
│  │   JSON          │  │   Protobuf      │              │
│  │   HTTP/1.1      │  │   HTTP/2        │              │
│  └─────────┬───────┘  └─────────┬───────┘              │
└────────────┼────────────────────┼──────────────────────┘
             │                    │
             ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                Application Layer                        │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   HTTP          │  │   gRPC          │              │
│  │   Handlers      │  │   Server        │              │
│  └─────────┬───────┘  └─────────┬───────┘              │
└────────────┼────────────────────┼──────────────────────┘
             │                    │
             ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                Business Logic                           │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   CreateEvent   │  │   GetEvent      │              │
│  │   UpdateEvent   │  │   DeleteEvent   │              │
│  │   ListEvents    │  │   ...           │              │
│  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

## 📊 Детальное сравнение

| Аспект | HTTP API | gRPC API |
|--------|----------|----------|
| **Протокол** | HTTP/1.1 | HTTP/2 |
| **Формат данных** | JSON | Protocol Buffers |
| **Типизация** | Слабая | Строгая |
| **Скорость** | Медленнее | Быстрее |
| **Размер данных** | Больше | Меньше |
| **Сложность** | Простая | Средняя |
| **Поддержка браузеров** | Полная | Ограниченная |
| **Инструменты** | Много | Меньше |
| **Отладка** | Легкая | Сложнее |

## 🔧 Реализация в проекте

### 1. **HTTP API - CreateEvent**

#### Запрос:
```bash
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
```

#### Обработчик:
```go
// internal/server/http/event_handlers.go
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    // Парсинг JSON
    var req CreateEventRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.logger.Error(fmt.Sprintf("Failed to decode request: %v", err))
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Валидация
    if req.Title == "" {
        h.logger.Error("CreateEvent: title is required")
        http.Error(w, "title is required", http.StatusBadRequest)
        return
    }
    
    // Генерация ID
    id := uuid.New().String()
    
    // Парсинг времени
    startTime, err := time.Parse(time.RFC3339, req.StartTime)
    if err != nil {
        h.logger.Error(fmt.Sprintf("Invalid start_time format: %v", err))
        http.Error(w, "invalid start_time format", http.StatusBadRequest)
        return
    }
    
    // Парсинг длительности
    duration, err := calendar_types.ParseDuration(req.Duration)
    if err != nil {
        h.logger.Error(fmt.Sprintf("Invalid duration format: %v", err))
        http.Error(w, "invalid duration format", http.StatusBadRequest)
        return
    }
    
    notifyBefore, err := calendar_types.ParseDuration(req.NotifyBefore)
    if err != nil {
        h.logger.Error(fmt.Sprintf("Invalid notify_before format: %v", err))
        http.Error(w, "invalid notify_before format", http.StatusBadRequest)
        return
    }
    
    // Вызов бизнес-логики
    err = h.application.CreateEvent(r.Context(), id, req.Title, req.Description, req.UserID, startTime, duration, notifyBefore)
    if err != nil {
        h.logger.Error(fmt.Sprintf("CreateEvent failed: %v", err))
        http.Error(w, "failed to create event", http.StatusInternalServerError)
        return
    }
    
    // Получение созданного события
    event, err := h.application.GetEventByID(r.Context(), id)
    if err != nil {
        h.logger.Error(fmt.Sprintf("Failed to get created event: %v", err))
        http.Error(w, "failed to get created event", http.StatusInternalServerError)
        return
    }
    
    // Маппинг в JSON
    response := mapStorageEventToJSONEvent(event)
    
    // Отправка ответа
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(response)
}
```

#### Ответ:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly team sync",
  "user_id": "user123",
  "start_time": "2025-07-15T12:00:00Z",
  "duration": "1h0m0s",
  "notify_before": "15m0s"
}
```

### 2. **gRPC API - CreateEvent**

#### Запрос:
```bash
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
```

#### Обработчик:
```go
// internal/server/grpc/server.go
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent called with title: %s", req.Title))
    
    // Валидация (строгая типизация)
    if req.Title == "" {
        s.logger.Error("CreateEvent: title is required")
        return nil, status.Error(codes.InvalidArgument, "title is required")
    }
    
    if req.UserId == "" {
        s.logger.Error("CreateEvent: user_id is required")
        return nil, status.Error(codes.InvalidArgument, "user_id is required")
    }
    
    // Генерация ID
    id := uuid.New().String()
    
    // Конвертация времени (автоматическая)
    startTime := req.StartTime.AsTime()
    
    // Конвертация длительности (автоматическая)
    duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
    notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())
    
    // Вызов бизнес-логики (та же самая!)
    err := s.application.CreateEvent(ctx, id, req.Title, req.Description, req.UserId, startTime, duration, notifyBefore)
    if err != nil {
        s.logger.Error(fmt.Sprintf("CreateEvent failed: %v", err))
        return nil, status.Error(codes.Internal, "failed to create event")
    }
    
    // Получение созданного события
    event, err := s.application.GetEventByID(ctx, id)
    if err != nil {
        s.logger.Error(fmt.Sprintf("Failed to get created event: %v", err))
        return nil, status.Error(codes.Internal, "failed to get created event")
    }
    
    // Маппинг в Protobuf (автоматический)
    protoEvent := mapStorageEventToProtoEvent(event)
    
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent completed successfully, event ID: %s", id))
    
    return &api.CreateEventResponse{Event: protoEvent}, nil
}
```

#### Ответ:
```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly team sync",
    "user_id": "user123",
    "start_time": "2025-07-15T12:00:00Z",
    "duration": "3600s",
    "notify_before": "900s"
  }
}
```

## 🔄 Маппинг данных

### 1. **HTTP DTO**
```go
// internal/server/http/eventdto.go
type CreateEventRequest struct {
    Title        string `json:"title"`
    Description  string `json:"description"`
    UserID       string `json:"user_id"`
    StartTime    string `json:"start_time"`    // RFC3339 формат
    Duration     string `json:"duration"`      // "1h30m" формат
    NotifyBefore string `json:"notify_before"` // "15m" формат
}

type EventResponse struct {
    ID           string `json:"id"`
    Title        string `json:"title"`
    Description  string `json:"description"`
    UserID       string `json:"user_id"`
    StartTime    string `json:"start_time"`
    Duration     string `json:"duration"`
    NotifyBefore string `json:"notify_before"`
}
```

### 2. **gRPC Protobuf**
```protobuf
// api/EventService.proto
message CreateEventRequest {
  string title = 1;
  string description = 2;
  string user_id = 3;
  google.protobuf.Timestamp start_time = 4;
  google.protobuf.Duration duration = 5;
  google.protobuf.Duration notify_before = 6;
}

message Event {
  string id = 1;
  string title = 2;
  string description = 3;
  string user_id = 4;
  google.protobuf.Timestamp start_time = 5;
  google.protobuf.Duration duration = 6;
  google.protobuf.Duration notify_before = 7;
}
```

## 🎯 Ключевые различия

### 1. **Парсинг данных**

#### HTTP (ручной парсинг):
```go
// Парсинг времени
startTime, err := time.Parse(time.RFC3339, req.StartTime)
if err != nil {
    return err
}

// Парсинг длительности
duration, err := calendar_types.ParseDuration(req.Duration)
if err != nil {
    return err
}
```

#### gRPC (автоматический):
```go
// Автоматическая конвертация
startTime := req.StartTime.AsTime()
duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
```

### 2. **Валидация**

#### HTTP (ручная):
```go
if req.Title == "" {
    http.Error(w, "title is required", http.StatusBadRequest)
    return
}
```

#### gRPC (структурированная):
```go
if req.Title == "" {
    return nil, status.Error(codes.InvalidArgument, "title is required")
}
```

### 3. **Обработка ошибок**

#### HTTP:
```go
http.Error(w, "event not found", http.StatusNotFound)
```

#### gRPC:
```go
return nil, status.Error(codes.NotFound, "event not found")
```

## 📈 Производительность

### 1. **Размер данных**
```json
// HTTP JSON (примерно 200 байт)
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly team sync",
  "user_id": "user123",
  "start_time": "2025-07-15T12:00:00Z",
  "duration": "1h0m0s",
  "notify_before": "15m0s"
}
```

```protobuf
// gRPC Protobuf (примерно 80 байт)
08 80 80 80 80 80 80 80 80 80 01 12 0B 54 65 61 6D 20 4D 65 65 74 69 6E 67 1A 10 57 65 65 6B 6C 79 20 74 65 61 6D 20 73 79 6E 63 22 07 75 73 65 72 31 32 33
```

### 2. **Скорость сериализации**
```go
// HTTP JSON
jsonData, _ := json.Marshal(event)  // ~1000ns

// gRPC Protobuf
protoData, _ := proto.Marshal(event)  // ~200ns
```

## 🛠️ Инструменты

### 1. **HTTP API**
```bash
# Тестирование с curl
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","user_id":"user1"}'

# Тестирование с Postman/Bruno
# Визуальный интерфейс для тестирования API
```

### 2. **gRPC API**
```bash
# Тестирование с grpcurl
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"title":"Test","user_id":"user1"}' \
  localhost:50051 calendar.EventService/CreateEvent

# Тестирование с BloomRPC/grpcui
# Визуальные инструменты для gRPC
```

## 🎯 Когда использовать что?

### HTTP API подходит для:
- ✅ Веб-приложения и браузеры
- ✅ Простые интеграции
- ✅ Отладка и тестирование
- ✅ RESTful архитектура
- ✅ Человекочитаемые данные

### gRPC API подходит для:
- ✅ Микросервисная архитектура
- ✅ Высокая производительность
- ✅ Строгая типизация
- ✅ Автогенерация клиентов
- ✅ Внутренние API

## 🔄 Единая бизнес-логика

### Ключевое преимущество нашей архитектуры:
```go
// Один и тот же код используется в HTTP и gRPC!
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    // ... парсинг и валидация ...
    
    // Единая бизнес-логика
    err := h.application.CreateEvent(r.Context(), id, req.Title, req.Description, req.UserID, startTime, duration, notifyBefore)
    
    // ... маппинг ответа ...
}

func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    // ... парсинг и валидация ...
    
    // Та же самая бизнес-логика!
    err := s.application.CreateEvent(ctx, id, req.Title, req.Description, req.UserId, startTime, duration, notifyBefore)
    
    // ... маппинг ответа ...
}
```

## 🎓 Заключение

В нашем проекте:

### HTTP API:
- **Простой** для понимания и отладки
- **Универсальный** - работает везде
- **Человекочитаемый** - JSON формат
- **Медленнее** - текстовый протокол

### gRPC API:
- **Быстрый** - бинарный протокол
- **Типобезопасный** - строгая типизация
- **Эффективный** - HTTP/2 мультиплексирование
- **Автогенерируемый** - клиентский код

### Общее:
- **Единая бизнес-логика** - один код для обоих API
- **Одинаковая функциональность** - все операции доступны
- **Консистентное логирование** - единый подход
- **Гибкость выбора** - клиенты могут использовать любой API

Это позволяет удовлетворить различные потребности клиентов, от простых веб-интерфейсов до высокопроизводительных микросервисов. 