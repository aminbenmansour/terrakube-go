#!/usr/bin/env bash
#
# Development services for Terrakube Go using plain podman (no compose
# provider required). For compose users, deploy/compose/compose.yaml is
# equivalent.
#
# Usage: scripts/dev-services.sh [up|down|status|wipe]

set -euo pipefail

PODMAN=${PODMAN:-podman}
NETWORK=${TERRAKUBE_NETWORK:-terrakube-go}

PG_NAME=${PG_NAME:-terrakube-postgres}
REDIS_NAME=${REDIS_NAME:-terrakube-redis}
MINIO_NAME=${MINIO_NAME:-terrakube-minio}

PG_VOLUME=${PG_VOLUME:-terrakube-go-postgres}
REDIS_VOLUME=${REDIS_VOLUME:-terrakube-go-redis}
MINIO_VOLUME=${MINIO_VOLUME:-terrakube-go-minio}

PG_PORT=${PG_PORT:-5432}
REDIS_PORT=${REDIS_PORT:-6379}
MINIO_PORT=${MINIO_PORT:-9000}
MINIO_CONSOLE_PORT=${MINIO_CONSOLE_PORT:-9001}

ensure_network() {
  if ! "$PODMAN" network exists "$NETWORK" >/dev/null 2>&1; then
    "$PODMAN" network create "$NETWORK" >/dev/null
  fi
}

ensure_volume() {
  local volume=$1
  if ! "$PODMAN" volume exists "$volume" >/dev/null 2>&1; then
    "$PODMAN" volume create "$volume" >/dev/null
  fi
}

container_exists() {
  "$PODMAN" container exists "$1" >/dev/null 2>&1
}

up() {
  ensure_network

  ensure_volume "$PG_VOLUME"
  if container_exists "$PG_NAME"; then
    "$PODMAN" start "$PG_NAME" >/dev/null
  else
    "$PODMAN" run -d \
      --name "$PG_NAME" \
      --network "$NETWORK" \
      -e POSTGRES_USER=terrakube \
      -e POSTGRES_HOST_AUTH_METHOD=trust \
      -e POSTGRES_DB=terrakube \
      -p "${PG_PORT}:5432" \
      -v "${PG_VOLUME}:/var/lib/postgresql/data" \
      docker.io/library/postgres:16-alpine >/dev/null
  fi

  ensure_volume "$REDIS_VOLUME"
  if container_exists "$REDIS_NAME"; then
    "$PODMAN" start "$REDIS_NAME" >/dev/null
  else
    "$PODMAN" run -d \
      --name "$REDIS_NAME" \
      --network "$NETWORK" \
      -p "${REDIS_PORT}:6379" \
      -v "${REDIS_VOLUME}:/data" \
      docker.io/library/redis:7-alpine redis-server --appendonly yes >/dev/null
  fi

  ensure_volume "$MINIO_VOLUME"
  if container_exists "$MINIO_NAME"; then
    "$PODMAN" start "$MINIO_NAME" >/dev/null
  else
    "$PODMAN" run -d \
      --name "$MINIO_NAME" \
      --network "$NETWORK" \
      -p "${MINIO_PORT}:9000" \
      -p "${MINIO_CONSOLE_PORT}:9001" \
      -v "${MINIO_VOLUME}:/data" \
      docker.io/pgsty/minio server /data --console-address ":9001" >/dev/null
  fi

  echo "Waiting for Postgres..."
  for _ in $(seq 1 30); do
    if "$PODMAN" exec "$PG_NAME" pg_isready -U terrakube -d terrakube >/dev/null 2>&1; then
      break
    fi
    sleep 1
  done

  echo
  echo "Postgres : postgres://terrakube@localhost:${PG_PORT}/terrakube?sslmode=disable"
  echo "Redis    : redis://localhost:${REDIS_PORT}/0"
  echo "MinIO    : http://localhost:${MINIO_PORT} (console http://localhost:${MINIO_CONSOLE_PORT}; credentials in container logs)"
  echo
  echo "Next: make db-migrate && make run-api"
}

down() {
  for name in "$MINIO_NAME" "$REDIS_NAME" "$PG_NAME"; do
    if container_exists "$name"; then
      "$PODMAN" stop "$name" >/dev/null
    fi
  done
  echo "Services stopped. Volumes were kept; use 'make services-wipe' to delete them."
}

status() {
  "$PODMAN" ps -a --filter "name=terrakube-" --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}'
}

wipe() {
  for name in "$MINIO_NAME" "$REDIS_NAME" "$PG_NAME"; do
    if container_exists "$name"; then
      "$PODMAN" rm -f "$name" >/dev/null
    fi
  done
  for volume in "$MINIO_VOLUME" "$REDIS_VOLUME" "$PG_VOLUME"; do
    "$PODMAN" volume rm -f "$volume" >/dev/null 2>&1 || true
  done
  echo "Services and volumes removed."
}

case "${1:-up}" in
  up) up ;;
  down) down ;;
  status) status ;;
  wipe) wipe ;;
  *)
    echo "usage: $0 [up|down|status|wipe]" >&2
    exit 2
    ;;
esac
