#!/usr/bin/env bash

set -euo pipefail

LOG_FILE="/tmp/github-ci.log"
ERROR_TITLE="ci command failed"
SUMMARY_TITLE="ci failed summary"
CONTEXT_TITLE="ci failure contexts"
CONTEXT_LINES=10
FAILURE_PATTERN='(^--- FAIL: )|(^FAIL[[:space:]])|(^panic: )|(\[build failed\])'
STREAM_OUTPUT=0
COMMAND=()

function show_help() {
  cat <<'EOF'
[INFO] 任意コマンドの標準出力/標準エラーをログ保存し、失敗時に要約を表示します。

使用方法:
  ./scripts/ops/github_ci_logging.sh [options] -- <command> [args...]

options:
  --log-file PATH         ログ保存先 (default: /tmp/github-ci.log)
  --error-title TEXT      GitHub annotation title (default: ci command failed)
  --summary-title TEXT    失敗要約グループ名 (default: ci failed summary)
  --context-title TEXT    前後行表示グループ名 (default: ci failure contexts)
  --context-lines NUMBER  一致行の前後表示行数 (default: 10)
  --failure-pattern REGEX 失敗抽出用の正規表現
  --stream-output         実行中ログを逐次表示する defaultはオフ
  --help                  ヘルプ表示
EOF
}

function check_requirements() {
  local required_commands=(tee grep sed wc nl cut)
  local cmd

  for cmd in "${required_commands[@]}"; do
    if ! command -v "${cmd}" >/dev/null 2>&1; then
      echo "required command not found: ${cmd}" >&2
      exit 1
    fi
  done
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
      --error-title=*)
        ERROR_TITLE="${1#*=}"
        shift
        ;;
      --error-title)
        ERROR_TITLE="${2:-}"
        shift 2
        ;;
      --summary-title=*)
        SUMMARY_TITLE="${1#*=}"
        shift
        ;;
      --summary-title)
        SUMMARY_TITLE="${2:-}"
        shift 2
        ;;
      --context-title=*)
        CONTEXT_TITLE="${1#*=}"
        shift
        ;;
      --context-title)
        CONTEXT_TITLE="${2:-}"
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
      --failure-pattern=*)
        FAILURE_PATTERN="${1#*=}"
        shift
        ;;
      --failure-pattern)
        FAILURE_PATTERN="${2:-}"
        shift 2
        ;;
      --stream-output)
        STREAM_OUTPUT=1
        shift
        ;;
      --)
        shift
        COMMAND=("$@")
        break
        ;;
      *)
        echo "unknown option: $1" >&2
        show_help
        exit 1
        ;;
    esac
  done

  if [[ ${#COMMAND[@]} -eq 0 ]]; then
    echo "command is required. use -- <command> [args...]" >&2
    exit 1
  fi

  if [[ -z "${LOG_FILE}" || -z "${ERROR_TITLE}" || -z "${SUMMARY_TITLE}" || -z "${CONTEXT_TITLE}" || -z "${CONTEXT_LINES}" || -z "${FAILURE_PATTERN}" || -z "${STREAM_OUTPUT}" ]]; then
    echo "empty option value is not allowed" >&2
    exit 1
  fi

  if ! [[ "${CONTEXT_LINES}" =~ ^[0-9]+$ ]]; then
    echo "--context-lines must be a non-negative integer" >&2
    exit 1
  fi

  if [[ "${STREAM_OUTPUT}" != "0" && "${STREAM_OUTPUT}" != "1" ]]; then
    echo "--stream-output must be either 0 or 1 internally" >&2
    exit 1
  fi
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

function run_with_logging() {
  local command_status=0
  local first_failure_line=""

  set +e
  if [[ "${STREAM_OUTPUT}" == "1" ]]; then
    "${COMMAND[@]}" 2>&1 | tee "${LOG_FILE}"
    command_status=${PIPESTATUS[0]}
  else
    "${COMMAND[@]}" > "${LOG_FILE}" 2>&1
    command_status=$?
  fi
  set -e

  if [[ "${command_status}" -ne 0 ]]; then
    first_failure_line="$(grep -m 1 -E "${FAILURE_PATTERN}" "${LOG_FILE}" || true)"
    if [[ -n "${first_failure_line}" ]]; then
      echo "::error title=${ERROR_TITLE}::${first_failure_line}"
    else
      echo "::error title=${ERROR_TITLE}::command exited with status ${command_status}"
    fi

    echo "::group::${SUMMARY_TITLE}"
    grep -E "${FAILURE_PATTERN}" "${LOG_FILE}" || true
    echo "::endgroup::"

    echo "::group::${CONTEXT_TITLE} (+/- ${CONTEXT_LINES} lines)"
    print_failure_contexts
    echo "::endgroup::"
    exit "${command_status}"
  fi
}

function main() {
  parse_args "$@"
  check_requirements
  run_with_logging
}

main "$@"
