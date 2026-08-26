# auth

Auth service with registration, authorization, and session management.

## Run

```bash
go run .
```

## Commands

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
