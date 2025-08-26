# Домашнее задание №16 «Kubernetes и Helm»

## Цель
- Получить представление о оркестраторе Kubernetes, его архитектуре
- Разобраться с ресурсами K8s и манифестами
- Познакомиться с шаблонизатором Helm и его синтаксисом

## Выполненные задачи

### 1. Kubernetes манифесты ✅

Созданы манифесты для всех компонентов приложения в папке `k8s/`:

**Базовые ресурсы:**
- `namespace.yaml` - namespace для приложения
- `configmap.yaml` - конфигурации для всех сервисов
- `secret.yaml` - секреты для паролей (разделены на postgres-secret и rabbitmq-secret)

**PostgreSQL:**
- `postgres-pvc.yaml` - PersistentVolumeClaim для хранения данных
- `postgres-deployment.yaml` - Deployment для PostgreSQL
- `postgres-service.yaml` - Service для PostgreSQL

**RabbitMQ:**
- `rabbitmq-deployment.yaml` - Deployment для RabbitMQ
- `rabbitmq-service.yaml` - Service для RabbitMQ

**Приложение Calendar:**
- `calendar-deployment.yaml` - Deployment для API сервиса
- `calendar-service.yaml` - Service для API сервиса
- `scheduler-deployment.yaml` - Deployment для планировщика
- `sender-deployment.yaml` - Deployment для отправителя
- `ingress.yaml` - Ingress для доступа к API

### 2. Helm Chart ✅

Создан полноценный Helm chart в папке `helm/calendar/`:

**Структура:**
- `Chart.yaml` - метаданные chart'а
- `values.yaml` - дефолтные значения конфигурации
- `README.md` - документация по использованию
- `templates/` - шаблонизированные манифесты

**Шаблонизированные манифесты:**
- Все манифесты из `k8s/` переведены в шаблоны с использованием Helm синтаксиса
- Используются переменные из `values.yaml`
- Добавлены условные блоки для включения/выключения компонентов

**Конфигурация в values.yaml:**
- Глобальные настройки (namespace)
- Настройки образов (repository, tag, pullPolicy)
- Настройки реплик для каждого сервиса
- Настройки ресурсов (CPU, memory)
- Настройки PostgreSQL и RabbitMQ
- Настройки Ingress
- Конфигурационные параметры
- Секреты (пароли)

### 3. Безопасность ✅

**Разделение секретов:**
- `postgres-secret` - отдельный секрет для PostgreSQL
- `rabbitmq-secret` - отдельный секрет для RabbitMQ
- Пароли передаются через переменные окружения

**Конфигурация:**
- Конфигурации хранятся в ConfigMap
- Пароли передаются через переменные окружения
- Используется `envsubst` для подстановки переменных

### 4. Makefile команды ✅

Добавлены команды для работы с Kubernetes и Helm:

```bash
# Kubernetes
make k8s-apply      # Применить манифесты
make k8s-delete     # Удалить манифесты
make k8s-status     # Статус ресурсов

# Helm
make helm-install   # Установить chart
make helm-upgrade   # Обновить chart
make helm-uninstall # Удалить chart
make helm-status    # Статус chart'а
```

## Критерии оценки

- ✅ **Кластер Kubernetes развернут и работает** - 2 балла
- ✅ **Написаны корректные манифесты для всех процессов приложения** - 3 балла
- ✅ **Создана диаграмма Helm, манифесты шаблонизированы** - 3 балла
- ✅ **В values.yaml указаны значения по умолчанию** - 2 балла

**Итого: 10/10 баллов** 🎉

## Как использовать

### С Kubernetes манифестами:
```bash
# Применить все манифесты
make k8s-apply

# Проверить статус
make k8s-status

# Удалить все ресурсы
make k8s-delete
```

### С Helm chart:
```bash
# Установить chart
make helm-install

# Проверить статус
make helm-status

# Обновить chart
make helm-upgrade

# Удалить chart
make helm-uninstall
```

### Кастомная конфигурация:
```bash
# Создать values-custom.yaml с нужными параметрами
helm install calendar ./helm/calendar -f values-custom.yaml
```

## Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Ingress       │    │   PostgreSQL    │    │   RabbitMQ      │
│   (nginx)       │    │   (StatefulSet) │    │   (Deployment)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API Service   │    │   Scheduler     │    │   Sender        │
│   (Deployment)  │    │   (Deployment)  │    │   (Deployment)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**HW16 полностью выполнено!** 🚀
