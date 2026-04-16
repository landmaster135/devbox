#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
DOCKERFILE_PATH="${SCRIPT_DIR}/Dockerfile.ci"

GO_VERSION="${GO_VERSION:-1.25}"
COV_FILE="${COV_FILE:-coverage.out}"
IMAGE_TAG="${CI_TEST_IMAGE_TAG:-devbox-ci-go${GO_VERSION}}"

function show_help() {
  cat <<EOF
[INFO] Dockerfile.ci でテスト用コンテナを構築し、コンテナ内で go test を実行します。

使用方法:
  ./scripts/ops/ci/ci_test.sh [options]

options:
  --go-version VERSION  Goバージョン (default: ${GO_VERSION})
  --cov-file PATH       coverprofile 出力先 (default: ${COV_FILE})
  --image-tag TAG       Docker image tag (default: ${IMAGE_TAG})
  --help                ヘルプ表示
EOF
}

function check_requirements() {
  if ! command -v docker >/dev/null 2>&1; then
    echo "docker コマンドが見つかりません。" >&2
    exit 1
  fi
}

function parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --help)
        show_help
        exit 0
        ;;
      --go-version=*)
        GO_VERSION="${1#*=}"
        shift
        ;;
      --go-version)
        GO_VERSION="${2:-}"
        shift 2
        ;;
      --cov-file=*)
        COV_FILE="${1#*=}"
        shift
        ;;
      --cov-file)
        COV_FILE="${2:-}"
        shift 2
        ;;
      --image-tag=*)
        IMAGE_TAG="${1#*=}"
        shift
        ;;
      --image-tag)
        IMAGE_TAG="${2:-}"
        shift 2
        ;;
      *)
        echo "unknown option: $1" >&2
        show_help
        exit 1
        ;;
    esac
  done

  if [[ -z "${GO_VERSION}" || -z "${COV_FILE}" || -z "${IMAGE_TAG}" ]]; then
    echo "empty option value is not allowed" >&2
    exit 1
  fi
}

function build_image() {
  docker build \
    --build-arg "GO_VERSION=${GO_VERSION}" \
    --file "${DOCKERFILE_PATH}" \
    --tag "${IMAGE_TAG}" \
    "${SCRIPT_DIR}"
}

function run_tests() {
  local uid gid
  local container_cmd
  local home_dir tmp_dir xdg_config_dir
  local -a env_flags
  local -a mount_flags
  local -a passthrough_vars
  local -a path_vars
  local path_var path_value path_dir
  declare -A mounted_dirs

  uid="$(id -u)"
  gid="$(id -g)"
  home_dir="${ROOT_DIR}"
  tmp_dir="${ROOT_DIR}/tmp/ci-test"
  xdg_config_dir="${ROOT_DIR}/tmp/xdg-config"

  env_flags=(
    -e "HOME=${home_dir}"
    -e "TMPDIR=${tmp_dir}"
    -e "XDG_CONFIG_HOME=${xdg_config_dir}"
    -e "GOTELEMETRY=off"
    -e "PATH=/usr/local/go/bin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    -e "GOCACHE=${ROOT_DIR}/.cache/go-build"
    -e "GOMODCACHE=${ROOT_DIR}/.cache/gomod"
  )
  mount_flags=(
    -v "${ROOT_DIR}:${ROOT_DIR}"
  )
  mounted_dirs["${ROOT_DIR}"]=1

  passthrough_vars=(
    GOOGLE_APPLICATION_CREDENTIALS
    CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE
    GOOGLE_GHA_CREDS_PATH
    GOOGLE_CLOUD_PROJECT
    GCLOUD_PROJECT
    CLOUDSDK_CORE_PROJECT
    GCLOUD_PROJECT_NUMBER
    GCLOUD_POOL_ID
    GCLOUD_PROVIDER_ID
    GCLOUD_SERVICE_ACCOUNT_EMAIL
  )
  path_vars=(
    GOOGLE_APPLICATION_CREDENTIALS
    CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE
    GOOGLE_GHA_CREDS_PATH
  )

  for path_var in "${path_vars[@]}"; do
    path_value="${!path_var:-}"
    if [[ -n "${path_value}" ]]; then
      path_dir="$(dirname "${path_value}")"
      # ROOT_DIR は既に丸ごと mount 済みなので、配下パスの重複 mount を避ける
      if [[ "${path_dir}" == "${ROOT_DIR}" || "${path_dir}" == "${ROOT_DIR}/"* ]]; then
        continue
      fi
      if [[ -d "${path_dir}" && -z "${mounted_dirs["${path_dir}"]:-}" ]]; then
        mount_flags+=(-v "${path_dir}:${path_dir}:ro")
        mounted_dirs["${path_dir}"]=1
      fi
    fi
  done

  for path_var in "${passthrough_vars[@]}"; do
    if [[ -n "${!path_var:-}" ]]; then
      env_flags+=(-e "${path_var}")
    fi
  done

  printf -v container_cmd \
    'mkdir -p %q %q %q %q %q && go test -v ./... -coverpkg=./... -covermode=count -coverprofile=%q' \
    "${home_dir}" \
    "${tmp_dir}" \
    "${xdg_config_dir}" \
    "${ROOT_DIR}/.cache/go-build" \
    "${ROOT_DIR}/.cache/gomod" \
    "${COV_FILE}"

  docker run --rm \
    --user "${uid}:${gid}" \
    "${env_flags[@]}" \
    "${mount_flags[@]}" \
    --workdir "${ROOT_DIR}" \
    "${IMAGE_TAG}" \
    bash -c "${container_cmd}"
}

function main() {
  parse_args "$@"
  check_requirements
  build_image
  run_tests
}

main "$@"
