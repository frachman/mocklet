# Current state

- Milestone: M1 implementation.
- Repository: empty workspace bootstrapped locally; remote verification/push pending authorization and connectivity.
- Working: Go repository/API source, SQL schema, public-repo safety files, documentation harness, token hashing and endpoint validation tests.
- Verified: PostgreSQL migration, create → public runtime → private management smoke test (201/200/401), Go tests, and Next.js production build.
- Known limitation: current runtime selects the single endpoint stored for a mock; M2 multi-endpoint management is not implemented.
- Environment note: sibling `../mikrolyt-ecosystem` exists but is empty and outside this session's writable root, so it was not initialized or modified.
