# auth

Auth service with registration, authorization, and session management.

## О проекте

Это микросервис аутентификации, который обеспечивает:

- **Регистрацию пользователей** — создание новых аккаунтов с email и паролем
- **Авторизацию** — безопасный вход в систему с валидацией учётных данных
- **Управление сессиями** — создание и поддержка пользовательских сессий

Сервис предназначен для использования в составе распределённой архитектуры, где аутентификация вынесена в отдельный сервис. Это позволяет:
- Централизовать логику авторизации
- Упростить интеграцию новых сервисов (все работают с единым auth-сервисом)
- Изолировать безопасность пользовательских данных
- Масштабировать независимо от других бизнес-сервисов

### Технологический стек

- **Go 1.22** — основной язык
- **PostgreSQL 16** — хранение пользователей и сессий
- **FastHTTP** — высокопроизводительный HTTP-сервер
- **SQL-migrate** — управление миграциями базы данных
- **Zap** — структурированное логирование

Команды

- `server` — start auth server
- `migrate up` — apply migrations
- `migrate down` — rollback migrations

## Endpoints

- `POST /api/v1/users` — registration
- `POST /api/v1/sessions` — login

## Config

| Env | Description |
|-----|-------------|
| AUTH_API_PORT | Server port |
| AUTH_API_MAXREQUESTBODYSIZE | Max request size |
| AUTH_API_TTL | Request TTL |
| AUTH_DB_HOST | Database host |
| AUTH_DB_PORT | Database port |
| AUTH_DB_USER | Database user |
| AUTH_DB_PASSWORD | Database password |
| AUTH_DB_DATABASENAME | Database name |
| AUTH_DB_SSL | SSL mode |
| AUTH_DEBUG | Debug mode |
