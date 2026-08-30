# Usage reporting

Mocklet records privacy-safe daily aggregates in PostgreSQL. It does not store IP addresses, user agents, tokens, request payloads, referrers, or persistent visitor identifiers.

Tracked metrics:

- `landing_views`: page-view beacons received by the homepage. This is total views, not unique people.
- `mocks_created`: successful mock creation requests.
- `runtime_requests`: requests to public mock endpoints.
- `management_requests`: authenticated management API requests.
- `rate_limited_requests`: runtime requests rejected by the process-local limiter.

Historical page views before this instrumentation are unavailable and must be reported as `not measured`, not estimated. The first meaningful conversion baseline starts after the instrumentation release is deployed.

Run the production report from the private host:

```bash
sudo /usr/local/sbin/mocklet-usage-report.sh 2026-08-30 2026-09-05
```

The report shows daily values, totals, and `mocks_created / landing_views` conversion. A visitor count cannot be derived from these aggregates; unique-user reporting requires a separate privacy and retention decision.
