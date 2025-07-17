# 3. Protocol Buffers (Protobuf)

## 🤔 Что такое Protocol Buffers?

**Protocol Buffers** (Protobuf) - это язык описания интерфейсов (IDL) и формат сериализации данных, разработанный Google. Это основа gRPC.

## 📝 Синтаксис Protobuf

### Базовые типы данных
```protobuf
// Числовые типы
int32, int64, uint32, uint64, sint32, sint64
fixed32, fixed64, sfixed32, sfixed64
float, double

// Строковые типы
string, bytes

// Логический тип
bool
```

### Сообщения (Messages)
```protobuf
message Event {
  string id = 1;           // Поле с номером 1
  string title = 2;        // Поле с номером 2
  string description = 3;  // Поле с номером 3
  string user_id = 4;      // Поле с номером 4
  int64 start_time = 5;    // Unix timestamp
  int64 duration = 6;      // Длительность в секундах
  int64 notify_before = 7; // Уведомление за N секунд
}
```

## 🏗️ Наш Protobuf файл

### Полное описание API
```protobuf
syntax = "proto3";  // Версия протокола

package calendar;   // Пакет для организации

import "google/protobuf/timestamp.proto";  // Импорт стандартных типов

// Описание сервиса
service EventService {
  // Создание события
  rpc CreateEvent(CreateEventRequest) returns (CreateEventResponse);
  
  // Получение события по ID
  rpc GetEvent(GetEventRequest) returns (GetEventResponse);
  
  // Обновление события
  rpc UpdateEvent(UpdateEventRequest) returns (UpdateEventResponse);
  
  // Удаление события
  rpc DeleteEvent(DeleteEventRequest) returns (DeleteEventResponse);
  
  // Получение событий за день
  rpc ListEventsForDay(ListEventsForDayRequest) returns (ListEventsResponse);
  
  // Получение событий за неделю
  rpc ListEventsForWeek(ListEventsForWeekRequest) returns (ListEventsResponse);
  
  // Получение событий за месяц
  rpc ListEventsForMonth(ListEventsForMonthRequest) returns (ListEventsResponse);
}

// Структура события
message Event {
  string id = 1;
  string title = 2;
  string description = 3;
  string user_id = 4;
  google.protobuf.Timestamp start_time = 5;
  google.protobuf.Duration duration = 6;
  google.protobuf.Duration notify_before = 7;
}

// Запросы и ответы для каждого метода
message CreateEventRequest {
  string title = 1;
  string description = 2;
  string user_id = 3;
  google.protobuf.Timestamp start_time = 4;
  google.protobuf.Duration duration = 5;
  google.protobuf.Duration notify_before = 6;
}

message CreateEventResponse {
  Event event = 1;
}

message GetEventRequest {
  string id = 1;
}

message GetEventResponse {
  Event event = 1;
}

// ... остальные запросы и ответы
```

## 🔧 Генерация кода

### 1. **Установка инструментов**
```bash
# Установка protoc компилятора
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 2. **Команда генерации**
```bash
protoc --go_out=. --go-grpc_out=. api/EventService.proto
```

### 3. **Автоматизация в Go**
```go
// api/event_service.go
//go:generate protoc --go_out=. --go-grpc_out=. EventService.proto

package api
```

```bash
# Запуск генерации
go generate ./api/
```

### 3.1. Зачем нужен файл api/event_service.go?

Файл [`api/event_service.go`](../../api/event_service.go) не содержит никакой бизнес-логики. Его единственная задача — содержать директиву:

```go
//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative EventService.proto
```

Это позволяет запускать команду `go generate ./...` (или `make generate`), чтобы автоматически сгенерировать Go-код для gRPC/Protobuf API на основе файла `EventService.proto`. Такой подход — стандартная практика для автоматизации генерации кода, чтобы не смешивать ручной и автогенерируемый код.

Сам файл не участвует в работе приложения и не содержит исполняемого кода.

## 📦 Автогенерированные файлы

### 1. **EventService.pb.go** - структуры данных
```go
// Автогенерированные структуры
type Event struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Id           string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
    Title        string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
    Description  string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
    UserId       string                 `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
    StartTime    *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
    Duration     *durationpb.Duration   `protobuf:"bytes,6,opt,name=duration,proto3" json:"duration,omitempty"`
    NotifyBefore *durationpb.Duration   `protobuf:"bytes,7,opt,name=notify_before,json=notifyBefore,proto3" json:"notify_before,omitempty"`
}

type CreateEventRequest struct {
    state         protoimpl.MessageState
    sizeCache     protoimpl.SizeCache
    unknownFields protoimpl.UnknownFields

    Title        string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
    Description  string                 `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
    UserId       string                 `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
    StartTime    *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
    Duration     *durationpb.Duration   `protobuf:"bytes,5,opt,name=duration,proto3" json:"duration,omitempty"`
    NotifyBefore *durationpb.Duration   `protobuf:"bytes,6,opt,name=notify_before,json=notifyBefore,proto3" json:"notify_before,omitempty"`
}
```

