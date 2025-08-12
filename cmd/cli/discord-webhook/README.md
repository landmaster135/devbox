# Discord Webhook通知CLIツール

Discord WebhookでメッセージやEmbed付き通知を送信するためのCLIツールです。

## 機能

- Discord Webhookへのメッセージ送信
- Embedなしの簡単な通知
- VSCode風のEmbed付き通知（フッターにVSCodeアイコン表示）
- カスタマイズ可能な色設定
- リンク付きタイトル対応

## インストール

```bash
cd /home/nov/devbox
go build -o bin/discord-webhook ./cmd/cli/discord-webhook
```

## 使用方法

### 基本的な通知（Embedなし）

```bash
./bin/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "Hello, Discord!" \
  -embed-type none
```

### VSCode風Embed付き通知

```bash
./bin/discord-webhook \
  -webhook-url "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN" \
  -content-text "デプロイ完了" \
  -embed-type vscode \
  -embed-text "アプリケーションが正常にデプロイされました"
```

### フルオプション

```bash
./bin/discord-webhook \
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
| `-embed-type` | `-et` | Embedのタイプ (none, vscode) |

### 任意オプション

| オプション | 短縮形 | 説明 |
|-----------|--------|------|
| `-embed-text` | `-et-text` | Embedのタイトル |
| `-embed-color` | `-ec` | Embedの色 |
| `-embed-url-linked-text` | `-eult` | EmbedタイトルのリンクURL |
| `-help` | `-h` | ヘルプを表示 |

### Embed色の選択肢

- `green` - 緑色
- `red` - 赤色
- `blue` - 青色（デフォルト）
- `yellow` - 黄色
- `orange` - オレンジ色
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

## 使用例

### 1. シンプルな通知

```bash
./bin/discord-webhook -wu "YOUR_WEBHOOK_URL" -ct "Hello World!" -et none
```

### 2. 成功通知

```bash
./bin/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "ビルドが完了しました" \
  -et vscode \
  -et-text "ビルド成功" \
  -ec "green"
```

### 3. エラー通知

```bash
./bin/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "エラーが発生しました" \
  -et vscode \
  -et-text "ビルド失敗" \
  -ec "red" \
  -eult "https://github.com/your-repo/actions"
```

### 4. デプロイ通知

```bash
./bin/discord-webhook \
  -wu "YOUR_WEBHOOK_URL" \
  -ct "新しいバージョンがデプロイされました" \
  -et vscode \
  -et-text "デプロイ完了 v1.2.3" \
  -ec "blue" \
  -eult "https://your-app.com"
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
