#!/usr/bin/env bash
set -Eeuo pipefail

config_file=${MOCKLET_REPORT_CONFIG:-/etc/mocklet/report.env}
# shellcheck disable=SC1090
source "$config_file"

from_date=${1:-$(date -u -d '6 days ago' +%F)}
to_date=${2:-$(date -u +%F)}

docker compose --env-file "$MOCKLET_ENV_FILE" -f "$MOCKLET_COMPOSE_FILE" exec -T postgres \
  psql -U mocklet -d mocklet -v from_date="$from_date" -v to_date="$to_date" <<'SQL'
\pset format aligned
\pset pager off
SELECT event_date, landing_views, mocks_created, runtime_requests,
       management_requests, rate_limited_requests,
       CASE WHEN landing_views = 0 THEN 0
            ELSE round(mocks_created::numeric / landing_views * 100, 2)
       END AS create_conversion_percent
FROM usage_daily
WHERE event_date BETWEEN :'from_date'::date AND :'to_date'::date
ORDER BY event_date;

SELECT COALESCE(sum(landing_views), 0) AS landing_views,
       COALESCE(sum(mocks_created), 0) AS mocks_created,
       COALESCE(sum(runtime_requests), 0) AS runtime_requests,
       COALESCE(sum(management_requests), 0) AS management_requests,
       COALESCE(sum(rate_limited_requests), 0) AS rate_limited_requests,
       CASE WHEN COALESCE(sum(landing_views), 0) = 0 THEN 0
            ELSE round(sum(mocks_created)::numeric / sum(landing_views) * 100, 2)
       END AS create_conversion_percent
FROM usage_daily
WHERE event_date BETWEEN :'from_date'::date AND :'to_date'::date;
SQL
