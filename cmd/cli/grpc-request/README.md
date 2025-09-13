# gRPC Request CLI Tool

gRPCサーバーにリクエストを送信するためのコマンドラインツールです。

## 概要

このツールは、gRPCサーバーに対してリクエストを送信し、レスポンスを取得するためのCLIツールです。gRPCリフレクションAPIを使用して、動的にサービスとメソッドの情報を取得し、JSONファイルからリクエストデータを読み込んでgRPCリクエストを実行します。

## 機能

- gRPCサーバーへの接続テスト
- 利用可能なサービス一覧の表示
- JSONファイルを使用したgRPCリクエストの送信
- TLS接続のサポート
- 認証トークンの送信（メタデータとして）
- リクエストタイムアウトの設定
- レスポンスの整形表示

## インストール

```bash
go build -o grpc-request ./cmd/cli/grpc-request
```

## 使用方法

### 基本的な使用方法

```bash
grpc-request [オプション]
```

### オプション

| オプション | 説明 | 必須 | デフォルト値 |
|-----------|------|------|-------------|
| `-server` | gRPCサーバーのアドレス（例: localhost:50051） | * | - |
| `-method` | 呼び出すメソッド（例: package.Service/Method） | △ | - |
| `-data` | リクエストデータのJSONファイルパス | △ | - |
| `-tls` | TLS接続を使用する | - | false |
| `-token` | 認証トークン（メタデータとして送信） | - | - |
| `-timeout` | リクエストタイムアウト | - | 30s |
| `-test` | 接続テストのみを実行する | - | false |
| `-list` | 利用可能なサービス一覧を表示する | - | false |

※ `-method`と`-data`は通常のリクエスト送信時に必須です。

## 使用例

### 1. 接続テスト

```bash
grpc-request -server localhost:50051 -test
```

### 2. サービス一覧表示

```bash
grpc-request -server localhost:50051 -list
```

### 3. 基本的なgRPCリクエスト送信

```bash
grpc-request -server localhost:50051 -method package.Service/Method -data request.json
```

### 4. TLS接続でリクエスト送信

```bash
grpc-request -server example.com:443 -method package.Service/Method -data request.json -tls
```

### 5. 認証トークン付きでリクエスト送信

```bash
grpc-request -server localhost:50051 -method package.Service/Method -data request.json -token your_token
```

### 6. タイムアウト設定付きでリクエスト送信

```bash
grpc-request -server localhost:50051 -method package.Service/Method -data request.json -timeout 60s
```

## リクエストデータファイル

リクエストデータはJSON形式で指定します。

### 例: request.json

```json
{
  "name": "John Doe",
  "age": 30,
  "email": "john.doe@example.com",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "zipcode": "10001"
  }
}
```

## レスポンス形式

ツールは以下の形式でレスポンスを表示します：

```
Status: 0 - OK
Duration: 150ms

Metadata:
  content-type: application/grpc
  authorization: Bearer token

Response Data:
{
  "id": "12345",
  "message": "Success",
  "timestamp": "2025-01-14T12:00:00Z"
}
```

## エラーハンドリング

- 接続エラー
- 認証エラー
- タイムアウトエラー
- gRPCステータスエラー
- JSONファイル読み込みエラー

各エラーは適切なエラーメッセージと共に表示されます。

## 制限事項

- gRPCリフレクションAPIが有効になっているサーバーでのみ動作します
- 現在の実装では、一部のgRPCリフレクション機能が簡略化されています
- ストリーミングRPCには対応していません

## 技術仕様

### アーキテクチャ

- **Domain Layer**: gRPCリクエスト/レスポンスのモデル定義
- **Repository Layer**: gRPC接続とリクエスト実行
- **Service Layer**: ビジネスロジックとデータ変換
- **CLI Layer**: コマンドライン引数の解析とメイン処理

### 依存関係

- `google.golang.org/grpc`: gRPCクライアント
- `google.golang.org/protobuf`: Protocol Buffersサポート
- `github.com/stretchr/testify`: テストフレームワーク

### ディレクトリ構造

```
internal/grpc_request/
├── config/                 # 設定管理
├── domain/models/          # ドメインモデル
├── interfaces/repositories/ # リポジトリインターフェース
└── usecases/               # サービス層
```

## 開発

### テスト実行

```bash
go test -v ./internal/grpc_request/...
```

### ビルド

```bash
go build -o grpc-request ./cmd/cli/grpc-request
```

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
