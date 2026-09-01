# Security baseline

Mock definitions and request traffic are untrusted. Management tokens are returned once at creation and stored only as SHA-256 hashes. Public keys are separate from management capability. Resources expire after 24 hours. Management input, imported documents, response size, and artificial delay are bounded.

The service does not execute user code, proxy requests, render response bodies as trusted HTML, or support permanent resources. Endpoint response content types are allowlisted to JSON, plain text, XML, form-encoded, and binary content; HTML and other active text types are rejected. Public runtime responses include `X-Content-Type-Options: nosniff` and a sandbox content-security policy.

Anonymous landing-view telemetry requires the exact JSON sentinel `{"source":"landing"}`. It contains no visitor identifier, cookie, or payload data. Usage counters are buffered in memory and flushed periodically, so the public runtime path does not wait for a database write on every request.

Scenario names are data values only. They select stored deterministic responses and cannot execute code, expressions, templates, or outbound requests. Unknown scenario names fall back to the endpoint default response.

The initial browser API surface allows credential-free CORS (`Access-Control-Allow-Origin: *`); management authorization remains an explicit token header and is not cookie-based.
