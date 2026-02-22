# web-clipper

Web記事要約用のMarkdownに、指定URLのリンク行を挿入して出力するCLIツールです。

## 機能

- `--operation=patch-markdown`
  - `--top-heading-level` で指定した見出しレベルの最初の見出し直後に `- [target-title](target-url)` を追加します。
  - 入力は `--src-markdown-content` または `--src-markdown-file` で指定します（同時指定不可）。
  - 処理後のMarkdownを `--out-file-path` に出力します。
  - 入力本文に見出しレベル4以上（`####` 以降）が含まれる場合はエラーになります。

## フラグ

- `--operation`（必須）: 実行する操作。現状は `patch-markdown` のみ対応。
- `--target-title`（必須）: 追加するリンク表示テキスト。
- `--target-url`（必須）: 追加するリンクURL。
- `--src-markdown-content`: 処理対象のMarkdown本文（`--src-markdown-file` と同時指定不可）。
- `--src-markdown-file`: 処理対象のMarkdownファイルパス（`--src-markdown-content` と同時指定不可）。
- `--out-file-path`（必須）: 出力先Markdownファイルパス。
- `--top-heading-level`（必須）: 追加位置の基準となる見出しレベル（`1` 以上）。
- `-help`, `-h`: ヘルプ表示。

## 使用方法

```bash
go run cmd/cli/web-clipper/main.go \
  --operation=patch-markdown \
  --target-title="OpenAI Blog" \
  --target-url="https://openai.com/blog" \
  --src-markdown-content=$'## 記事タイトル 要約\n\n### 見出し1\n本文\n' \
  --out-file-path=./tmp/out.md \
  --top-heading-level=2
```

```bash
go run cmd/cli/web-clipper/main.go \
  --operation=patch-markdown \
  --target-title="OpenAI Blog" \
  --target-url="https://openai.com/blog" \
  --src-markdown-file=./tmp/in.md \
  --out-file-path=./tmp/out.md \
  --top-heading-level=2
```

## 変換例

Before:

```md
## <記事タイトル> 要約

### <見出し1: 固有名詞を含める>
<内容>
```

After:

```md
## <記事タイトル> 要約
- [target-title](target-url)

### <見出し1: 固有名詞を含める>
<内容>
```

## 出力例

成功時（標準出力）:

```text
出力しました: ./tmp/out.md
```

失敗時（標準エラー出力）:

```text
エラー: 見出しレベル4以上（#### 以降）は使用できません
```
