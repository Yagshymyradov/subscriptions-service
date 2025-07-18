# Subscriptions service

REST-сервис для агрегации данных об он-лайн подписках пользователей

## Стек
- Go 1.22+
- PostgreSQL 16 (Docker)
- chi - HTTP-роутер
- zap - структурное логирование
- golang-migrate - миграции
- swag - авто-генерация Swagger-документации

---

## Быстрый старт (локально)

```bash
# 1. Клонируем репозиторий
git clone https://github.com/Yagshymyradov/subscriptions-service.git
cd subscriptions-service

# 2. Поднимаем базу
docker compose up -d db

# 3. Применяем миграции (если нужен CLI migrate → brew install golang-migrate)
migrate -path migrations \
        -database "postgres://postgres:postgres@localhost:5678/subscriptions?sslmode=disable" \
        up

# 4. Запускаем приложение
APP_DB_DSN="postgres://postgres:postgres@localhost:5678/subscriptions?sslmode=disable" \
go run ./cmd/subscriptions-service
```

После запуска сервис слушает `http://localhost:8080`.

---

## Конфигурация

| Переменная          | Описание                          | Значение по умолчанию |
| ------------------- | --------------------------------- | --------------------- |
| `APP_HTTP_PORT`     | порт HTTP-сервера                 | `8080`                |
| `APP_DB_DSN`        | строка подключения к PostgreSQL   | – (обязательно)       |
| `APP_LOG_LEVEL`     | `debug` / `info` / `warn` / …     | `info`                |

Пример файла `.env` (не коммитить):

```env
APP_DB_DSN=postgres://postgres:postgres@localhost:5678/subscriptions?sslmode=disable
APP_HTTP_PORT=8080
APP_LOG_LEVEL=debug
```

---

## API

Полная интерактивная спецификация доступна по `/swagger/index.html`.

Кратко:

| Метод | Путь                            | Описание                              |
| ----- | ------------------------------  | ------------------------------------- |
| POST  | `/subscriptions`               | создать подписку                      |
| GET   | `/subscriptions/{id}`          | получить подписку по ID               |
| GET   | `/subscriptions?userID=...`    | список подписок пользователя          |
| PUT   | `/subscriptions/{id}`          | обновить подписку                     |
| DELETE| `/subscriptions/{id}`          | удалить подписку                      |
| GET   | `/subscriptions/total`         | суммарная стоимость (фильтры — ниже)  |

### /subscriptions/total

GET /subscriptions/total?userID=<uuid>&month=<1-12>&year=<YYYY>[&serviceFilter=<подстрока>]

Возвращает JSON: `{ "total": 800 }`

---

## Тестирование

1. Unit-тесты репозитория — `go test ./internal/repository/...`  
2. HTTP-хендлеры — `httptest` + мок-репозиторий (`go test ./internal/handlers/...`)

---

## TODO / Возможные улучшения

- Валидация входных JSON (`go-playground/validator`)  
- Graceful-shutdown (`srv.Shutdown(ctx)`)  
- CI workflow (lint, test, swag init)  
- Dockerfile + docker-compose с самим сервисом