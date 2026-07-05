# web-clipper

Web記事要約用のMarkdownに、指定URLのリンク行を挿入したり、添付ファイルを要約用の命名規則へリネームしたりするCLIツールです。

## 機能

- `--operation=patch-markdown`
  - `--top-heading-level` で指定した見出しレベルの最初の見出し直後に `- [target-title](target-url)` を追加します。
  - 入力は `--src-markdown-content` または `--src-markdown-file` で指定します（同時指定不可）。
  - 処理後のMarkdownを `--out-file-path` に出力します。
  - 入力本文に見出しレベル4以上（`####` 以降）が含まれる場合はエラーになります。
- `--operation=rename-attachments`
  - `--src-dir` で指定したディレクトリ直下の通常ファイルを `web-summary-YYYYMMDD-hhmmss-<slug>_<serial_number>.<extension>` へリネームします。
  - `--time` 指定時は更新時刻の昇順、`--name` 指定時はファイル名の昇順で採番します（同時指定不可）。
  - `--slug` は英小文字、数字、半角ハイフンのみ使用できます。
  - ディレクトリは対象外です。リネーム先が既に存在する場合はエラーになります。

## フラグ

| フラグ | 必須 | 説明 |
| --- | --- | --- |
| `--operation` | はい | 実行する操作。`patch-markdown`, `rename-attachments` に対応。 |
| `--target-title` | `patch-markdown` で必須 | 追加するリンク表示テキスト。 |
| `--target-url` | `patch-markdown` で必須 | 追加するリンクURL。 |
| `--src-markdown-content` | いいえ | 処理対象のMarkdown本文（`--src-markdown-file` と同時指定不可）。 |
| `--src-markdown-file` | いいえ | 処理対象のMarkdownファイルパス（`--src-markdown-content` と同時指定不可）。 |
| `--out-file-path` | `patch-markdown` で必須 | 出力先Markdownファイルパス。 |
| `--top-heading-level` | `patch-markdown` で必須 | 追加位置の基準となる見出しレベル（`1` 以上）。 |
| `--src-dir` | `rename-attachments` で必須 | リネーム対象ディレクトリ。 |
| `--slug` | `rename-attachments` で必須 | 出力ファイル名に含めるslug。英小文字、数字、半角ハイフンのみ使用可能。 |
| `--start` | `rename-attachments` で必須 | 採番開始番号（`0` 以上）。 |
| `--digits` | いいえ | 採番のゼロ埋め桁数（デフォルト: `2`）。 |
| `--time` | `rename-attachments` でどちらか必須 | 更新時刻の昇順でリネーム（`--name` と同時指定不可）。 |
| `--name` | `rename-attachments` でどちらか必須 | ファイル名の昇順でリネーム（`--time` と同時指定不可）。 |
| `--json` | いいえ | JSON形式で出力。 |
| `--verbose` | いいえ | 詳細を出力。 |
| `-help`, `-h` | いいえ | ヘルプ表示。 |

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

`rename-attachments` on default.
```bash
go run cmd/cli/web-clipper/main.go \
	--operation=rename-attachments \
	--src-dir=./tmp/attachments \
	--slug=openai-blog \
	--start=11 \
	--digits=2 \
	--time \
	--json
```

```bash
go run cmd/cli/web-clipper/main.go \
  --operation=rename-attachments \
  --src-dir=./tmp/attachments \
  --slug=openai-blog \
  --start=1 \
  --name \
  --verbose
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

```text
リネームしました: 2件
./tmp/attachments/a.png -> ./tmp/attachments/web-summary-20260705-171234-openai-blog_01.png
./tmp/attachments/b.jpg -> ./tmp/attachments/web-summary-20260705-171234-openai-blog_02.jpg
```

JSON出力:

```json
{
  "renamed_count": 1,
  "files": [
    {
      "from": "./tmp/attachments/a.png",
      "to": "./tmp/attachments/web-summary-20260705-171234-openai-blog_01.png"
    }
  ]
}
```

失敗時（標準エラー出力）:

```text
エラー: 見出しレベル4以上（#### 以降）は使用できません
```

```text
エラー: --slug は英小文字、数字、半角ハイフンのみ使用できます
```
