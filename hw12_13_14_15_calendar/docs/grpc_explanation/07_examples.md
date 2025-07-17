# 7. Практические примеры

## 🚀 Быстрый старт

### 1. **Запуск окружения**
```bash
# Запуск PostgreSQL и RabbitMQ
docker compose -f env/docker-compose.yml up -d

# Проверка статуса
docker compose -f env/docker-compose.yml ps
```

### 2. **Запуск сервера**
```bash
# Сборка и запуск
go build ./cmd/calendar/
./calendar

# Или напрямую
go run ./cmd/calendar/main.go
```

### 3. **Проверка работы**
```bash
# Проверка HTTP API
curl http://localhost:8080/health

# Проверка gRPC API
grpcurl -plaintext localhost:50051 list
```

## 📝 Примеры HTTP API

### 1. **Создание события**
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync meeting",
    "user_id": "john.doe",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "1h0m0s",
    "notify_before": "15m0s"
  }'
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly team sync meeting",
  "user_id": "john.doe",
  "start_time": "2025-07-15T14:00:00Z",
  "duration": "1h0m0s",
  "notify_before": "15m0s"
}
```

### 2. **Получение события по ID**
```bash
curl "http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000"
```

**Ответ:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Team Meeting",
  "description": "Weekly team sync meeting",
  "user_id": "john.doe",
  "start_time": "2025-07-15T14:00:00Z",
  "duration": "1h0m0s",
  "notify_before": "15m0s"
}
```

### 3. **Обновление события**
```bash
curl -X PUT http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Team Meeting",
    "description": "Updated weekly team sync",
    "user_id": "john.doe",
    "start_time": "2025-07-15T15:00:00Z",
    "duration": "2h0m0s",
    "notify_before": "30m0s"
  }'
```

### 4. **Получение событий за день**
```bash
curl "http://localhost:8080/events?day=2025-07-15"
```

**Ответ:**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly team sync meeting",
    "user_id": "john.doe",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "1h0m0s",
    "notify_before": "15m0s"
  },
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "title": "Client Call",
    "description": "Discussion with client",
    "user_id": "john.doe",
    "start_time": "2025-07-15T16:00:00Z",
    "duration": "30m0s",
    "notify_before": "5m0s"
  }
]
```

### 5. **Получение событий за неделю**
```bash
curl "http://localhost:8080/events?week=2025-07-15"
```

### 6. **Получение событий за месяц**
```bash
curl "http://localhost:8080/events?month=2025-07"
```

### 7. **Удаление события**
```bash
curl -X DELETE http://localhost:8080/events/550e8400-e29b-41d4-a716-446655440000
```

## 🔧 Примеры gRPC API

### 1. **Создание события**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{
    "title": "Team Meeting",
    "description": "Weekly team sync meeting",
    "user_id": "john.doe",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "3600s",
    "notify_before": "900s"
  }' \
  localhost:50051 calendar.EventService/CreateEvent
```

**Ответ:**
```json
{
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Team Meeting",
    "description": "Weekly team sync meeting",
    "user_id": "john.doe",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "3600s",
    "notify_before": "900s"
  }
}
```

### 2. **Получение события по ID**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 calendar.EventService/GetEvent
```

### 3. **Обновление события**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Updated Team Meeting",
    "description": "Updated weekly team sync",
    "user_id": "john.doe",
    "start_time": "2025-07-15T15:00:00Z",
    "duration": "7200s",
    "notify_before": "1800s"
  }' \
  localhost:50051 calendar.EventService/UpdateEvent
```

### 4. **Получение событий за день**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"date": "2025-07-15T00:00:00Z"}' \
  localhost:50051 calendar.EventService/ListEventsForDay
```

**Ответ:**
```json
{
  "events": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Team Meeting",
      "description": "Weekly team sync meeting",
      "user_id": "john.doe",
      "start_time": "2025-07-15T14:00:00Z",
      "duration": "3600s",
      "notify_before": "900s"
    }
  ]
}
```

### 5. **Получение событий за неделю**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"date": "2025-07-15T00:00:00Z"}' \
  localhost:50051 calendar.EventService/ListEventsForWeek
```

### 6. **Получение событий за месяц**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"date": "2025-07-01T00:00:00Z"}' \
  localhost:50051 calendar.EventService/ListEventsForMonth
