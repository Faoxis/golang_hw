#!/bin/sh

# Заполняем шаблон конфигурации из переменных окружения
cat > config.yml << EOF
logger:
  level: ${LOG_LEVEL:-debug}

server:
  host: ${SERVER_HOST:-localhost}
  port: ${SERVER_PORT:-8888}
  grpc_port: ${GRPC_PORT:-:6523}

storage:
  type: ${STORAGE_TYPE:-database}
  host: ${POSTGRES_HOST:-localhost}
  port: ${POSTGRES_PORT:-5432}
  user: ${POSTGRES_USER:-calendar_user}
  password: ${POSTGRES_PASSWORD:-calendar_pass}
  database: ${POSTGRES_DB:-calendar}
  sslmode: ${POSTGRES_SSLMODE:-disable}
EOF

# Запускаем приложение
exec ./calendar -config ./config.yml
