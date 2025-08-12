# Weather Notificator CLI

指定した都市の天気予報をDiscord Webhookを通じて通知するCLIツールです。

## 機能

- OpenWeather APIから天気予報データを取得
- 指定した日数分（1-5日）の天気予報を取得
- 各日の天気情報を個別のDiscord embedとして送信
- 天気に応じた色分けとアイコン表示
- 詳細な天気情報（気温、湿度、気圧、風速、3時間毎の予報）

## 必要な準備

### 1. OpenWeather API キーの取得

1. [OpenWeatherMap](https://openweathermap.org/api)にアクセス
2. アカウントを作成してログイン
3. API Keysページでフリープランのキーを取得

### 2. Discord Webhook URLの取得

1. Discordサーバーの設定を開く
2. 「連携サービス」→「ウェブフック」を選択
3. 「新しいウェブフック」を作成
4. ウェブフックURLをコピー

## 使用方法

### 基本的な使用方法

```bash
./weather-notificator -api-key YOUR_API_KEY -city Tokyo -max-days 3 -webhook-url YOUR_WEBHOOK_URL
```

### パラメータ

| パラメータ | 必須 | 説明 | 例 |
|-----------|------|------|-----|
| `-api-key` | ✅ | OpenWeather API キー | `abc123def456` |
| `-city` | ✅ | 都市名（英語または日本語） | `Tokyo`, `Osaka`, `"New York"` |
| `-max-days` | ✅ | 取得する日数（1-5日） | `3` |
| `-webhook-url` | ✅ | Discord Webhook URL | `https://discord.com/api/webhooks/...` |
| `-help`, `-h` | ❌ | ヘルプを表示 | - |

## 使用例

### 東京の3日間天気予報

```bash
./weather-notificator \
  -api-key your_openweather_api_key \
  -city Tokyo \
  -max-days 3 \
  -webhook-url https://discord.com/api/webhooks/123456789/abcdefg
```

### 大阪の5日間天気予報

```bash
./weather-notificator \
  -api-key your_openweather_api_key \
  -city Osaka \
  -max-days 5 \
  -webhook-url https://discord.com/api/webhooks/123456789/abcdefg
```

### スペースを含む都市名の場合

```bash
./weather-notificator \
  -api-key your_openweather_api_key \
  -city "New York" \
  -max-days 2 \
  -webhook-url https://discord.com/api/webhooks/123456789/abcdefg
```

## 出力形式

各日の天気予報は個別のDiscord embedとして送信されます：

### Embedの内容

- **タイトル**: 天気アイコン + 都市名 + 日付 + 進捗（例：☀️ Tokyo の天気予報 (1/3日目)）
- **色**: 天気に応じた色分け（晴れ=黄色、雨=青色、曇り=灰色など）
- **本文**:
  - 📅 日付
  - 🌡️ 最低・最高気温
  - ☁️ 天気と詳細
  - 💧 湿度
  - 🌪️ 気圧
  - 💨 風速
  - ⏰ 3時間毎の詳細予報

### 天気アイコンと色の対応

| 天気 | アイコン | 色 |
|------|----------|-----|
| 晴れ (Clear) | ☀️ | 黄色 |
| 曇り (Clouds) | ☁️ | 灰色 |
| 雨 (Rain) | ☔️ | 青色 |
| 霧雨 (Drizzle) | 🌦️ | シアン |
| 雷雨 (Thunderstorm) | ⛈️ | 紫色 |
| 雪 (Snow) | ❄️ | 白色 |
| 霧 (Mist/Fog) | 🌫️ | 灰色 |
| 砂塵 (Dust/Sand) | 🌪️ | 茶色 |
| 突風 (Squall/Tornado) | 🌪️ | 赤色 |

## エラーハンドリング

### よくあるエラーと対処法

API キーが無効

```
エラー: 天気予報の取得に失敗しました: APIエラー: ステータスコード 401
```

**対処法**: OpenWeather APIキーが正しいか確認してください。

都市名が見つからない

```
エラー: 天気予報の取得に失敗しました: APIエラー: ステータスコード 404
```

**対処法**: 都市名のスペルを確認するか、英語名で試してください。

Webhook URLが無効

```
エラー: Discord通知の送信に失敗しました
```

**対処法**: Discord Webhook URLが正しいか確認してください。

日数制限エラー

```
エラー: 最大日数は5日以下である必要があります（OpenWeather API制限）
```

**対処法**: `-max-days`を1-5の範囲で指定してください。

## 技術仕様

### 依存関係

- OpenWeather API (5 Day / 3 Hour Forecast)
- Discord Webhook API
- Go 1.19以上

### アーキテクチャ

```
cmd/cli/weather-notificator/
├── main.go                              # CLIエントリーポイント
└── README.md                           # このファイル

internal/weather_notificator/
├── config/
│   └── config.go                       # 設定とフラグ解析
└── usecases/
    ├── services.go                     # メインビジネスロジック
    └── services_test.go               # テストコード

internal/open_weather_map/usecases/     # OpenWeather APIクライアント
internal/discord_webhook/usecases/      # Discord Webhook クライアント
```

### 制限事項

- OpenWeather APIの無料プランでは5日間の予報まで
- Discord Webhookの送信レート制限に配慮して500ms間隔で送信
- 都市名は英語での指定を推奨

## ライセンス

このプロジェクトのライセンスに従います。

## 貢献

バグ報告や機能要望は、プロジェクトのIssueトラッカーまでお願いします。
