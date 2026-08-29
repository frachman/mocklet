CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS mock_apis (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  public_key TEXT NOT NULL UNIQUE,
  management_token_hash TEXT NOT NULL,
  owner_user_id TEXT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS mock_endpoints (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  mock_api_id UUID NOT NULL REFERENCES mock_apis(id) ON DELETE CASCADE,
  method TEXT NOT NULL,
  path TEXT NOT NULL,
  status_code INTEGER NOT NULL,
  headers JSONB NOT NULL DEFAULT '{}'::jsonb,
  body TEXT NOT NULL DEFAULT '',
  content_type TEXT NOT NULL DEFAULT 'application/json',
  delay_ms INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(mock_api_id, method, path)
);

