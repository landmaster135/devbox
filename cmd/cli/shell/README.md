# Shell CLI

Codexの`shell`ツール互換でコマンドを安全に実行するCLIです。`command: Vec<String>`の形で引数を受け取り、ベースディレクトリやサンドボックス権限を厳密に制御しながらコマンドを起動できます。

## 主な特徴

- `command`配列をそのまま`execvp`に渡す設計で、シェル展開を伴わない安全な実行
- `--base-dir`と`--workdir`で実行ディレクトリを制限し、ベース外のアクセスを拒否
- `--env KEY=VALUE`を複数指定して環境変数を補完
- 実行結果はJSONで出力（stdout/stderr/exit_code/timed_outなどを含む）
- Codex MCP版と同じ許可コマンド一覧を`list_allowed`操作で取得可能

## 主要フラグ

| フラグ | 説明 | 必須 |
| ------ | ---- | ---- |
| `-operation` | 実行する操作。`execute`（コマンド実行）か`list_allowed`（許可コマンド一覧）を指定 | mandatory |
| `-command` (複数指定可) | `command: Vec<String>`の各要素を追加。`--`以降の位置引数も同様に扱われる | `execute`時のみ必須 |
| `-workdir`, `-cwd` | 実行ディレクトリ。`-base-dir`配下に限定され、相対・絶対どちらも可 | optional |
| `-base-dir`, `-basedir` | 許可されたルートディレクトリ。未指定時は現在の作業ディレクトリが適用 | optional |
| `-timeout-ms`, `-timeout` | ミリ秒単位のタイムアウト。0の場合は内部デフォルト60秒が利用される | optional |
| `-env KEY=VALUE` (複数指定可) | 追加の環境変数を設定。1回の指定で1組のKEY=VALUEを受け付ける | optional |
| `-sandbox-permissions`, `-sandbox` | `use_default` または `require_escalated` を指定し、必要に応じてサンドボックス外実行を要求 | optional |
| `-justification`, `-reason` | `-sandbox-permissions=require_escalated`を選んだ場合の理由。承認フローへ伝える1行コメント | `-sandbox-permissions=require_escalated`時に必須 |
| `-help`, `-h` | ヘルプを表示し即時終了 | optional |

※ `--` を挟んだ後ろの位置引数は、そのまま `command` 配列に連結されます（例: `-- bash -lc "ls"`）。

### エスカレーション要求の例

```bash
go run ./cmd/cli/shell \
  -operation=execute \
  -sandbox-permissions=require_escalated \
  -justification="Need docker build outside sandbox" \
  -- docker build .
```

`require_escalated`を指定すると、結果JSONの`escalation_requested`が`true`になり、CLIは同じく実コマンドを実行します。別の承認フローと連携する際はこの値をトリガーに利用できます。

## ビルド

```bash
# プロジェクトルートから
./scripts/build_shell.sh  # 無い場合は go build でも可
go build -o bin/shell ./cmd/cli/shell
```

## 使用例

### execute操作

```bash
# -- の後ろにVec<String>として渡す
go run ./cmd/cli/shell -operation=execute -workdir=app -- bash -lc "ls -a"

# パイプラインを用いた記述
go run ./cmd/cli/shell -operation=execute -workdir=. -- bash -lc "ls -a | grep .git"

# --commandを複数指定することも可能
go run ./cmd/cli/shell -operation=execute \
  -command bash -command -lc -command "ls -a" \
  -base-dir=. -workdir=.
```

出力例:

```json
{
  "command": ["bash", "-lc", "ls"],
  "base_dir": "/home/user/devbox",
  "workdir": "/home/user/devbox",
  "stdout": "cmd\ninternal\n",
  "stderr": "",
  "exit_code": 0,
  "success": true,
  "timed_out": false,
  "duration_ms": 45,
  "timeout_ms": 60000,
  "sandbox_permissions": "use_default",
  "escalation_requested": false
}
```

### list_allowed操作

```bash
go run ./cmd/cli/shell -operation=list_allowed
```

```json
{
  "commands": ["awk", "bash", "bun", "cargo", ... ]
}
```

## 注意事項

- `-workdir`が`-base-dir`外を指すとエラーになります。
- タイムアウトを超えると`exit_code=-1`, `timed_out=true`でJSONが返り、CLIの終了コードは1になります。
- コマンド終了コードが非0の場合、CLIも同じ終了コードを返すので、他のスクリプトから容易に検出できます。

## 参考資料

- `cmd/cli/arithmetic-calculator/README.md` & `cmd/cli/service-implementing-viewer/README.md` をベースにしたCLIドキュメント構成
- Codex `shell` 実装 (`/home/user/package_references/codex/codex-rs`)の`SandboxPermissions`仕様
