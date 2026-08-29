# Current state

- Milestone: M1 complete; M2 next.
- Repository: empty workspace bootstrapped locally; remote verification/push pending authorization and connectivity.
- Working: Go repository/API source, SQL schema, public-repo safety files, documentation harness, token hashing and endpoint validation tests.
- Verified: PostgreSQL migration, create → public runtime → private management smoke test (201/200/401), Go tests, and Next.js production build.
- Browser integration: credential-free CORS and preflight handling are enabled for the separate local web origin.
- M2: multi-endpoint CRUD, five-route cap, collision handling, response headers, bounded delay, expiry cleanup, rate limiting, containers, CI, and readiness endpoint are implemented.
- Known limitation: current runtime selects the single endpoint stored for a mock; M2 multi-endpoint management is not implemented.
- Environment note: sibling `../mikrolyt-ecosystem` exists but is empty and outside this session's writable root, so it was not initialized or modified.
