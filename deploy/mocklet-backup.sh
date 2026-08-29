#!/usr/bin/env bash
set -Eeuo pipefail

config_file=${MOCKLET_BACKUP_CONFIG:-/etc/mocklet/backup.env}
[[ ${EUID} -eq 0 ]] || { echo 'error: backup must run as root' >&2; exit 1; }
# shellcheck disable=SC1090
source "$config_file"

for key in MOCKLET_COMPOSE_FILE MOCKLET_ENV_FILE MOCKLET_BACKUP_DIR MOCKLET_GPG_HOME MOCKLET_GPG_RECIPIENT MOCKLET_SFTP_KEY MOCKLET_SFTP_KNOWN_HOSTS MOCKLET_SFTP_USER MOCKLET_SFTP_HOST MOCKLET_SFTP_REMOTE_DIR MOCKLET_RETENTION_DAYS; do
  [[ -n ${!key:-} ]] || { echo "error: missing $key" >&2; exit 1; }
done

for command_name in docker gpg sftp sha256sum find; do
  command -v "$command_name" >/dev/null || { echo "error: required command unavailable: $command_name" >&2; exit 1; }
done

install -d -o root -g root -m 0700 "$MOCKLET_BACKUP_DIR"
tmp_dir=$(mktemp -d "$MOCKLET_BACKUP_DIR/.tmp.XXXXXX")
cleanup() { rm -rf "$tmp_dir"; }
trap cleanup EXIT

backup_id="mocklet-$(date -u +%Y%m%dT%H%M%SZ)"
artifact_name="${backup_id}.dump.gpg"
artifact_path="$tmp_dir/$artifact_name"
checksum_path="$tmp_dir/$artifact_name.sha256"

docker compose --env-file "$MOCKLET_ENV_FILE" -f "$MOCKLET_COMPOSE_FILE" exec -T postgres pg_dump -U mocklet -d mocklet -Fc \
  | gpg --homedir "$MOCKLET_GPG_HOME" --batch --yes --trust-model always \
      --recipient "$MOCKLET_GPG_RECIPIENT" --output "$artifact_path" --encrypt
(cd "$tmp_dir" && sha256sum "$artifact_name" > "$checksum_path")

sftp -q -b - -i "$MOCKLET_SFTP_KEY" -o BatchMode=yes -o IdentitiesOnly=yes \
  -o StrictHostKeyChecking=yes -o UserKnownHostsFile="$MOCKLET_SFTP_KNOWN_HOSTS" \
  "$MOCKLET_SFTP_USER@$MOCKLET_SFTP_HOST" <<SFTP
cd $MOCKLET_SFTP_REMOTE_DIR
put $artifact_path $artifact_name.part
put $checksum_path $artifact_name.sha256.part
rename $artifact_name.part $artifact_name
rename $artifact_name.sha256.part $artifact_name.sha256
SFTP

mv "$artifact_path" "$MOCKLET_BACKUP_DIR/$artifact_name"
mv "$checksum_path" "$MOCKLET_BACKUP_DIR/$artifact_name.sha256"
chmod 0600 "$MOCKLET_BACKUP_DIR/$artifact_name" "$MOCKLET_BACKUP_DIR/$artifact_name.sha256"
find "$MOCKLET_BACKUP_DIR" -maxdepth 1 -type f -name 'mocklet-*.dump.gpg' -mtime "+$MOCKLET_RETENTION_DAYS" -print0 \
  | while IFS= read -r -d '' file_path; do rm -f -- "$file_path" "${file_path}.sha256"; done

printf '%s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > "$MOCKLET_BACKUP_DIR/last-success"
chmod 0600 "$MOCKLET_BACKUP_DIR/last-success"
echo "backup succeeded: $artifact_name"
