#!/bin/bash -x
#
# Startup script to install Chrome remote desktop and a desktop environment.
#
# See environmental variables at then end of the script for configuration
#

function send_discord_notification() {
  # --helpオプションの確認
  if [ "$1" == "--help" ] || [ "$1" == "-h" ]; then
    cat <<EOF
Usage: ${FUNCNAME[0]} <通知テキスト> <Embedのテキスト> <Embedのフッターテキスト> <EmbedのアイコンのURL> <Embedの色> [WEBHOOK_URL]

Parameters:
  通知テキスト           : 通知内容のテキスト
  Embedのテキスト        : Embedメッセージの内容
  Embedのフッターテキスト : Embedフッターに表示するテキスト
  EmbedのアイコンのURL   : Embedフッターに表示するアイコンのURL
  Embedの色             : Embedの色（10進数の数値または以下の文字列指定が可能）
  WEBHOOK_URL            : (任意) DiscordのWebhook URL。省略した場合は環境変数 DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD を使用します。

Color Samples:
  green     : 4569935   (0x45BB4F)
  red       : 16711680  (0xFF0000)
  sky_blue  : 52479     (0x00CCFF)
  orange    : 14177041  (0xD85311)
  white     : 16777215  (0xFFFFFF)
  blue      : 39423     (0x0099FF)
  yellow    : 16770560  (0xFFE600)
  pink      : 16711833  (0xFF0099)
  purple    : 10494192  (0xA020F0)
  gray_blue : 9212588   (0x8C92AC)
  black     : 3355443   (0x333333)

Examples:
  # Using environment variable DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD (5 arguments)
  ${FUNCNAME[0]} "通知テキスト" "Embedのテキスト" "Embedのフッターテキスト" "https://example.com/footer_icon.png" green

  # Explicitly specifying WEBHOOK_URL as the last argument (6 arguments)
  ${FUNCNAME[0]} "通知テキスト" "Embedのテキスト" "Embedのフッターテキスト" "https://example.com/footer_icon.png" red "https://discord.com/api/webhooks/your_webhook_id/your_webhook_token"
EOF
    return 0
  fi

  local message embed_text embed_footer_text embed_icon_url embed_color webhook_url

  # Embedなし
  # 引数の個数チェック
  if [ "$#" -eq 1 ]; then
    message="$1"
    webhook_url="${DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD}"
    if [ -z "$webhook_url" ]; then
      echo "[WARN] DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD is empty. skip notification."
      return 0
    fi
    # JSONペイロードの作成
    local payload=$(cat <<EOF
{
  "content": "$message"
}
EOF
    )
    # curlでDiscordのWebhookにPOSTリクエストを送信
    curl -fsS -H "Content-Type: application/json" \
      -X POST \
      -d "$payload" \
      "$webhook_url"
    return 0
  fi

  # Embedあり
  # 引数の個数チェック
  if [ "$#" -eq 5 ]; then
    message="$1"
    embed_text="$2"
    embed_footer_text="$3"
    embed_icon_url="$4"
    embed_color="$5"
    webhook_url="${DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD}"
  elif [ "$#" -eq 6 ]; then
    message="$1"
    embed_text="$2"
    embed_footer_text="$3"
    embed_icon_url="$4"
    embed_color="$5"
    webhook_url="$6"
  else
    echo "Error: 引数の数が正しくありません。"
    cat <<EOF
Usage: ${FUNCNAME[0]} <通知テキスト> <Embedのテキスト> <Embedのフッターテキスト> <EmbedのアイコンのURL> <Embedの色> [WEBHOOK_URL];
EOF
    return 1
  fi

  # いずれかの引数が空文字だった場合のエラーハンドリング
  if [ -z "$message" ] || [ -z "$embed_text" ] || [ -z "$embed_icon_url" ] || [ -z "$embed_color" ] || [ -z "$webhook_url" ]; then
    echo "Error: 引数に空文字が含まれています。"
    cat <<EOF
Usage: ${FUNCNAME[0]} <通知テキスト> <Embedのテキスト> <Embedのフッターテキスト> <EmbedのアイコンのURL> <Embedの色> [WEBHOOK_URL];
EOF
    return 1
  fi

  # Embedの色が文字列の場合、対応する10進数の値に変換
  case "$embed_color" in
    green)
      embed_color=4569935
      ;;
    red)
      embed_color=16711680
      ;;
    sky_blue)
      embed_color=52479
      ;;
    orange)
      embed_color=14177041
      ;;
    white)
      embed_color=16777215
      ;;
    blue)
      embed_color=39423
      ;;
    yellow)
      embed_color=16770560
      ;;
    pink)
      embed_color=16711833
      ;;
    purple)
      embed_color=10494192
      ;;
    gray_blue)
      embed_color=9212588
      ;;
    black)
      embed_color=3355443
      ;;
    *)
      # 数値が直接入力されているものとみなす
      ;;
  esac

  # JSONペイロードの作成
  local payload=$(cat <<EOF
{
  "content": "$message",
  "embeds": [
    {
      "description": "$embed_text",
      "color": $embed_color,
      "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
      "footer": {
          "text": "$embed_footer_text",
          "icon_url": "$embed_icon_url"
      }
    }
  ]
}
EOF
  )

  # curlでDiscordのWebhookにPOSTリクエストを送信
  curl -fsS -H "Content-Type: application/json" \
    -X POST \
    -d "$payload" \
    "$webhook_url"
}

