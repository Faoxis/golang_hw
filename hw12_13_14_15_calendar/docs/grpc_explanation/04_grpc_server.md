# 4. gRPC сервер

## 🏗️ Архитектура gRPC сервера

```
┌─────────────────────────────────────────────────────────┐
│                    gRPC Server                          │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   HTTP/2        │  │   gRPC          │              │
│  │   Transport     │  │   Protocol      │              │
│  └─────────┬───────┘  └─────────┬───────┘              │
└────────────┼────────────────────┼──────────────────────┘
             │                    │
             ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                EventServiceServer                       │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   CreateEvent   │  │   GetEvent      │              │
│  │   UpdateEvent   │  │   DeleteEvent   │              │
│  │   ListEvents    │  │   ...           │              │
│  └─────────┬───────┘  └─────────┬───────┘              │
└────────────┼────────────────────┼──────────────────────┘
             │                    │
             ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                  Application Layer                      │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   Business      │  │   Validation    │              │
│  │   Logic         │  │   & Mapping     │              │
│  └─────────┬───────┘  └─────────┬───────┘              │
└────────────┼────────────────────┼──────────────────────┘
             │                    │
             ▼                    ▼
┌─────────────────────────────────────────────────────────┐
│                   Storage Layer                         │
│  ┌─────────────────┐  ┌─────────────────┐              │
│  │   PostgreSQL    │  │   Memory        │              │
│  │   Database      │  │   Storage       │              │
│  └─────────────────┘  └─────────────────┘              │
└─────────────────────────────────────────────────────────┘
```

## 🔧 Реализация сервера

### 1. **Структура сервера**
```go
// internal/server/grpc/server.go
type Server struct {
    api.UnimplementedEventServiceServer  // Встраиваем базовый интерфейс
    application app.Application          // Бизнес-логика
    logger      logger.Logger           // Логирование
    host        string                  // Хост сервера
    port        int                     // Порт сервера
}
```

### 2. **Создание сервера**
```go
func NewServer(logger logger.Logger, host string, port int, application app.Application) *Server {
    return &Server{
        application: application,
        logger:      logger,
        host:        host,
        port:        port,
    }
}
```

### 3. **Запуск сервера**
```go
func (s *Server) Start(ctx context.Context) error {
    // Создаем TCP listener
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
    if err != nil {
        return fmt.Errorf("failed to listen: %w", err)
    }
    
    // Создаем gRPC сервер
    grpcServer := grpc.NewServer()
    
    // Регистрируем наш сервис
    api.RegisterEventServiceServer(grpcServer, s)
    
    s.logger.Info(fmt.Sprintf("gRPC server starting on %s:%d", s.host, s.port))
    
    // Запускаем сервер
    if err := grpcServer.Serve(listener); err != nil {
        return fmt.Errorf("failed to serve: %w", err)
    }
    
    return nil
}
```

## 🎯 Реализация методов

### 1. **CreateEvent - создание события**
```go
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent called with title: %s", req.Title))
    
    // Валидация входных данных
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
    
    // Конвертация времени
    startTime := req.StartTime.AsTime()
    
    // Конвертация длительности
    duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
    notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())
    
    // Вызов бизнес-логики
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
    
    // Маппинг в Protobuf формат
    protoEvent := mapStorageEventToProtoEvent(event)
    
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent completed successfully, event ID: %s", id))
    
    return &api.CreateEventResponse{Event: protoEvent}, nil
}
```

### 2. **GetEvent - получение события**
```go
func (s *Server) GetEvent(ctx context.Context, req *api.GetEventRequest) (*api.GetEventResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC GetEvent called with ID: %s", req.Id))
    
    // Валидация
    if req.Id == "" {
        s.logger.Error("GetEvent: id is required")
        return nil, status.Error(codes.InvalidArgument, "id is required")
    }
    
    // Получение события
    event, err := s.application.GetEventByID(ctx, req.Id)
    if err != nil {
        s.logger.Error(fmt.Sprintf("GetEvent failed: %v", err))
        return nil, status.Error(codes.NotFound, "event not found")
    }
    
    // Маппинг в Protobuf формат
    protoEvent := mapStorageEventToProtoEvent(event)
    
    s.logger.Info(fmt.Sprintf("gRPC GetEvent completed successfully, event ID: %s", req.Id))
    
    return &api.GetEventResponse{Event: protoEvent}, nil
}
```

