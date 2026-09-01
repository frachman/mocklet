# Routes and responses

Each route defines an HTTP method, path, status code, response body, content
type, optional response headers, and optional delay.

## Methods and paths

Supported methods are `GET`, `POST`, `PUT`, `PATCH`, and `DELETE`. Paths must be
absolute, such as `/users/{id}`. A path segment enclosed in braces matches one
segment at runtime; other segments match literally. The first matching route
for the request method is used.

An anonymous mock can contain up to five routes. A duplicate method and path is
rejected instead of silently replacing the existing route.

## Response fields

```json
{
  "method": "GET",
  "path": "/users/{id}",
  "status_code": 200,
  "content_type": "application/json",
  "headers": {"X-Fixture": "mocklet"},
  "body": "{\"id\":\"demo\"}",
  "delay_ms": 0
}
```

Status codes must be between 100 and 599. Supported content types are JSON,
plain text, XML, form-encoded, and binary content. Response bodies are limited
to 1 MiB and delays to 10 seconds.

Use `PUT` to update a route and `DELETE` to remove it. Management operations
require the mock's management token.
