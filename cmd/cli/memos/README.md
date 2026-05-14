# Memos CLI

Memos API（`/api/v1`）を操作するCLIツールです。

## 機能

- `create-memo`: メモを作成
- `get-memo`: 単一メモを取得
- `delete-memo`: 単一メモを削除
- `list-memos`: メモ一覧を取得
- `list-attachments`: 添付一覧を取得
- `update-memo`: 既存メモを更新（UpdateMemo）
- `update-tag`: 既存タグを新しいタグへ一括置換
- `patch-files`: ローカルファイルを添付として作成し、メモ添付を更新
- `list-memo-relations`: 対象メモのリレーション一覧を取得
- `add-memo-relations`: 対象メモのリレーションを追加/置換して更新

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
| `-operation` | 実行する操作（`create-memo`, `get-memo`, `delete-memo`, `list-memos`, `list-attachments`, `update-memo`, `update-tag`, `patch-files`, `list-memo-relations`, `add-memo-relations`） | 必須 |
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

### delete-memo

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | 削除対象の memo 識別子（例: `memo-123` / `memos/memo-123`） | 必須 |
| `-force` | `true` のとき強制削除を要求する（デフォルト: `false`） | 任意 |

### list-memos

| オプション | 説明 | 必須 |
|---|---|---|
| `-page-size` | 取得件数（デフォルト: 20） | 任意 |
| `-page-token` | ページトークン | 任意 |
| `-state` | 状態フィルタ（`NORMAL`, `ARCHIVED`） | 任意 |
| `-order-by` | 並び順（例: `update_time desc`） | 任意 |
| `-filter` | フィルタ条件（CEL形式。例: `visibility == "PUBLIC"`。`created_ts` / `updated_ts` の RFC3339 比較値は内部で Unix 秒へ変換） | 任意 |
| `-any-contents` | `content.contains` 検索キーワード（カンマ区切り。例: `meeting,study`） | 任意 |

補足:
- `-any-contents` を指定すると、キーワードごとに `content.contains("<keyword>")` を評価して結果を統合します。
- 複数キーワードで同一メモがヒットした場合は、メモID単位で重複排除して返却します。
- `-filter` と併用した場合は、`(<filter>) && content.contains("<keyword>")` をキーワードごとに評価します。

**filter で利用可能な主なフィールド**

Memos の `filter` は CEL 形式です。公式ドキュメント上で確認できる主なフィールドは以下です（インスタンスのバージョン差で一部差異あり）。

| フィールド | 型 | 例 |
|---|---|---|
| `content` | `string` | `content.contains("meeting")` |
| `visibility` | `string` | `visibility == "PUBLIC"` |
| `tag` / `tags` | `string` / `[]string` | `tag in ["work","project"]`, `"work" in tags` |
| `created_ts` | `int` (Unix timestamp) | `created_ts > 1700000000`, `created_ts > "2023-01-01T13:00:00Z"` |
| `updated_ts` | `int` (Unix timestamp) | `updated_ts > 1700000000`, `updated_ts >= "2023-01-01T13:00:00+09:00"` |
| `pinned` | `bool` | `pinned` |
| `has_task_list` | `bool` | `has_task_list` |
| `has_incomplete_tasks` | `bool` | `has_incomplete_tasks` |
| `has_link` | `bool` | `has_link` |
| `has_code` | `bool` | `has_code` |
| `creator_id` | `int` | `creator_id == 101` |

