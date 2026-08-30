CREATE TABLE IF NOT EXISTS usage_daily (
  event_date DATE PRIMARY KEY DEFAULT CURRENT_DATE,
  landing_views BIGINT NOT NULL DEFAULT 0,
  mocks_created BIGINT NOT NULL DEFAULT 0,
  runtime_requests BIGINT NOT NULL DEFAULT 0,
  management_requests BIGINT NOT NULL DEFAULT 0,
  rate_limited_requests BIGINT NOT NULL DEFAULT 0
);
