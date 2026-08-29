# Deployment readiness

Mocklet is independently deployable as a Go API, Next.js web app, and PostgreSQL database. Production secrets must be supplied by the deployment environment; never commit `.env`.

Required configuration:

- `DATABASE_URL`
- `API_PORT`
- `NEXT_PUBLIC_API_ORIGIN` set at web build time
- `MOCKLET_BIND_ADDRESS` and port variables should bind staging/production services to loopback behind a reverse proxy.

Minimum release gate:

1. CI passes API tests and web build.
2. Apply `db/migrations/001_init.sql` against the target Mocklet database.
3. Verify `GET /healthz` and `GET /readyz`.
4. Create a synthetic mock through the management API and verify public runtime status/body and invalid-token rejection.
5. Confirm reverse proxy limits, TLS, backups, logs, and a rollback image/tag before exposure.

The current rate limiter is process-local and the current cleanup loop runs every 15 minutes. A multi-instance deployment therefore requires an explicit operational decision before scaling horizontally. No production host, domain, or credentials are encoded here.

Current readiness: suitable for controlled testing; container build/rehearsal, M2 HTTP coverage, dependency audit, and visual QA pass. The intended production target is `mocklet.mikrolyt.com` on `root@192.129.240.188`, with Caddy routing `/api/*` and `/m/*` to the private API and the web root to the private web container. It is not yet approved for public exposure until DNS, encrypted off-host backup delivery, rollback rehearsal, and target-specific ingress checks pass.

Homelab staging is verified on `192.168.1.5` with a separate Compose project and loopback-only ports. The staging backup/restore rehearsal passed using synthetic data. Production backup must still use the approved encrypted off-host destination and remain outside Git/chat.
