# Deployment readiness

Mocklet is independently deployable as a Go API, Next.js web app, and PostgreSQL database. Production secrets must be supplied by the deployment environment; never commit `.env`.

Required configuration:

- `DATABASE_URL`
- `API_PORT`
- `NEXT_PUBLIC_API_ORIGIN` set at web build time

Minimum release gate:

1. CI passes API tests and web build.
2. Apply `db/migrations/001_init.sql` against the target Mocklet database.
3. Verify `GET /healthz` and `GET /readyz`.
4. Create a synthetic mock through the management API and verify public runtime status/body and invalid-token rejection.
5. Confirm reverse proxy limits, TLS, backups, logs, and a rollback image/tag before exposure.

The current rate limiter is process-local and the current cleanup loop runs every 15 minutes. A multi-instance deployment therefore requires an explicit operational decision before scaling horizontally. No production host, domain, or credentials are encoded here.

Current readiness: suitable for controlled testing; container build/rehearsal and visual QA pass. Not yet approved for public production exposure until backup/rollback rehearsal and target-specific ingress checks pass.
