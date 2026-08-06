---
name: hooks-craftsman
description: Codex hooks の作成・更新・検証を行うスキル。`~/.codex/hooks.json`、`features.codex_hooks`、PreToolUse/PostToolUse/SessionStart/UserPromptSubmit/Stop、hook スクリプト、hook が発火しない原因調査、Codex hooks の schema や実装情報の確認が必要な場面で使用する。
---

# Hooks Craftsman

## 基本方針

- 回答・記述は日本語で行う。
- hook は「通常作業を過剰に止めない」ことを優先し、ブロック条件を具体的に絞る。
- hook スクリプトは必ず単体テストしてから、実際の Codex tool 実行で確認する。
- `hooks.json` と hook スクリプトを分け、複雑な判定は shell script 側へ寄せる。
- スクリプト変更時は `bash -n` と代表 payload テストを実施する。

## Codex hooks 情報の調査方法

まずローカルの Codex 実装・設定を確認する。

```bash
codex --version
CODEX_HOME="${CODEX_HOME:-$HOME/.codex}"
CODEX_REFS="${CODEX_REFS:-$HOME/package_references/codex}"
rg -n "codex_hooks|hooks.json|PreToolUse|PostToolUse|SessionStart|UserPromptSubmit|Stop" "$CODEX_HOME" "$CODEX_REFS" -g '*.md' -g '*.rs' -g '*.json' -g '*.toml'
```

主な確認先。

- `~/.codex/config.toml`: `[features] codex_hooks = true` の有無
- `~/.codex/hooks.json`: hooks 設定。symlink の場合は `readlink -f ~/.codex/hooks.json` で実体を確認する
- `$CODEX_REFS/docs/config.md`: 公式 config docs への入口
- `$CODEX_REFS/codex-rs/hooks/src/engine/config.rs`: `hooks.json` の構造
- `$CODEX_REFS/codex-rs/hooks/src/engine/discovery.rs`: hooks の探索・matcher の扱い
- `$CODEX_REFS/codex-rs/hooks/src/events/pre_tool_use.rs`: `PreToolUse` の入力・ブロック処理
- `$CODEX_REFS/codex-rs/hooks/schema/generated/*.json`: 各 hook event の入出力 JSON schema

公式 docs が必要な場合は OpenAI/Codex 公式情報を優先する。ただし hooks の詳細はローカルの `openai/codex` 実装と generated schema が最も具体的な一次情報になりやすい。

## 設定ファイル

`~/.codex/config.toml` に feature を有効化する。

```toml
[features]
codex_hooks = true
```

`~/.codex/hooks.json` を配置する。管理場所を別にする場合は symlink でよい。

```bash
ln -s /path/to/hooks.json ~/.codex/hooks.json
```

