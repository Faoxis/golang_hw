#!/bin/sh

# Заполняем шаблон конфигурации из переменных окружения
envsubst < config.template.yml > config.yml

# Запускаем приложение
exec ./calendar -config ./config.yml
