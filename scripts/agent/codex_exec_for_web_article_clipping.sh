#!/usr/bin/env bash
set -euo pipefail

function show_help() {
  cat <<'EOF'
Usage:
  codex_exec_for_web_article_clipping.sh -t <target_dir> -f <article_list_file>
  codex_exec_for_web_article_clipping.sh <target_dir> <article_list_file>

Description:
  Execute the following command sequentially for each line in article_list_file:
    codex exec '$web-clipping-summary <target_dir>, <embedded_article>'

article_list_file format:
  One article per line, e.g.
    [title](https://example.com/article)

Notes:
  - Empty lines are ignored.
  - Lines starting with '#' are treated as comments and ignored.
EOF
}

function parse_args() {
  local -n out_target_dir_ref="$1"
  local -n out_article_list_file_ref="$2"
  shift 2

  local -a positional=()

  while [[ $# -gt 0 ]]; do
    case "$1" in
      -t|--target-dir)
        if [[ $# -lt 2 ]]; then
          echo "[ERROR] Missing value for $1" >&2
          exit 1
        fi
        out_target_dir_ref="$2"
        shift 2
        ;;
      -f|--article-list-file)
        if [[ $# -lt 2 ]]; then
          echo "[ERROR] Missing value for $1" >&2
          exit 1
        fi
        out_article_list_file_ref="$2"
        shift 2
        ;;
      -h|--help)
        show_help
        exit 0
        ;;
      --)
        shift
        while [[ $# -gt 0 ]]; do
          positional+=("$1")
          shift
        done
        ;;
      -*)
        echo "[ERROR] Unknown option: $1" >&2
        show_help >&2
        exit 1
        ;;
      *)
        positional+=("$1")
        shift
        ;;
    esac
  done

  if [[ ${#positional[@]} -gt 2 ]]; then
    echo "[ERROR] Too many positional arguments." >&2
    show_help >&2
    exit 1
  fi

  if [[ ${#positional[@]} -ge 1 ]]; then
    out_target_dir_ref="${positional[0]}"
  fi

  if [[ ${#positional[@]} -ge 2 ]]; then
    out_article_list_file_ref="${positional[1]}"
  fi
}

function main() {
  local target_dir=""
  local article_list_file=""
  local -a article_list=()

  parse_args target_dir article_list_file "$@"

  if [[ -z "$target_dir" ]]; then
    echo "[ERROR] target_dir is required." >&2
    show_help >&2
    exit 1
  fi

  if [[ -z "$article_list_file" ]]; then
    echo "[ERROR] article_list_file is required." >&2
    show_help >&2
    exit 1
  fi

  if [[ ! -f "$article_list_file" ]]; then
    echo "[ERROR] article_list_file not found: $article_list_file" >&2
    exit 1
  fi

  mapfile -t article_list < "$article_list_file"

  for embedded_article in "${article_list[@]}"; do
    if [[ "$embedded_article" =~ ^[[:space:]]*$ ]]; then
      continue
    fi

    if [[ "$embedded_article" =~ ^[[:space:]]*# ]]; then
      continue
    fi

    if [[ -z "$embedded_article" ]]; then
      continue
    fi

    local prompt
    local header
    header="clipped_at_date\\ttitle\\turl"
    prompt="\$web-clipping-summary ${target_dir}\\n\\n${header}\\n${embedded_article}"

    echo "[INFO] Running: codex exec '<prompt with article>'"
    codex exec "$prompt"
  done
}

main "$@"
