#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DEFAULT_ENV_FILE="${PROJECT_ROOT}/env.yml"
ENV_FILE="${1:-$DEFAULT_ENV_FILE}"
DOCKERFILE_PATH="${PROJECT_ROOT}/pkg/docker/Dockerfile.cron"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"

declare -A VITE_ENV=()
declare -a VITE_KEYS=()

trim() {
  local value="$1"
  value="${value#${value%%[![:space:]]*}}"
  value="${value%${value##*[![:space:]]}}"
  printf '%s' "$value"
}

rtrim() {
  local value="$1"
  value="${value%${value##*[![:space:]]}}"
  printf '%s' "$value"
}

strip_inline_comment() {
  local value="$1"
  local length=${#value}
  local in_single=0
  local in_double=0
  local result=""
  local prev_char=""

  for ((i=0; i<length; i++)); do
    local char="${value:i:1}"
    if [[ "$char" == "'" && $in_double -eq 0 ]]; then
      result+="$char"
      if ((in_single)); then
        in_single=0
      else
        in_single=1
      fi
      prev_char="$char"
      continue
    fi

    if [[ "$char" == '"' && $in_single -eq 0 ]]; then
      result+="$char"
      if [[ "$prev_char" != "\\" ]]; then
        if ((in_double)); then
          in_double=0
        else
          in_double=1
        fi
      fi
      prev_char="$char"
      continue
    fi

    if [[ "$char" == "#" && $in_single -eq 0 && $in_double -eq 0 ]]; then
      break
    fi

    result+="$char"
    prev_char="$char"
  done

  rtrim "$result"
}

strip_quotes() {
  local value
  value="$(trim "$1")"
  local length=${#value}
  if ((length >= 2)); then
    local first="${value:0:1}"
    local last="${value: -1}"
    if [[ "$first" == "$last" && ( "$first" == '"' || "$first" == "'" ) ]]; then
      value="${value:1:$((length-2))}"
    fi
  fi
  printf '%s' "$value"
}

parse_env_file() {
  local file="$1"
  while IFS= read -r line || [[ -n "$line" ]]; do
    local trimmed
    trimmed="$(trim "$line")"
    if [[ -z "$trimmed" || "${trimmed:0:1}" == "#" ]]; then
      continue
    fi
    if [[ "$line" =~ ^[[:space:]]*([A-Za-z0-9_]+)[[:space:]]*:(.*)$ ]]; then
      local key="${BASH_REMATCH[1]}"
      local raw_value="${BASH_REMATCH[2]}"
      local without_comment
      without_comment="$(strip_inline_comment "$raw_value")"
      local value
      value="$(strip_quotes "$(trim "$without_comment")")"
      if [[ "$key" == VITE_* ]]; then
        if [[ -z "${VITE_ENV[$key]+x}" ]]; then
          VITE_KEYS+=("$key")
        fi
        VITE_ENV[$key]="$value"
      fi
    fi
  done < "$file"
}

format_command_preview() {
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

if [[ ! -f "$ENV_FILE" ]]; then
  echo "[docker:build:frontend] envファイルが見つかりません: $ENV_FILE" >&2
  exit 1
fi

parse_env_file "$ENV_FILE"

if [[ ! -f "$DOCKERFILE_PATH" ]]; then
  echo "[docker:build:frontend] Dockerfileが見つかりません: $DOCKERFILE_PATH" >&2
  exit 1
fi

DOCKER_ARGS=(build -f "$DOCKERFILE_PATH" -t "$IMAGE_TAG")
for key in "${VITE_KEYS[@]}"; do
  DOCKER_ARGS+=(--build-arg)
  DOCKER_ARGS+=("${key}=${VITE_ENV[$key]}")
done
DOCKER_ARGS+=("${PROJECT_ROOT}")

preview=$(format_command_preview "${DOCKER_ARGS[@]}")
echo "[docker:build:frontend] $preview"

docker "${DOCKER_ARGS[@]}"
