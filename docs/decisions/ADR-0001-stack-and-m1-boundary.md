# ADR-0001: Small M1 stack and boundary

## Decision

Use Go's HTTP standard library with PostgreSQL through pgx, Next.js/TypeScript for the management UI, and SQL migrations. Deliver one endpoint per mock before implementing multi-route management.

## Rationale

This keeps the request-serving path lightweight and the first proof reversible. PostgreSQL preserves ownership and expiry boundaries without introducing an ORM or extra infrastructure.

## Consequences

M1 is intentionally incomplete relative to the M2 MVP definition. Multi-endpoint behavior must be added before calling the product an MVP.

