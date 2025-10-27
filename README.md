# Subscriptions Service (Go + PostgreSQL)

Мини‑сервис для CRUD по подпискам + подсчёт суммы за период.
Запуск: `docker compose up --build`

## Быстрый старт
1. Скопируйте `.env.example` → `.env` (или задайте переменные окружения).
2. `docker compose up --build`
3. API: `http://localhost:8080/api/v1`
4. Swagger UI: `http://localhost:8081` (читает `openapi/openapi.yaml`).

## Переменные окружения
- `APP_PORT` — порт HTTP (по умолчанию 8080).
- `DATABASE_URL` — DSN для Postgres.
- `LOG_LEVEL` — debug|info|warn|error (по умолчанию info).

## Миграции
Контейнер `migrate/migrate` применит `migrations/*.sql` при старте стека.

## Форматы дат
`start_date` и `end_date` принимают **месяц и год** в формате `YYYY-MM` или `MM-YYYY`. В БД хранится дата первым числом месяца.

## Подсчёт стоимости
Суммируются **полные месяцы пересечения** подписки и заданного периода, границы включительно.# touch test
