# notion-to-memos-markdown

`notion-to-memos-markdown` は、Notion 由来の JSON（Content / Artifact）と `con_id.md` 群をもとに Markdown を加工する CLI ツールです。

## 機能

- `--operation=distribute-files`
  - `--page-type=content` の Content JSON を読み込み
  - `--src-body-dir` 内の `con_id.md` を探索
  - `--out-dir/<category>/` へコピー
  - `category` が空の場合は `uncategorized` に振り分け
- `--operation=craft-markdown`
  - `--page-type=content` / `artifact` / `task` を対象に処理
  - `con_id` の数値範囲（`--con_number_start` から `--con_number_end`）を対象に処理
  - 入力 body は `src-body-dir/<con_id>.md` を優先し、無い場合は `src-body-dir/<con_id>_*.md`（分割ファイル）を対象にする
  - `--skips-no-src-body=true` 指定時はコピー元Markdownがないページをスキップ（未指定/false時は空Markdownを作成して加工）
  - `content`:
    - `--category` 指定時は一致する category の Content のみ処理
    - コピー元 Markdown に front matter が既に存在する場合、H1 見出しとタグの新規追加を行わない（front matter へのキー追加は実施）
    - `page_title` を H1 見出しとして先頭に追加
    - front matter（`title`, `bought_at`, `score_of_100`, `price_yen`, `con_id`, `url`）を追加
    - `tags.md` の `## Frequent Tags` セクションを使ってタグを解決し、分類タグ群を追加
  - `artifact`:
    - `output_url` を本文先頭に追加した上で、`add-tags` と `add-heading1` 相当の加工を実行
    - `tags.md` の `## Frequent Tags` と `## Artifact` を使ってタグを解決し、`#91-backup/tool-migration/202602-notion` を必ず付与
  - `task`:
    - 出力は `out-dir` 直下に作成
    - `--src-resource-dir` 配下を再帰走査し、対象 Task の添付ファイルを取得する
    - 添付ファイルは `out-resource-dir` 直下へ複製する
    - `src-body-dir/<con_id>.md` がある場合は見出し分割せず、1入力につき1ファイル（`<con_id>_01.md`）を加工対象にする
    - `src-body-dir/<con_id>.md` が無い場合のみ `src-body-dir/<con_id>_*.md` を加工対象にする
    - どの入力でも各出力ファイルに `add-heading1(page_title)` と `add-tags` 相当の加工を適用
    - `priority` と `status_id` から Task用タグを付与し、`#91-backup/tool-migration/202602_notion` を必ず付与
    - `powered_artifacts[].page_title` が定義済みマップに一致した場合、対応する追加タグを付与
    - `status_id` に応じて `done_at_start` または `updated_at` を採用し、`yyyyMMddhhmmss_<index>.md` へリネーム
    - 添付ファイル名は `con_id_<number>(.<ext>)` または `con_id_<index>_<number>(.<ext>)` のみ受け付ける（規約外はエラー）
    - 添付ファイルは対応する Markdown の `yyyyMMddhhmmss_<index>` を使って `yyyyMMddhhmmss_<index>_<number>(.<ext>)` へリネームする（拡張子なしも許容）
    - 添付ファイルの同名衝突は回避せず、出力先が既に存在する場合はエラーにする
    - リネーム先が `out-dir` 内で衝突する場合は、タイムスタンプの秒（`ss`）を 1 秒ずつ増やして非衝突名を採用
- `--operation=check-body-length`
  - `--src-body-dir` 配下を再帰的に走査して全ファイルを対象にする
  - ファイル本文の文字数をルーン数でカウントする
  - `--threshold` より大きい文字数のファイルを列挙する
  - 閾値超過ファイルの総数を表示する
- `--operation=grep-str`
  - `--src-body-dir` 配下を再帰的に走査して全ファイルを対象にする
  - `--target-str` を含むファイルを列挙する
  - 該当ファイルの総数を表示する
- `--operation=rename-bodies-by-category-id`
  - `--page-type=content` の Content JSON を読み込み、`--con_number_start`〜`--con_number_end` の範囲にある Content の `category_id -> con_id` マップを作成
  - `--src-resource-dir` 配下を再帰的に走査して全ファイルを対象にする
  - ファイル名の先頭プレフィックスが `category_id` と一致する場合、同じ部分を `con_id` へ置換してリネームする
  - 処理件数（対象Content件数、対象ファイル総数、リネーム成功数、スキップ数）を表示する
- `--operation=migrate-to-memos`
  - `--src-body-dir` 配下を再帰走査して Markdown ファイル（`.md`）を対象にする
  - body ファイルごとに Memos API でメモを作成し、レスポンスから `memo.name` を取得する
  - body ファイル名（拡張子除去）を `con_id` として扱う
  - `--src-resource-dir` 配下を再帰走査し、ファイル名が同じ `con_id` プレフィックスのリソースを作成・添付する
  - 一致リソースが 0 件の場合は添付 API を呼ばずにスキップする
  - 処理中の進捗（`con_id` ごとの開始/完了/スキップ）を標準エラー出力へ逐次表示する

