# Current state

- Milestone: M2 MVP implementation complete; production-readiness gate next.
- Repository: empty workspace bootstrapped locally; remote verification/push pending authorization and connectivity.
- Working: Go repository/API source, SQL schema, public-repo safety files, documentation harness, token hashing and endpoint validation tests.
- Verified: PostgreSQL migration, create → public runtime → private management smoke test (201/200/401), Go tests, and Next.js production build.
- Browser integration: credential-free CORS and preflight handling are enabled for the separate local web origin.
- M2: multi-endpoint CRUD, five-route cap, collision handling, response headers, bounded delay, expiry cleanup, rate limiting, containers, CI, and readiness endpoint are implemented.
- Known limitation: process-local rate limiting is not suitable for horizontal scaling without a shared limiter; production target, backup rehearsal, proxy limits, and visual browser QA remain unverified.
- Environment note: sibling `../mikrolyt-ecosystem` exists but is empty and outside this session's writable root, so it was not initialized or modified.
