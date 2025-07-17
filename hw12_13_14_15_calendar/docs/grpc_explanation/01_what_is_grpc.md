# 1. Что такое gRPC?

## 🤔 gRPC - что это?

**gRPC** (Google Remote Procedure Call) - это современный фреймворк для создания API, разработанный Google. Он позволяет клиентам и серверам общаться как будто они находятся в одном процессе.

## 🔄 Как работает gRPC?

### Традиционный REST API:
```
Клиент → HTTP запрос → Сервер
Клиент ← HTTP ответ ← Сервер
```

### gRPC:
```
Клиент → Вызов функции → Сервер
Клиент ← Результат функции ← Сервер
```

## 🏗️ Основные компоненты gRPC

### 1. **Protocol Buffers (Protobuf)**
- Язык описания интерфейсов (IDL)
- Бинарный формат сериализации
- Строгая типизация
- Автогенерация кода

### 2. **HTTP/2**
- Мультиплексирование (несколько запросов в одном соединении)
- Сжатие заголовков
- Server Push
- Бинарный протокол

### 3. **Строгая типизация**
- Контракт между клиентом и сервером
- Автоматическая валидация
- IDE поддержка (автодополнение, проверка типов)

## 📊 Сравнение с REST

| Аспект | REST | gRPC |
|--------|------|------|
| **Протокол** | HTTP/1.1 | HTTP/2 |
| **Формат данных** | JSON/XML | Protocol Buffers |
| **Типизация** | Слабая | Строгая |
| **Скорость** | Медленнее | Быстрее |
| **Размер данных** | Больше | Меньше |
| **Поддержка потоков** | Ограниченная | Полная |
| **Сложность** | Простая | Средняя |

## 🎯 Преимущества gRPC

### 1. **Производительность**
```protobuf
// Protobuf - компактный бинарный формат
message Event {
  string id = 1;
  string title = 2;
  int64 start_time = 3;
}
```

### 2. **Строгая типизация**
```go
// Автогенерированный код с типами
func (s *Server) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
    // req.Title - строго типизированная строка
    // req.StartTime - строго типизированное время
}
```

### 3. **Мультиплексирование**
```go
// Один HTTP/2 соединение для всех запросов
conn, _ := grpc.Dial("localhost:50051", grpc.WithInsecure())
client := api.NewEventServiceClient(conn)

// Все эти вызовы используют одно соединение
client.CreateEvent(ctx, req1)
client.GetEvent(ctx, req2)
client.ListEvents(ctx, req3)
```

### 4. **Автогенерация кода**
```bash
# Генерируем код из .proto файла
protoc --go_out=. --go-grpc_out=. EventService.proto

# Получаем готовые структуры и методы
# - EventServiceClient (для клиентов)
# - EventServiceServer (для серверов)
# - Все структуры запросов/ответов
```

## 🔧 Типы gRPC вызовов

### 1. **Unary RPC** (как обычная функция)
```protobuf
rpc CreateEvent(CreateEventRequest) returns (CreateEventResponse);
```
```go
// Клиент
response, err := client.CreateEvent(ctx, request)

// Сервер
func (s *Server) CreateEvent(ctx context.Context, req *CreateEventRequest) (*CreateEventResponse, error) {
    // Обработка
    return response, nil
}
```

### 2. **Server Streaming RPC** (сервер отправляет поток данных)
```protobuf
rpc ListEvents(ListEventsRequest) returns (stream Event);
```

### 3. **Client Streaming RPC** (клиент отправляет поток данных)
```protobuf
rpc UploadEvents(stream Event) returns (UploadResponse);
```

### 4. **Bidirectional Streaming RPC** (двусторонний поток)
```protobuf
rpc Chat(stream Message) returns (stream Message);
```

## 🚀 Почему gRPC в нашем проекте?

### 1. **Микросервисная архитектура**
- Строгие контракты между сервисами
- Автогенерация клиентов
- Эффективная коммуникация

### 2. **Производительность**
- Бинарный протокол быстрее JSON
- HTTP/2 мультиплексирование
- Сжатие данных

### 3. **Типобезопасность**
- Ошибки на этапе компиляции
- Автодополнение в IDE
- Валидация данных

## 📝 Пример из нашего проекта

### Protobuf определение:
```protobuf
service EventService {
  rpc CreateEvent(CreateEventRequest) returns (CreateEventResponse);
  rpc GetEvent(GetEventRequest) returns (GetEventResponse);
  rpc ListEvents(ListEventsRequest) returns (ListEventsResponse);
}
```

### Автогенерированный код:
```go
// Клиент
client := api.NewEventServiceClient(conn)
response, err := client.CreateEvent(ctx, &api.CreateEventRequest{
    Title: "Meeting",
    UserId: "user123",
})

// Сервер
type Server struct {
    api.UnimplementedEventServiceServer
}

func (s *Server) CreateEvent(ctx context.Context, req *api.CreateEventRequest) (*api.CreateEventResponse, error) {
    // Бизнес-логика
    return &api.CreateEventResponse{Event: event}, nil
}
```

## 🎓 Заключение

gRPC - это современный подход к созданию API, который:
- Обеспечивает высокую производительность
- Предоставляет строгую типизацию
- Автоматически генерирует код
- Поддерживает различные типы вызовов

В нашем календаре gRPC используется параллельно с HTTP API, что позволяет клиентам выбирать наиболее подходящий протокол для своих нужд. 