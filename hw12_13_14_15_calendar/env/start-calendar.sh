#!/bin/sh

# Подставляем переменные окружения в конфигурацию
envsubst < ./configs/config.yml > ./configs/config-resolved.yml

# Запускаем приложение
exec ./calendar -config ./configs/config-resolved.yml
