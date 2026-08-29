# Current state

- Milestone: M2 MVP implementation complete; production-readiness gate next.
- Repository: empty workspace bootstrapped locally; remote verification/push pending authorization and connectivity.
- Working: Go repository/API source, SQL schema, public-repo safety files, documentation harness, token hashing and endpoint validation tests.
- Verified: PostgreSQL migration, create → public runtime → private management smoke test (201/200/401), Go tests, and Next.js production build.
- Browser integration: credential-free CORS and preflight handling are enabled for the separate local web origin.
- M2: multi-endpoint CRUD, five-route cap, collision handling, response headers, bounded delay, expiry cleanup, rate limiting, containers, CI, and readiness endpoint are implemented.
- Verified: API and web Docker images build successfully after tightening the Docker context; M2 HTTP smoke covered two routes, configured status/body/header, collision 409, and authenticated listing.
- Security review: `npm audit --omit=dev` reports zero vulnerabilities after upgrading to Next.js 16.3.3; `gitleaks` is unavailable in the local environment, so the fallback scan found only documented synthetic development credentials.
- Final QA: clean Compose stack on isolated ports passed migration, `/healthz`, `/readyz`, web HTTP 200, browser create mock, browser add endpoint, and zero new browser console errors.
- Controlled M2 testing: five routes succeeded; sixth route and duplicate route returned 409; update returned runtime 299; delete returned 204 and runtime 404; rate burst produced 429; expired resource returned 404.
- Homelab staging: source deployed to an isolated private staging host; loopback-only ports, host health/readiness, create/runtime/management, browser tunnel QA, and synthetic PostgreSQL backup/restore rehearsal passed.
- Production is live on the approved public domain; encrypted backup delivery, checksum, scheduled freshness check, ingress smoke, and existing-site regression passed. Isolated decrypt/restore using the private recovery key remains pending.
- Known limitation: process-local rate limiting is not suitable for horizontal scaling without a shared limiter.
- Environment note: sibling `../mikrolyt-ecosystem` exists but is empty and outside this session's writable root, so it was not initialized or modified.
