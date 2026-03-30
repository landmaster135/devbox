#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
CI_TEST_SCRIPT="${SCRIPT_DIR}/ci_test.sh"

LOG_FILE="${ROOT_DIR}/.agents/tmp/go-test.log"
GO_VERSION="${GO_VERSION:-1.25}"
COV_FILE="${COV_FILE:-coverage.out}"
IMAGE_TAG="${CI_TEST_IMAGE_TAG:-devbox-ci-go${GO_VERSION}}"
CONTEXT_LINES=10
FAILURE_PATTERN='(^--- FAIL: )|(^FAIL$)|(^FAIL[[:space:]])|(^panic: )|(\[build failed\])'

function show_help() {
  cat <<EOF
[INFO] ci_test.sh の実行ログを保存し、失敗時に要約と前後文脈を表示します。

使用方法:
  ./scripts/ops/ci/ci_test_with_logging.sh [options]

options:
  --log-file PATH       ログ保存先 (default: ${LOG_FILE})
  --go-version VERSION  Goバージョン (default: ${GO_VERSION})
  --cov-file PATH       coverprofile 出力先 (default: ${COV_FILE})
  --image-tag TAG       Docker image tag (default: ${IMAGE_TAG})
  --context-lines N     一致行の前後表示行数 (default: ${CONTEXT_LINES})
  --help                ヘルプ表示
EOF
}

function parse_args() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --help)
        show_help
        exit 0
        ;;
      --log-file=*)
        LOG_FILE="${1#*=}"
        shift
        ;;
      --log-file)
        LOG_FILE="${2:-}"
        shift 2
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
      --context-lines=*)
        CONTEXT_LINES="${1#*=}"
        shift
        ;;
      --context-lines)
        CONTEXT_LINES="${2:-}"
        shift 2
        ;;
      *)
        echo "unknown option: $1" >&2
        show_help
        exit 1
        ;;
    esac
  done

  if [[ -z "${LOG_FILE}" || -z "${GO_VERSION}" || -z "${COV_FILE}" || -z "${IMAGE_TAG}" ]]; then
    echo "empty option value is not allowed" >&2
    exit 1
  fi

  if ! [[ "${CONTEXT_LINES}" =~ ^[0-9]+$ ]]; then
    echo "--context-lines must be a non-negative integer" >&2
    exit 1
  fi
}

function ensure_log_path() {
  mkdir -p "$(dirname "${LOG_FILE}")"
}

function print_marker() {
  local label="$1"
  echo "==================== ${label} ===================="
}

function print_failure_index() {
  echo "[failure-tests]"
  grep -nE '^--- FAIL: ' "${LOG_FILE}" || echo "none"

  echo "[failure-packages]"
  grep -nE '^FAIL([[:space:]]+[^[:space:]]+)?([[:space:]]+\[build failed\])?$' "${LOG_FILE}" || echo "none"

  echo "[panic-lines]"
  grep -nE '^panic: ' "${LOG_FILE}" || echo "none"

  echo "[build-failed-lines]"
  grep -nE '\[build failed\]' "${LOG_FILE}" || echo "none"
}

function print_failure_contexts() {
  local total_lines=0
  local line_no=0
  local start_line=0
  local end_line=0
  local range_start=-1
  local range_end=-1
  local match_lines=()

  total_lines="$(wc -l < "${LOG_FILE}")"
  mapfile -t match_lines < <(grep -nE "${FAILURE_PATTERN}" "${LOG_FILE}" | cut -d: -f1 || true)
  if [[ ${#match_lines[@]} -eq 0 ]]; then
    echo "no lines matched failure pattern"
    return
  fi

  for line_no in "${match_lines[@]}"; do
    start_line=$((line_no - CONTEXT_LINES))
    end_line=$((line_no + CONTEXT_LINES))
    if ((start_line < 1)); then
      start_line=1
    fi
    if ((end_line > total_lines)); then
      end_line=${total_lines}
    fi

    if ((range_start == -1)); then
      range_start=${start_line}
      range_end=${end_line}
      continue
    fi

    if ((start_line <= range_end + 1)); then
      if ((end_line > range_end)); then
        range_end=${end_line}
      fi
      continue
    fi

    echo "--- context ${range_start}-${range_end} ---"
    sed -n "${range_start},${range_end}p" "${LOG_FILE}" | nl -ba -v "${range_start}" -w 1 -s ': '
    range_start=${start_line}
    range_end=${end_line}
  done

  if ((range_start != -1)); then
    echo "--- context ${range_start}-${range_end} ---"
    sed -n "${range_start},${range_end}p" "${LOG_FILE}" | nl -ba -v "${range_start}" -w 1 -s ': '
  fi
}

function run_ci_test() {
  local test_status=0
  local first_failure_line=""

  ensure_log_path

  set +e
  bash "${CI_TEST_SCRIPT}" \
    --go-version="${GO_VERSION}" \
    --cov-file="${COV_FILE}" \
    --image-tag="${IMAGE_TAG}" 2>&1 | tee "${LOG_FILE}"
  test_status=${PIPESTATUS[0]}
  set -e

  if [[ "${test_status}" -ne 0 ]]; then
    first_failure_line="$(grep -m 1 -E "${FAILURE_PATTERN}" "${LOG_FILE}" || true)"
    if [[ -n "${first_failure_line}" ]]; then
      echo "::error title=go test failed::${first_failure_line}"
    else
      echo "::error title=go test failed::command exited with status ${test_status}"
    fi

    print_marker "FAILURE SUMMARY START"
    print_marker "FAILURE INDEX START"
    print_failure_index
    print_marker "FAILURE INDEX END"

    echo "::group::go test failed summary"
    grep -E "${FAILURE_PATTERN}" "${LOG_FILE}" || true
    echo "::endgroup::"
    print_marker "FAILURE SUMMARY END"

    print_marker "FAILURE CONTEXT START"
    echo "::group::go test failure contexts (+/- ${CONTEXT_LINES} lines)"
    print_failure_contexts
    echo "::endgroup::"
    print_marker "FAILURE CONTEXT END"

    echo "full test log: ${LOG_FILE}"
    exit "${test_status}"
  fi

  echo "full test log: ${LOG_FILE}"
}

function main() {
  parse_args "$@"
  run_ci_test
}

main "$@"
