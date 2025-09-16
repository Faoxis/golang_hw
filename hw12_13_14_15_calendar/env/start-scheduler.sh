#!/bin/sh

# Заполняем шаблон конфигурации из переменных окружения
cat > scheduler_config.yaml << EOF
logger:
  level: ${LOG_LEVEL:-debug}

storage:
  type: database
  host: ${POSTGRES_HOST:-localhost}
  port: ${POSTGRES_PORT:-5432}
  user: ${POSTGRES_USER:-calendar_user}
  password: ${POSTGRES_PASSWORD:-calendar_pass}
  database: ${POSTGRES_DB:-calendar}
  sslmode: ${POSTGRES_SSLMODE:-disable}

rabbit:
  url: amqp://${RABBITMQ_USER:-calendar_user}:${RABBITMQ_PASS:-calendar_pass}@${RABBITMQ_HOST:-localhost}:${RABBITMQ_PORT:-5672}/
  username: ${RABBITMQ_USER:-calendar_user}
  password: ${RABBITMQ_PASS:-calendar_pass}

event-queue:
  name: ${EVENT_QUEUE_NAME:-events}
  exchange: ${EVENT_QUEUE_EXCHANGE:-events}

scheduler:
  check_interval: ${SCHEDULER_CHECK_INTERVAL:-10s}
EOF

# Запускаем приложение
exec ./calendar_scheduler -config ./scheduler_config.yaml
