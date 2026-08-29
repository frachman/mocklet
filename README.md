# Mocklet

Build against the API contract before the API is ready.

Mocklet is an anonymous-first, disposable mock API service. The first vertical slice creates one mock API with one configured endpoint, persists it in PostgreSQL, serves it from a public URL, and protects management access with a private token.

## Status

M0/M1 implementation is in progress. The current API supports one endpoint per mock. Multi-endpoint management, scenarios, and OpenAPI import are later milestones.

## Local API

```bash
cp .env.example .env
docker compose up -d postgres
psql "$DATABASE_URL" -f db/migrations/001_init.sql
go run ./apps/api/cmd/mocklet
```

Create a mock:

```bash
curl -X POST http://localhost:8080/api/v1/mocks \
  -H 'Content-Type: application/json' \
  -d '{"name":"Checkout","method":"GET","path":"/users","status_code":200,"body":"{\"ok\":true}"}'
```

Call `http://localhost:8080/m/{public_key}/users`. Keep `management_token` private and send it as `X-Management-Token` when reading management details.

## Boundaries

- Anonymous resources expire after 24 hours.
- Request bodies are limited to 64 KiB for management input; response bodies to 1 MiB.
- Delay is bounded to 10 seconds.
- No user-provided code, scripting, proxying, accounts, billing, AI, or OpenAPI import exists yet.
- Every public repository file must be safe to disclose.

