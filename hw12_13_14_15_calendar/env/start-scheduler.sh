#!/bin/sh

# Заполняем шаблон конфигурации из переменных окружения
envsubst < scheduler_config.template.yaml > scheduler_config.yaml

# Запускаем приложение
exec ./calendar_scheduler -config ./scheduler_config.yaml
