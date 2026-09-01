#!/usr/bin/env bash
set -eu

image_tag="$1"
cd /srv/mocklet
backup="/srv/mocklet/.env.pre-${image_tag#sha-}"
cp .env "$backup"
rollback() {
  cp "$backup" .env
  docker compose --env-file .env -f docker-compose.yml up -d api web >/dev/null 2>&1 || true
}
trap rollback ERR
sed -i "s#^MOCKLET_IMAGE_TAG=.*#MOCKLET_IMAGE_TAG=$image_tag#" .env
docker compose --env-file .env -f docker-compose.yml pull api web
docker compose --env-file .env -f docker-compose.yml up -d api web
for attempt in $(seq 1 30); do
  api_state=$(docker inspect --format '{{.State.Status}}' mocklet-api-1)
  web_state=$(docker inspect --format '{{.State.Status}}' mocklet-web-1)
  if [ "$api_state" = running ] && [ "$web_state" = running ]; then
    if ! curl --fail --silent --show-error --max-time 10 https://mocklet.mikrolyt.com/ >/dev/null; then
      sleep 5
      continue
    fi
    smoke=$(mktemp)
    trap 'rm -f "$smoke"' EXIT
    curl --fail --silent --show-error --max-time 10 -H 'Content-Type: application/json' -X POST https://mocklet.mikrolyt.com/api/v1/mocks -d '{"name":"production-smoke","method":"GET","path":"/health","status_code":200,"body":"{\"ok\":true}","content_type":"application/json"}' > "$smoke"
    public_key=$(jq -er '.public_key' "$smoke")
    curl --fail --silent --show-error --max-time 10 "https://mocklet.mikrolyt.com/m/$public_key/health" | jq -e '.ok == true' >/dev/null
    trap - ERR
    echo "production deployment verified: $image_tag"
    exit 0
  fi
  sleep 5
done
echo "production deployment did not become ready" >&2
exit 1
