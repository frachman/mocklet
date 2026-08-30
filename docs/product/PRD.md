# Mocklet Product Requirements

The canonical product requirements are maintained in the project handoff and PRD supplied for this build. This repository copy records the active scope: anonymous-first disposable mocks, independent deployability, Go + Next.js + PostgreSQL, and staged delivery from M1 to M4.

M1 proves create → configure one endpoint → persist → call public URL → receive response → authenticate management access. M2 adds multiple endpoints, editing/deletion, headers, delay, route conflicts, clearer UI, privacy-safe aggregate telemetry, and SEO essentials.

## M3 — Scenarios and failure simulation

An endpoint can have named response scenarios for deliberate success, empty, error, unauthorized, and delayed states. The selector contract is stable:

- `X-Mocklet-Scenario: <name>` selects a scenario.
- `?__scenario=<name>` is the query-string alternative.
- No selector uses the endpoint default response.
- An unknown scenario safely falls back to the default response.

Scenarios are managed with authenticated endpoint-scoped CRUD routes. The runtime remains deterministic and performs no AI or user-code execution.

## M4 — Bounded OpenAPI import

M4 accepts OpenAPI 3.x YAML/JSON, previews generated routes and scenarios, and requires human review before activation. It supports documented response status codes, examples, local references, and deterministic schema-derived fallback payloads. Composite generation, remote references, `oneOf`/`anyOf`, server handling, security enforcement, and request-body schema matching are explicitly unsupported.
