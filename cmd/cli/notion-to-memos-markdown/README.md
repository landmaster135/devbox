# notion-to-memos-markdown

`notion-to-memos-markdown` は、Content JSON と `con_id.md` 群をもとに Markdown を加工する CLI ツールです。

## 機能

- `--operation=distribute-files`
  - `--page-type=content` の Content JSON を読み込み
  - `--src-body-dir` 内の `con_id.md` を探索
  - `--out-dir/<category>/` へコピー
  - `category` が空の場合は `uncategorized` に振り分け
- `--operation=craft-markdown`
  - `con_id` の数値範囲（`--con_number_start` から `--con_number_end`）を対象に処理
  - `--category` 指定時は一致する category の Content のみ処理
  - `--skips-no-src-body=true` 指定時はコピー元Markdownがない Content をスキップ（未指定/false時は空Markdownを作成して加工）
  - `page_title` を H1 見出しとして先頭に追加
  - front matter（`bought_at`, `score_of_100`, `price_yen`, `con_id`, `url`）を追加
  - `tags.md` の `## Frequent Tags` セクションを使ってタグを解決し、分類タグ群を追加
- `--operation=check-body-length`
  - `--src-body-dir` 配下を再帰的に走査して全ファイルを対象にする
  - ファイル本文の文字数をルーン数でカウントする
  - `--threshold` より大きい文字数のファイルを列挙する
  - 閾値超過ファイルの総数を表示する

## フラグ

| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `--operation` | 必須 | 操作タイプ。`distribute-files`、`craft-markdown`、`check-body-length` |
| `--page-type` | `distribute-files`/`craft-markdown`で必須 | ページタイプ。`content` を指定 |
| `--category` | `craft-markdown`で任意 | 対象 category。指定時は一致する Content のみ処理 |
| `--skips-no-src-body` | `craft-markdown`で任意 | `true` ならコピー元Markdown未存在の Content をスキップ（デフォルト: `false`） |
| `--con_number_start` | `craft-markdown`で必須 | 対象 `con_id` 範囲の開始番号（1以上） |
| `--con_number_end` | `craft-markdown`で必須 | 対象 `con_id` 範囲の終了番号（1以上） |
| `--threshold` | `check-body-length`で必須 | 文字数の閾値（0以上） |
| `--src-json-file` | `distribute-files`/`craft-markdown`で必須 | Content JSON ファイルのパス |
| `--src-body-dir` | 必須 | 入力ディレクトリ |
| `--out-dir` | `distribute-files`/`craft-markdown`で必須 | 出力先ルートディレクトリ |
| `-help`, `-h` | 任意 | ヘルプ表示 |

## 使用例

`distribute-files`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=distribute-files \
  --page-type=content \
  --src-json-file=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/contents.json \
  --src-body-dir=$HOME/path/to/dir \
  --out-dir=/tmp/notion-memos-out
```

`craft-markdown`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=craft-markdown \
  --page-type=content \
  --category=software \
  --skips-no-src-body=true \
  --con_number_start=2000 \
  --con_number_end=2200 \
  --src-json-file=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/contents.json \
  --src-body-dir=$HOME/path/to/dir \
  --out-dir=/tmp/notion-memos-crafted
```

`check-body-length`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=check-body-length \
  --src-body-dir=$HOME/path/to/dir \
  --threshold=1000
```

## 出力例

成功時:

```text
処理完了
JSON基準: 総件数=120, コピー成功=98, 未検出=15, スキップ=7
src-body-dir基準: 総md件数=130, JSON対応=110, JSON未対応=20
```

```text
処理完了
対象件数=42, 加工成功=42
```

```text
処理完了
対象ファイル総数=4
閾値超過ファイル総数=2
閾値超過ファイル一覧:
/tmp/notion-body/CO0001.md: 1201
/tmp/notion-body/nested/notes.txt: 1490
```

エラー時:

```text
エラー: 未対応のpage-typeです: artifact
```
