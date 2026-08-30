CREATE TABLE IF NOT EXISTS mock_scenarios (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  endpoint_id UUID NOT NULL REFERENCES mock_endpoints(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  is_default BOOLEAN NOT NULL DEFAULT FALSE,
  status_code INTEGER NOT NULL,
  headers JSONB NOT NULL DEFAULT '{}'::jsonb,
  body TEXT NOT NULL DEFAULT '',
  content_type TEXT NOT NULL DEFAULT 'application/json',
  delay_ms INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(endpoint_id, name)
);

CREATE UNIQUE INDEX IF NOT EXISTS mock_scenarios_one_default
  ON mock_scenarios(endpoint_id) WHERE is_default;
