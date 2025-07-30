# Valkey CLI Tool

Valkey データベースを操作するためのコマンドラインツールです。高性能なキー・バリューストアであるValkeyに対して、包括的なデータ操作機能を提供します。

## 機能

- **キー管理**: パターンマッチングによるキー検索・取得
- **データ操作**: 値の取得・設定・削除
- **型情報**: キーのデータ型確認
- **一括操作**: 複数キーの同時処理
- **条件検索**: 柔軟な条件でのデータ選択
- **認証対応**: パスワード認証とデータベース選択
- **安全削除**: ドライラン機能による事前確認
- **JSON出力**: 構造化されたデータの見やすい表示
- **短縮オプション**: 全てのオプションに短縮形を提供
- **エラーハンドリング**: 詳細なエラーメッセージとガイダンス

## インストール

```bash
# プロジェクトルートから
go build -o bin/valkey ./cmd/cli/valkey
```

## 使用方法

### 基本的な使用方法

```bash
./bin/valkey -operation get-keys -pattern "user:*"
```

出力例:
```
パターン 'user:*' に一致するキー (3件):
  user:123
  user:456
  user:789
```

## オプション

### 基本オプション

| オプション | 短縮形 | 説明 | デフォルト | 例 |
|-----------|--------|------|-----------|-----|
| `-operation` | `-o` | Valkey操作 (get-keys, get-value, get-type, set-value, delete-key, delete-keys, select-keys, get-all-values, delete-data) | | `-o get-keys` |
| `-help` | | ヘルプを表示 | | `-help` |

### データ操作オプション

| オプション | 短縮形 | 説明 | 例 |
|-----------|--------|------|-----|
| `-key` | `-k` | 単一キー | `-k "user:123"` |
| `-keys` | | 複数キー（カンマ区切り） | `-keys "user:123,user:456"` |
| `-pattern` | `-p` | パターン | `-p "user:*"` |
| `-value` | `-v` | 値 | `-v '{"name":"John"}'` |
| `-all` | `-a` | 全件フラグ | `-a` |
| `-dry-run` | `-dr` | ドライランフラグ | `-dr` |

### 接続オプション

| オプション | 短縮形 | デフォルト | 説明 | 例 |
|-----------|--------|-----------|------|-----|
| `-host` | `-h` | `localhost` | Valkeyホスト | `-h "redis.example.com"` |
| `-port` | | `6379` | Valkeyポート | `-port 6380` |
| `-password` | `-pass` | | Valkeyパスワード | `-pass "secret"` |
| `-database` | `-db` | `0` | データベース番号 | `-db 1` |

## 使用例

### キー取得 (get-keys)

```bash
# 基本的なパターンマッチング
./bin/valkey -operation get-keys -pattern "user:*"

# 全キーを取得
./bin/valkey -operation get-keys -pattern "*"

# 特定のプレフィックスで検索
./bin/valkey -operation get-keys -pattern "session:*"

# 短縮形を使用
./bin/valkey -o get-keys -p "cache:*"
```

出力例:
```
パターン 'user:*' に一致するキー (3件):
  user:123
  user:456
  user:789
```

### 値取得 (get-value)

```bash
# 基本的な値取得
./bin/valkey -operation get-value -key "user:123"

# 短縮形を使用
./bin/valkey -o get-value -k "session:abc123"
```

出力例:
```
キー 'user:123' の値: {"name":"John","age":30,"email":"john@example.com"}
```

### 型取得 (get-type)

```bash
# キーの型を確認
./bin/valkey -operation get-type -key "user:123"

# 短縮形を使用
./bin/valkey -o get-type -k "counter:visits"
```

出力例:
```
キー 'user:123' の型: string
```

### 値設定 (set-value)

```bash
# JSON形式の値を設定
./bin/valkey -operation set-value -key "user:123" -value '{"name":"John","age":30}'

# 単純な文字列値を設定
./bin/valkey -operation set-value -key "status" -value "active"

# 短縮形を使用
./bin/valkey -o set-value -k "counter:visits" -v "1000"
```

出力例:
```
キー 'user:123' に値を設定しました
```

### キー削除 (delete-key)

```bash
# 単一キーを削除
./bin/valkey -operation delete-key -key "temp:data"

# 短縮形を使用
./bin/valkey -o delete-key -k "cache:expired"
```

出力例:
```
キー 'temp:data' を削除しました
```