```

### 7. **Удаление события**
```bash
grpcurl -plaintext -proto api/EventService.proto \
  -d '{"id": "550e8400-e29b-41d4-a716-446655440000"}' \
  localhost:50051 calendar.EventService/DeleteEvent
```

## 🐍 Python клиент для gRPC

### 1. **Установка зависимостей**
```bash
pip install grpcio grpcio-tools
```

### 2. **Генерация Python кода**
```bash
python -m grpc_tools.protoc \
  --python_out=./python_client \
  --grpc_python_out=./python_client \
  --proto_path=./api \
  EventService.proto
```

### 3. **Python клиент**
```python
import grpc
import api.EventService_pb2 as pb2
import api.EventService_pb2_grpc as pb2_grpc
from google.protobuf.timestamp_pb2 import Timestamp
from google.protobuf.duration_pb2 import Duration
import datetime

def create_event():
    # Подключение к серверу
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = pb2_grpc.EventServiceStub(channel)
        
        # Создание запроса
        start_time = Timestamp()
        start_time.FromDatetime(datetime.datetime(2025, 7, 15, 14, 0, 0))
        
        duration = Duration()
        duration.FromTimedelta(datetime.timedelta(hours=1))
        
        notify_before = Duration()
        notify_before.FromTimedelta(datetime.timedelta(minutes=15))
        
        request = pb2.CreateEventRequest(
            title="Python Team Meeting",
            description="Meeting created from Python",
            user_id="python.user",
            start_time=start_time,
            duration=duration,
            notify_before=notify_before
        )
        
        # Вызов gRPC метода
        response = stub.CreateEvent(request)
        
        print(f"Created event: {response.event.title}")
        print(f"Event ID: {response.event.id}")

def get_event(event_id):
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = pb2_grpc.EventServiceStub(channel)
        
        request = pb2.GetEventRequest(id=event_id)
        response = stub.GetEvent(request)
        
        print(f"Event: {response.event.title}")
        print(f"Description: {response.event.description}")

if __name__ == "__main__":
    create_event()
    get_event("550e8400-e29b-41d4-a716-446655440000")
```

## 🔄 JavaScript клиент для HTTP

### 1. **Node.js клиент**
```javascript
const axios = require('axios');

const API_BASE = 'http://localhost:8080';

async function createEvent() {
    try {
        const response = await axios.post(`${API_BASE}/events`, {
            title: 'JavaScript Team Meeting',
            description: 'Meeting created from JavaScript',
            user_id: 'js.user',
            start_time: '2025-07-15T14:00:00Z',
            duration: '1h0m0s',
            notify_before: '15m0s'
        });
        
        console.log('Created event:', response.data);
        return response.data.id;
    } catch (error) {
        console.error('Error creating event:', error.response?.data);
    }
}

async function getEvent(eventId) {
    try {
        const response = await axios.get(`${API_BASE}/events/${eventId}`);
        console.log('Event:', response.data);
    } catch (error) {
        console.error('Error getting event:', error.response?.data);
    }
}

async function listEventsForDay(date) {
    try {
        const response = await axios.get(`${API_BASE}/events?day=${date}`);
        console.log('Events for day:', response.data);
    } catch (error) {
        console.error('Error listing events:', error.response?.data);
    }
}

async function main() {
    const eventId = await createEvent();
    if (eventId) {
        await getEvent(eventId);
        await listEventsForDay('2025-07-15');
    }
}

main();
```

### 2. **Browser JavaScript**
```html
<!DOCTYPE html>
<html>
<head>
    <title>Calendar API Test</title>
</head>
<body>
    <h1>Calendar API Test</h1>
    <button onclick="createEvent()">Create Event</button>
    <button onclick="listEvents()">List Events</button>
    <div id="result"></div>

    <script>
        const API_BASE = 'http://localhost:8080';

        async function createEvent() {
            try {
                const response = await fetch(`${API_BASE}/events`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        title: 'Browser Team Meeting',
                        description: 'Meeting created from browser',
                        user_id: 'browser.user',
                        start_time: '2025-07-15T14:00:00Z',
                        duration: '1h0m0s',
                        notify_before: '15m0s'
                    })
                });
                
                const data = await response.json();
                document.getElementById('result').innerHTML = 
                    `<pre>${JSON.stringify(data, null, 2)}</pre>`;
            } catch (error) {
                console.error('Error:', error);
            }
        }

        async function listEvents() {
            try {
                const response = await fetch(`${API_BASE}/events?day=2025-07-15`);
                const data = await response.json();
                document.getElementById('result').innerHTML = 
                    `<pre>${JSON.stringify(data, null, 2)}</pre>`;
            } catch (error) {
                console.error('Error:', error);
            }
        }
    </script>