function send_discord_notification_about_gce() {
  # 第一引数以降のパラメータは send_discord_notification と同じ順序
  # 通知テキストに [GCE] プレフィックスを追加する例
  local message="$1"
  local embed_text="$2"
  local embed_color="$3"
  local embed_footer_text="GoogleComputeEngine"
  local embed_icon_url=$GCE_ICON_URL
  local webhook_url=$DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD

  # webhook_url が空でないかチェックし、あれば最後の引数として渡す
  if [ -n "$webhook_url" ]; then
    send_discord_notification "$message" "$embed_text" "$embed_footer_text" "$embed_icon_url" "$embed_color" "$webhook_url"
  else
    send_discord_notification "$message" "$embed_text" "$embed_footer_text" "$embed_icon_url" "$embed_color"
  fi
}

function set_env_var_from_custom_metadata() {
  local key=$1
  local value=""
  value=$(curl -fsS -H "Metadata-Flavor: Google" \
    "http://metadata.google.internal/computeMetadata/v1/instance/attributes/$key" 2>/dev/null || true)
  if [ -z "$value" ]; then
    value=$(curl -fsS -H "Metadata-Flavor: Google" \
      "http://metadata.google.internal/computeMetadata/v1/project/attributes/$key" 2>/dev/null || true)
  fi

  # Remove wrapping single/double quotes for values like "https://..." or ''.
  if [[ "$value" =~ ^\".*\"$ ]] || [[ "$value" =~ ^\'.*\'$ ]]; then
    value="${value:1:${#value}-2}"
  fi
  export "$key"="$value"
}

function install_chrome_remote_desktop() {
  curl --fail --show-error --location --retry 3 --max-time 120 https://dl.google.com/linux/linux_signing_key.pub \
    | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/chrome-remote-desktop.gpg
  echo "deb [arch=amd64] https://dl.google.com/linux/chrome-remote-desktop/deb stable main" \
    | sudo tee /etc/apt/sources.list.d/chrome-remote-desktop.list
  sudo apt-get update
  sudo DEBIAN_FRONTEND=noninteractive \
    apt-get install --assume-yes chrome-remote-desktop
  echo "[INFO] Setup for Chrome Remote Desktop completed successfully."
}

function install_xfce_desktop_env() {
  # Xfce デスクトップ環境の導入
  sudo DEBIAN_FRONTEND=noninteractive \
    apt install --assume-yes xfce4 desktop-base dbus-x11 xscreensaver
  # Xfce デスクトップ環境をデフォルトへ設定
  sudo bash -c 'echo "exec /etc/X11/Xsession /usr/bin/xfce4-session" > /etc/chrome-remote-desktop-session'
  # Chrome ブラウザをインストール
  curl --fail --show-error --location --retry 3 --max-time 300 -o google-chrome-stable_current_amd64.deb \
  https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
  sudo DEBIAN_FRONTEND=noninteractive apt install --assume-yes --fix-broken ./google-chrome-stable_current_amd64.deb
  # Color Manager を停止
  sudo systemctl stop colord
  sudo systemctl disable colord
  echo "[INFO] Setup for Xfce (X Windows System Desktop Environment) completed successfully."
}

function install_desktop_env_with_xfce() {
  install_chrome_remote_desktop
  install_xfce_desktop_env
}

function install_lxqt_desktop_env() {
  # LXQt デスクトップ環境の導入
  sudo DEBIAN_FRONTEND=noninteractive \
    apt install --assume-yes lxqt dbus-x11
  # LXQt デスクトップ環境をデフォルトへ設定
  echo "exec startlxqt" > /etc/chrome-remote-desktop-session
  # Chrome ブラウザをインストール
  curl --fail --show-error --location --retry 3 --max-time 300 -o google-chrome-stable_current_amd64.deb \
    https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
  sudo DEBIAN_FRONTEND=noninteractive apt install --assume-yes --fix-broken ./google-chrome-stable_current_amd64.deb
  # Color Manager を停止
  sudo systemctl stop colord 2>/dev/null
  sudo systemctl disable colord 2>/dev/null
  echo "[INFO] Setup for LXQt (X Windows System Desktop Environment) completed successfully."
}

function install_desktop_env_with_lxqt() {
  install_chrome_remote_desktop
  install_lxqt_desktop_env
}

function download_and_install { # args URL FILENAME
  if [[ -e "$2" ]] ; then
    echo "cannot download $1 to $2 - file exists"
    return 1;
  fi
  curl -L -o "$2" "$1" && \
    apt-get install --assume-yes --fix-broken "$2" && \
    rm "$2"
}

function is_installed {  # args PACKAGE_NAME
  dpkg-query --list "$1" | grep -q "^ii" 2>/dev/null
  return $?
}

function install_japanese_locale() {
  # install japanese locale
  sudo apt -y install locales
  sudo localectl set-locale LANG=ja_JP.UTF-8 LANGUAGE="ja_JP:ja"
  source /etc/default/locale
  echo $LANG

  # install japanese locale for desktop
  sudo apt -y install task-japanese-desktop
}

function setup_japanese_timezone() {
  sudo chmod 777 /etc/timezone
  sudo rm /etc/localtime
  echo Asia/Tokyo > /etc/timezone
  sudo chmod 644 /etc/timezone
  sudo dpkg-reconfigure -f noninteractive tzdata
}

function install_japanese_input_method() {
  # install input methods
  sudo apt update
  sudo apt install fcitx5-mozc -y

  touch ~/.profile
  inserting=$(cat <<EOF
while true; do
  dbus-update-activation-environment --systemd DBUS_SESSION_BUS_ADDRESS DISPLAY XAUTHORITY 2> /dev/null && break
done

export GTK_IM_MODULE=fcitx
export QT_IM_MODULE=fcitx
export XMODIFIERS="@im=fcitx"
if [ $SHLVL = 1 ] ; then
  (fcitx5 --disable=wayland -d --verbose '*'=0 &)
  xset -r 49  > /dev/null 2>&1
fi
EOF
  )
  echo $inserting >> ~/.profile

  im-config -n fcitx5
}

function setup_dev_resources() {
  # Setup development resources
  # local home_dir="$DEV_HOME"
  sudo apt update
  sudo apt install git -y
}

function setup_vscode() {
  # Install VSCode
  wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
  sudo install -D -o root -g root -m 644 packages.microsoft.gpg /etc/apt/keyrings/packages.microsoft.gpg
  sudo sh -c 'echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/keyrings/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" > /etc/apt/sources.list.d/vscode.list'
  rm -f packages.microsoft.gpg
  sudo apt install apt-transport-https
  sudo apt update
  sudo apt install code
}

function install_docker() {
  if command -v docker >/dev/null 2>&1; then
    echo "[INFO] Docker is already installed. skip installation."
    return 0
  fi

  local docker_repo_os
  docker_repo_os="$(. /etc/os-release && echo "${ID}")"
  if [ "$docker_repo_os" != "debian" ] && [ "$docker_repo_os" != "ubuntu" ]; then
    docker_repo_os="debian"
  fi

  sudo DEBIAN_FRONTEND=noninteractive apt-get update
  sudo DEBIAN_FRONTEND=noninteractive \
    apt-get install --assume-yes ca-certificates curl gnupg lsb-release

  sudo install -m 0755 -d /etc/apt/keyrings
  sudo rm -f /etc/apt/keyrings/docker.gpg
  curl -fsSL "https://download.docker.com/linux/${docker_repo_os}/gpg" \
    | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
  sudo chmod a+r /etc/apt/keyrings/docker.gpg

  echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/${docker_repo_os} \
    $(lsb_release -cs) stable" \
    | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

  sudo DEBIAN_FRONTEND=noninteractive apt-get update
  sudo DEBIAN_FRONTEND=noninteractive \
    apt-get install --assume-yes docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
  echo "[INFO] Docker installation completed successfully."
}

# Configure the following environmental variables as required:
INSTALL_XFCE=yes
INSTALL_CINNAMON=yes
INSTALL_CHROME=yes
INSTALL_FULL_DESKTOP=yes

# Any additional packages that should be installed on startup can be added here
EXTRA_PACKAGES="less bzip2 zip unzip tasksel wget"

set_env_var_from_custom_metadata VSC_PROFILE_URL
set_env_var_from_custom_metadata DISCORD_WEBHOOK_URL_FOR_IAC_ON_GCLOUD
set_env_var_from_custom_metadata GCE_ICON_URL
send_discord_notification "VMのカスタムメタデータを環境変数に反映したよ！"

send_discord_notification "VMのスタートアップスクリプトを実行するよ！"

apt-get update

# Install X Windows desktop system
if ! is_installed chrome-remote-desktop; then
  # install_desktop_env_with_xfce
  if ! install_desktop_env_with_lxqt; then
    send_discord_notification_about_gce "失敗…" "VMのデスクトップ環境の設定に失敗したよ…" "red"
    echo "[ERROR] Desktop environment setup failed"
    exit 1
  fi
fi

# install_desktop_env_with_xfce
# install_desktop_env_with_lxqt
send_discord_notification "VMのデスクトップ環境の設定が完了したよ！"

# [[ "$INSTALL_CHROME" = "yes" ]] && ! is_installed google-chrome-stable && \
#   download_and_install \
#     https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
#     /tmp/google-chrome-stable_current_amd64.deb

echo "[INFO] Chrome remote desktop installation completed"
send_discord_notification "VMのChromeリモデの設定が完了したよ！"

install_japanese_locale
echo "[INFO] Install japanese locale completed"
send_discord_notification "VMのロケールの設定が完了したよ！"

setup_japanese_timezone
echo "[INFO] Setup timezone completed"
send_discord_notification "VMのタイムゾーンの設定が完了したよ！"

install_japanese_input_method
echo "[INFO] Pre-installation of japanese input methods completed"
send_discord_notification "VMのIMEの設定が完了したよ！"

setup_dev_resources
echo "[INFO] Setup development resources completed"
send_discord_notification "VMの開発リソースの設定が完了したよ！"

if ! install_docker; then
  send_discord_notification_about_gce "失敗…" "VMのDockerの設定に失敗したよ…" "red"
  echo "[ERROR] Docker installation failed"
  exit 1
fi
send_discord_notification "VMのDockerの設定が完了したよ！"

setup_vscode
echo "[INFO] Setup for VSCode completed"
send_discord_notification "VMのVSCodeの設定が完了したよ！"
send_discord_notification_about_gce "終わった！" "VMのスタートアップスクリプトが完了したよ！" "green"