注意:
- `create_time_after(...)` や `visibilities` は、少なくとも一部環境では未サポートで `undeclared reference` エラーになります。
- `created_ts` / `updated_ts` の日時文字列は RFC3339/RFC3339Nano（タイムゾーン必須）で指定してください。例: `2023-01-01T13:00:00Z`
- CEL 形式記法は [CEL documentation](https://github.com/google/cel-spec) を参照してください。

### list-attachments

| オプション | 説明 | 必須 |
|---|---|---|
| `-page-size` | 取得件数（デフォルト: 20） | 任意 |
| `-page-token` | ページトークン | 任意 |
| `-order-by` | 並び順（例: `create_time desc`） | 任意 |
| `-filter` | フィルタ条件（CEL形式） | 任意 |

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
| `-updates-time` | `true` のとき `displayTime` を現在日時（UTC/RFC3339）で更新（結果として `updateTime` も更新される） | 任意（デフォルト: `false`） |

### update-tag

| オプション | 説明 | 必須 |
|---|---|---|
| `-src-tag` | 置換元タグ（例: `work` または `#work`） | 必須 |
| `-dest-tag` | 置換先タグ（例: `project` または `#project`） | 必須 |

補足:
- `update-tag` は `ListMemos` で `"<src-tag>" in tags` を満たすメモを取得し、本文中の `#src-tag` を `#dest-tag` に置換して `UpdateMemo`（`updateMask=content`）で更新します。
- `#tag-xxx` のように `src-tag` が接頭辞として含まれる別タグは置換しません（完全一致のみ置換）。

### patch-files

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | 添付対象の memo 識別子 | 必須 |
| `-files` | 添付するローカルファイルのパスをカンマ区切りで指定 | 必須 |
| `-replaces` | `true` なら新規添付のみで置換。`false`（デフォルト）なら既存添付を `ListMemoAttachments` で取得して保持したまま追加 | 任意 |

補足:
- `patch-files` は、指定された全ファイルの読み込みと MIME type 判定が成功した場合にのみ `CreateAttachment` を開始します。
- 1件でもファイル読み込みまたは MIME type 判定に失敗した場合、添付作成は一切行わずに中断します。

### list-memo-relations

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | リレーション取得対象の memo 識別子 | 必須 |

補足:
- `list-memo-relations` は内部で全ページを取得し、対象メモの既存リレーションをまとめて返します。
- 任意のメモにおけるレスポンスボディにおいて、自身がリンクしているメモだと`relations[n].memo`で自身のID、自身がリンクされているメモだと`relations[n].relatedMemo`で自身のIDが返却されます。

### add-memo-relations

| オプション | 説明 | 必須 |
|---|---|---|
| `-memo` | リレーション更新対象の memo 識別子 | 必須 |
| `-related-memos` | 追加対象の related memo 識別子をカンマ区切りで指定 | 必須 |
| `-replaces` | `true` なら既存リレーションを破棄して置換。`false`（デフォルト）なら既存を保持して追加 | 任意 |

補足:
- `add-memo-relations` はまず `ListMemoRelations` で既存リレーションを取得し、`SetMemoRelations` で更新します。
- 結果には `discardedRelations`（破棄された関係）と `addedRelations`（追加された関係）を含みます。

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

メモ削除
```bash
go run ./cmd/cli/memos \
  -operation=delete-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -force=true
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

メモ一覧（CELフィルタ）
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -filter='visibility == "PUBLIC"'
```

メモ一覧（CELフィルタ: 作成時刻＋公開範囲）
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -filter='created_ts > "2023-01-01T13:00:00Z" && visibility == "PUBLIC"'
```

メモ一覧（CELフィルタ: 複数タグのAND検索）
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -filter='"work" in tags && "project" in tags'
```

メモ一覧（`-any-contents`: 複数キーワード検索＋重複排除）
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -any-contents='meeting,study'
```

メモ一覧（`-filter` と `-any-contents` の併用）
```bash
go run ./cmd/cli/memos \
  -operation=list-memos \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -filter='visibility == "PUBLIC"' \
  -any-contents='meeting,study'
```

添付一覧
```bash
go run ./cmd/cli/memos \
  -operation=list-attachments \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -page-size=50 \
  -order-by="create_time desc"
```

補足（Windows のクォート）:
- `cmd.exe` では `'` が文字列クォートとして機能しないため、`-filter` の値は `"` で囲って指定してください。
- 例: `-filter="visibility == 'PUBLIC'"`、`-filter="created_ts > '2023-01-01T13:00:00Z'"`、`-filter="visibility == \"PUBLIC\""`

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

メモ更新（`displayTime` を現在日時へ設定し、`updateTime` も更新させる）
```bash
go run ./cmd/cli/memos \
  -operation=update-memo \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -content="更新後の本文" \
  -updates-time=true
```

タグ更新（`#work` を `#project` に一括置換）
```bash
go run ./cmd/cli/memos \
  -operation=update-tag \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -src-tag=work \
  -dest-tag=project
```

添付を追加（既存添付を保持: デフォルト `-replaces=false`）
```bash
go run ./cmd/cli/memos \
  -operation=patch-files \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -files="./a.png,./b.pdf"
```

添付を置換（新規添付だけにする）
```bash
go run ./cmd/cli/memos \
  -operation=patch-files \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -files="./a.png,./b.pdf" \
  -replaces=true
```

メモの既存リレーションを確認
```bash
go run ./cmd/cli/memos \
  -operation=list-memo-relations \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123
```

メモへリレーションを追加（既存保持）
```bash
go run ./cmd/cli/memos \
  -operation=add-memo-relations \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -related-memos="memo-456,memo-789"
```

メモのリレーションを置換（既存破棄）
```bash
go run ./cmd/cli/memos \
  -operation=add-memo-relations \
  -base-url=$MEMOS_BASE_URL \
  -api-token=$MEMOS_TOKEN \
  -memo=memo-123 \
  -related-memos="memo-456,memo-789" \
  -replaces=true
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
- DeleteMemo: https://usememos.com/docs/api/memoservice/DeleteMemo
- ListMemos: https://usememos.com/docs/api/memoservice/ListMemos
- Filtering (ListMemos query): https://www.usememos.com/docs/api/v1#tag/memo/GET/api/v1/memos
- Filter fields reference: https://www.usememos.com/docs/api/v1#tag/shortcut/field/filter
- UpdateMemo: https://usememos.com/docs/api/memoservice/UpdateMemo
- CreateAttachment: https://usememos.com/docs/api/attachmentservice/CreateAttachment
- ListAttachments: https://usememos.com/docs/api/attachmentservice/ListAttachments
- ListMemoAttachments: https://usememos.com/docs/api/memoservice/ListMemoAttachments
- SetMemoAttachments: https://usememos.com/docs/api/memoservice/SetMemoAttachments
- ListMemoRelations: https://usememos.com/docs/api/memoservice/ListMemoRelations
- SetMemoRelations: https://usememos.com/docs/api/memoservice/SetMemoRelations
