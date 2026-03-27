#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"
OUTPUT_DIR="${ROOT_DIR}/pkg/bin"

SOURCE_DIRS=(
  "${ROOT_DIR}/pkg/bin/cli"
  "${ROOT_DIR}/pkg/bin/mcp"
  "${ROOT_DIR}/pkg/taskfile"
)

function show_help() {
  cat <<EOF
[INFO] コンパイル済みツール群を1つの zip にまとめて pkg/bin 配下へ出力します。
出力ファイル名: compiled_tools_yyyyMMddhhmmss.zip

使用方法:
  ./scripts/ops/zip_compiled_tools.sh

例:
  ./scripts/ops/zip_compiled_tools.sh

--help を指定するとこのメッセージを表示します。
EOF
}

function check_requirements() {
  if ! command -v zip >/dev/null 2>&1; then
    echo "zip コマンドが見つかりません。" >&2
    exit 1
  fi
}

function ensure_source_dirs() {
  for src_dir in "${SOURCE_DIRS[@]}"; do
    if [ ! -d "${src_dir}" ]; then
      echo "ディレクトリが見つかりません: ${src_dir}" >&2
      exit 1
    fi
  done
}

function zip_compiled_tools() {
  local relative_paths=()
  local timestamp archive_path

  mkdir -p "${OUTPUT_DIR}"
  ensure_source_dirs

  for src_dir in "${SOURCE_DIRS[@]}"; do
    relative_paths+=("${src_dir#${ROOT_DIR}/}")
  done

  timestamp="$(date '+%Y%m%d%H%M%S')"
  archive_path="${OUTPUT_DIR}/compiled_tools_${timestamp}.zip"

  (
    cd "${ROOT_DIR}"
    zip -r -db -dc "${archive_path}" "${relative_paths[@]}"
  )

  echo "作成: ${archive_path}"
}

function main() {
  if [[ "${1:-}" == "--help" ]]; then
    show_help
    exit 0
  fi

  check_requirements
  zip_compiled_tools
}

main "$@"
