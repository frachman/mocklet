# Local development

Requirements: Go 1.26+, Node.js 24+, Docker, and PostgreSQL client tools.

```bash
docker compose up -d postgres
PGPASSWORD=mocklet_dev psql -h localhost -U mocklet -d mocklet -f db/migrations/001_init.sql
DATABASE_URL='postgres://mocklet:mocklet_dev@localhost:5432/mocklet?sslmode=disable' go run ./apps/api/cmd/mocklet
```

For a full container rehearsal, apply the migration and start the API/web services:

```bash
docker compose up -d --build
docker compose exec -T postgres psql -U mocklet -d mocklet < db/migrations/001_init.sql
```

If the default ports are occupied, use an isolated mapping, for example `MOCKLET_DB_PORT=55432 MOCKLET_API_PORT=18080 MOCKLET_WEB_PORT=13000 NEXT_PUBLIC_API_ORIGIN=http://localhost:18080 docker compose up -d --build`.

In another terminal:

```bash
npm install --prefix apps/web
NEXT_PUBLIC_API_ORIGIN=http://localhost:8080 npm run dev --prefix apps/web
```

The API has `GET /healthz`, `POST /api/v1/mocks`, and public runtime URLs under `/m/{public_key}/{path}`. Management inspection is `GET /api/v1/mocks/{public_key}` with `X-Management-Token`.
