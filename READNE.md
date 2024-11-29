# Сервис курсов USDT

Сервис для получения и хранения курсов USDT с биржи Garantex.  Использует gRPC, PostgreSQL, и поддерживает graceful shutdown.

## Быстрый старт

1. **Установите зависимости:** `go mod tidy`
2. **Соберите приложение:** `make build`
3. **Запустите с помощью Docker Compose:** `docker compose up -d`
4. **Остановите:** `docker compose down`

## Настройка

Настройте подключение к базе данных через переменные окружения:

* `DB_USER` (default: `postgres`)
* `DB_PASSWORD` (default: `postgres`)
* `DB_HOST` (default: `db`)
* `DB_PORT` (default: `5432`)
* `DB_DATABASE` (default: `postgres`)

## API (gRPC)

* `/GetRates`:  Получение курса USDT.  Аргумент: `target_currency` (например, "USD").
* `/HealthCheck`: Проверка работоспособности.


## Тестирование

Запустите тесты: `make test`


## Дополнительная информация

* Логирование: `zap`
* Проверка кода: `golangci-lint`