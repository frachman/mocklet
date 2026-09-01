# Quick start

Create a mock through the web app or the management API, then call its runtime
URL.

## Create a mock

```bash
curl -sS -X POST 'https://mocklet.mikrolyt.com/api/v1/mocks' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Checkout",
    "method": "GET",
    "path": "/checkout/session",
    "status_code": 200,
    "body": "{\"status\":\"ready\"}",
    "content_type": "application/json"
  }'
```

The response contains a `public_key`, a one-time `management_token`, an
`expires_at` timestamp, and the first endpoint. Save the management token
securely; it is not a runtime credential.

## Call the runtime

Use the returned public key as the base URL:

```bash
curl 'https://mocklet.mikrolyt.com/m/PUBLIC_KEY/checkout/session'
```

Change or add routes through the authenticated management API, then repeat the
client request against the same runtime URL.
