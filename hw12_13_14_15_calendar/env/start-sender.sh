#!/bin/sh

# Заполняем шаблон конфигурации из переменных окружения
cat > sender_config.yaml << EOF
logger:
  level: ${LOG_LEVEL:-debug}

rabbit:
  url: amqp://${RABBITMQ_USER:-calendar_user}:${RABBITMQ_PASS:-calendar_pass}@${RABBITMQ_HOST:-localhost}:${RABBITMQ_PORT:-5672}/
  username: ${RABBITMQ_USER:-calendar_user}
  password: ${RABBITMQ_PASS:-calendar_pass}

event-queue:
  name: ${EVENT_QUEUE_NAME:-events}
  exchange: ${EVENT_QUEUE_EXCHANGE:-events}
EOF

# Запускаем приложение
exec ./calendar_sender -config ./sender_config.yaml
