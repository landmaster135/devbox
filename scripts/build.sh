#!/bin/bash

function run_all_sh_scripts() {
  # ローカル変数の宣言
  local func_name="${FUNCNAME[0]}"
  local target_dir=""
  local exit_code=0
  local script_count=0
  local failed_scripts=()
  local skipped_scripts=()
  local skip_list=()

  # ヘルプメッセージの表示
  if [[ "$1" == "--help" ]]; then
    echo "${func_name} - 指定されたディレクトリ内のすべての.shスクリプトを実行します"
    echo ""
    echo "使用方法:"
    echo "  ${func_name} <ディレクトリパス> [スキップするファイル名...]"
    echo ""
    echo "例:"
    echo "  ${func_name} /path/to/scripts"
    echo "  ${func_name} . skip_this.sh also_skip.sh"
    echo ""
    echo "説明:"
    echo "  第1引数にディレクトリパスを指定します。"
    echo "  第2引数以降にスキップしたいスクリプトのファイル名を指定できます。"
    echo "  スキップするファイル名はパスではなく、ファイル名のみを指定してください。"
    return 0
  fi

  # 引数のチェック
  if [[ $# -lt 1 ]]; then
    echo "[ERROR] ${func_name}: ディレクトリパスを最低1つ指定してください" >&2
    echo "[ERROR] ${func_name}: 使用方法については '${func_name} --help' を実行してください" >&2
    return 1
  fi

  target_dir="$1"
  shift

  # スキップリストの作成
  if [[ $# -gt 0 ]]; then
    skip_list=("$@")
    echo "[INFO] ${func_name}: 以下のスクリプトはスキップします: ${skip_list[*]}"
  fi

  # ディレクトリの存在チェック
  if [[ ! -d "${target_dir}" ]]; then
    echo "[ERROR] ${func_name}: ディレクトリ '${target_dir}' が存在しません" >&2
    return 1
  fi

  echo "[INFO] ${func_name}: ディレクトリ '${target_dir}' 内の.shスクリプトを検索しています..."

  # .shファイルの検索
  local sh_files=()
  while IFS= read -r file; do
    sh_files+=("$file")
  done < <(find "${target_dir}" -type f -name "*.sh" -print)

  # スクリプトが見つからない場合
  if [[ ${#sh_files[@]} -eq 0 ]]; then
    echo "[INFO] ${func_name}: ディレクトリ '${target_dir}' 内に.shスクリプトが見つかりませんでした"
    return 0
  fi

  echo "[INFO] ${func_name}: ${#sh_files[@]}個の.shスクリプトが見つかりました"

  # 各スクリプトの実行
  for script in "${sh_files[@]}"; do
    # ファイル名のみを取得
    local script_basename=$(basename "${script}")

    # スキップリストに含まれているかチェック
    local should_skip=0
    for skip_file in "${skip_list[@]}"; do
      if [[ "${script_basename}" == "${skip_file}" ]]; then
        should_skip=1
        skipped_scripts+=("${script}")
        echo "[INFO] ${func_name}: スクリプト '${script}' はスキップリストに含まれているためスキップします"
        break
      fi
    done

    # スキップする場合は次のスクリプトへ
    if [[ ${should_skip} -eq 1 ]]; then
      continue
    fi

    ((script_count++))

    # 実行権限のチェック
    if [[ ! -x "${script}" ]]; then
      echo "[INFO] ${func_name}: スクリプト '${script}' に実行権限を付与します"
      echo "[INFO] ${func_name}: 実行コマンド: chmod +x \"${script}\""
      chmod +x "${script}"

      if [[ $? -ne 0 ]]; then
        echo "[ERROR] ${func_name}: スクリプト '${script}' に実行権限を付与できませんでした" >&2
        failed_scripts+=("${script}")
        continue
      fi
    fi

    echo "[INFO] ${func_name}: (${script_count}/$((${#sh_files[@]} - ${#skipped_scripts[@]}))) スクリプト '${script}' を実行します"
    echo "[INFO] ${func_name}: 実行コマンド: bash \"${script}\""

    # スクリプトの実行
    bash "${script}"
    local script_status=$?

    if [[ ${script_status} -ne 0 ]]; then
      echo "[ERROR] ${func_name}: スクリプト '${script}' はエラーコード ${script_status} で終了しました" >&2
      failed_scripts+=("${script}")
      exit_code=1
    else
      echo "[INFO] ${func_name}: スクリプト '${script}' が正常に実行されました"
    fi

    echo ""
  done

  # 実行結果のサマリー
  echo "[INFO] ${func_name}: 全スクリプトの処理が完了しました (成功: $((script_count - ${#failed_scripts[@]})), 失敗: ${#failed_scripts[@]}, スキップ: ${#skipped_scripts[@]})"

  if [[ ${#skipped_scripts[@]} -gt 0 ]]; then
    echo "[INFO] ${func_name}: 以下のスクリプトはスキップされました:"
    for skipped in "${skipped_scripts[@]}"; do
      echo "[INFO] ${func_name}:   - ${skipped}"
    done
  fi

  if [[ ${#failed_scripts[@]} -gt 0 ]]; then
    echo "[ERROR] ${func_name}: 以下のスクリプトの実行に失敗しました:"
    for failed in "${failed_scripts[@]}"; do
      echo "[ERROR] ${func_name}:   - ${failed}"
    done
    return ${exit_code}
  fi

  return 0
}

run_all_sh_scripts ./scripts build.sh create_project_files.sh build_mcp_tools.sh
