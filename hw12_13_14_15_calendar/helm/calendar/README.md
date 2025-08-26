# Calendar Helm Chart

Helm chart для развертывания приложения Calendar в Kubernetes.

## Установка

```bash
# Установка с дефолтными значениями
helm install calendar ./helm/calendar

# Установка с кастомными значениями
helm install calendar ./helm/calendar -f values-custom.yaml

# Обновление
helm upgrade calendar ./helm/calendar

# Удаление
helm uninstall calendar
```

## Конфигурация

Основные параметры конфигурации в `values.yaml`:

- `global.namespace` - namespace для развертывания
- `api.replicas` - количество реплик API сервиса
- `postgres.enabled` - включить/выключить PostgreSQL
- `rabbitmq.enabled` - включить/выключить RabbitMQ
- `ingress.enabled` - включить/выключить Ingress

## Компоненты

- **calendar-api** - основной API сервис
- **calendar-scheduler** - планировщик уведомлений
- **calendar-sender** - отправитель уведомлений
- **postgres** - база данных PostgreSQL
- **rabbitmq** - брокер сообщений RabbitMQ