### 2. **EventService_grpc.pb.go** - клиент и сервер
```go
// Интерфейс сервера
type EventServiceServer interface {
    CreateEvent(context.Context, *CreateEventRequest) (*CreateEventResponse, error)
    GetEvent(context.Context, *GetEventRequest) (*GetEventResponse, error)
    UpdateEvent(context.Context, *UpdateEventRequest) (*UpdateEventResponse, error)
    DeleteEvent(context.Context, *DeleteEventRequest) (*DeleteEventResponse, error)
    ListEventsForDay(context.Context, *ListEventsForDayRequest) (*ListEventsResponse, error)
    ListEventsForWeek(context.Context, *ListEventsForWeekRequest) (*ListEventsResponse, error)
    ListEventsForMonth(context.Context, *ListEventsForMonthRequest) (*ListEventsResponse, error)
    mustEmbedUnimplementedEventServiceServer()
}

// Клиент
type EventServiceClient interface {
    CreateEvent(ctx context.Context, in *CreateEventRequest, opts ...grpc.CallOption) (*CreateEventResponse, error)
    GetEvent(ctx context.Context, in *GetEventRequest, opts ...grpc.CallOption) (*GetEventResponse, error)
    // ... остальные методы
}

type eventServiceClient struct {
    cc grpc.ClientConnInterface
}
```

## 🔄 Маппинг типов

### 1. **Время и длительность**
```go
// Protobuf → Go
import (
    "google.golang.org/protobuf/types/known/timestamppb"
    "google.golang.org/protobuf/types/known/durationpb"
)

// Go время → Protobuf
startTime := time.Now()
protoStartTime := timestamppb.New(startTime)

// Go длительность → Protobuf
duration := time.Hour
protoDuration := durationpb.New(duration)

// Protobuf → Go время
goStartTime := protoStartTime.AsTime()

// Protobuf → Go длительность
goDuration := protoDuration.AsDuration()
```

### 2. **Маппинг в нашем коде**
```go
// internal/server/grpc/server.go
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

## 🎯 Преимущества Protobuf

### 1. **Компактность**
```json
// JSON
{
  "id": "123",
  "title": "Meeting",
  "start_time": "2025-07-15T12:00:00Z"
}
```
```protobuf
// Protobuf (бинарный, примерно в 2-3 раза меньше)
08 7B 12 06 4D 65 65 74 69 6E 67 2A 08 08 80 80 80 80 80 80 80 80 80 01
```

### 2. **Скорость сериализации**
```go
// JSON сериализация
jsonData, _ := json.Marshal(event)  // ~1000ns

// Protobuf сериализация
protoData, _ := proto.Marshal(event)  // ~200ns
```

### 3. **Строгая типизация**
```go
// Ошибка компиляции при неправильном типе
event := &api.Event{
    StartTime: "not a timestamp",  // ❌ Ошибка компиляции
}

// Правильно
event := &api.Event{
    StartTime: timestamppb.New(time.Now()),  // ✅ Правильный тип
}
```

### 4. **Обратная совместимость**
```protobuf
// Версия 1
message Event {
  string id = 1;
  string title = 2;
}

// Версия 2 (добавили поле)
message Event {
  string id = 1;
  string title = 2;
  string description = 3;  // Новое поле
}
```
Старые клиенты продолжают работать, просто не получают новое поле.

## 🔧 Работа с Protobuf в проекте

### 1. **Изменение API**
```protobuf
// 1. Редактируем EventService.proto
message Event {
  string id = 1;
  string title = 2;
  string description = 3;
  string user_id = 4;
  google.protobuf.Timestamp start_time = 5;
  google.protobuf.Duration duration = 6;
  google.protobuf.Duration notify_before = 7;
  // Добавляем новое поле
  string location = 8;  // Новое поле
}
```

### 2. **Регенерация кода**
```bash
go generate ./api/
```

### 3. **Обновление сервера**
```go
func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    // Теперь req.Location доступно
    if req.Location != "" {
        // Обработка локации
    }
    
    // ... остальная логика
}
```

## 🚀 Лучшие практики

### 1. **Нумерация полей**
```protobuf
message Event {
  string id = 1;           // ✅ Начинаем с 1
  string title = 2;        // ✅ Последовательная нумерация
  string description = 3;  // ✅ Не пропускаем номера
  // string old_field = 4; // ❌ Не удаляем, а помечаем как deprecated
  string new_field = 5;    // ✅ Новые поля в конце
}
```

### 2. **Именование**
```protobuf
message CreateEventRequest {  // ✅ PascalCase для сообщений
  string title = 1;          // ✅ snake_case для полей
  string user_id = 2;        // ✅ snake_case для полей
}

service EventService {        // ✅ PascalCase для сервисов
  rpc CreateEvent(...) returns (...);  // ✅ PascalCase для методов
}
```

### 3. **Документация**
```protobuf
// Event представляет событие в календаре
message Event {
  // Уникальный идентификатор события
  string id = 1;
  
  // Заголовок события (обязательное поле)
  string title = 2;
  
  // Описание события (опциональное)
  string description = 3;
}
```

## 🎓 Заключение

Protocol Buffers обеспечивают:
- **Эффективность** - компактный бинарный формат
- **Скорость** - быстрая сериализация/десериализация
- **Типобезопасность** - строгая типизация на уровне компиляции
- **Совместимость** - обратная совместимость версий
- **Автогенерацию** - автоматическое создание клиентского и серверного кода

В нашем проекте Protobuf используется как основа для gRPC API, обеспечивая эффективную коммуникацию между клиентами и сервером. 