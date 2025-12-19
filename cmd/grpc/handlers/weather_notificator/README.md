# Weather Notificator gRPC Handler

このディレクトリには、天気通知サービスのgRPCハンドラーが含まれています。

## 概要

`WeatherNotificatorHandler`は、gRPCプロトコルを通じて天気通知機能を提供するハンドラーです。既存の`WeatherNotificatorService`を利用して、gRPCリクエストを処理し、指定された都市の天気予報をDiscordに通知します。

## アーキテクチャ

```
gRPCクライアント
    ↓
WeatherNotificatorHandler (gRPCハンドラー)
    ↓
WeatherNotificatorService (ビジネスロジック)
    ↓
┌─────────────────────┬─────────────────────┐
│  WeatherService     │  DiscordService     │
│  (OpenWeatherMap)   │  (Discord Webhook)  │
└─────────────────────┴─────────────────────┘
```

## 機能

### SendWeatherNotification

指定された都市の天気予報を取得し、Discordに通知を送信します。

**gRPCメソッド:**
```protobuf
rpc SendWeatherNotification(WeatherNotificationRequest) returns (WeatherNotificationResponse);
```

**リクエストパラメータ:**
- `api_key` (string): OpenWeatherMap APIキー
  - 必須項目
  - OpenWeatherMapから取得したAPIキーを指定
- `city` (string): 都市名
  - 必須項目
  - 例: "Tokyo", "Osaka", "New York", "London"
  - 英語での都市名を推奨
- `max_days` (int32): 予報日数
  - 必須項目
  - 範囲: 1-5日
  - 指定した日数分の天気予報を取得
- `webhook_url` (string): Discord Webhook URL
  - 必須項目
  - Discord チャンネルのWebhook URL
  - 形式: `https://discord.com/api/webhooks/...`

**レスポンス:**
- `success` (bool): 処理成功フラグ
  - `true`: 正常に処理完了
  - `false`: エラーが発生
- `message` (string): 成功時のメッセージ
  - 成功時に返される詳細メッセージ
  - 例: "✅ Tokyoの3日間天気予報をDiscordに送信しました"
- `error` (string): エラー時のエラーメッセージ
  - エラー発生時の詳細なエラー内容

## バリデーション

リクエストは以下の条件でバリデーションされます：

1. **APIキー**: 空文字列でないこと
2. **都市名**: 空文字列でないこと
3. **最大日数**: 1-5の範囲内であること
4. **Webhook URL**: 空文字列でないこと

バリデーションエラーが発生した場合、gRPCステータスコード `InvalidArgument` が返されます。

## エラーハンドリング

### gRPCステータスコード

- `codes.InvalidArgument`: リクエストパラメータのバリデーションエラー
- `codes.Internal`: 内部処理エラー（天気データ取得失敗、Discord通知失敗など）

### エラーレスポンス例

```json
{
  "success": false,
  "message": "",
  "error": "APIキーが指定されていません"
}
```

## 使用例

### Go クライアント

```go
package main

import (
  "context"
  "log"
  
  pb "github.com/landmaster135/devbox/cmd/grpc/proto/weather_notificator"
  "google.golang.org/grpc"
)

func main() {
  // gRPCサーバーに接続
  conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
  if err != nil {
    log.Fatalf("接続失敗: %v", err)
  }
  defer conn.Close()
  
  // クライアントを作成
  client := pb.NewWeatherNotificatorServiceClient(conn)
  
  // リクエストを作成
  req := &pb.WeatherNotificationRequest{
    ApiKey:     "your-openweathermap-api-key",
    City:       "Tokyo",
    MaxDays:    3,
    WebhookUrl: "https://discord.com/api/webhooks/...",
  }
  
  // gRPCメソッドを呼び出し
  resp, err := client.SendWeatherNotification(context.Background(), req)
  if err != nil {
    log.Fatalf("リクエスト失敗: %v", err)
  }
  
  if resp.Success {
    log.Printf("成功: %s", resp.Message)
  } else {
    log.Printf("エラー: %s", resp.Error)
  }
}
```

### grpcurl を使用したテスト

```bash
# サーバーのサービス一覧を確認
grpcurl -plaintext localhost:50051 list

# メソッドの詳細を確認
grpcurl -plaintext localhost:50051 describe weather_notificator.WeatherNotificatorService

# 天気通知を送信
grpcurl -plaintext -d '{
  "api_key": "your-api-key",
  "city": "Tokyo",
  "max_days": 3,
  "webhook_url": "https://discord.com/api/webhooks/..."
}' localhost:50051 weather_notificator.WeatherNotificatorService/SendWeatherNotification
```

## 実装詳細

### ハンドラー構造体

