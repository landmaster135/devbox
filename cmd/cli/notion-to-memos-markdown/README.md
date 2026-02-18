# notion-to-memos-markdown

`notion-to-memos-markdown` は、Content JSON と `con_id.md` 群をもとに、`category` 単位で Markdown を仕分けコピーする CLI ツールです。

## 機能

- `--operation=distribute-files` をサポート
- `--page-type=content` の Content JSON を読み込み
- `--src-body-dir` 内の `con_id.md` を探索
- `--out-dir/<category>/` へコピー
- `category` が空の場合は `uncategorized` に振り分け

## フラグ

| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `--operation` | 必須 | 操作タイプ。`distribute-files` を指定 |
| `--page-type` | 必須 | ページタイプ。`content` を指定 |
| `--src-json-path` | 必須 | Content JSON ファイルのパス |
| `--src-body-dir` | 必須 | `con_id.md` があるディレクトリ |
| `--out-dir` | 必須 | カテゴリ別の出力先ルートディレクトリ |
| `-help`, `-h` | 任意 | ヘルプ表示 |

## 使用例

```bash
go run ./cmd/cli/notion-to-memos-markdown \
  --operation=distribute-files \
  --page-type=content \
  --src-json-path=$HOME/devbox/cmd/cli/notion-to-memos-markdown/tmp/contents.json \
  --src-body-dir=$HOME/path/to/dir \
  --out-dir=/tmp/notion-memos-out
```

## 出力例

成功時:

```text
処理完了
JSON基準: 総件数=120, コピー成功=98, 未検出=15, スキップ=7
src-body-dir基準: 総md件数=130, JSON対応=110, JSON未対応=20
```

エラー時:

```text
エラー: 未対応のpage-typeです: artifact
```
