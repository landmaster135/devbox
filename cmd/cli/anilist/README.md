# AniList CLI ツール

AniListからアニメ情報を取得するためのコマンドラインツールです。

## 概要

このツールは、AniList GraphQL APIを使用してユーザーのアニメリスト情報を取得し、JSON形式またはテーブル形式で出力します。環境変数は使用せず、すべてのパラメータはコマンドライン引数で指定します。

## 機能

- **query-anime**: 指定したユーザーのアニメリスト情報を取得
- **出力形式**: JSON形式またはテーブル形式での出力
- **フィルタリング**: ステータス別のフィルタリング
- **制限**: 取得件数の制限

## インストール

```bash
cd devbox
go build -o bin/anilist ./cmd/cli/anilist
```

## 使用例

### 基本的な使用方法

```bash
# ユーザー名を指定してアニメリストを取得
go run ./cmd/cli/anilist -operation query-anime -username your_username

# ユーザーIDを指定してアニメリストを取得
go run ./cmd/cli/anilist -operation query-anime -user-id 123456

# 短縮形を使用
go run ./cmd/cli/anilist -o query-anime -u your_username
```

### 出力形式の指定

```bash
# JSON形式で出力（デフォルト）
go run ./cmd/cli/anilist -o query-anime -u your_username -format json

# テーブル形式で出力
go run ./cmd/cli/anilist -o query-anime -u your_username -format table
```

### フィルタリング

```bash
# 完了したアニメのみを取得
go run ./cmd/cli/anilist -o query-anime -u your_username -status COMPLETED

# 現在視聴中のアニメのみを取得
go run ./cmd/cli/anilist -o query-anime -u your_username -status CURRENT
```

### 取得件数の制限

```bash
# 最新の10件のみを取得
go run ./cmd/cli/anilist -o query-anime -u your_username -limit 10

# 最新の5件のみを取得
go run ./cmd/cli/anilist -o query-anime -u your_username -limit 5
```

### ファイル出力

```bash
# 結果をファイルに保存
go run ./cmd/cli/anilist -o query-anime -u your_username -output-dir ./output

# 出力ディレクトリを指定（短縮形）
go run ./cmd/cli/anilist -o query-anime -u your_username -d ./results

# 出力形式とディレクトリを組み合わせ
go run ./cmd/cli/anilist -o query-anime -u your_username -format table -output-dir ./output
```

## オプション

| オプション | 短縮形 | 説明 | デフォルト値 |
|-----------|--------|------|-------------|
| `-operation` | `-o` | 操作タイプ (query-anime) | - |
| `-username` | `-u` | AniListユーザー名 | - |
| `-user-id` | `-i` | AniListユーザーID | - |
| `-format` | `-f` | 出力形式 (json, table) | json |
| `-limit` | `-l` | 取得件数の制限 (0は無制限) | 0 |
| `-status` | `-s` | ステータスフィルタ | - |
| `-output-dir` | `-d` | 出力ディレクトリ (指定時はファイルに保存) | - |
| `-help` | `-h` | ヘルプを表示 | - |

### ステータスの種類

- `CURRENT`: 現在視聴中
- `PLANNING`: 視聴予定
- `COMPLETED`: 完了
- `DROPPED`: 中断
- `PAUSED`: 一時停止
- `REPEATING`: 再視聴中

## 出力例

### JSON形式

```json
[
  {
    "id": 21,
    "title": "ONE PIECE",
    "score": 85,
    "status": "CURRENT",
    "progress": 1000,
    "completed_at": "",
    "notes": "長編アニメの代表作",
    "cover_image_url": "https://s4.anilist.co/file/anilistcdn/media/anime/cover/large/bx21-YCDoj1EkAxFn.jpg",
    "site_url": "https://anilist.co/anime/21",
    "studio": "Toei Animation",
    "updated_at": "2024-01-15T10:30:00Z"
  }
]
```

### テーブル形式

```
ID	タイトル	ステータス	スコア	進行状況	完了日	スタジオ
---	---	---	---	---	---	---
21	ONE PIECE	CURRENT	85	1000		Toei Animation
```

## エラーハンドリング

- ユーザー名またはユーザーIDのいずれかが必須です
- 無効な操作タイプが指定された場合はエラーメッセージを表示します
- AniList APIからエラーが返された場合は詳細なエラーメッセージを表示します
- ネットワークエラーの場合は適切なエラーメッセージを表示します

## 開発者向け情報

### アーキテクチャ

このツールは以下の層で構成されています：

- **Domain層**: データモデルとビジネスルールの定義
- **Infrastructure層**: AniList APIクライアントの実装
- **Usecases層**: ビジネスロジックの実装
- **Config層**: コマンドライン引数の解析と設定管理
- **CLI層**: メインエントリーポイント

### ディレクトリ構造

```
devbox/
├── cmd/cli/anilist/
│   ├── main.go
│   └── README.md
└── internal/anilist/
    ├── config/
    │   ├── config.go
    │   └── flag_parser.go
    ├── domain/
    │   └── models.go
    ├── infrastructure/
    │   └── anilist_client.go
    └── usecases/
        └── services.go
```

### テスト実行

```bash
cd devbox
go test ./internal/anilist/...
```

## 注意事項

- このツールは認証を必要としないパブリックなAniList APIを使用します
- レート制限に注意してください（AniListのAPIレート制限に従います）
- 大量のデータを取得する場合は、適切な制限値を設定してください

## ライセンス

このプロジェクトのライセンスに従います。

## 貢献

バグ報告や機能要求は、プロジェクトのIssueトラッカーまでお願いします。
