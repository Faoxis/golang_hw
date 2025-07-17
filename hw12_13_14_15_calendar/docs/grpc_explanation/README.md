# 📚 Документация по gRPC и архитектуре календаря

## Содержание

1. [Что такое gRPC?](./01_what_is_grpc.md)
2. [Архитектура проекта](./02_architecture.md)
3. [Protobuf спецификация](./03_protobuf.md)
4. [gRPC сервер](./04_grpc_server.md)
5. [HTTP vs gRPC](./05_http_vs_grpc.md)
6. [Тестирование](./06_testing.md)
7. [Практические примеры](./07_examples.md)

## 🎯 Цель документации

Эта документация объясняет:
- Как работает gRPC
- Как устроена архитектура календаря
- Как HTTP и gRPC API взаимодействуют
- Как тестировать API
- Практические примеры использования

## 🚀 Быстрый старт

### Запуск сервера
```bash
# Запуск окружения (PostgreSQL + RabbitMQ)
docker compose -f env/docker-compose.yml up -d

# Сборка и запуск сервера
go build ./cmd/calendar/
./calendar
```

### Тестирование HTTP API
```bash
# Создание события
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{"title":"Test","start_time":"2025-07-15T12:00:00Z","duration":"1h0m0s","user_id":"user1"}'

# Получение событий
curl "http://localhost:8080/events?day=2025-07-15"
```

### Тестирование gRPC API
```bash
# Установка grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Создание события
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"title":"Test","start_time":"2025-07-15T12:00:00Z","duration":"3600s","user_id":"user1"}' \
  localhost:50051 calendar.EventService/CreateEvent
```

## 📋 Требования

- Go 1.23+
- Docker и Docker Compose
- grpcurl (для тестирования gRPC)

## 🔗 Полезные ссылки

- [gRPC официальная документация](https://grpc.io/docs/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [grpcurl](https://github.com/fullstorydev/grpcurl) 