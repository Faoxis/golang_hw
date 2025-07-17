# 2. Архитектура проекта

## 🏗️ Общая архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │    │   gRPC Client   │    │   Web Browser   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          │ HTTP/1.1             │ HTTP/2 + gRPC        │ HTTP/1.1
          │ JSON                 │ Protocol Buffers     │ JSON
          ▼                      ▼                      ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Calendar Service                             │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │   HTTP Server   │  │   gRPC Server   │  │   Application   │ │
│  │   (port 8080)   │  │   (port 50051)  │  │   (Business     │ │
│  │                 │  │                 │  │    Logic)       │ │
│  └─────────┬───────┘  └─────────┬───────┘  └─────────┬───────┘ │
└────────────┼────────────────────┼────────────────────┼─────────┘
             │                    │                    │
             │                    │                    │
             ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Storage Layer                            │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐ │
│  │  Memory Storage │  │  SQL Storage    │  │   Logger        │ │
│  │   (in-memory)   │  │   (PostgreSQL)  │  │                 │ │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## 📁 Структура проекта

```
hw12_13_14_15_calendar/
├── api/                          # Protobuf и автогенерированный код
│   ├── EventService.proto        # Описание gRPC API
│   ├── EventService.pb.go        # Автогенерированные структуры
│   ├── EventService_grpc.pb.go   # Автогенерированные клиент/сервер
│   └── event_service.go          # go:generate директива
├── cmd/calendar/                 # Точка входа приложения
│   ├── main.go                   # Основной файл
│   ├── config.go                 # Конфигурация
│   └── version.go                # Версионирование
├── internal/                     # Внутренний код приложения
│   ├── app/                      # Бизнес-логика
│   │   └── app.go                # Основная логика календаря
│   ├── server/                   # Серверы
│   │   ├── http/                 # HTTP сервер
│   │   │   ├── server.go         # HTTP сервер
│   │   │   ├── event_handlers.go # HTTP обработчики
│   │   │   ├── eventdto.go       # HTTP DTO
│   │   │   └── middleware.go     # HTTP middleware
│   │   └── grpc/                 # gRPC сервер
│   │       └── server.go         # gRPC сервер
│   ├── storage/                  # Слой хранения
│   │   ├── event.go              # Модель события
│   │   ├── memory/               # In-memory хранилище
│   │   └── sql/                  # SQL хранилище
│   └── logger/                   # Логирование
├── configs/                      # Конфигурационные файлы
│   └── config.yml                # Основной конфиг
├── env/                          # Docker окружение
│   └── docker-compose.yml        # PostgreSQL + RabbitMQ
└── migrations/                   # SQL миграции
```

## 🔄 Поток данных

### 1. **HTTP запрос**
```
HTTP Client → HTTP Server → Application → Storage → Database
HTTP Client ← HTTP Server ← Application ← Storage ← Database
```

### 2. **gRPC запрос**
```
gRPC Client → gRPC Server → Application → Storage → Database
gRPC Client ← gRPC Server ← Application ← Storage ← Database
```

## 🎯 Ключевые принципы архитектуры

### 1. **Разделение ответственности**
- **HTTP Server** - обработка HTTP запросов, сериализация JSON
- **gRPC Server** - обработка gRPC запросов, сериализация Protobuf
- **Application** - бизнес-логика (не зависит от протоколов)
- **Storage** - абстракция хранения данных
- **Logger** - централизованное логирование

### 2. **Инверсия зависимостей**
```go
// Application не знает о HTTP/gRPC
type Application interface {
    CreateEvent(ctx context.Context, id, title, description, userID string, startTime time.Time, duration, notifyBefore calendar_types.CalendarDuration) error
    GetEventByID(ctx context.Context, id string) (storage.Event, error)
    // ...
}

// HTTP и gRPC серверы используют Application
type Server struct {
    application Application  // Зависимость от интерфейса
    logger      Logger
}
```

### 3. **Единая бизнес-логика**
```go
// Один и тот же код используется в HTTP и gRPC
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    // Валидация
    if req.Title == "" {
        return nil, status.Error(codes.InvalidArgument, "title is required")
    }
    
    // Вызов бизнес-логики
    err := s.application.CreateEvent(ctx, id, req.Title, req.Description, req.UserId, startTime, duration, notifyBefore)
    
    // Маппинг ответа
    return &api.CreateEventResponse{Event: mapStorageEventToProtoEvent(event)}, nil
}
```

