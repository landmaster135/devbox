#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../../.." && pwd)"
CI_TEST_SCRIPT="${SCRIPT_DIR}/ci_test.sh"

LOG_FILE="${ROOT_DIR}/.agents/tmp/go-test.log"
GO_VERSION="${GO_VERSION:-1.25}"
COV_FILE="${COV_FILE:-coverage.out}"
COVERAGE_REPORT_FILE="${COVERAGE_REPORT_FILE:-.agents/tmp/coverage-report.txt}"
IMAGE_TAG="${CI_TEST_IMAGE_TAG:-devbox-ci-go${GO_VERSION}}"
RUN_CONTEXT="${RUN_CONTEXT:-auto}"
CONTEXT_LINES=10
FAILURE_PATTERN='(^--- FAIL: )|(^FAIL$)|(^FAIL[[:space:]])|(^panic: )|(\[build failed\])'
LOCAL_ALLOWED_FAILURE_TESTS=(
  "TestGetClientOptions_Normal"
  "TestGetClientOptions_Normal/WithServiceAccount_Normal"
)
LOCAL_FAILED_TESTS_ALLOWED=()
LOCAL_FAILED_TESTS_BLOCKED=()

function show_help() {
  cat <<EOF
[INFO] ci_test.sh の実行ログを保存し、失敗時に要約と前後文脈を表示します。

使用方法:
  ./scripts/ops/ci/ci_test_with_logging.sh [options]

options:
  --log-file PATH       ログ保存先 (default: ${LOG_FILE})
  --go-version VERSION  Goバージョン (default: ${GO_VERSION})
  --cov-file PATH       coverprofile 出力先 (default: ${COV_FILE})
  --coverage-report-file PATH
                       go tool cover -func レポート出力先 (default: ${COVERAGE_REPORT_FILE})
  --image-tag TAG       Docker image tag (default: ${IMAGE_TAG})
  --run-context VALUE   実行コンテキスト local|github-actions|auto (default: ${RUN_CONTEXT})
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
      --run-context=*)
        RUN_CONTEXT="${1#*=}"
        shift
        ;;
      --run-context)
        RUN_CONTEXT="${2:-}"
        shift 2
        ;;
      --context-lines=*)
        CONTEXT_LINES="${1#*=}"
        shift
        ;;
      --coverage-report-file=*)
        COVERAGE_REPORT_FILE="${1#*=}"
        shift
        ;;
      --coverage-report-file)
        COVERAGE_REPORT_FILE="${2:-}"
        shift 2
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

  if [[ -z "${LOG_FILE}" || -z "${GO_VERSION}" || -z "${COV_FILE}" || -z "${COVERAGE_REPORT_FILE}" || -z "${IMAGE_TAG}" || -z "${RUN_CONTEXT}" ]]; then
    echo "empty option value is not allowed" >&2
    exit 1
  fi

  if ! [[ "${CONTEXT_LINES}" =~ ^[0-9]+$ ]]; then
    echo "--context-lines must be a non-negative integer" >&2
    exit 1
  fi
}

function resolve_run_context() {
  if [[ "${RUN_CONTEXT}" == "auto" ]]; then
    if [[ "${GITHUB_ACTIONS:-}" == "true" ]]; then
      RUN_CONTEXT="github-actions"
    else
      RUN_CONTEXT="local"
    fi
  fi

  case "${RUN_CONTEXT}" in
    local|github-actions)
      ;;
    *)
      echo "--run-context must be one of local, github-actions, auto" >&2
      exit 1
      ;;
  esac
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

function collect_failed_test_names() {
  grep -E '^[[:space:]]*--- FAIL: ' "${LOG_FILE}" \
    | sed -E 's/^[[:space:]]*--- FAIL: ([^[:space:]]+).*/\1/' \
    | sort -u
}

function is_local_allowed_failure_test() {
  local test_name="$1"
  local allowed_test=""

  for allowed_test in "${LOCAL_ALLOWED_FAILURE_TESTS[@]}"; do
    if [[ "${test_name}" == "${allowed_test}" ]]; then
      return 0
    fi
  done

  return 1
}

function classify_failed_tests_for_local() {
  local test_name=""
  local -a failed_tests=()

  mapfile -t failed_tests < <(collect_failed_test_names || true)
  for test_name in "${failed_tests[@]}"; do
    if is_local_allowed_failure_test "${test_name}"; then
      LOCAL_FAILED_TESTS_ALLOWED+=("${test_name}")
    else
      LOCAL_FAILED_TESTS_BLOCKED+=("${test_name}")
    fi
  done
}