## フラグ

| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `--operation` | 必須 | 操作タイプ。`distribute-files`、`craft-markdown`、`check-body-length`、`grep-str`、`rename-bodies-by-category-id`、`migrate-to-memos` |
| `--page-type` | `distribute-files`/`craft-markdown`/`rename-bodies-by-category-id`/`migrate-to-memos`で必須 | ページタイプ。`distribute-files`/`rename-bodies-by-category-id`/`migrate-to-memos` は `content`、`craft-markdown` は `content` / `artifact` / `task` |
| `--base-url` | `migrate-to-memos`で必須 | Memos API のベースURL（例: `https://memos.example.com`） |
| `--api-token` | `migrate-to-memos`で必須 | Memos API の Bearer トークン |
| `--category` | `craft-markdown`で任意 | 対象 category。指定時は一致するページのみ処理 |
| `--skips-no-src-body` | `craft-markdown`で任意 | `true` ならコピー元Markdown未存在の Content をスキップ（デフォルト: `false`） |
| `--con_number_start` | `craft-markdown`/`rename-bodies-by-category-id`で必須 | 対象 `con_id` 範囲の開始番号（1以上） |
| `--con_number_end` | `craft-markdown`/`rename-bodies-by-category-id`で必須 | 対象 `con_id` 範囲の終了番号（1以上） |
| `--threshold` | `check-body-length`で必須 | 文字数の閾値（0以上） |
| `--target-str` | `grep-str`で必須 | 検索対象の文字列 |
| `--src-json-file` | `distribute-files`/`craft-markdown`/`rename-bodies-by-category-id`で必須 | 入力 JSON ファイルのパス（`craft-markdown` では `content`/`artifact` の両方を受け付ける。`--src-json-path` も利用可能） |
| `--src-body-dir` | `distribute-files`/`craft-markdown`/`check-body-length`/`grep-str`/`migrate-to-memos`で必須 | 入力ディレクトリ |
| `--src-resource-dir` | `craft-markdown`(`task`)/`rename-bodies-by-category-id`/`migrate-to-memos`で必須 | リソース入力ディレクトリ |
| `--out-dir` | `distribute-files`/`craft-markdown`で必須 | 出力先ルートディレクトリ |
| `--out-resource-dir` | `craft-markdown`(`task`)で必須 | Task添付リソースの出力先ディレクトリ |
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

`craft-markdown` (artifact):

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=craft-markdown \
  --page-type=artifact \
  --con_number_start=100 \
  --con_number_end=300 \
  --src-json-path=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/artifacts.json \
  --src-body-dir=$HOME/path/to/dir \
  --out-dir=/tmp/notion-artifacts-crafted
```

`craft-markdown` (task):

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=craft-markdown \
  --page-type=task \
  --con_number_start=2000 \
  --con_number_end=2200 \
  --src-json-path=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/tasks.json \
  --src-body-dir=$HOME/path/to/dir \
  --src-resource-dir=$HOME/path/to/resource_dir \
  --out-dir=/tmp/notion-tasks-crafted \
  --out-resource-dir=/tmp/notion-tasks-crafted-resources
```

`check-body-length`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=check-body-length \
  --src-body-dir=$HOME/path/to/dir \
  --threshold=1000
```

`grep-str`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=grep-str \
  --src-body-dir=$HOME/path/to/dir \
  --target-str=TODO
```

`rename-bodies-by-category-id`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=rename-bodies-by-category-id \
  --page-type=content \
  --con_number_start=2000 \
  --con_number_end=2200 \
  --src-json-file=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/contents.json \
  --src-resource-dir=$HOME/path/to/resource_dir
```

`migrate-to-memos`:

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=migrate-to-memos \
  --page-type=content \
  --base-url=https://memos.example.com \
  --api-token=$MEMOS_TOKEN \
  --src-body-dir=$HOME/path/to/body_dir \
  --src-resource-dir=$HOME/path/to/resource_dir
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
対象件数=42, 加工成功=39, スキップ=3
```

```text
処理完了
対象ファイル総数=4
閾値超過ファイル総数=2
閾値超過ファイル一覧:
/tmp/notion-body/CO0001.md: 1201
/tmp/notion-body/nested/notes.txt: 1490
```

```text
処理完了
対象ファイル総数=4
該当ファイル総数=2
該当ファイル一覧:
/tmp/notion-body/CO0001.md
/tmp/notion-body/nested/notes.txt
```

```text
処理完了
対象Content件数=23
対象ファイル総数=75
リネーム成功=31
スキップ(プレフィックスなし)=2
スキップ(マップ未対応)=42
```

```text
処理完了
対象body件数=120
メモ作成成功=120
添付ファイル総数=245
添付スキップ(リソースなし)=18
```

エラー時:

```text
エラー: 未対応のpage-typeです: unknown
```
