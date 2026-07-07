#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
BUILD_SCRIPT="${PROJECT_ROOT}/pkg/docker/build_image.sh"
DOCKERFILE_PATH="${PROJECT_ROOT}/pkg/docker/Dockerfile.http"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"
TAR_PATH="${PROJECT_ROOT}/${IMAGE_NAME}-image.tar"
ZIP_PATH="${PROJECT_ROOT}/${IMAGE_NAME}-image.zip"

function log() {
  echo "[docker:pkg] $1"
}

function build_frontend() {
  log "フロントエンドイメージをビルド"
  "$BUILD_SCRIPT" "$DOCKERFILE_PATH"
}

function verify_image() {
  log "ビルド済みイメージを確認"
  if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${IMAGE_TAG}$"; then
    echo "[docker:pkg] ${IMAGE_TAG} が見つかりません" >&2
    exit 1
  fi
}

function verify_pg_dump() {
  log "イメージ内の pg_dump を確認"
  docker run --rm --entrypoint /bin/sh "$IMAGE_TAG" -c 'which pg_dump && pg_dump --version'
}

function archive_image() {
  log "既存のアーカイブをクリーンアップ"
  rm -f "$TAR_PATH" "$ZIP_PATH"

  log "Dockerイメージを保存: $TAR_PATH"
  docker save -o "$TAR_PATH" "$IMAGE_TAG"
  ls -lh "$TAR_PATH"

  local owner_user="${USER:-$(id -un)}"
  local owner_group="${GROUP:-$(id -gn)}"
  if command -v sudo >/dev/null 2>&1; then
    log "sudo chown ${owner_user}:${owner_group} を実行"
    sudo chown "${owner_user}:${owner_group}" "$TAR_PATH"
  else
    log "sudo が見つからないため chown をスキップ"
  fi

  log "アーカイブをzip化: $ZIP_PATH"
  (
    cd "$PROJECT_ROOT"
    zip -q "${ZIP_PATH}" "$(basename "$TAR_PATH")"
  )
  ls -lh "$ZIP_PATH"
}

build_frontend
verify_image
verify_pg_dump
archive_image

log "パッケージ用イメージの作成が完了しました"
