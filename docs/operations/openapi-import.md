# OpenAPI import

Mocklet supports a bounded OpenAPI 3.x preview and import workflow.

1. `POST /api/v1/import/openapi/preview` with a YAML or JSON document.
2. Review the generated routes, response examples, and status scenarios.
3. Submit the reviewed endpoint list to `POST /api/v1/mocks/import`.

The preview uses `github.com/getkin/kin-openapi`, accepts at most 256 KiB, and generates at most five routes. Explicit response examples take precedence, followed by media-type examples, schema examples/defaults, and deterministic schema-derived values. The import creates a disposable mock and returns its management token once.

Unsupported features are rejected or surfaced rather than silently approximated: remote references, exotic composite schemas, `oneOf`/`anyOf` generation, server selection, security-scheme enforcement, and request-body schema matching. Human edits to the preview are the source of truth for activation; re-import never overwrites an existing mock.
