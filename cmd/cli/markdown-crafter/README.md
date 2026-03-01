# markdown-crafter

Markdown ファイルを加工する CLI ツールです。`--operation` で下記の 5 操作を切り替えます。

- `split-headings`: 指定見出しレベルごとに分割して複数ファイル出力
- `add-front-matter`: `--kv` で指定した key-value をフロントマターへ追加
- `add-tags`: `--tags` で指定したタグを `#tag1 #tag2` 形式で追加（`--file-path` または `--dir-path`）
- `delete-empty-files`: 指定ディレクトリ内の条件一致 Markdown ファイルを削除
- `add-heading1`: 本文の先頭または末尾へ見出し1を追加

## フラグ一覧

| フラグ | 必須 | デフォルト | 説明 |
| --- | --- | --- | --- |
| `--operation` | 必須 | なし | `split-headings` / `add-front-matter` / `add-tags` / `delete-empty-files` / `add-heading1` |
| `--file-path` | `split-headings` / `add-front-matter` / `add-heading1` で必須、`add-tags` で `--dir-path` と排他で必須 | なし | 対象の Markdown ファイル |
| `--dir-path` | `add-tags` で `--file-path` と排他で必須 | なし | 対象の Markdown ディレクトリ（直下の `.md` を処理） |
| `--directory-path` | `delete-empty-files` で必須 | なし | 対象の Markdown ディレクトリ |
| `--heading-level` | `split-headings` で必須 | `0` | 分割対象の見出しレベル（1-6） |
| `--heading-text` | `add-heading1` で必須 | なし | 追加する見出し1のテキスト |
| `--heading-position` | `add-heading1` で必須 | なし | 追加位置（`head` または `tail`） |
| `--output-dir` | `split-headings` で必須 | なし | 分割ファイルの出力先ディレクトリ |
| `--kv` | `add-front-matter` で必須（1件以上） | なし | 追加する key-value（`key=value`）複数指定可 |
| `--tags` | `add-tags` で必須 | なし | カンマ区切りタグ（例: `go,markdown`） |
| `--help`, `-h` | 任意 | `false` | ヘルプ表示 |

## 使用方法

```bash
go run ./cmd/cli/markdown-crafter --operation <operation> [--file-path <path> | --dir-path <dir>] [flags...]
```

## 使用例

### split-headings

```bash
go run ./cmd/cli/markdown-crafter \
  --operation split-headings \
  --file-path ./sample.md \
  --heading-level 2 \
  --output-dir ./out
```

`sample.md`:

```md
# aaa

## bbb

test1

## ccc

test2
```

出力:

- `./out/001.md`
- `./out/002.md`

### add-front-matter

```bash
go run ./cmd/cli/markdown-crafter \
  --operation add-front-matter \
  --file-path ./sample.md \
  --kv title=example \
  --kv author=alice
```

### add-tags

```bash
go run ./cmd/cli/markdown-crafter \
  --operation add-tags \
  --file-path ./sample.md \
  --tags go,markdown
```

```bash
go run ./cmd/cli/markdown-crafter \
  --operation add-tags \
  --dir-path ./notes \
  --tags go,markdown
```

### delete-empty-files

```bash
go run ./cmd/cli/markdown-crafter \
  --operation delete-empty-files \
  --directory-path ./notes
```

次のいずれかに一致する `.md` ファイルを削除します。

- 空文字
- `# Miscellaneous notes\n- `

### add-heading1

```bash
go run ./cmd/cli/markdown-crafter \
  --operation add-heading1 \
  --file-path ./sample.md \
  --heading-text 概要 \
  --heading-position head
```

## 出力例

### 成功時

```text
split-headings: 2 ファイルを出力しました
- out/001.md
- out/002.md
```

```text
add-front-matter: ./sample.md を更新しました (3 キー)
```

```text
add-tags: ./sample.md にタグを追加しました (#go #markdown)
```

```text
add-tags: 2 ファイルにタグを追加しました (#go #markdown)
- ./notes/a.md
- ./notes/b.md
```

```text
delete-empty-files: 2 ファイルを削除しました
- ./notes/a.md
- ./notes/b.md
```

```text
add-heading1: ./sample.md に見出しを追加しました (head)
```

### エラー時

```text
エラー: --heading-level は 1 から 6 の範囲で指定してください
```

```text
エラー: --kv の形式が不正です: title (key=value 形式で指定してください)
```

```text
エラー: --heading-position には head, tail のいずれかを指定してください
```
