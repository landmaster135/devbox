# Memos Utility CLI

`web-summary-*` / `movie-summary-*` の Markdown から、Memos にメモを作成するための CLI です。

## 機能

- `create-web-clip`: `web-summary-YYYYMMDD-hhmmss-<slug>.md` からメモ作成
- `create-movie-clip`: `movie-summary-YYYYMMDD-hhmmss-<slug>.md` からメモ作成

作成時の固定値:
- `visibility=PRIVATE`
- `state=NORMAL`
- `pinned=false`

`content-file` のファイル名から `displayTime`（`+09:00`）を自動生成します。

`-attachments` を指定した場合は、メモ作成前に全ファイルの存在・読込可能性を検証します。1件でも不正なパスがある場合、メモは作成しません。

## インストール

```bash
go build -o bin/memos-utility ./cmd/cli/memos-utility
```

## 使用方法

```bash
go run ./cmd/cli/memos-utility \
  -operation=<operation> \
  -base-url=<https://memos.example.com> \
  -api-token=<token> \
  -content-file=<path> \
  [-attachments=<path1,path2,...>]
```

## フラグ一覧

| フラグ | 説明 | 必須 | デフォルト |
|---|---|---|---|
| `-operation` | 実行操作（`create-web-clip`, `create-movie-clip`） | 必須 | - |
| `-base-url` | Memos のベースURL | 必須 | - |
| `-api-token` | Memos API の Bearer トークン | 必須 | - |
| `-content-file` | メモ本文を読み込む Markdown ファイルパス | 必須 | - |
| `-attachments` | 添付するローカルファイルパス（カンマ区切り） | 任意 | 空 |
| `-timeout` | HTTP タイムアウト秒 | 任意 | `30` |
| `-help`, `-h` | ヘルプ表示 | 任意 | `false` |

## content-file の命名制約

- `-operation=create-web-clip`
  - `web-summary-YYYYMMDD-hhmmss-<slug>.md` のみ受け付けます。
- `-operation=create-movie-clip`
  - `movie-summary-YYYYMMDD-hhmmss-<slug>.md` のみ受け付けます。

## 使用例

Web クリップを作成:

```bash
go run ./cmd/cli/memos-utility \
  -operation=create-web-clip \
  -base-url="$MEMOS_BASE_URL" \
  -api-token="$MEMOS_TOKEN" \
  -content-file=/tmp/web-summary-20240719-231059-palworld-steam-dedicated-server.md
```

Movie クリップを作成し、添付を追加:

```bash
go run ./cmd/cli/memos-utility \
  -operation=create-movie-clip \
  -base-url="$MEMOS_BASE_URL" \
  -api-token="$MEMOS_TOKEN" \
  -content-file=/tmp/movie-summary-20260319-055716-trump-masako-diplomacy.md \
  -attachments=/tmp/shot1.png,/tmp/shot2.png
```

## 出力例

成功時:

```json
{
  "operation": "create-web-clip",
  "displayTime": "2024-07-19T23:10:59+09:00",
  "memo": {
    "name": "memos/abc123",
    "displayTime": "2024-07-19T23:10:59+09:00",
    "visibility": "PRIVATE",
    "state": "NORMAL"
  },
  "attachments": [
    "/tmp/shot1.png",
    "/tmp/shot2.png"
  ],
  "setMemoAttachments": {
    "name": "memos/abc123"
  }
}
```

エラー時:

```text
エラー: create-web-clip の content-file は web-summary-YYYYMMDD-hhmmss-<slug>.md 形式のみ指定できます: invalid.md
```
