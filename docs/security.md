# Security baseline

Mock definitions and request traffic are untrusted. Management tokens are returned once at creation and stored only as SHA-256 hashes. Public keys are separate from management capability. Resources expire after 24 hours. Management input, response size, and artificial delay are bounded.

The initial service does not execute user code, proxy requests, render response bodies as trusted HTML, or support permanent resources. Rate limiting and asynchronous cleanup remain known hardening work.