## 🔌 Интерфейсы и абстракции

### 1. **Application Interface**
```go
type Application interface {
    CreateEvent(ctx context.Context, id, title, description, userID string, startTime time.Time, duration, notifyBefore calendar_types.CalendarDuration) error
    UpdateEvent(ctx context.Context, id, title, description, userID string, startTime time.Time, duration, notifyBefore calendar_types.CalendarDuration) error
    DeleteEvent(ctx context.Context, id string) error
    GetEventByID(ctx context.Context, id string) (storage.Event, error)
    ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
    ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
    ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}
```

### 2. **Storage Interface**
```go
type Storage interface {
    AddEvent(ctx context.Context, e storage.Event) error
    UpdateEvent(ctx context.Context, e storage.Event) error
    DeleteEvent(ctx context.Context, id string) error
    GetEventByID(ctx context.Context, id string) (storage.Event, error)
    ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
    ListEventsForWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
    ListEventsForMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
    Close() error
}
```

### 3. **Logger Interface**
```go
type Logger interface {
    Debug(msg string)
    Info(msg string)
    Error(msg string)
    Warn(msg string)
}
```

## 🚀 Запуск и конфигурация

### 1. **Конфигурация**
```yaml
# configs/config.yml
logger:
  level: debug

server:
  host: "localhost"
  port: 8080          # HTTP порт
  grpc_port: 50051    # gRPC порт

storage:
  type: database      # memory или database
  host: localhost
  port: 5434
  user: calendar_user
  password: calendar_pass
  database: calendar
  sslmode: disable
```

### 2. **Запуск**
```go
// main.go
func main() {
    // Загрузка конфигурации
    config, err := LoadConfig(configFile)
    
    // Создание логгера
    logg := logger.New(config.Logger.Level)
    
    // Создание хранилища
    var storage app.Storage
    switch config.Storage.Type {
    case "memory":
        storage = memorystorage.New(logg)
    case "database":
        storage, err = sqlstorage.New(config.Storage.GetPostgresDSN(), logg)
    }
    
    // Создание приложения
    calendar := app.New(logg, storage)
    
    // Создание серверов
    httpServer := internalhttp.NewServer(logg, config.Server.Host, config.Server.Port, calendar)
    grpcServer := grpcserver.NewServer(logg, config.Server.Host, config.Server.GRPCPort, calendar)
    
    // Запуск серверов
    go httpServer.Start(ctx)
    grpcServer.Start(ctx)
}
```

## 🔄 Жизненный цикл запроса

### HTTP запрос:
1. **Клиент** отправляет HTTP запрос на `localhost:8080/events`
2. **HTTP Server** получает запрос через Chi router
3. **Middleware** логирует запрос
4. **Handler** парсит JSON в структуру
5. **Application** выполняет бизнес-логику
6. **Storage** сохраняет данные в PostgreSQL
7. **Response** возвращается клиенту в JSON

### gRPC запрос:
1. **Клиент** вызывает `CreateEvent` через gRPC
2. **gRPC Server** получает запрос через HTTP/2
3. **Server** парсит Protobuf в структуру
4. **Application** выполняет ту же бизнес-логику
5. **Storage** сохраняет данные в PostgreSQL
6. **Response** возвращается клиенту в Protobuf

## 🎓 Преимущества архитектуры

### 1. **Модульность**
- Каждый компонент можно заменить независимо
- Легко добавлять новые протоколы (WebSocket, GraphQL)
- Простое тестирование отдельных компонентов

### 2. **Масштабируемость**
- HTTP и gRPC серверы работают независимо
- Можно запускать на разных портах/серверах
- Легко добавлять балансировщики нагрузки

### 3. **Тестируемость**
- Моки для всех интерфейсов
- Интеграционные тесты с реальными серверами
- Изолированное тестирование бизнес-логики

### 4. **Гибкость**
- Поддержка разных типов хранилищ
- Конфигурируемые порты и настройки
- Легкое переключение между протоколами 