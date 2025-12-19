#!/bin/bash
function generate_service_implementing_table(){
  ./pkg/bin/cli/linux_amd64/service-implementing-viewer -root-dir=$HOME/devbox/cmd -target-dirs=cli,mcp,grpc/handlers,http/handlers -operation=output
}

function show_help() {
  local FUNC="${FUNCNAME[0]}"
  cat <<EOF
[INFO] [$FUNC] 各サービスが実装されているツールの種別の早見表をMarkdownテーブル形式で自動生成します。

使用方法:
  ./pkg/bash/generate_service_implementing_table.sh

例:
  ./pkg/bash/generate_service_implementing_table.sh

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

  generate_service_implementing_table
}

main
