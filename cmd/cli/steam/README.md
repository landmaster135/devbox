# Steam API CLI Tool

Steam Web APIを使用してユーザーのゲーム情報を取得し、JSONファイルに出力するCLIツールです。

## 機能

- ユーザーの所有ゲーム一覧の取得
- ゲーム情報の詳細取得（名前、ID、アイコン、サムネイル、プレイ時間など）
- 全ゲームの統計情報と実績情報の一括取得
- JSONファイルへの出力
- 統計情報の表示

## 取得できるゲーム情報

### games操作で取得できる情報

各ゲームについて以下の情報を取得できます：

- **name**: ゲーム名
- **id**: Steam App ID
- **icon**: ゲームアイコンのURL
- **thumbnail**: ゲームサムネイル（ライブラリ画像）のURL
- **playtime_recent_2_weeks**: 最近2週間のプレイ時間（分）
- **playtime_disconnected**: オフラインプレイ時間（分）
- **playtime_forever**: 総プレイ時間（分）
- **recent_time_last_played**: 最後にプレイした時刻（Unix timestamp）
- **achievements_can_be_retrieved**: 実績情報が取得可能かどうか
- **stats**: 統計情報が取得可能かどうか

### game-stats操作で取得できる情報

全ゲームの詳細な統計情報と実績情報を一括取得できます：

**統計情報（stats）**:
- ゲーム内の各種統計データ（キル数、デス数、プレイ時間など）
- ゲームごとに異なる統計項目
- 数値データとして取得

**実績情報（achievements）**:
- 実績名と表示名
- 実績の説明
- 達成状況（achieved: true/false）
- 達成日時（Unix timestamp）
- 実績アイコンのURL

## 必要な環境

- Go 1.19以上
- Steam Web API キー

## Steam Web API キーの取得

