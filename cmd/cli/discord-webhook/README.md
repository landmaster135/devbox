# Discord Webhook通知CLIツール

Discord WebhookでメッセージやEmbed付き通知を送信するためのCLIツールです。

## 機能

- Discord Webhookへのメッセージ送信
- Embedなしの簡単な通知
- VSCode風のEmbed付き通知（フッターにVSCodeアイコン表示）
- PostgreSQLダンプ結果のEmbed付き通知（VSCodeレイアウトを流用）
- OpenWeatherMap風のEmbed付き通知（専用カラーとアイコンを使用）
- Google Cloud各サービス（Compute Engine / Secret Manager / Cloud Runなど）の成功・失敗通知Embed
- カスタマイズ可能な色設定
- リンク付きタイトル対応

## インストール

```bash
go build -o bin/discord-webhook ./cmd/cli/discord-webhook
```

## 使用方法

### 基本的な通知（Embedなし）

```bash
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -bot-name "テスト用のボット" \
  -content-text "Hello, Discord!" \
  -embed-type none
```

### VSCode用Embed付き通知

```bash
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "デプロイ完了" \
  -embed-type vscode \
  -embed-text "アプリケーションが正常にデプロイされました"
```

### PostgreSQLダンプ通知

```bash
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "PostgreSQLのダンプが完了しました" \
  -embed-type postgres \
  -embed-text "バックアップ完了"
```

### OpenWeatherMap用Embed付き通知

```bash
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "東京の天気予報です" \
  -embed-type open-weather-map \
  -embed-text "今日の天気予報"
```

### Google Cloudサービスの結果通知

Google Cloudの各サービス名と成功/失敗に応じたEmbedタイプを指定します。Embedテキストや色を省略すると、成功時は緑、失敗時は赤で送信されます。

```bash
# Google Compute Engineの操作成功
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "GCEインスタンスの作成に成功しました" \
  -embed-type google-compute-engine-success

# Cloud Runのデプロイ失敗（赤色で通知）
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "Cloud Runデプロイでエラーが発生しました" \
  -embed-type google-cloud-run-failed \
  -embed-text "デプロイジョブが失敗しました"
```

### フルオプション

```bash
go run ./cmd/cli/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "通知" \
  -embed-type vscode \
  -embed-text "タイトル" \
  -embed-color "green" \
  -embed-url-linked-text "https://example.com"
```

## オプション

### 必須オプション

| オプション | 短縮形 | 説明 |
|-----------|--------|------|
| `-webhook-url` | `-wu` | Discord WebhookのURL |
| `-content-text` | `-ct` | メッセージの本文 |
| `-embed-type` | `-et` | Embedのタイプ (none, vscode, postgres, open-weather-map, google-*-success, google-*-failed) |

### 任意オプション

| オプション | 短縮形 | 説明 |
|-----------|--------|------|
| `-bot-name` | `-bn` | ボットの名前 |
| `-embed-text` | `-et-text` | Embedのタイトル |
| `-embed-color` | `-ec` | Embedの色 |
| `-embed-url-linked-text` | `-eult` | EmbedタイトルのリンクURL |
| `-help`, `-h` | ヘルプを表示 |

### Embed色の選択肢

- `green` - 緑色
- `red` - 赤色
- `blue` - 青色（VSCodeモードのデフォルト）
- `yellow` - 黄色
- `orange` - オレンジ色（OpenWeatherMapモードのデフォルト）
- `purple` - 紫色
- `pink` - ピンク色
- `sky_blue` - 空色
- `gray_blue` - グレー青
- `white` - 白色
- `black` - 黒色

## Embed-typeについて

### none
- Embedを使用せず、content-textのみを送信
- シンプルなテキストメッセージ

### vscode
- VSCode風のEmbedを使用
- フッターにVSCodeアイコンを表示
- タイムスタンプ付き
- カスタマイズ可能な色とリンク

### postgres
- VSCode風レイアウトでPostgreSQL関連イベントを通知
- フッターにPostgreSQL公式アイコンを表示
- デフォルトで紫色・「PostgreSQLダンプ」タイトルを設定
- `-embed-text` / `-embed-color` / `-embed-url-linked-text` で上書き可能

### open-weather-map
- OpenWeatherMap風のEmbedを使用
- デフォルトでオレンジ色・天気予報タイトル・専用ロゴを適用
- `-embed-text` / `-embed-color` / `-embed-url-linked-text` で上書き可能

### gcloud系
`google-*-success` / `google-*-failed`の形式で指定
- Google Cloudの各サービス（`compute-engine`, `secret-manager`, `cloud-storage`, `cloud-scheduler`, `cloud-iam`, `cloud-run`, `cloud-run-function`）に対応
- `-success`指定時はデフォルトで緑色、`-failed`指定時は赤色で通知
- Embedのフッターにサービス名とGoogle Cloudアイコンを表示
- `-embed-text` / `-embed-color` / `-embed-url-linked-text` で上書き可能

## 使用例

### 1. シンプルな通知

```bash
go run ./cmd/cli/discord-webhook -wu "YOUR_WEBHOOK_URL" -ct "Hello World!" -et none
```

### 2. 成功通知

```bash
go run ./cmd/cli/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "ビルドが完了しました" \
  -et vscode \
  -et-text "ビルド成功" \
  -ec "green"
```

### 3. エラー通知

```bash
go run ./cmd/cli/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "エラーが発生しました" \
  -et vscode \
  -et-text "ビルド失敗" \
  -ec "red" \
  -eult "https://github.com/your-repo/actions"
```

### 4. 天気予報通知（カスタム設定）

```bash
go run ./cmd/cli/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "大阪の天気予報をお届けします" \
  -et open-weather-map \
  -et-text "大阪の週末予報" \
  -ec "sky_blue" \
  -eult "https://example.com/weather"
```

## エラーハンドリング

- 必須パラメータが不足している場合、エラーメッセージとヘルプが表示されます
- 無効なembed-typeを指定した場合、有効な値が表示されます
- 無効なWebhook URLの場合、適切なエラーメッセージが表示されます
- ネットワークエラーや Discord API エラーも適切に処理されます

## 技術仕様

- Go言語で実装
- Clean Architecture採用
- テスト可能な設計
- 既存のinfrastructureリソースを活用

## ファイル構成

```
devbox/
├── cmd/cli/discord-webhook/
│   ├── main.go              # CLIエントリーポイント
│   └── README.md            # このファイル
├── internal/discord_webhook/
│   ├── config/              # 設定とフラグ解析
│   │   ├── config.go
│   │   ├── flag_parser.go
│   │   └── interfaces.go
│   ├── usecases/            # ビジネスロジック
│   │   └── services.go
│   ├── infrastructure/      # 外部リソース
│   │   └── discord/         # Discord API クライアント
│   └── asset/
│       └── vscode.webp      # VSCodeアイコン
```

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのライセンスファイルを参照してください。
