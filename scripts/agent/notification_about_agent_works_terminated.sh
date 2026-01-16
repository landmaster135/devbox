#!/bin/bash
function notification_about_agent_works_terminated(){
  $HOME/devbox/pkg/bin/cli/linux_amd64/discord-webhook \
    -webhook-url "$1" \
    -content-text "作業が終わったみたいです！" \
    -embed-type vscode \
    -embed-text "$2"
}

function show_help() {
  local FUNC="${FUNCNAME[0]}"
  cat <<EOF
[INFO] [$FUNC] AIエージェントの作業が完了した時にDiscordに通知します。

使用方法:
  ./pkg/bash/notification_about_agent_works_terminated.sh

例:
  ./pkg/bash/notification_about_agent_works_terminated.sh

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

  local webhook_url="$1"
  local embed_text="$2"

  notification_about_agent_works_terminated $webhook_url $embed_text
}

main "$@"