### 3. **UpdateEvent - обновление события**
```go
func (s *Server) UpdateEvent(ctx context.Context, req *api.UpdateEventRequest) (*api.UpdateEventResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC UpdateEvent called with ID: %s", req.Id))
    
    // Валидация
    if req.Id == "" {
        s.logger.Error("UpdateEvent: id is required")
        return nil, status.Error(codes.InvalidArgument, "id is required")
    }
    
    if req.Title == "" {
        s.logger.Error("UpdateEvent: title is required")
        return nil, status.Error(codes.InvalidArgument, "title is required")
    }
    
    // Конвертация данных
    startTime := req.StartTime.AsTime()
    duration := calendar_types.CalendarDuration(req.Duration.AsDuration())
    notifyBefore := calendar_types.CalendarDuration(req.NotifyBefore.AsDuration())
    
    // Обновление события
    err := s.application.UpdateEvent(ctx, req.Id, req.Title, req.Description, req.UserId, startTime, duration, notifyBefore)
    if err != nil {
        s.logger.Error(fmt.Sprintf("UpdateEvent failed: %v", err))
        return nil, status.Error(codes.Internal, "failed to update event")
    }
    
    // Получение обновленного события
    event, err := s.application.GetEventByID(ctx, req.Id)
    if err != nil {
        s.logger.Error(fmt.Sprintf("Failed to get updated event: %v", err))
        return nil, status.Error(codes.Internal, "failed to get updated event")
    }
    
    // Маппинг в Protobuf формат
    protoEvent := mapStorageEventToProtoEvent(event)
    
    s.logger.Info(fmt.Sprintf("gRPC UpdateEvent completed successfully, event ID: %s", req.Id))
    
    return &api.UpdateEventResponse{Event: protoEvent}, nil
}
```

### 4. **DeleteEvent - удаление события**
```go
func (s *Server) DeleteEvent(ctx context.Context, req *api.DeleteEventRequest) (*api.DeleteEventResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC DeleteEvent called with ID: %s", req.Id))
    
    // Валидация
    if req.Id == "" {
        s.logger.Error("DeleteEvent: id is required")
        return nil, status.Error(codes.InvalidArgument, "id is required")
    }
    
    // Удаление события
    err := s.application.DeleteEvent(ctx, req.Id)
    if err != nil {
        s.logger.Error(fmt.Sprintf("DeleteEvent failed: %v", err))
        return nil, status.Error(codes.Internal, "failed to delete event")
    }
    
    s.logger.Info(fmt.Sprintf("gRPC DeleteEvent completed successfully, event ID: %s", req.Id))
    
    return &api.DeleteEventResponse{}, nil
}
```

### 5. **ListEvents - получение списка событий**
```go
func (s *Server) ListEventsForDay(ctx context.Context, req *api.ListEventsForDayRequest) (*api.ListEventsResponse, error) {
    s.logger.Info(fmt.Sprintf("gRPC ListEventsForDay called for date: %s", req.Date.AsTime().Format("2006-01-02")))
    
    // Получение событий
    events, err := s.application.ListEventsForDay(ctx, req.Date.AsTime())
    if err != nil {
        s.logger.Error(fmt.Sprintf("ListEventsForDay failed: %v", err))
        return nil, status.Error(codes.Internal, "failed to list events")
    }
    
    // Маппинг событий
    protoEvents := make([]*api.Event, len(events))
    for i, event := range events {
        protoEvents[i] = mapStorageEventToProtoEvent(event)
    }
    
    s.logger.Info(fmt.Sprintf("gRPC ListEventsForDay completed, found %d events", len(events)))
    
    return &api.ListEventsResponse{Events: protoEvents}, nil
}
```

## 🔄 Маппинг данных

### 1. **Storage Event → Protobuf Event**
```go
func mapStorageEventToProtoEvent(event storage.Event) *api.Event {
    return &api.Event{
        Id:           event.ID,
        Title:        event.Title,
        Description:  event.Description,
        UserId:       event.UserID,
        StartTime:    timestamppb.New(event.StartTime),
        Duration:     durationpb.New(event.Duration.Duration()),
        NotifyBefore: durationpb.New(event.NotifyBefore.Duration()),
    }
}
```

