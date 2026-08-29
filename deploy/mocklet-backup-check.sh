#!/usr/bin/env bash
set -Eeuo pipefail

config_file=${MOCKLET_BACKUP_CONFIG:-/etc/mocklet/backup.env}
# shellcheck disable=SC1090
source "$config_file"
max_age_minutes=${MOCKLET_BACKUP_MAX_AGE_MINUTES:-1500}
latest_file=$(find "$MOCKLET_BACKUP_DIR" -maxdepth 1 -type f -name 'mocklet-*.dump.gpg' -mmin "-$max_age_minutes" -print -quit)
[[ -n ${latest_file:-} ]] || { echo "error: no fresh encrypted backup within ${max_age_minutes} minutes" >&2; exit 1; }
checksum_file="${latest_file}.sha256"
expected=$(awk 'NR == 1 {print $1}' "$checksum_file")
actual=$(sha256sum "$latest_file" | awk '{print $1}')
[[ "$expected" == "$actual" ]] || { echo 'error: checksum mismatch' >&2; exit 1; }
echo "backup check succeeded: $(basename "$latest_file")"