1. [Steam Web API Key](https://steamcommunity.com/dev/apikey)にアクセス
2. Steamアカウントでログイン
3. ドメイン名を入力（例: localhost）
4. API キーを取得

## インストール・ビルド

```bash
# リポジトリのルートディレクトリから
cd devbox
go build -o steam-cli ./cmd/cli/steam/
```

## 使用方法

### 基本的な使用方法

```bash
go run ./cmd/cli/steam --operation games --steam-api-key YOUR_API_KEY --steam-id STEAM_ID
```

### ショートハンドオプション

```bash
go run ./cmd/cli/steam -o games -k YOUR_API_KEY -s STEAM_ID
```

### パラメータ

| パラメータ | ショートハンド | 必須 | 説明 |
|-----------|---------------|------|------|
| `--operation` | `-o` | * | 実行する操作（`games`, `game-stats`） |
| `--steam-api-key` | `-k` | * | Steam Web API キー |
| `--steam-id` | `-s` | * | 対象ユーザーの17桁のSteam ID |
| `--game-id` | `-g` | - | ゲームID（game-stats操作で使用） |
| `--output-dir` | `-d` | - | JSONファイルの出力ディレクトリ（デフォルト: カレントディレクトリ） |

### Steam IDの確認方法

Steam IDは以下の方法で確認できます：

1. **Steam プロフィールURL**から：
  - `https://steamcommunity.com/profiles/76561198000000000/` の数字部分
   
2. **カスタムURL**の場合：
  - [SteamID Finder](https://steamidfinder.com/)などのツールを使用

3. **Steam クライアント**から：
  - プロフィール → プロフィールを編集 → カスタムURL

## 出力

### JSONファイル

実行すると、操作に応じて以下の形式でJSONファイルが生成されます：

**games操作の出力**

ファイル名: `steam_games_{STEAM_ID}_{TIMESTAMP}.json`

```json
{
  "steam_id": "76561198000000000",
  "generated_at": "2025-01-19T12:34:56Z",
  "total_games": 150,
  "games": [
    {
      "name": "Counter-Strike 2",
      "id": 730,
      "icon": "https://media.steampowered.com/steamcommunity/public/images/apps/730/69f7ebe2735c366c65c0b33dae00e12dc40edbe4.jpg",
      "thumbnail": "https://shared.fastly.steamstatic.com/store_item_assets/steam/apps/730/library_600x900.jpg",
      "playtime_recent_2_weeks": 120,
      "playtime_disconnected": 0,
      "playtime_forever": 5400,
      "recent_time_last_played": 1705123456,
      "achievements_can_be_retrieved": true,
      "stats": true
    }
  ]
}
```

**game-stats操作の出力**

ファイル名: `steam_games_stats_{STEAM_ID}_{TIMESTAMP}.json`

```json
[
  {
    "game_name": "Counter-Strike 2",
    "game_id": 730,
    "stats": [
      {
        "name": "total_kills",
        "value": 12345
      },
      {
        "name": "total_deaths", 
        "value": 9876
      }
    ],
    "achievements": [
      {
        "name": "First Kill",
        "display_name": "First Blood",
        "description": "Get your first kill",
        "achieved": true,
        "unlock_time": 1705123456
      }
    ]
  }
]
```

### コンソール出力

実行時には以下の情報がコンソールに表示されます：

- 取得したゲーム数
- 統計情報（総プレイ時間、実績対応ゲーム数など）
- 最もプレイしたゲームのトップ5

## 使用例

```bash
# ゲーム一覧の取得
go run ./cmd/cli/steam --operation games --steam-api-key ABCD1234567890 --steam-id 76561198000000000

# ゲーム統計・実績情報の取得
go run ./cmd/cli/steam --operation game-stats --steam-api-key ABCD1234567890 --steam-id 76561198000000000

# ショートハンドを使用（ゲーム一覧）
go run ./cmd/cli/steam -o games -k ABCD1234567890 -s 76561198000000000

# ショートハンドを使用（統計・実績情報）
go run ./cmd/cli/steam -o game-stats -k ABCD1234567890 -s 76561198000000000

# 出力ディレクトリを指定してゲーム一覧を取得
go run ./cmd/cli/steam --operation games --steam-api-key ABCD1234567890 --steam-id 76561198000000000 --output-dir /path/to/output

# ショートハンドで出力ディレクトリを指定
go run ./cmd/cli/steam -o games -k ABCD1234567890 -s 76561198000000000 -d /path/to/output

# 相対パスで出力ディレクトリを指定
go run ./cmd/cli/steam -o games -k ABCD1234567890 -s 76561198000000000 -d ./data/steam
```

## エラーハンドリング

以下の場合にエラーが発生します：

- 必須パラメータが不足している場合
- Steam API キーが無効な場合
- Steam IDが存在しない、または形式が正しくない場合
- プロフィールがプライベートに設定されている場合
- Steam APIのレート制限に達した場合

## 制限事項

- Steam APIのレート制限により、大量のゲームを持つユーザーの場合、処理に時間がかかる場合があります
- プライベートプロフィールの場合、一部の情報が取得できない場合があります
- 実績・統計情報の取得可能性は、ゲームの設定やプライバシー設定に依存します

## トラブルシューティング

### よくある問題

1. **"invalid Steam ID format" エラー**
  - Steam IDが17桁の数字であることを確認してください

2. **"failed to get owned games" エラー**
  - Steam API キーが正しいことを確認してください
  - 対象ユーザーのプロフィールが公開されていることを確認してください

3. **"HTTP error: 403" エラー**
  - Steam API キーが無効、または期限切れの可能性があります
  - 新しいAPI キーを取得してください

### デバッグ

詳細なエラー情報が必要な場合は、以下のようにログレベルを上げて実行できます：

```bash
# 詳細なログ出力
STEAM_DEBUG=1 go run ./cmd/cli/steam -o games -k YOUR_KEY -s STEAM_ID
```

## ライセンス

このプロジェクトはMITライセンスの下で公開されています。

## 貢献

バグ報告や機能要求は、GitHubのIssueでお知らせください。

## 関連リンク

- [Steam Web API Documentation](https://developer.valvesoftware.com/wiki/Steam_Web_API)
- [Steam Web API Key](https://steamcommunity.com/dev/apikey)
- [SteamID Finder](https://steamidfinder.com/)