`hooks.json` の基本形。

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "^Bash$",
        "hooks": [
          {
            "type": "command",
            "command": "/absolute/path/to/script.sh",
            "timeout": 5,
            "statusMessage": "Checking policy"
          }
        ]
      }
    ]
  }
}
```

注意点。

- `matcher` は regex。`^Bash$` は Bash tool のみ対象にする。
- `type: "command"` の hook が実用対象。`prompt` / `agent` は実装上 skip される場合がある。
- `command` に script path を直接書くなら実行権限が必要。
- 実行権限を避けるなら `command` を `bash /path/to/script.sh` にする。

## PreToolUse の入力

`PreToolUse` command hook は stdin で JSON payload を受け取る。Bash の場合、最低限次を使う。

```bash
payload="$(cat)"
cwd="$(jq -r '.cwd // ""' <<<"$payload")"
cmd="$(jq -r '.tool_input.command // ""' <<<"$payload")"
```

主なフィールド。

- `.cwd`: 実行時 cwd
- `.tool_name`: 通常 `Bash`
- `.tool_input.command`: 実行予定の shell command
- `.model`, `.permission_mode`, `.session_id`, `.turn_id`: 必要に応じてログ・判定に使う

## ブロック方法

`PreToolUse` のブロック方法は2つある。

1. stderr に理由を書き、exit code `2` で終了する。

```bash
echo "Blocked: reason" >&2
exit 2
```

2. stdout に JSON を出して exit code `0` で終了する。

```bash
jq -n --arg reason "Blocked: reason" '{decision:"block", reason:$reason}'
exit 0
```

単純な guard では exit code `2` が扱いやすい。

## hook スクリプト作成手順

1. `hooks.json` の対象 event と matcher を決める。
2. hook script を absolute path で作成する。
3. script 先頭で `payload="$(cat)"` を必ず行う。
4. `cwd` から project root を解決する。
5. ブロック条件を最小化する。
6. `bash -n` で構文検査する。
7. 代表 payload を stdin に流して単体テストする。
8. 実際の Codex tool 実行で hook が発火するか確認する。

project root 解決の基本形。

```bash
project_root="$(
  git -C "$cwd" rev-parse --show-toplevel 2>/dev/null || realpath -m "$cwd"
)"
allowed_root="$(realpath -m "$project_root/agent-work")"
```

`agent-work` 配下判定。

```bash
is_under_allowed_root() {
  local path="$1"
  local resolved

  if [[ "$path" = /* ]]; then
    resolved="$(realpath -m "$path")"
  else
    resolved="$(realpath -m "$cwd/$path")"
  fi

  [[ "$resolved" == "$allowed_root" || "$resolved" == "$allowed_root/"* ]]
}
```

## よくある guard パターン

### mkdir/touch の tmp/cache 引数をブロックする

`mkdir` / `touch` の引数に `tmp` または `cache` が含まれる場合だけブロックする。通常の実装用ディレクトリ作成を過剰に止めない。

```bash
block_tmp_cache_args_for_command() {
  local command_name="$1"

  while read -r path; do
    [[ -z "$path" ]] && continue

    local lower_path
    lower_path="$(tr '[:upper:]' '[:lower:]' <<<"$path")"

    if [[ "$lower_path" == *tmp* || "$lower_path" == *cache* ]]; then
      echo "Blocked: ${command_name} arguments must not contain tmp or cache." >&2
      echo "Rejected argument: $path" >&2
      exit 2
    fi
  done < <(
    grep -oE "(^|[;&|[:space:]])${command_name}[[:space:]][^;&|]+" <<<"$cmd" \
      | sed -E "s/(^|[;&|[:space:]])${command_name}[[:space:]]+//" \
      | tr ' ' '\n' \
      | grep -v '^-' \
      || true
  )
}
```

### Go 一時環境変数を guard する

`TMPDIR`, `GOTMPDIR`, `GOCACHE` は Go 系コマンドの時だけ検査する。

対象例。

- `go test`
- `go test -coverprofile=...`
- `go tool cover ...`
- `go build`
- `go run`
- `go vet`
- `go generate`
- `go list`

```bash
is_go_command() {
  [[ "$cmd" =~ (^|[[:space:];&|])go[[:space:]]+(test|tool|build|run|vet|generate|list)([[:space:]]|$) ]]
}

guard_go_temp_envs() {
  is_go_command || return 0

  for key in TMPDIR GOTMPDIR GOCACHE; do
    value="$(
      grep -oE "(^|[[:space:]])${key}=[^[:space:]]+" <<<"$cmd" \
        | head -1 \
        | sed -E "s/^[[:space:]]*${key}=//" \
        || true
    )"

    [[ -z "$value" ]] && continue
    block_if_outside_agents "$value"
  done
}
```

### `go test -c` の出力先を強制する

`go test -c` は test binary を生成するため、`-o` を必須にし、出力先を project `agent-work` 配下へ限定する。

```bash
extract_go_test_compile_output() {
  awk '
    {
      for (i = 1; i <= NF; i++) {
        if ($i == "-o") {
          if (i + 1 <= NF) {
            print $(i + 1)
          }
          exit
        }
        if ($i ~ /^-o=.+/) {
          sub(/^-o=/, "", $i)
          print $i
          exit
        }
      }
    }
  ' <<<"$cmd"
}

guard_go_test_compile_output() {
  [[ "$cmd" =~ (^|[[:space:];&|])go[[:space:]]+test([[:space:]][^;&|]*)?[[:space:]]-c([[:space:]]|$) ]] || return 0

  local output_path
  output_path="$(extract_go_test_compile_output)"

  if [[ -z "$output_path" ]]; then
    echo "Blocked: go test -c must specify -o under project agent-work directory." >&2
    echo "Project root: $project_root" >&2
    echo "Allowed root: $allowed_root" >&2
    exit 2
  fi

  block_if_outside_agents "$output_path"
}
```

## 検証手順

### 1. 構文検査

```bash
bash -n /path/to/hook-script.sh
```

### 2. 単体 payload テスト

ブロックされる例。

```bash
printf '%s\n' '{"cwd":"<repo>","tool_input":{"command":"mkdir -p .codex-cache/hook-test"}}' \
  | /path/to/hook-script.sh
```

許可される例。

```bash
printf '%s\n' '{"cwd":"<repo>","tool_input":{"command":"mkdir -p backend/cmd/rpc/connect_handlers/prefixes/hook_test_dir"}}' \
  | /path/to/hook-script.sh
```

`go test -c` のブロック例。

```bash
printf '%s\n' '{"cwd":"<repo>","tool_input":{"command":"go test -c ./cmd/rpc/foo"}}' \
  | /path/to/hook-script.sh
```

`go test -c` の許可例。

```bash
printf '%s\n' '{"cwd":"<repo>","tool_input":{"command":"go test -c -o agent-work/cache/foo.test ./cmd/rpc/foo"}}' \
  | /path/to/hook-script.sh
```

### 3. 実 hook 経路テスト

Codex の shell tool で実際に対象コマンドを実行する。

```bash
mkdir -p .codex-cache/hook-test
go test -c ./cmd/rpc/connect_handlers/medication_records/query_medication_records
go test -c -o backend/foo.test ./cmd/rpc/connect_handlers/medication_records/query_medication_records
```

期待値は `Command blocked by PreToolUse hook:` が表示され、対象ファイル・ディレクトリが作成されないこと。

## トラブルシュート

- hook が発火しない:
  - `codex --version` を確認する。
  - `~/.codex/config.toml` に `[features] codex_hooks = true` があるか確認する。
  - `~/.codex/hooks.json` が存在するか、symlink が正しいか確認する。
  - Codex セッション起動後に設定した場合は再起動を疑う。
  - 利用中の実行経路が Codex CLI の hook runtime を通っているか確認する。
- hook は走るがブロックされない:
  - script が stdin から `payload="$(cat)"` しているか確認する。
  - `.tool_input.command` を見ているか確認する。
  - matcher が `^Bash$` など実際の tool name と一致しているか確認する。
- hook が失敗扱いになる:
  - `set -euo pipefail` と `grep` の exit `1` に注意する。
  - パイプ末尾に `|| true` を入れるべき箇所を確認する。
  - `bash -n` と `bash -x` で切り分ける。
- script path で permission denied:
  - `chmod +x /path/to/script.sh` を実行する。
  - もしくは `hooks.json` の `command` を `bash /path/to/script.sh` にする。

## 最終報告

最終報告では次を簡潔に示す。

- 更新した `hooks.json` / script のパス
- 追加・変更した guard 条件
- 単体 payload テストの結果
- 実 hook 経路テストの結果
- 残る制約や誤爆リスク
