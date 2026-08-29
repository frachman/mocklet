# Local development

Requirements: Go 1.26+, Node.js 24+, Docker, and PostgreSQL client tools.

```bash
docker compose up -d postgres
PGPASSWORD=mocklet_dev psql -h localhost -U mocklet -d mocklet -f db/migrations/001_init.sql
DATABASE_URL='postgres://mocklet:mocklet_dev@localhost:5432/mocklet?sslmode=disable' go run ./apps/api/cmd/mocklet
```

In another terminal:

```bash
npm install --prefix apps/web
NEXT_PUBLIC_API_ORIGIN=http://localhost:8080 npm run dev --prefix apps/web
```

The API has `GET /healthz`, `POST /api/v1/mocks`, and public runtime URLs under `/m/{public_key}/{path}`. Management inspection is `GET /api/v1/mocks/{public_key}` with `X-Management-Token`.

