#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

if [[ $# -lt 1 || -z "${1:-}" ]]; then
  echo "[docker:build:frontend] DOCKERFILE_PATH 引数が必要です" >&2
  echo "使い方: $0 <Dockerfileパス>" >&2
  exit 1
fi

DOCKERFILE_PATH="$1"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"

function format_command_preview() {
  local preview="docker"
  for arg in "$@"; do
    if [[ "$arg" =~ [^A-Za-z0-9._/:=-] ]]; then
      local escaped="${arg//\\/\\\\}"
      escaped="${escaped//\"/\\\"}"
      preview+=" \"$escaped\""
    else
      preview+=" $arg"
    fi
  done
  printf '%s' "$preview"
}

if [[ ! -f "$DOCKERFILE_PATH" ]]; then
  echo "[docker:build:frontend] Dockerfileが見つかりません: $DOCKERFILE_PATH" >&2
  exit 1
fi

DOCKER_ARGS=(build -f "$DOCKERFILE_PATH" -t "$IMAGE_TAG" "$PROJECT_ROOT")

preview=$(format_command_preview "${DOCKER_ARGS[@]}")
echo "[docker:build:frontend] $preview"

docker "${DOCKER_ARGS[@]}"
