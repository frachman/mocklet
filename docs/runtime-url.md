# Use the runtime URL

The runtime URL has this shape:

```text
https://mocklet.mikrolyt.com/m/{public_key}/{route-path}
```

The `public_key` identifies the disposable mock, but it does not authorize
management operations. It is suitable as a development or test API base URL;
do not treat it as a production service or a secret.

```javascript
const response = await fetch(
  `${process.env.MOCKLET_URL}/checkout/session`,
);
const checkout = await response.json();
```

Mocklet returns the configured status, headers, content type, and body. A
configured delay is applied before the response. Unknown routes return `404`,
and an expired mock is no longer available.

## Scenarios

An endpoint may have named scenarios. Select one with the `X-Mocklet-Scenario`
header or the `__scenario` query parameter. If neither is supplied, the default
scenario is used. Scenario names select stored data; they do not execute code,
templates, or outbound requests.
