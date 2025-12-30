#!/bin/bash

function main() {
  # ヘルプ表示
  if [[ "$1" == "-h" || "$1" == "--help" || -z "$1" ]]; then
    cat << 'EOF'
Usage: script_generator_to_build.sh <tool_name>

Description:
  指定されたツール名でビルドスクリプトを生成します。

Arguments:
  tool_name      生成するツール名（必須）

Options:
  -h, --help     このヘルプメッセージを表示

Examples:
  script_generator_to_build.sh mytool
  script_generator_to_build.sh user-cli
  script_generator_to_build.sh "build_helper"
EOF
    return 0
  fi

  local tool_name="$1"

  # 実行するコマンドを表示
  echo "Executing: ./pkg/bin/cli/linux_amd64/script-generator-to-build \"$tool_name\""

  # コマンド実行
  ./pkg/bin/cli/linux_amd64/script-generator-to-build "$tool_name"
}

main "$@"
