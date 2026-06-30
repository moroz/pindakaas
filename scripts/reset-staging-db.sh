#!/usr/bin/env -S bash -euo pipefail

# Reset the staging database: drop it, run migrations, then apply seeds.
#
# This script ships inside the release tarball and is meant to be run on the
# server. Trigger it over SSH, e.g.:
#
#   ssh deploy@pindakaas.virtualq.run sudo /usr/local/lib/server/current/reset-staging-db.sh
#
# The current database is backed up before it is dropped, so a mistaken run can
# be recovered from the timestamped .bak file next to the database.

if [[ "${EUID}" -ne 0 ]]; then
  echo "This script must be run as root (use sudo)." >&2
  exit 1
fi

APP_BASE_DIR="/usr/local/lib/server"
APP_USER="app"
SERVICE="pindakaas"
ENV_FILE="${APP_BASE_DIR}/${SERVICE}.env"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SEEDS_FILE="${SCRIPT_DIR}/seeds.sql"

# Load DATABASE_URL and the GOOSE_* settings. The systemd EnvironmentFile is
# plain KEY="value" lines with no shell expansion, so it is safe to source.
set -a
# shellcheck disable=SC1090
source "${ENV_FILE}"
set +a

echo "==> Stopping ${SERVICE} (releases the open SQLite file)"
systemctl stop "${SERVICE}"

if [[ -f "${DATABASE_URL}" ]]; then
  BACKUP_DIR="${APP_BASE_DIR}/backups"
  BACKUP_FILE="${BACKUP_DIR}/$(basename "${DATABASE_URL}").$(date -u +%Y%m%dT%H%M%S).bak"
  echo "==> Backing up current database to ${BACKUP_FILE}"
  # Create the backup as ${APP_USER} so the directory and file are owned by it.
  # .backup checkpoints the WAL into a single consistent file.
  sudo -u "${APP_USER}" mkdir -p "${BACKUP_DIR}"
  sudo -u "${APP_USER}" sqlite3 "${DATABASE_URL}" ".backup '${BACKUP_FILE}'"
else
  echo "==> No existing database at ${DATABASE_URL}, skipping backup"
fi

echo "==> Deleting staging database at ${DATABASE_URL}"
rm -f "${DATABASE_URL}" "${DATABASE_URL}-wal" "${DATABASE_URL}-shm"

echo "==> Running migrations"
# Run as ${APP_USER} so the recreated database file is owned by the service user.
sudo -u "${APP_USER}" env \
  GOOSE_DRIVER="${GOOSE_DRIVER}" \
  GOOSE_DBSTRING="${GOOSE_DBSTRING}" \
  GOOSE_MIGRATION_DIR="${GOOSE_MIGRATION_DIR}" \
  goose up

echo "==> Applying seeds"
sudo -u "${APP_USER}" sqlite3 "${DATABASE_URL}" < "${SEEDS_FILE}"

echo "==> Starting ${SERVICE}"
systemctl start "${SERVICE}"

echo "==> Done: staging database reset, migrated, and seeded."
