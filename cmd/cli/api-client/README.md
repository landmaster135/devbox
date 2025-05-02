# APIクライアント

任意のURLのAPIにリクエストを送信するコマンドラインツールです。

## 機能

- 様々なHTTPメソッド（GET, POST, PUT, DELETE, etc.）をサポート
- JSONファイルをリクエストボディとして送信可能
- レスポンスのステータスコード、ヘッダー、ボディを表示
- JSONレスポンスを自動的に整形して表示

## ビルド方法

```bash
# リポジトリのルートディレクトリで実行
./scripts/build_api_client.sh
```

## 使用方法

### GETリクエスト

```bash
./bin/api-client -url https://example.com/api
```

### POSTリクエスト（JSONファイル使用）

```bash
./bin/api-client -url https://example.com/api -method POST -json path/to/data.json
```

### PUTリクエスト（JSONファイル使用）

```bash
./bin/api-client -url https://example.com/api -method PUT -json path/to/data.json
```

### DELETEリクエスト

```bash
./bin/api-client -url https://example.com/api -method DELETE
```

## オプション

- `-url`: リクエスト先のURL（必須）
- `-method`: HTTPメソッド（デフォルト: GET）
- `-json`: リクエストボディとして送信するJSONファイルのパス（POST/PUT/PATCHの場合は必須）
- `-token`: 認証トークン（指定した場合、「Bearer <token>」形式でAuthorizationヘッダーに設定されます）

## 使用例

### JSONPlaceholderのAPIを使用した例

```bash
# ユーザー情報の取得（GET）
./bin/api-client -url https://jsonplaceholder.typicode.com/users/1

# 新しい投稿の作成（POST）
./bin/api-client -url https://jsonplaceholder.typicode.com/posts -method POST -json testdata/sample_request.json

# 認証トークンを使用したリクエスト
./bin/api-client -url https://api.example.com/secured-endpoint -token your-api-token

# 認証トークンとJSONファイルを使用したPOSTリクエスト
./bin/api-client -url https://api.example.com/secured-endpoint -method POST -json testdata/sample_request.json -token your-api-token
```

## エラー処理

- URLが指定されていない場合はエラーメッセージを表示
- POST/PUT/PATCHメソッドでJSONファイルが指定されていない場合はエラーメッセージを表示
- JSONファイルが存在しない、または無効なJSON形式の場合はエラーメッセージを表示
- APIリクエストが失敗した場合はエラーメッセージを表示
