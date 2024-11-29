#!/bin/bash

if ! docker network inspect usdt_network &> /dev/null; then
  echo "Сеть usdt_network не существует. Создаем..."
  docker network create usdt_network
fi

docker compose up -d