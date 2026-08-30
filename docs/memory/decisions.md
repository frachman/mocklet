# Decisions

1. M1 is the active delivery boundary; M2 is required before claiming the full MVP.
2. Go standard HTTP + pgx + SQL migrations are the initial implementation choices.
3. Anonymous resources expire after 24 hours, with policy values kept easy to change.
4. Public keys identify runtime resources; management tokens are separate and hashed.
5. Runtime and aggregate telemetry events are buffered in memory and flushed periodically; request handling must not synchronously write usage data to PostgreSQL.
6. Public response content types are restricted to non-HTML allowlisted types, with browser MIME-sniffing protection enabled.
