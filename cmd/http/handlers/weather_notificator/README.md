# Weather Notificator HTTP Handler

天気予報をDiscordに通知するHTTPハンドラです。

## 概要

このハンドラは、OpenWeather APIから天気予報を取得し、Discord Webhookを使用してDiscordチャンネルに通知を送信するHTTPエンドポイントを提供します。

## 機能

- 指定した都市の天気予報を取得
- 1〜5日間の天気予報をサポート
- Discord Webhookを使用した通知送信
- リクエストのバリデーション
- 適切なエラーハンドリング
- JSON形式のレスポンス

## エンドポイント

### POST /weather-notification

指定した都市の天気予報をDiscordに送信します。

### リクエスト

**Content-Type:** `application/json`

**パラメータ:**

| パラメータ | 型 | 必須 | 説明 |
|-----------|---|------|------|
| `api_key` | string | * | OpenWeather API キー |
| `city` | string | * | 都市名（例: "Tokyo", "New York"） |
| `max_days` | int | * | 最大日数（1-5日） |
| `webhook_url` | string | * | Discord Webhook URL |

### リクエスト例

```json
{
  "api_key": "your_openweather_api_key",
  "city": "Tokyo",
  "max_days": 3,
  "webhook_url": "https://discord.com/api/webhooks/123456789/abcdefghijklmnop"
}
```

### レスポンス

**成功時 (200 OK):**

```json
{
  "success": true,
  "message": "✅ Tokyoの3日間天気予報をDiscordに送信しました"
}
```

**エラー時 (400 Bad Request / 500 Internal Server Error):**

```json
{
  "success": false,
  "error": "エラーメッセージ"
}
```

## エラーコード

| HTTPステータス | 説明 |
|---------------|------|
| 200 | 成功 |
| 400 | リクエストエラー（バリデーション失敗、不正なJSON等） |
| 405 | メソッドが許可されていない（POST以外） |
| 500 | サーバー内部エラー（API呼び出し失敗等） |

## バリデーション

以下の条件でリクエストをバリデーションします：

- `api_key`: 空文字列でないこと
- `city`: 空文字列でないこと
- `max_days`: 1以上5以下の整数であること
- `webhook_url`: 空文字列でないこと
- Content-Type: `application/json`であること
- HTTPメソッド: `POST`であること

## 使用例

### cURLを使用した例

```bash
curl -X POST http://localhost:8080/weather-notification \
  -H "Content-Type: application/json" \
  -d '{
    "api_key": "your_openweather_api_key",
    "city": "Tokyo",
    "max_days": 3,
    "webhook_url": "https://discord.com/api/webhooks/123456789/abcdefghijklmnop"
  }'
```

### Goクライアントの例

```go
package main

import (
  "bytes"
  "encoding/json"
  "fmt"
  "net/http"
)

type WeatherRequest struct {
  APIKey     string `json:"api_key"`
  City       string `json:"city"`
  MaxDays    int    `json:"max_days"`
  WebhookURL string `json:"webhook_url"`
}

func main() {
  req := WeatherRequest{
    APIKey:     "your_openweather_api_key",
    City:       "Tokyo",
    MaxDays:    3,
    WebhookURL: "https://discord.com/api/webhooks/123456789/abcdefghijklmnop",
  }

  jsonData, _ := json.Marshal(req)
  
  resp, err := http.Post(
    "http://localhost:8080/weather-notification",
    "application/json",
    bytes.NewBuffer(jsonData),
  )
  if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
  }
  defer resp.Body.Close()

  fmt.Printf("Status: %s\n", resp.Status)
}
```

## 依存関係

このハンドラは以下のサービスに依存しています：

- `github.com/landmaster135/devbox/internal/weather_notificator/usecases.WeatherNotificatorService`
- OpenWeather API
- Discord Webhook API

## ログ出力

ハンドラは以下の情報をログに出力します：

- リクエストの受信
- バリデーションエラー
- サービス実行エラー
- 成功時のメッセージ
- レスポンス送信エラー

## セキュリティ考慮事項

- API キーは環境変数や設定ファイルで管理することを推奨
- Discord Webhook URLは機密情報として扱う
- 適切なレート制限の実装を検討
- HTTPS通信の使用を推奨

## テスト

ハンドラのテストを実行するには：

```bash
go test ./cmd/http/handlers/weather_notificator/...
```

## トラブルシューティング

## よくある問題

1. **400 Bad Request: API キーが指定されていません**
   - `api_key`フィールドが空または未設定

2. **400 Bad Request: 最大日数は1以上である必要があります**
   - `max_days`が0以下の値

3. **500 Internal Server Error: 天気通知の実行に失敗しました**
   - OpenWeather APIキーが無効
   - Discord Webhook URLが無効
   - ネットワーク接続の問題

4. **405 Method Not Allowed**
   - POST以外のHTTPメソッドを使用

## デバッグ方法

1. ログを確認してエラーの詳細を把握
2. OpenWeather APIキーの有効性を確認
3. Discord Webhook URLの正確性を確認
4. ネットワーク接続を確認

## 関連ドキュメント

- [OpenWeather API Documentation](https://openweathermap.org/api)
- [Discord Webhook Documentation](https://discord.com/developers/docs/resources/webhook)
- [Weather Notificator CLI README](../../cli/weather-notificator/README.md)
