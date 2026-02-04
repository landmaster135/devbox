#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
COMPOSE_PATH="${PROJECT_ROOT}/docker-compose.yml"
ENV_PATH="${PROJECT_ROOT}/env.yml"
BUILD_SCRIPT="${PROJECT_ROOT}/pkg/docker/build_image.sh"
DOCKERFILE_PATH="${PROJECT_ROOT}/pkg/docker/Dockerfile.cron"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"
PORT_KEY="CRON_URL_PORT"
VOLUME_KEY="MOUNT_VOLUME"
USER_KEY="HOST_ID"
SERVICE_NAME="devbox"

log() {
  echo "[docker:pkg-dev] $1"
}

run_env_sync() {
  log "env.yml を docker-compose.yml に反映"
  (
    cd "$PROJECT_ROOT"
    go run ./cmd/cli/docker/main.go \
      --operation=env-into-compose \
      --compose-path="$COMPOSE_PATH" \
      --env-yaml-path="$ENV_PATH"
  )
}

run_port_sync() {
  log "ports 情報を docker-compose.yml に反映"
  (
    cd "$PROJECT_ROOT"
    go run ./cmd/cli/docker/main.go \
      --operation=ports-into-compose \
      --compose-path="$COMPOSE_PATH" \
      --env-yaml-path="$ENV_PATH" \
      --port-key="$PORT_KEY" \
      --service="$SERVICE_NAME"
  )
}

run_volume_sync() {
  log "volumes 情報を docker-compose.yml に反映"
  (
    cd "$PROJECT_ROOT"
    go run ./cmd/cli/docker/main.go \
      --operation=volumes-into-compose \
      --compose-path="$COMPOSE_PATH" \
      --env-yaml-path="$ENV_PATH" \
      --volume-key="$VOLUME_KEY" \
      --service="$SERVICE_NAME"
  )
}

run_user_sync() {
  log "user 情報を docker-compose.yml に反映"
  (
    cd "$PROJECT_ROOT"
    go run ./cmd/cli/docker/main.go \
      --operation=user-into-compose \
      --compose-path="$COMPOSE_PATH" \
      --env-yaml-path="$ENV_PATH" \
      --user-key="$USER_KEY" \
      --service="$SERVICE_NAME"
  )
}

build_frontend() {
  log "フロントエンドイメージをビルド"
  "$BUILD_SCRIPT" "$ENV_PATH" "$DOCKERFILE_PATH"
}

verify_image() {
  log "ビルド済みイメージを確認"
  if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_TAG}$"; then
    echo "[docker:pkg-dev] ${IMAGE_TAG} が見つかりません" >&2
    exit 1
  fi
}

compose_up() {
  log "docker compose up -d を実行"
  docker compose -f "$COMPOSE_PATH" up -d
}

run_env_sync
run_volume_sync
run_user_sync
run_port_sync
build_frontend
verify_image
compose_up

log "開発用デプロイが完了しました"
