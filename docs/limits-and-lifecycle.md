# Limits and lifecycle

Mocklet is intentionally disposable.

| Resource or setting | Current behavior |
| --- | --- |
| Resource lifetime | 24 hours |
| Routes per anonymous mock | 5 |
| Management request body | 64 KiB |
| OpenAPI document/import body | 256 KiB |
| Response body | 1 MiB |
| Artificial delay | 0–10,000 ms |
| Runtime rate limit | 120 requests per minute per tracked client |

When a mock expires, its runtime and management operations are no longer
available. Create a new mock for a new test session. These values describe the
current service limits and may change as the product evolves.