### 2. **Protobuf Event → Storage Event**
```go
func mapProtoEventToStorageEvent(event *api.Event) storage.Event {
    return storage.Event{
        ID:           event.Id,
        Title:        event.Title,
        Description:  event.Description,
        UserID:       event.UserId,
        StartTime:    event.StartTime.AsTime(),
        Duration:     calendar_types.CalendarDuration(event.Duration.AsDuration()),
        NotifyBefore: calendar_types.CalendarDuration(event.NotifyBefore.AsDuration()),
    }
}
```

## 🚨 Обработка ошибок

### 1. **gRPC коды ошибок**
```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Различные типы ошибок
return nil, status.Error(codes.InvalidArgument, "title is required")     // 400
return nil, status.Error(codes.NotFound, "event not found")              // 404
return nil, status.Error(codes.Internal, "failed to create event")       // 500
return nil, status.Error(codes.AlreadyExists, "event already exists")    // 409
```

### 2. **Логирование ошибок**
```go
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    // Логируем начало операции
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent called with title: %s", req.Title))
    
    // Валидация
    if req.Title == "" {
        s.logger.Error("CreateEvent: title is required")  // Логируем ошибку
        return nil, status.Error(codes.InvalidArgument, "title is required")
    }
    
    // Бизнес-логика
    err := s.application.CreateEvent(ctx, id, req.Title, req.Description, req.UserId, startTime, duration, notifyBefore)
    if err != nil {
        s.logger.Error(fmt.Sprintf("CreateEvent failed: %v", err))  // Логируем ошибку
        return nil, status.Error(codes.Internal, "failed to create event")
    }
    
    // Логируем успешное завершение
    s.logger.Info(fmt.Sprintf("gRPC CreateEvent completed successfully, event ID: %s", id))
    
    return response, nil
}
```

## 🔧 Конфигурация и запуск

### 1. **В main.go**
```go
func main() {
    // ... загрузка конфигурации и создание зависимостей
    
    // Создание gRPC сервера
    grpcServer := grpcserver.NewServer(logg, config.Server.Host, config.Server.GRPCPort, calendar)
    
    // Запуск в горутине
    go func() {
        if err := grpcServer.Start(ctx); err != nil {
            logg.Error(fmt.Sprintf("gRPC server failed: %v", err))
            cancel()
        }
    }()
    
    // ... запуск HTTP сервера
}
```

### 2. **Graceful shutdown**
```go
func (s *Server) Start(ctx context.Context) error {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
    if err != nil {
        return fmt.Errorf("failed to listen: %w", err)
    }
    
    grpcServer := grpc.NewServer()
    api.RegisterEventServiceServer(grpcServer, s)
    
    s.logger.Info(fmt.Sprintf("gRPC server starting on %s:%d", s.host, s.port))
    
    // Graceful shutdown
    go func() {
        <-ctx.Done()
        s.logger.Info("gRPC server shutting down...")
        grpcServer.GracefulStop()
    }()
    
    if err := grpcServer.Serve(listener); err != nil {
        return fmt.Errorf("failed to serve: %w", err)
    }
    
    return nil
}
```

## 🎯 Преимущества реализации

### 1. **Разделение ответственности**
- gRPC сервер отвечает только за протокол
- Бизнес-логика изолирована в Application
- Маппинг данных вынесен в отдельные функции

### 2. **Единообразие**
- Одинаковая обработка ошибок
- Единый стиль логирования
- Консистентная валидация

### 3. **Тестируемость**
- Легко мокать Application
- Изолированное тестирование методов
- Интеграционные тесты с реальным сервером

### 4. **Масштабируемость**
- Независимый от HTTP сервера
- Возможность запуска на разных портах
- Легкое добавление middleware

## 🎓 Заключение

gRPC сервер в нашем проекте:
- **Реализует** все методы из Protobuf спецификации
- **Использует** единую бизнес-логику с HTTP API
- **Обеспечивает** строгую типизацию и валидацию
- **Поддерживает** логирование и обработку ошибок
- **Интегрируется** с существующей архитектурой

Это позволяет клиентам выбирать между HTTP и gRPC API в зависимости от их потребностей. 