### 複数キー削除 (delete-keys)

```bash
# 複数のキーを一度に削除
./bin/valkey -operation delete-keys -keys "temp:1,temp:2,temp:3"

# セッションキーを一括削除
./bin/valkey -operation delete-keys -keys "session:abc,session:def,session:ghi"
```

出力例:
```
キー 'temp:1' を削除しました
キー 'temp:2' を削除しました
キー 'temp:3' は存在しませんでした
合計 2 個のキーを削除しました
```

### 値選択 (select-keys)

```bash
# 単一キーでの選択
./bin/valkey -operation select-keys -key "user:123"

# 複数キーでの選択
./bin/valkey -operation select-keys -keys "user:123,user:456"

# パターンでの選択
./bin/valkey -operation select-keys -pattern "session:*"

# 全データの選択
./bin/valkey -operation select-keys -all

# 短縮形を使用
./bin/valkey -o select-keys -p "cache:*"
```

出力例:
```json
選択された値:
{
  "key": "user:123",
  "value": "{\"name\":\"John\",\"age\":30}",
  "type": "string"
}
```

### 全値取得 (get-all-values)

```bash
# 指定されたキーの全値を取得
./bin/valkey -operation get-all-values -keys "user:123,user:456,user:789"

# 短縮形を使用
./bin/valkey -o get-all-values -keys "session:abc,session:def"
```

出力例:
```json
取得された全値:
{
  "values": [
    {
      "key": "user:123",
      "value": "{\"name\":\"John\",\"age\":30}",
      "type": "string"
    },
    {
      "key": "user:456",
      "value": "{\"name\":\"Alice\",\"age\":25}",
      "type": "string"
    }
  ],
  "count": 2
}
```

### データ削除 (delete-data)

```bash
# ドライラン（実際には削除しない）
./bin/valkey -operation delete-data -pattern "temp:*" -dry-run

# 実際の削除
./bin/valkey -operation delete-data -pattern "temp:*"

# 単一キーの削除
./bin/valkey -operation delete-data -key "old:data"

# 複数キーの削除
./bin/valkey -operation delete-data -keys "cache:1,cache:2"

# 短縮形を使用
./bin/valkey -o delete-data -p "expired:*" -dr
```

ドライラン出力例:
```json
削除対象データ（ドライラン）:
{
  "keys": [
    "temp:data1",
    "temp:data2",
    "temp:data3"
  ],
  "count": 3,
  "deleted": 0,
  "message": "ドライランモードで実行しました。実際のデータ削除は行われていません。"
}
```

実際の削除出力例:
```json
削除結果:
{
  "keys": [
    "temp:data1",
    "temp:data2"
  ],
  "results": {
    "temp:data1": true,
    "temp:data2": true
  },
  "count": 2,
  "deleted": 2
}
```

### 認証とデータベース選択

```bash
# パスワード認証
./bin/valkey -operation get-keys -pattern "*" -password "mypassword"

# 特定のデータベースに接続
./bin/valkey -operation get-keys -pattern "*" -database 1

# ホスト・ポート指定
./bin/valkey -operation get-keys -pattern "*" -host "redis.example.com" -port 6380

# 全オプション組み合わせ
./bin/valkey -operation get-keys -pattern "user:*" \
  -host "redis.example.com" \
  -port 6380 \
  -password "secret" \
  -database 2

# 短縮形を使用
./bin/valkey -o get-keys -p "*" -h "localhost" -pass "secret" -db 1
```

## パターンマッチング

Valkeyは以下のパターンマッチング構文をサポートします：

| パターン | 説明 | 例 | マッチするキー |
|---------|------|-----|---------------|
| `*` | 0文字以上の任意の文字 | `user:*` | `user:123`, `user:abc`, `user:` |
| `?` | 1文字の任意の文字 | `user:?` | `user:1`, `user:a` |
| `[abc]` | 指定した文字のいずれか | `user:[123]` | `user:1`, `user:2`, `user:3` |
| `[a-z]` | 範囲内の文字 | `user:[0-9]` | `user:0`, `user:5`, `user:9` |

### パターンマッチングの例

```bash
# notion で始まるキー
./bin/valkey -o get-keys -p "notion*"

# notion を含むキー
./bin/valkey -o get-keys -p "*notion*"

# 3桁の数字で終わるキー
./bin/valkey -o get-keys -p "*[0-9][0-9][0-9]"

# session: で始まり、1文字の英数字が続くキー
./bin/valkey -o get-keys -p "session:[a-zA-Z0-9]"
```