```go
type WeatherNotificatorHandler struct {
  pb.UnimplementedWeatherNotificatorServiceServer
  service *usecases.WeatherNotificatorService
}
```

### コンストラクタ

- `NewWeatherNotificatorHandler()`: デフォルトの依存関係でハンドラーを作成
- `NewWeatherNotificatorHandlerWithService(service)`: 依存関係を注入してハンドラーを作成（テスト用）

### 主要メソッド

- `SendWeatherNotification(ctx, req)`: gRPCメソッドの実装
- `validateRequest(req)`: リクエストバリデーション
- `GetService()`: サービスインスタンスの取得（テスト用）

## ログ出力

ハンドラーは以下の情報をログに出力します：

- **リクエスト受信**: リクエストパラメータの概要
- **バリデーションエラー**: バリデーション失敗時の詳細
- **処理エラー**: 天気データ取得やDiscord通知の失敗
- **処理成功**: 正常完了時のメッセージ

### ログ例

```
2025/09/14 07:30:00 天気通知成功: ✅ Tokyoの3日間天気予報をDiscordに送信しました
2025/09/14 07:30:01 リクエストバリデーションエラー: APIキーが指定されていません
2025/09/14 07:30:02 天気通知処理エラー: 天気予報の取得に失敗しました: invalid API key
```

## 依存関係

### 内部依存関係

- `github.com/landmaster135/devbox/internal/weather_notificator/usecases`: ビジネスロジック
- `github.com/landmaster135/devbox/cmd/grpc/proto/weather_notificator`: Protocol Buffers生成コード

### 外部依存関係

- `google.golang.org/grpc`: gRPCフレームワーク
- `google.golang.org/grpc/codes`: gRPCステータスコード
- `google.golang.org/grpc/status`: gRPCステータス管理

## テスト

### 単体テスト

```go
func TestWeatherNotificatorHandler_SendWeatherNotification(t *testing.T) {
  // モックサービスを作成
  mockService := &MockWeatherNotificatorService{}
  handler := NewWeatherNotificatorHandlerWithService(mockService)
  
  // テストケースを実行
  req := &pb.WeatherNotificationRequest{
    ApiKey:     "test-api-key",
    City:       "Tokyo",
    MaxDays:    3,
    WebhookUrl: "https://discord.com/api/webhooks/test",
  }
  
  resp, err := handler.SendWeatherNotification(context.Background(), req)
  
  assert.NoError(t, err)
  assert.True(t, resp.Success)
}
```

### 統合テスト

```bash
# gRPCサーバーを起動
go run cmd/grpc/main.go

# 別ターミナルでテスト実行
grpcurl -plaintext -d '{"api_key":"test","city":"Tokyo","max_days":1,"webhook_url":"test"}' \
  localhost:50051 weather_notificator.WeatherNotificatorService/SendWeatherNotification
```

## パフォーマンス

### 推奨設定

- **同時接続数**: 100接続まで
- **リクエストタイムアウト**: 30秒
- **レスポンスサイズ**: 最大1MB

### 制限事項

- OpenWeatherMap APIの制限に依存（無料プランは1000リクエスト/日）
- Discord Webhook APIの制限に依存（レート制限あり）

## セキュリティ

### 推奨事項

1. **APIキーの管理**: 環境変数やシークレット管理システムを使用
2. **TLS暗号化**: 本番環境ではTLSを有効化
3. **認証・認可**: 必要に応じてgRPCインターセプターで実装
4. **入力検証**: 悪意のある入力に対する適切な検証

### セキュリティ設定例

```go
// TLS設定
creds, err := credentials.LoadTLSConfig(&tls.Config{
    ServerName: "your-server.com",
})
conn, err := grpc.Dial("your-server.com:443", grpc.WithTransportCredentials(creds))
```

## トラブルシューティング

### よくある問題

1. **接続エラー**: サーバーが起動しているか確認
2. **APIキーエラー**: OpenWeatherMap APIキーが有効か確認
3. **Discord通知失敗**: Webhook URLが正しいか確認
4. **都市名エラー**: 英語の都市名を使用しているか確認

### デバッグ方法

```bash
# サーバーログを確認
tail -f /var/log/grpc-server.log

# gRPCリフレクションでサービス確認
grpcurl -plaintext localhost:50051 list
```

## 関連ドキュメント

- [Protocol Buffers定義](../proto/weather_notificator/weather_notificator.proto)
- [gRPCサーバー設定](../server/README.md)
- [天気通知サービス](../../internal/weather_notificator/README.md)
- [OpenWeatherMap API](https://openweathermap.org/api)
- [Discord Webhook API](https://discord.com/developers/docs/resources/webhook)