</body>
</html>
```

## 🧪 Тестирование сценариев

### 1. **Создание нескольких событий**
```bash
#!/bin/bash

# Создаем несколько событий
echo "Creating events..."

# Утренняя встреча
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Morning Standup",
    "description": "Daily team standup",
    "user_id": "team.lead",
    "start_time": "2025-07-15T09:00:00Z",
    "duration": "30m0s",
    "notify_before": "5m0s"
  }'

# Обеденный перерыв
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Lunch Break",
    "description": "Team lunch",
    "user_id": "team.lead",
    "start_time": "2025-07-15T12:00:00Z",
    "duration": "1h0m0s",
    "notify_before": "0s"
  }'

# Встреча с клиентом
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Client Meeting",
    "description": "Discussion with client about requirements",
    "user_id": "team.lead",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "1h30m0s",
    "notify_before": "15m0s"
  }'

echo "Events created. Listing events for day:"
curl "http://localhost:8080/events?day=2025-07-15"
```

### 2. **Тестирование ошибок**
```bash
#!/bin/bash

echo "Testing validation errors..."

# Пустой заголовок
echo "1. Testing empty title:"
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Test description",
    "user_id": "test.user",
    "start_time": "2025-07-15T14:00:00Z",
    "duration": "1h0m0s",
    "notify_before": "15m0s"
  }'

echo -e "\n\n2. Testing invalid time format:"
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Event",
    "user_id": "test.user",
    "start_time": "invalid-time",
    "duration": "1h0m0s",
    "notify_before": "15m0s"
  }'

echo -e "\n\n3. Testing non-existent event:"
curl "http://localhost:8080/events/non-existent-id"
```

## 📊 Мониторинг и логи

### 1. **Просмотр логов сервера**
```bash
# Запуск сервера с подробными логами
LOG_LEVEL=debug go run ./cmd/calendar/main.go
```

### 2. **Проверка базы данных**
```bash
# Подключение к PostgreSQL
docker exec -it hw12_13_14_15_calendar-postgres-1 psql -U calendar_user -d calendar

# Просмотр событий
SELECT * FROM events;

# Подсчет событий по пользователям
SELECT user_id, COUNT(*) FROM events GROUP BY user_id;
```

### 3. **Мониторинг производительности**
```bash
# Тест производительности HTTP API
ab -n 1000 -c 10 -T application/json -p event.json http://localhost:8080/events

# Тест производительности gRPC API
grpcurl -plaintext -proto api/EventService.proto \
  -d @event.json \
  localhost:50051 calendar.EventService/CreateEvent
```

## 🎯 Полезные команды

### 1. **Очистка данных**
```bash
# Очистка базы данных
docker exec -it hw12_13_14_15_calendar-postgres-1 psql -U calendar_user -d calendar -c "DELETE FROM events;"
```

### 2. **Проверка здоровья сервиса**
```bash
# HTTP health check
curl http://localhost:8080/health

# gRPC health check (если реализован)
grpcurl -plaintext localhost:50051 grpc.health.v1.Health/Check
```

### 3. **Экспорт событий**
```bash
# Экспорт в JSON
curl "http://localhost:8080/events?day=2025-07-15" | jq '.' > events.json

# Экспорт в CSV (если реализован)
curl "http://localhost:8080/events?day=2025-07-15&format=csv" > events.csv
```

## 🎓 Заключение

Эти примеры демонстрируют:

### **HTTP API:**
- ✅ Простота использования с curl
- ✅ Универсальность для веб-приложений
- ✅ Человекочитаемые запросы и ответы
- ✅ Легкая интеграция с JavaScript

### **gRPC API:**
- ✅ Строгая типизация
- ✅ Высокая производительность
- ✅ Автогенерация клиентов
- ✅ Эффективная сериализация

### **Общие возможности:**
- ✅ Полный CRUD для событий
- ✅ Фильтрация по времени
- ✅ Валидация данных
- ✅ Обработка ошибок
- ✅ Логирование операций

Это позволяет разработчикам выбирать наиболее подходящий API для своих нужд и интегрировать календарь в различные приложения. 