## 出力フォーマット

### 通常の出力

```
パターン 'user:*' に一致するキー (3件):
  user:123
  user:456
  user:789
```

### JSON出力

構造化されたデータは見やすいJSON形式で出力されます：

```json
{
  "values": [
    {
      "key": "user:123",
      "value": "{\"name\":\"John\",\"age\":30}",
      "type": "string"
    }
  ],
  "count": 1
}
```

### 削除操作の出力

```json
{
  "keys": ["temp:1", "temp:2"],
  "results": {
    "temp:1": true,
    "temp:2": false
  },
  "count": 2,
  "deleted": 1
}
```

## エラーハンドリング

### 操作タイプの指定忘れ

```bash
./bin/valkey
```

```
エラー: 操作タイプが指定されていません
Valkey CLIツール

使用方法:
  キー取得:
    ./bin/valkey -operation get-keys -pattern "user:*"
...
```

### 必須パラメータの不足

```bash
./bin/valkey -operation get-value
```

```
エラー: get-value操作にはkeyが必要です
```

### 無効な操作タイプ

```bash
./bin/valkey -operation invalid
```

```
エラー: 未対応の操作タイプです: invalid
```

### 接続エラー

```bash
./bin/valkey -operation get-keys -pattern "*" -host "nonexistent.host"
```

```
エラー: リポジトリの初期化に失敗しました: failed to parse valkey URL: ...
```

### JSON変換エラー

内部的にJSON変換でエラーが発生した場合：

```
エラー: JSON変換エラー: invalid character 'x' looking for beginning of value
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **依存性注入**: テスタビリティを考慮した設計
- **レイヤー分離**: 設定、サービス、リポジトリの明確な分離

### ディレクトリ構造

```
internal/valkey/
├── config/                    # 設定管理
│   ├── config.go             # 設定構造体とURL構築
│   ├── flag_parser.go        # フラグ解析実装
│   └── interfaces.go         # インターフェース定義
├── usecases/                 # ビジネスロジック
│   └── data_service.go       # データ操作サービス
└── infrastructure/           # インフラストラクチャ層
    ├── logger/               # ログ機能
    │   └── repository/       # ログリポジトリ
    └── valkey/               # Valkey接続
        ├── data_store.go     # データストア実装
        └── repository/       # データリポジトリ
```

### 使用技術

- **Go**: プログラミング言語
- **valkey-go**: Valkeyクライアントライブラリ
- **標準ライブラリ**: `encoding/json`, `flag`, `fmt`, `os`, `context`など
- **テストフレームワーク**: Go標準のtestingパッケージ

### 接続管理

- **URL形式**: `valkey://[user:password@]host:port[/database]`
- **認証**: パスワード認証をサポート
- **データベース選択**: 0-15のデータベース番号を指定可能
- **接続プール**: valkey-goライブラリによる効率的な接続管理

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/valkey ./cmd/cli/valkey

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/valkey-linux ./cmd/cli/valkey
GOOS=windows GOARCH=amd64 go build -o bin/valkey.exe ./cmd/cli/valkey
GOOS=darwin GOARCH=amd64 go build -o bin/valkey-mac ./cmd/cli/valkey
```

### テスト

```bash
# 単体テスト
go test ./internal/valkey/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/valkey/...
go tool cover -html=coverage.out -o coverage.html

# 統合テスト（Valkeyサーバーが必要）
go test -tags=integration ./internal/valkey/...
```

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です：
- Model Context Protocol経由で同じValkey操作機能にアクセス可能

## 実用的な使用例

### 開発環境でのデバッグ

```bash
# 全キーを確認
./bin/valkey -o get-keys -p "*"

# セッションデータの確認
./bin/valkey -o get-value -k "session:user123"

# キャッシュの状態確認
./bin/valkey -o select-keys -p "cache:*"

# テストデータの設定
./bin/valkey -o set-value -k "test:user" -v '{"id":1,"name":"Test User"}'

# テストデータの削除（ドライラン）
./bin/valkey -o delete-data -p "test:*" -dr
```

### 本番環境での運用

```bash
# 認証付きでセッション情報を確認
./bin/valkey -o get-keys -p "session:*" -pass "production_password" -db 1

