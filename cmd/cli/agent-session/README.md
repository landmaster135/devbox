# agent-session

`agent-session` は、Codex のセッション一覧を取得して表示するCLIツールです。

## フラグ

| フラグ | 必須 | デフォルト | 説明 |
| --- | --- | --- | --- |
| `--operation` | 必須 | なし | 実行する操作。現在は `retrieve-session` のみ対応。 |
| `--agent-type` | 必須 | なし | エージェントタイプ。現在は `codex` のみ対応。 |
| `--limit` | 任意 | `200` | 取得件数の上限。1以上を指定。 |
| `--start-date` | 任意 | なし | 開始日。形式は `yyyyMMdd`。 |
| `--end-date` | 任意 | なし | 終了日。形式は `yyyyMMdd`。 |
| `--agent-home-dir` | 任意 | `~/.codex` | エージェントホームディレクトリ。セッションは `<agent-home-dir>/sessions` から取得する。 |

## 使用方法

```bash
go run ./cmd/cli/agent-session \
  --operation=retrieve-session \
  --agent-type=codex
```

## 使用例

期間指定なし（最新200件）:

```bash
go run ./cmd/cli/agent-session \
  --operation=retrieve-session \
  --agent-type=codex
```

期間指定あり:

```bash
go run ./cmd/cli/agent-session \
  --operation=retrieve-session \
  --agent-type=codex \
  --limit=50 \
  --start-date=20260301 \
  --end-date=20260331
```

ホームディレクトリを明示:

```bash
go run ./cmd/cli/agent-session \
  --operation=retrieve-session \
  --agent-type=codex \
  --agent-home-dir=$HOME/.codex
```

## 出力例

成功時:

```text
UUID                                  Created at      Updated at      Branch   CWD                Conversation
11111111-2222-4333-8444-555555555555  2 minutes ago   19 seconds ago  task     /workspace/devbox   Codexのセッション一覧ってどうやって取得出来るの？
aaaaaaaa-bbbb-4ccc-8ddd-eeeeeeeeeeee  54 minutes ago  53 minutes ago  feature  /workspace/dotfiles Create Git commit message in English.
```

対象なし:

```text
対象セッションは見つかりませんでした。
```

エラー時（例: 日付フォーマット不正）:

```text
エラー: 設定の初期化に失敗しました: --start-date の形式が不正です: parsing time "2026-03-01" as "20060102": cannot parse "-03-01" as "01"
```
