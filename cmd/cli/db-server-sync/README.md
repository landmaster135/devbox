# db-server-sync

AniListデータをデータベースサーバー用のリクエストボディ形式に変換するCLIツールです。

## 機能

- AniListからエクスポートしたJSONデータを指定されたリクエストボディ形式に変換
- 追加データファイルとの結合機能（anilist_idをキーとして）
- ISO8601形式のタイムスタンプをUnixタイムスタンプに変換

## 使用方法

### 基本的な使用方法

```bash
# AniListデータのみを変換
go run ./cmd/cli/db-server-sync -operation=append-anime -input-file-path=anilist.json -output-file-path=output.json

# 追加データと結合して変換
go run ./cmd/cli/db-server-sync -operation=append-anime -input-file-path=anilist.json -additional-input-file-path=additional.json -output-file-path=output.json
```

### パラメータ

| パラメータ | 必須 | 説明 |
|-----------|------|------|
| `-operation` | * | 実行する操作（現在は`append-anime`のみ対応） |
| `-input-file-path` | * | 入力ファイルのパス（AniListデータのJSONファイル） |
| `-output-file-path` | * | 出力ファイルのパス |
| `-additional-input-file-path` | | 追加入力ファイルのパス（tier情報等） |
| `-help` | | ヘルプを表示 |

## 入力データ形式

### メインファイル（AniListデータ）

```json
[
  {
    "id": 5114,
    "title": "鋼の錬金術師 FULLMETAL ALCHEMIST",
    "score": 94,
    "status": "COMPLETED",
    "progress": 64,
    "completed_at": "2022-01-13T00:00:00+09:00",
    "notes": "こういう終わり方もクールだと思います。",
    "cover_image_url": "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx5114-nSWCgQlmOMtj.jpg",
    "site_url": "https://anilist.co/anime/5114",
    "studio": "bones",
    "updated_at": "2024-01-22T16:31:54+09:00"
  }
]
```

### 追加ファイル（tier情報等）

```json
[
  {
    "anilist_id": 5114,
    "con_id": "AN0001",
    "cover_image_url_modified": "https://storage.googleapis.com/test/test_contents/AN0001_01.webp",
    "fleeting_tier": 0,
    "funny_tier": 2,
    "heartwarming_tier": 2,
    "motivating_tier": 1,
    "nihilistic_tier": 1,
    "tearjerking_tier": 2
  }
]
```

## 出力データ形式

```json
{
  "data": {
    "animes": [
      {
        "anilist_id": 5114,
        "completed_at": "2022-01-13T00:00:00+09:00",
        "con_id": "AN0001",
        "cover_image_url": "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx5114-nSWCgQlmOMtj.jpg",
        "cover_image_url_modified": "https://storage.googleapis.com/test/test_contents/AN0001_01.webp",
        "fleeting_tier": 0,
        "funny_tier": 2,
        "heartwarming_tier": 2,
        "motivating_tier": 1,
        "nihilistic_tier": 1,
        "notes": "こういう終わり方もクールだと思います。",
        "progress": 64,
        "score": 94,
        "site_url": "https://anilist.co/anime/5114",
        "status": "COMPLETED",
        "studio": "bones",
        "tearjerking_tier": 2,
        "title": "鋼の錬金術師 FULLMETAL ALCHEMIST",
        "updated_at": 1642777914
      }
    ]
  },
  "description": "Anime data from AniList",
  "name": "My Anime List"
}
```

## データ変換の詳細

### フィールドマッピング

| 入力フィールド | 出力フィールド | 変換内容 |
|---------------|---------------|----------|
| `id` | `anilist_id` | そのまま |
| `title` | `title` | そのまま |
| `score` | `score` | そのまま |
| `status` | `status` | そのまま |
| `progress` | `progress` | そのまま |
| `completed_at` | `completed_at` | そのまま（nullの場合は空文字） |
| `notes` | `notes` | そのまま |
| `cover_image_url` | `cover_image_url` | そのまま |
| `site_url` | `site_url` | そのまま |
| `studio` | `studio` | そのまま |
| `updated_at` | `updated_at` | ISO8601文字列 → Unixタイムスタンプ |

### 追加データのマージ

追加ファイルが指定された場合、`anilist_id`をキーとして以下のフィールドが追加されます：

- `con_id`
- `cover_image_url_modified`
- `fleeting_tier`（nullの場合は0）
- `funny_tier`（nullの場合は0）
- `heartwarming_tier`（nullの場合は0）
- `motivating_tier`（nullの場合は0）
- `nihilistic_tier`（nullの場合は0）
- `tearjerking_tier`（nullの場合は0）

## ビルド方法

```bash
# devboxディレクトリで実行
go build -o bin/db-server-sync ./cmd/cli/db-server-sync
```

## テスト実行

```bash
# 全体のテスト実行
go test ./internal/db_server_sync/...

# カバレッジ付きテスト実行
go test -coverprofile=coverage.out ./internal/db_server_sync/...
go tool cover -html=coverage.out -o coverage.html
```

## エラーハンドリング

- 必須パラメータが不足している場合はエラーメッセージとヘルプを表示
- ファイルが存在しない場合は適切なエラーメッセージを表示
- JSONの形式が不正な場合は詳細なエラーメッセージを表示
- タイムスタンプの変換に失敗した場合はエラーメッセージを表示

## 使用例

### 例1: 基本的な変換

```bash
go run ./cmd/cli/db-server-sync -operation=append-anime -input-file-path=/path/to/anilist_data.json -output-file-path=/path/to/output.json
```

### 例2: 追加データとの結合

```bash
go run ./cmd/cli/db-server-sync -operation=append-anime -input-file-path=/path/to/anilist_data.json -additional-input-file-path=/path/to/tier_data.json -output-file-path=/path/to/output.json
```

### 例3: ヘルプの表示

```bash
go run ./cmd/cli/db-server-sync -help
