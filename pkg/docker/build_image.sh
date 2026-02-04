#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"
DEFAULT_ENV_FILE="${PROJECT_ROOT}/env.yml"
ENV_FILE="${1:-$DEFAULT_ENV_FILE}"
if [[ $# -lt 2 || -z "${2:-}" ]]; then
  echo "[docker:build:frontend] DOCKERFILE_PATH 引数が必要です" >&2
  echo "使い方: $0 [env.ymlパス] <Dockerfileパス> [KEY1,KEY2,...]" >&2
  exit 1
fi
DOCKERFILE_PATH="$2"
BUILD_ARG_LIST="${3:-CRON_URL_PORT}"
IMAGE_NAME="devbox-cron"
IMAGE_TAG="${IMAGE_NAME}:latest"

declare -A BUILD_ARG_VALUES=()
declare -A REQUESTED_KEYS=()
declare -a BUILD_ARG_KEYS=()

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
      if [[ -n "${REQUESTED_KEYS[$key]+x}" ]]; then
        BUILD_ARG_VALUES[$key]="$value"
      fi
    fi
  done < "$file"
}

parse_key_list() {
  local csv="$1"
  local IFS=','
  read -ra raw_keys <<< "$csv"
  for raw_key in "${raw_keys[@]}"; do
    local key
    key="$(trim "$raw_key")"
    if [[ -z "$key" ]]; then
      continue
    fi
    if [[ -z "${REQUESTED_KEYS[$key]+x}" ]]; then
      BUILD_ARG_KEYS+=("$key")
      REQUESTED_KEYS[$key]=1
    fi
  done
}

ensure_required_args() {
  local missing=()
  for key in "${BUILD_ARG_KEYS[@]}"; do
    if [[ -z "${BUILD_ARG_VALUES[$key]+x}" || -z "${BUILD_ARG_VALUES[$key]}" ]]; then
      missing+=("$key")
    fi
  done

  if ((${#missing[@]} > 0)); then
    echo "[docker:build:frontend] envファイルに以下のキーがありません: ${missing[*]}" >&2
    exit 1
  fi
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

parse_key_list "$BUILD_ARG_LIST"
parse_env_file "$ENV_FILE"
ensure_required_args

if [[ ! -f "$DOCKERFILE_PATH" ]]; then
  echo "[docker:build:frontend] Dockerfileが見つかりません: $DOCKERFILE_PATH" >&2
  exit 1
fi

DOCKER_ARGS=(build -f "$DOCKERFILE_PATH" -t "$IMAGE_TAG")
for key in "${BUILD_ARG_KEYS[@]}"; do
  DOCKER_ARGS+=(--build-arg)
  DOCKER_ARGS+=("${key}=${BUILD_ARG_VALUES[$key]}")
done
DOCKER_ARGS+=("${PROJECT_ROOT}")

preview=$(format_command_preview "${DOCKER_ARGS[@]}")
echo "[docker:build:frontend] $preview"

docker "${DOCKER_ARGS[@]}"
