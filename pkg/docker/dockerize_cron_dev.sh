#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
COMPOSE_PATH="${PROJECT_ROOT}/docker-compose.yml"
BUILD_SCRIPT="${PROJECT_ROOT}/pkg/docker/build_image.sh"
DOCKERFILE_PATH="${PROJECT_ROOT}/pkg/docker/Dockerfile.cron"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"

function log() {
  echo "[docker:pkg-dev] $1"
}

function build_frontend() {
  log "フロントエンドイメージをビルド"
  "$BUILD_SCRIPT" "$DOCKERFILE_PATH"
}

function verify_image() {
  log "ビルド済みイメージを確認"
  if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_TAG}$"; then
    echo "[docker:pkg-dev] ${IMAGE_TAG} が見つかりません" >&2
    exit 1
  fi
}

function compose_up() {
  log "docker compose up -d を実行"
  docker compose -f "$COMPOSE_PATH" up -d
}

build_frontend
verify_image
compose_up

log "開発用デプロイが完了しました"
