# Deployment readiness

Mocklet is independently deployable as a Go API, Next.js web app, and PostgreSQL database. Production secrets must be supplied by the deployment environment; never commit `.env`.

Required configuration:

- `DATABASE_URL`
- `API_PORT`
- `NEXT_PUBLIC_API_ORIGIN` set at web build time
- `MOCKLET_BIND_ADDRESS` and port variables should bind staging/production services to loopback behind a reverse proxy.

Minimum release gate:

1. CI passes API tests and web build.
2. Apply migrations in order, including `db/migrations/003_scenarios.sql`, against the target Mocklet database.
3. Verify `GET /healthz` and `GET /readyz` from the private container/network context; these endpoints are not exposed through the public ingress.
4. Create a synthetic mock through the management API and verify public runtime status/body and invalid-token rejection.
5. Confirm reverse proxy limits, TLS, backups, logs, and a rollback image/tag before exposure.

The current rate limiter is process-local and the current cleanup loop runs every 15 minutes. A multi-instance deployment therefore requires an explicit operational decision before scaling horizontally. No production host, domain, or credentials are encoded here.

Current state: deployed to `mocklet.mikrolyt.com` on the approved production host. Caddy routes `/api/*` and `/m/*` to the private API and the web root to the private web container. Immutable image, HTTPS, API/runtime smoke, scenario CRUD/runtime, OpenAPI preview, no-public-API-or-DB-port, encrypted off-host delivery, checksum, and scheduled-backup checks pass. The isolated decrypt/restore rehearsal remains pending because the private restore key was not available on the audited homelab or production config path; do not represent disaster-recovery readiness as complete until that rehearsal is recorded.

Homelab staging is verified on the approved private staging host with a separate Compose project and loopback-only ports. The staging backup/restore rehearsal passed using synthetic data. Production backup must still use the approved encrypted off-host destination and remain outside Git/chat.
