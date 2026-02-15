#!/bin/bash
function generate_service_summary() {
  ./pkg/bin/cli/linux_amd64/service-implementing-viewer \
    -operation=aggregate-summary \
    -root-dir=$HOME/devbox/cmd \
    -target-dirs=cli,mcp,grpc/handlers,http/handlers \
    -write-file=docs/project_status/service_summary.md
}

function show_help() {
  local FUNC="${FUNCNAME[0]}"
  cat <<EOF
[INFO] [$FUNC] 各ツールのREADME概要をカテゴリ別に集約し、Markdownに書き出します。

使用方法:
  ./scripts/ops/generate_service_summary.sh

例:
  ./scripts/ops/generate_service_summary.sh

--help を指定するとこのメッセージを表示します。
EOF
}

# === 実行部 ===
function main() {
  local FUNC="${FUNCNAME[0]}"

  if [[ "$1" == "--help" ]]; then
    show_help
    exit 0
  fi

  generate_service_summary
}

main "$@"
