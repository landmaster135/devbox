# shellツールメモ
`openai/codex`のリポジトリ内の`codex-rs`ディレクトリ内で実装されているshellツールを参考に作るためのメモである。

## 役割
- Codex が shell ツール呼び出し時に `execvp()` へ渡すコマンドを厳密に制御し、必要に応じてサンドボックス外実行（エスカレーション）を要求できる。
- 通常は `"bash", "-lc", "<任意のシェル文字列>"` の 3 要素で実行する。

## パラメータ
- `command: Vec<String>` — 必須。実行するコマンド。ベクタの各要素がそのまま `execvp()` の argv に並ぶ。
- `workdir: Option<String>` — 任意。指定するとそのディレクトリで実行される（`cd` 不要）。
- `timeout_ms: Option<u64>` — 任意。ミリ秒単位の実行上限。超過時はプロセスが強制終了。
- `sandbox_permissions: Option<SandboxPermissions>` — 任意。`use_default`（既定サンドボックス）か `require_escalated`（サンドボックス外実行要求）を指定。
- `justification: Option<String>` — `sandbox_permissions` に `require_escalated` を指定した場合に必須となる 1 文の説明。

## 動作メモ
- Windows では PowerShell を推奨し、`["powershell.exe", "-Command", ...]` 形で指定する。
- UNIX 系 OS では `bash -lc` で文字列コマンドを渡すのが基本。`workdir` を忘れない。
- `require_escalated` の場合、MCP から人間承認や `.rules` 適用が入り、許可されればサンドボックス外で同じコマンドが実行される。
