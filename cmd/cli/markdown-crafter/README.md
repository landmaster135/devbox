# markdown-crafter

Markdown ファイルを加工する CLI ツールです。`--operation` で下記の 8 操作を切り替えます。

- `split-headings`: 指定見出しレベルごとに分割して複数ファイル出力
- `add-front-matter`: `--kv` で指定した key-value をフロントマターへ追加
- `add-tags`: `--tags` で指定したタグを `#tag1 #tag2` 形式で追加（`--file-path` または `--dir-path`）
- `delete-empty-files`: 指定ディレクトリ内の条件一致 Markdown ファイルを削除
- `add-heading1`: 本文の先頭または末尾へ見出し1を追加
- `replace-images`: Markdown 画像記法 `![alt](url)` を指定文字列へ置換
- `remove-heading-annotations`: 指定見出しレベルの `**見出し**` 注釈を `見出し` へ変換
- `remove-title-hash-tags`: 各ファイルの指定行範囲にあるハッシュタグ（例: `#Python`）を除去

## 内部構成

- CLI エントリーポイントは `cmd/cli/markdown-crafter/main.go` に集約。
- operation の振り分けは `internal/markdown_crafter/usecases/services.go` が担当。
- 各 operation の実処理は `internal/markdown_crafter/usecases/operations/*/service.go` に分離。
- front matter / tag などの共通ロジックは `internal/markdown_crafter/usecases/common/helpers.go` に集約。

## フラグ一覧

| フラグ | 必須 | デフォルト | 説明 |
| --- | --- | --- | --- |
| `--operation` | 必須 | なし | `split-headings` / `add-front-matter` / `add-tags` / `delete-empty-files` / `add-heading1` / `replace-images` / `remove-heading-annotations` / `remove-title-hash-tags` |
| `--file-path` | `split-headings` / `add-front-matter` / `add-heading1` / `replace-images` / `remove-heading-annotations` で必須、`add-tags` で `--dir-path` と排他で必須 | なし | 対象の Markdown ファイル |
| `--dir-path` | `add-tags` で `--file-path` と排他で必須、`delete-empty-files` / `remove-title-hash-tags` で必須 | なし | 対象の Markdown ディレクトリ（直下の `.md` を処理） |
| `--heading-level` | `split-headings` / `remove-heading-annotations` で必須 | `0` | 対象の見出しレベル（1-6） |
| `--start-line` | `remove-title-hash-tags` で必須 | `0` | ハッシュタグ除去の対象開始行（1始まり） |
| `--end-line` | `remove-title-hash-tags` で必須 | `0` | ハッシュタグ除去の対象終了行（1始まり、`--start-line` 以上） |
| `--heading-text` | `add-heading1` で必須 | なし | 追加する見出し1のテキスト |
| `--heading-position` | `add-heading1` で必須 | なし | 追加位置（`head` または `tail`） |
| `--replacement-text` | `replace-images` で必須 | なし | 画像記法を置換する文字列（例: `(添付画像)`） |
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
  --dir-path ./notes
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

### replace-images

```bash
go run ./cmd/cli/markdown-crafter \
  --operation replace-images \
  --file-path ./sample.md \
  --replacement-text "(添付画像)"
```

### remove-heading-annotations

```bash
go run ./cmd/cli/markdown-crafter \
  --operation remove-heading-annotations \
  --file-path ./sample.md \
  --heading-level 3
```

### remove-title-hash-tags

```bash
go run ./cmd/cli/markdown-crafter \
  --operation remove-title-hash-tags \
  --dir-path ./notes \
  --start-line 1 \
  --end-line 2
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

```text
replace-images: ./sample.md の画像記法 2 件を置換しました
```

```text
remove-heading-annotations: ./sample.md の見出し注釈 2 件を除去しました
```

```text
remove-title-hash-tags: 2 ファイルの 1 行目から 2 行目までのハッシュタグを除去しました
- ./notes/a.md
- ./notes/b.md
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

```text
エラー: --replacement-text は必須です (--operation=replace-images)
```

```text
エラー: --start-line は必須です (--operation=remove-title-hash-tags)
```
