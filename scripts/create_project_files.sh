#!/bin/bash

# === ヘルプ表示 ===
function show_help() {
  local FUNC="${FUNCNAME[0]}"
  cat <<EOF
[INFO] [$FUNC] プロジェクト用のGoファイル構成を自動生成します。

使用方法:
  ./scripts/create_project_files.sh <project-name>

例:
  ./scripts/create_project_files.sh image-renamer-for-screenshot
  ./scripts/create_project_files.sh kana-converter

--help を指定するとこのメッセージを表示します。
EOF
}

# === ベースディレクトリ取得 ===
function get_base_dir() {
  local FUNC="${FUNCNAME[0]}"
  echo "$(dirname "$(realpath "$0")")/.."
}

# === スネークケースに変換 ===
function to_snake_case() {
  local FUNC="${FUNCNAME[0]}"
  echo "$1" | tr '-' '_'
}

# === ファイル生成処理 ===
function create_file_with_content() {
  local FUNC="${FUNCNAME[0]}"
  local file_path="$1"
  local base_dir="$2"
  local full_path="${base_dir}/${file_path}"
  local dir_path
  dir_path=$(dirname "$full_path")

  echo "[INFO] [$FUNC] mkdir -p \"$dir_path\""
  mkdir -p "$dir_path" || {
    echo "[ERROR] [$FUNC] ディレクトリ作成に失敗しました: $dir_path"
    exit 1
  }

  echo "[INFO] [$FUNC] touch \"$full_path\"（内容の生成あり）"
  case "$file_path" in
    *main.go)
      cat <<EOF > "$full_path"
package main

func main() {
  // TODO: implement
}
EOF
      ;;
    *services.go)
      cat <<EOF > "$full_path"
package usecases

// TODO: implement service functions
EOF
      ;;
    *services_test.go)
      cat <<EOF > "$full_path"
package usecases

import "testing"

func TestSomething(t *testing.T) {
  // TODO: add tests
}
EOF
      ;;
    *)
      touch "$full_path" || {
        echo "[ERROR] [$FUNC] ファイル作成に失敗しました: $full_path"
        exit 1
      }
      ;;
  esac
}

# === プロジェクトファイル作成 ===
function create_project_files() {
  local FUNC="${FUNCNAME[0]}"
  local project_name="$1"

  if [[ -z "$project_name" ]]; then
    echo "[ERROR] [$FUNC] プロジェクト名が指定されていません。"
    show_help
    exit 1
  fi

  local base_dir
  base_dir=$(get_base_dir)
  local snake_case_name
  snake_case_name=$(to_snake_case "$project_name")

  local files=(
    "cmd/cli/${project_name}/main.go"
    "cmd/cli/${project_name}/README.md"
    "internal/independencies/${snake_case_name}/usecases/services_test.go"
    "internal/independencies/${snake_case_name}/usecases/services.go"
    "internal/independencies/${snake_case_name}/.gitkeep"
  )

  for file in "${files[@]}"; do
    create_file_with_content "$file" "$base_dir"
  done

  echo "[INFO] [$FUNC] ✅ プロジェクト「$project_name」のファイルとディレクトリが作成されました。"
}

# === 実行部 ===
function main() {
  local FUNC="${FUNCNAME[0]}"

  if [[ "$1" == "--help" ]]; then
    show_help
    exit 0
  fi

  if [[ $# -ne 1 ]]; then
    echo "[ERROR] [$FUNC] 引数が不正です。"
    show_help
    exit 1
  fi

  create_project_files "$1"
}

main "$@"
