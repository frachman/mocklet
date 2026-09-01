# Management token

The management token authorizes changes and reads for one disposable mock.
It is returned once when the mock is created or imported. Mocklet stores only a
SHA-256 hash, so the original token cannot be recovered from the service.

Send it with either header when calling management endpoints:

```bash
curl -sS \
  -H "X-Management-Token: $MANAGEMENT_TOKEN" \
  'https://mocklet.mikrolyt.com/api/v1/mocks/PUBLIC_KEY'
```

The `Authorization: Bearer ...` form is also accepted. Keep the token in a
local secret store or environment variable. Never put it in a browser bundle,
runtime URL, logs, screenshots, or source control.

The public key and management token have different roles: the public key
selects runtime traffic, while the management token grants access to the mock's
configuration.
