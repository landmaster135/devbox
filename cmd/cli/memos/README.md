# Memos CLI

Memos API（`/api/v1`）を操作するCLIツールです。

## 機能

- `create-memo`: メモを作成
- `get-memo`: 単一メモを取得
- `list-memos`: メモ一覧を取得
- `update-memo`: 既存メモを更新（UpdateMemo）

## インストール

```bash
go build -o bin/memos ./cmd/cli/memos
```

## 使用方法

```bash
go run ./cmd/cli/memos \
  -operation=<operation> \
  -base-url=<https://memos.example.com> \
  -api-token=<token> \
  [operationごとのオプション]
```

## 共通オプション

| オプション | 説明 | 必須 |
|---|---|---|
| `-operation` | 実行する操作（`create-memo`, `get-memo`, `list-memos`, `update-memo`） | 必須 |
| `-base-url` | Memos のベースURL（例: `https://memos.example.com`） | 必須 |
| `-api-token` | Bearer トークン | 必須 |
| `-timeout` | HTTPタイムアウト秒（デフォルト: 30） | 任意 |
| `-help`, `-h` | ヘルプ表示 | 任意 |

## operation別オプション

### create-memo

| オプション | 説明 | 必須 |
|---|---|---|
| `-content` | メモ本文（`-content-file` とどちらか一方必須） | 条件付き必須 |
| `-content-file` | メモ本文を読み込むファイルパス（`-content` とどちらか一方必須） | 条件付き必須 |
| `-memo-id` | 作成時に指定する memoId | 任意 |
| `-visibility` | 公開範囲（`PRIVATE`, `PROTECTED`, `PUBLIC`） | 任意 |
| `-state` | 状態（`NORMAL`, `ARCHIVED`） | 任意 |
| `-pinned` | ピン留め（`true/false`） | 任意 |
| `-display-time` | 表示日時（RFC3339） | 任意 |

### get-memo

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | 取得対象の memo 識別子（例: `memo-123` / `memos/memo-123`） | 必須 |

### list-memos

| オプション | 説明 | 必須 |
|---|---|---|
| `-page-size` | 取得件数（デフォルト: 20） | 任意 |
| `-page-token` | ページトークン | 任意 |
| `-state` | 状態フィルタ（`NORMAL`, `ARCHIVED`） | 任意 |
| `-order-by` | 並び順（例: `update_time desc`） | 任意 |

### update-memo

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | 更新対象の memo 識別子 | 必須 |
| `-content` | 更新後の本文（`-content-file` とどちらか一方必須） | 条件付き必須 |
| `-content-file` | 更新後本文を読み込むファイルパス（`-content` とどちらか一方必須） | 条件付き必須 |
| `-visibility` | 更新後の公開範囲 | 任意 |
| `-state` | 更新後の状態 | 任意 |
| `-pinned` | 更新後のピン留め（`true/false`） | 任意 |
| `-update-mask` | 更新対象フィールド（例: `content,visibility`） | 任意 |

## 使用例

メモ作成
```bash
go run ./cmd/cli/memos \
  -operation=create-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -content="今日のメモ" \
  -visibility=PRIVATE
```

メモ取得
```bash
go run ./cmd/cli/memos \
  -operation=get-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123
```

メモ一覧
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -page-size=20 \
  -state=NORMAL \
  -order-by="update_time desc"
```

メモ更新
```bash
go run ./cmd/cli/memos \
  -operation=update-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -content="更新後の本文" \
  -visibility=PUBLIC
```

メモ作成（改行を含む本文をファイルから渡す）
```bash
go run ./cmd/cli/memos \
  -operation=create-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -content-file="./memo.md" \
  -visibility=PRIVATE
```

## 出力

成功時は API レスポンスを整形済み JSON で `stdout` に出力します。失敗時は `stderr` に `エラー: ...` を出力し、終了コード `1` で終了します。

## 参考リンク

- CreateMemo: https://usememos.com/docs/api/memoservice/CreateMemo
- GetMemo: https://usememos.com/docs/api/memoservice/GetMemo
- ListMemos: https://usememos.com/docs/api/memoservice/ListMemos
- UpdateMemo: https://usememos.com/docs/api/memoservice/UpdateMemo