function should_tolerate_local_failures() {
  if grep -qE '^panic: |\[build failed\]' "${LOG_FILE}"; then
    return 1
  fi

  if [[ ${#LOCAL_FAILED_TESTS_ALLOWED[@]} -eq 0 ]]; then
    return 1
  fi

  if [[ ${#LOCAL_FAILED_TESTS_BLOCKED[@]} -gt 0 ]]; then
    return 1
  fi

  return 0
}

function print_local_filter_result() {
  local result="$1"
  print_marker "LOCAL FAILURE FILTER START"
  echo "[run-context] ${RUN_CONTEXT}"
  echo "[policy] known failures are tolerated only in local runs"
  echo "[filter-result] ${result}"
  echo "[tolerated-failed-tests]"
  if [[ ${#LOCAL_FAILED_TESTS_ALLOWED[@]} -eq 0 ]]; then
    echo "none"
  else
    printf '%s\n' "${LOCAL_FAILED_TESTS_ALLOWED[@]}"
  fi
  echo "[remaining-failed-tests]"
  if [[ ${#LOCAL_FAILED_TESTS_BLOCKED[@]} -eq 0 ]]; then
    echo "none"
  else
    printf '%s\n' "${LOCAL_FAILED_TESTS_BLOCKED[@]}"
  fi
  print_marker "LOCAL FAILURE FILTER END"
}

function print_final_result() {
  local overall_result="$1"

  print_marker "FINAL RESULT START"
  echo "[overall-result] ${overall_result}"
  echo "[run-context] ${RUN_CONTEXT}"
  print_marker "FINAL RESULT END"
}

function resolve_coverage_file_path() {
  if [[ "${COV_FILE}" = /* ]]; then
    echo "${COV_FILE}"
    return
  fi
  echo "${ROOT_DIR}/${COV_FILE}"
}

function resolve_coverage_report_path() {
  if [[ "${COVERAGE_REPORT_FILE}" = /* ]]; then
    echo "${COVERAGE_REPORT_FILE}"
    return
  fi
  echo "${ROOT_DIR}/${COVERAGE_REPORT_FILE}"
}

function resolve_module_path() {
  (
    cd "${ROOT_DIR}" && go list -m -f '{{.Path}}' 2>/dev/null
  ) || true
}

function build_module_coverage_profile() {
  local cov_path="$1"
  local module_path="$2"
  local out_path="$3"

  {
    head -n 1 "${cov_path}"
    rg "^${module_path}/" "${cov_path}" || true
  } > "${out_path}"
}

function generate_coverage_report() {
  local cov_path="$1"
  local report_path="$2"
  local module_path=""
  local source_cov_path=""
  local filtered_cov_path=""

  mkdir -p "$(dirname "${report_path}")"
  module_path="$(resolve_module_path)"
  source_cov_path="${cov_path}"

  if [[ -n "${module_path}" ]]; then
    filtered_cov_path="$(mktemp)"
    build_module_coverage_profile "${cov_path}" "${module_path}" "${filtered_cov_path}"
    if [[ "$(wc -l < "${filtered_cov_path}")" -gt 1 ]]; then
      source_cov_path="${filtered_cov_path}"
    fi
  fi

  if ! go tool cover -func="${source_cov_path}" > "${report_path}" 2>/dev/null; then
    if [[ -n "${filtered_cov_path}" ]]; then
      rm -f "${filtered_cov_path}"
    fi
    return 1
  fi

  if [[ -n "${filtered_cov_path}" ]]; then
    rm -f "${filtered_cov_path}"
  fi
}

function print_coverage_total() {
  local cov_path=""
  local report_path=""
  local total_percent=""

  cov_path="$(resolve_coverage_file_path)"
  report_path="$(resolve_coverage_report_path)"
  print_marker "COVERAGE TOTAL START"
  echo "[run-context] ${RUN_CONTEXT}"
  echo "[coverage-file] ${cov_path}"
  echo "[coverage-report-file] ${report_path}"

  if [[ ! -f "${cov_path}" ]]; then
    echo "[coverage-total] unavailable file not found"
    print_marker "COVERAGE TOTAL END"
    return
  fi

  if generate_coverage_report "${cov_path}" "${report_path}"; then
    total_percent="$(awk '/^total:/{print $NF}' "${report_path}" || true)"
  else
    total_percent=""
  fi

  if [[ -z "${total_percent}" ]]; then
    echo "[coverage-total] unavailable failed to parse total"
  else
    echo "[coverage-total] ${total_percent}"
  fi
  print_marker "COVERAGE TOTAL END"
}

function filter_console_output() {
  awk '
    /^[[:space:]]*coverage:[[:space:]][0-9.]+% of statements in \.\/\.\.\.$/ { next }
    /^[[:space:]]*ok[[:space:]].*coverage:[[:space:]][0-9.]+% of statements in \.\/\.\.\.$/ { next }
    { print }
  '
}

function run_ci_test() {
  local test_status=0
  local first_failure_line=""

  ensure_log_path

  set +e
  bash "${CI_TEST_SCRIPT}" \
    --go-version="${GO_VERSION}" \
    --cov-file="${COV_FILE}" \
    --image-tag="${IMAGE_TAG}" 2>&1 | tee "${LOG_FILE}" | filter_console_output
  test_status=${PIPESTATUS[0]}
  set -e

  if [[ "${test_status}" -ne 0 ]]; then
    if [[ "${RUN_CONTEXT}" == "local" ]]; then
      classify_failed_tests_for_local
      if should_tolerate_local_failures; then
        print_local_filter_result "PASS local known failures only"
        print_coverage_total
        print_final_result "SUCCESS"
        echo "full test log: ${LOG_FILE}"
        return 0
      fi
      print_local_filter_result "FAIL unknown failures remain"
    fi

    if [[ "${RUN_CONTEXT}" == "local" && ${#LOCAL_FAILED_TESTS_BLOCKED[@]} -gt 0 ]]; then
      first_failure_line="--- FAIL: ${LOCAL_FAILED_TESTS_BLOCKED[0]}"
    else
      first_failure_line="$(grep -m 1 -E "${FAILURE_PATTERN}" "${LOG_FILE}" || true)"
    fi

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

    print_coverage_total
    print_final_result "FAILURE"
    echo "full test log: ${LOG_FILE}"
    exit "${test_status}"
  fi

  print_coverage_total
  print_final_result "SUCCESS"
  echo "full test log: ${LOG_FILE}"
}

function main() {
  parse_args "$@"
  resolve_run_context
  run_ci_test
}

main "$@"
