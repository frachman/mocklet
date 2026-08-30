# Mocklet

## Build against the API contract first.

Mocklet is a disposable mock API service for frontend, QA, and integration work. Create a predictable HTTP endpoint in seconds, use it while the real backend is being built, then replace the base URL when the production API is ready.

Try Mocklet live at [mocklet.mikrolyt.com](https://mocklet.mikrolyt.com).

## Why Mocklet

- Unblock frontend work before the backend exists.
- Test success, error, empty, and delayed-response states.
- Share a temporary endpoint with teammates and integration tests.
- Keep the contract visible and easy to replace later.

## Status

The M2 MVP is implemented and deployed to the live URL. The service currently supports disposable mocks with up to five configured endpoints, runtime delivery, management-token protection, bounded delays and payloads, aggregate usage metrics, and crawlable SEO documentation.

Usage metrics are privacy-preserving aggregates. Mocklet does not store webhook payloads, headers, viewer tokens, full IP addresses, persistent visitor identifiers, or referrer data. The initial production baseline contains controlled smoke-test traffic; real adoption reporting will follow after the service has had time to receive organic usage.

See the [usage reporting guide](docs/operations/usage-report.md) for metric definitions and the reporting workflow.

OpenAPI 3.x contracts can also be previewed and imported through the [bounded OpenAPI workflow](docs/operations/openapi-import.md). Review generated routes before activation; unsupported features are called out explicitly.

## Quick start

Run the API locally:

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
- Delays are bounded to 10 seconds.
- Process-local runtime rate limiting is enabled.
- No user-provided code, scripting, proxying, accounts, billing, AI, or OpenAPI import exists yet.
- Every public repository file must be safe to disclose.

## Development

The API is written in Go and the web application uses Next.js. Run the repository checks with:

```bash
go test ./...
npm ci --prefix apps/web
npm run build --prefix apps/web
```

Contribution and branch-protection expectations are documented in [CONTRIBUTING.md](CONTRIBUTING.md).
