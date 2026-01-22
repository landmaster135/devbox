# HTTPリクエスト

任意のURLのAPIにリクエストを送信するコマンドラインツールです。

## 機能

- 様々なHTTPメソッド（GET, POST, PUT, DELETE, etc.）をサポート
- JSONファイルをリクエストボディとして送信可能
- レスポンスのステータスコード、ヘッダー、ボディを表示
- JSONレスポンスを自動的に整形して表示
- 文字エンコーディング対応（Shift_JIS、UTF-8、EUC-JP、自動検出）
- HTMLレスポンスの文字化け解消
- 主要ブラウザに近いUser-Agent/HTTPヘッダーを自動付与し、CloudflareなどのBot検出を回避しやすい
- gzip/deflate/brotli圧縮レスポンスを自動展開
- Cloudflareチャレンジを検出すると自動リトライと警告表示を行い、問題の切り分けを支援

## ビルド方法

```bash
# リポジトリのルートディレクトリで実行
./scripts/build_http_request.sh
```

## オプション

- `-url`: リクエスト先のURL（必須）
- `-method`: HTTPメソッド（デフォルト: GET）
- `-json`: リクエストボディとして送信するJSONファイルのパス（POST/PUT/PATCHの場合は必須）
- `-token`: 認証トークン（指定した場合、「Bearer <token>」形式でAuthorizationヘッダーに設定されます）
- `-encoding`: レスポンスの文字エンコーディング指定（shift_jis, utf-8, euc-jp, auto）

## 使用例

### JSONPlaceholderのAPIを使用した例

```bash
# ユーザー情報の取得（GET）
go run ./cmd/cli//http-request -url https://jsonplaceholder.typicode.com/users/1

# 新しい投稿の作成（POST）
go run ./cmd/cli//http-request -url https://jsonplaceholder.typicode.com/posts -method POST -json testdata/sample_request.json

# 認証トークンを使用したリクエスト
go run ./cmd/cli//http-request -url https://api.example.com/secured-endpoint -token your-api-token

# 認証トークンとJSONファイルを使用したPOSTリクエスト
go run ./cmd/cli//http-request -url https://api.example.com/secured-endpoint -method POST -json testdata/sample_request.json -token your-api-token

# Shift_JISエンコーディングのHTMLページを取得
go run ./cmd/cli//http-request -url http://abehiroshi.la.coocan.jp/ -encoding shift_jis

# エンコーディング自動検出でHTMLページを取得
go run ./cmd/cli//http-request -url http://example.com -encoding auto

# リソースの追加もしくは更新（PUT）
go run ./cmd/cli//http-request -url https://example.com/api -method PUT -json path/to/data.json

# 既存の投稿の削除（DELETE）
go run ./cmd/cli//http-request -url https://example.com/api -method DELETE
```

## エンコーディング機能

### 対応エンコーディング

- `shift_jis` (または `shift-jis`): Shift_JIS形式のレスポンスをUTF-8に変換
- `utf-8`: UTF-8形式（デフォルト、変換なし）
- `euc-jp`: EUC-JP形式のレスポンスをUTF-8に変換
- `auto`: Content-TypeヘッダーやHTMLメタタグから自動検出

### 自動検出の仕組み

`-encoding auto`を指定した場合、以下の順序で文字エンコーディングを検出します：

1. HTTPレスポンスのContent-Typeヘッダーからcharsetを抽出
2. HTMLの`<meta charset="...">`タグから検出
3. HTMLの`<meta http-equiv="Content-Type" content="text/html; charset=...">`タグから検出

### 注意事項

- エンコーディング変換はHTMLコンテンツ（Content-Type: text/html）にのみ適用されます
- JSONレスポンスなど、他のコンテンツタイプでは変換は行われません
- 不正なエンコーディング指定の場合、元のバイト列がそのまま表示されます

## エラー処理

- URLが指定されていない場合はエラーメッセージを表示
- POST/PUT/PATCHメソッドでJSONファイルが指定されていない場合はエラーメッセージを表示
- JSONファイルが存在しない、または無効なJSON形式の場合はエラーメッセージを表示
- APIリクエストが失敗した場合はエラーメッセージを表示

## Cloudflare対策と警告表示

- `GET/HEAD/OPTIONS`リクエストでは、403/429/503/520-524などのレスポンスを検知すると指数バックオフ付きで自動リトライします。
- Cloudflare固有のヘッダーや「Just a moment...」ページを検出した場合は、警告メッセージを`Warnings`欄に表示します。
- 警告が出力された場合は、一時的なブロックである可能性が高いため、時間を空けるかブラウザでクッキーを取得してから再度実行してください。