# 特定のユーザーデータを確認
./bin/valkey -o get-value -k "user:12345" -h "prod-redis.example.com" -port 6380 -pass "secret"

# キャッシュのクリア（ドライラン）
./bin/valkey -o delete-data -p "cache:expired:*" -h "prod-redis.example.com" -pass "secret" -dr

# 統計情報の取得
./bin/valkey -o get-all-values -keys "stats:daily,stats:weekly,stats:monthly" -h "prod-redis.example.com" -pass "secret"
```

### データ移行とバックアップ

```bash
# 特定のパターンのデータを確認
./bin/valkey -o select-keys -p "old_format:*" -db 0

# 新しい形式でデータを設定
./bin/valkey -o set-value -k "new_format:user:123" -v '{"version":2,"data":{"name":"John"}}' -db 1

# 古いデータの削除（ドライラン）
./bin/valkey -o delete-data -p "old_format:*" -db 0 -dr
```

### 監視とメンテナンス

```bash
# キー数の確認
./bin/valkey -o get-keys -p "*" | grep -c "  "

# 特定のパターンのキー数確認
./bin/valkey -o select-keys -p "session:*" | jq '.count'

# 期限切れキーの確認
./bin/valkey -o get-keys -p "temp:*"

# メモリ使用量の多いキーの確認
./bin/valkey -o get-all-values -keys "large:data1,large:data2,large:data3"
```

### パフォーマンステスト

```bash
# 大量データの作成
for i in {1..1000}; do
  ./bin/valkey -o set-value -k "perf:test:$i" -v "{\"id\":$i,\"data\":\"test data $i\"}"
done

# パターンマッチングのパフォーマンス確認
time ./bin/valkey -o get-keys -p "perf:test:*"

# 一括取得のパフォーマンス確認
time ./bin/valkey -o get-all-values -keys "perf:test:1,perf:test:2,perf:test:3"
```

## ベストプラクティス

### パターンマッチング

- **具体的なパターン**: `*` よりも具体的なパターンを使用してパフォーマンスを向上
- **プレフィックス活用**: `user:*` のようにプレフィックスを活用した効率的な検索
- **範囲指定**: `[0-9]` のような範囲指定で精密な検索

### セキュリティ

- **パスワード管理**: コマンドラインでのパスワード指定は履歴に残るため、環境変数の使用を推奨
- **接続制限**: 本番環境では適切なファイアウォール設定とアクセス制限
- **権限分離**: 読み取り専用操作と書き込み操作で異なる認証情報を使用

### パフォーマンス

- **バッチ処理**: 複数の操作は `delete-keys` や `get-all-values` でまとめて実行
- **ドライラン活用**: 削除操作は必ず `-dry-run` で事前確認
- **適切なデータベース選択**: 用途に応じてデータベースを分離

### 運用

- **ログ記録**: 重要な操作は実行前後でログを記録
- **バックアップ**: 削除操作前にはデータのバックアップを確認
- **監視**: 定期的なキー数やデータサイズの監視

### トラブルシューティング

- **接続確認**: エラー時はまずValkeyサーバーの状態を確認
- **権限確認**: 認証エラーの場合はパスワードとデータベース権限を確認
- **パターン確認**: 期待した結果が得られない場合はパターンマッチング構文を確認

## よくある質問

### Q: パターン `"notion"` で検索しても結果が0件になる

A: `"notion"` は完全一致検索です。`notion:sync:20250728003930` のようなキーを検索するには `"notion*"` を使用してください。

```bash
# 正しい方法
./bin/valkey -o get-keys -p "notion*"
```

### Q: JSON出力が見づらい

A: `jq` コマンドと組み合わせて使用できます：

```bash
./bin/valkey -o select-keys -all | jq .
```

### Q: 大量のキーがある場合のパフォーマンス

A: 具体的なパターンを使用し、必要に応じてデータベースを分離してください：

```bash
# 効率的
./bin/valkey -o get-keys -p "user:active:*"

# 非効率的
./bin/valkey -o get-keys -p "*"
```

### Q: 削除操作を取り消したい

A: 削除操作は元に戻せません。必ず `-dry-run` で事前確認してください：

```bash
# 安全な方法
./bin/valkey -o delete-data -p "temp:*" -dr
# 確認後に実行
./bin/valkey -o delete-data -p "temp:*"
```

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
