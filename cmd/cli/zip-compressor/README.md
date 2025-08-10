# Zip Compressor

ファイルやディレクトリをZip形式で圧縮・展開するCLIツールです。セキュリティを重視した設計で、パストラバーサル攻撃対策を含む安全な圧縮・展開機能を提供します。

## 機能

- **ファイル圧縮**: 単一ファイルをZip形式で圧縮。約70-80%の圧縮率
- **ディレクトリ圧縮**: ディレクトリ全体を再帰的にZip形式で圧縮
- **Zip展開**: Zipファイルを指定されたディレクトリに安全に展開
- **セキュリティ対策**: パストラバーサル攻撃を防ぐ安全な展開処理
- **柔軟な入力方式**: フラグ指定と位置引数の両方に対応
- **短縮オプション**: 全てのオプションに短縮形を提供
- **エラーハンドリング**: 詳細なエラーメッセージと適切な終了コード

## インストール

```bash
# プロジェクトルートから
go build -o bin/zip-compressor ./cmd/cli/zip-compressor
```

## 使用方法

### 基本的な使用方法

```bash
./bin/zip-compressor -operation compress -path /path/to/file_or_directory
```

出力例:
```
圧縮が完了しました: /path/to/file_or_directory.zip
```

### オプション

| オプション | 短縮形 | 説明 | 必須 | 例 |
|-----------|--------|------|------|-----|
| `-operation` | `-o` | 実行する操作 (compress, decompress) | ✓ | `-o compress` |
| `-path` | `-p` | 対象ファイル/ディレクトリのパス | ✓ | `-p /path/to/file.txt` |
| `-help` | `-h` | ヘルプを表示 | | `-h` |

## 使用例

### ファイル圧縮

```bash
# 単一ファイルの圧縮
./bin/zip-compressor -operation compress -path /home/user/document.txt
./bin/zip-compressor -o compress -p /home/user/document.txt
./bin/zip-compressor compress /home/user/document.txt  # 位置引数

# 出力: document.txt.zip が作成される
```

### ディレクトリ圧縮

```bash
# ディレクトリ全体の圧縮
./bin/zip-compressor -operation compress -path /home/user/my_folder
./bin/zip-compressor -o compress -p /home/user/my_folder
./bin/zip-compressor compress /home/user/my_folder  # 位置引数

# 出力: my_folder.zip が作成される
```

### Zipファイル展開

```bash
# Zipファイルの展開
./bin/zip-compressor -operation decompress -path /home/user/archive.zip
./bin/zip-compressor -o decompress -p /home/user/archive.zip
./bin/zip-compressor decompress /home/user/archive.zip  # 位置引数

# 出力: archive_decompressed/ ディレクトリに展開される
```

### ヘルプ表示

```bash
./bin/zip-compressor -help
./bin/zip-compressor -h
```

## 出力フォーマット

### 圧縮成功時の出力

```
圧縮が完了しました: document.txt.zip
圧縮が完了しました: my_folder.zip
```

### 展開成功時の出力

```
展開が完了しました: archive_decompressed
展開が完了しました: backup_decompressed
```

### ヘルプメッセージ

```
Zip圧縮CLIツール

使用方法:
  ファイル/ディレクトリ圧縮:
    ./bin/zip-compressor -operation compress -path /path/to/file_or_directory
    ./bin/zip-compressor -o compress -p /path/to/file_or_directory
    ./bin/zip-compressor compress /path/to/file_or_directory

  Zipファイル展開:
    ./bin/zip-compressor -operation decompress -path /path/to/archive.zip
    ./bin/zip-compressor -o decompress -p /path/to/archive.zip
    ./bin/zip-compressor decompress /path/to/archive.zip

オプション:
  -operation, -o    操作タイプ (compress, decompress)
  -path, -p         対象ファイル/ディレクトリのパス
  -help, -h         このヘルプを表示

例:
  # ファイル圧縮
  ./bin/zip-compressor compress /home/user/document.txt
  # → document.txt.zip が作成される

  # ディレクトリ圧縮
  ./bin/zip-compressor compress /home/user/my_folder
  # → my_folder.zip が作成される

  # Zipファイル展開
  ./bin/zip-compressor decompress /home/user/archive.zip
  # → archive_decompressed/ ディレクトリに展開される
```

## エラーハンドリング

### 無効な操作タイプ

```bash
./bin/zip-compressor -operation invalid
```

```
エラー: 未対応の操作タイプです: invalid
```

### 必須パラメータの不足

```bash
./bin/zip-compressor
```

```
エラー: 操作タイプが指定されていません
```

### 存在しないファイル/ディレクトリ

```bash
./bin/zip-compressor compress /nonexistent/file.txt
```

```
エラー: 指定されたパスが存在しません: /nonexistent/file.txt
```

### 存在しないZipファイル

```bash
./bin/zip-compressor decompress /nonexistent/archive.zip
```

```
エラー: 指定されたZipファイルが存在しません: /nonexistent/archive.zip
```

### 無効なZipファイル

```bash
./bin/zip-compressor decompress /path/to/textfile.txt
```

```
エラー: 指定されたファイルはZipファイルではありません: /path/to/textfile.txt
```

### セキュリティエラー（パストラバーサル攻撃）

```bash
# 悪意のあるZipファイルを展開しようとした場合
./bin/zip-compressor decompress malicious.zip
```

```
エラー: 不正なパスが検出されました: ../../../etc/passwd
```

### 権限エラー

```bash
./bin/zip-compressor compress /root/protected_file.txt
```

```
エラー: ファイルを開けませんでした: /root/protected_file.txt
```

## 技術仕様

### アーキテクチャ

- **Clean Architecture**: ドメイン、ユースケース、インフラストラクチャの分離
- **SOLID原則**: インターフェースを活用した疎結合な設計
- **TDD**: テスト駆動開発によるテストファースト実装
- **セキュリティファースト**: パストラバーサル攻撃対策を組み込んだ設計

### ディレクトリ構造

```
internal/zip_compressor/
├── config/          # 設定管理
│   ├── config.go    # 設定構造体とパーサー
│   ├── flag_parser.go # フラグ解析
│   └── interfaces.go  # インターフェース定義
└── usecases/        # ビジネスロジック
    └── services.go  # サービス層
```

### 使用技術

- **Go**: プログラミング言語
- **標準ライブラリ**: `archive/zip`, `flag`, `fmt`, `os`, `path/filepath`など
- **テストフレームワーク**: Go標準のtestingパッケージ

### セキュリティ機能

- **パストラバーサル攻撃対策**: 展開時に不正なパス（`../`を含むパス）を検出・拒否
- **安全なパス処理**: `filepath.Clean`を使用した正規化処理
- **展開先制限**: 指定されたディレクトリ外への展開を防止

## 開発者向け情報

### ビルド

```bash
# 開発用ビルド
go build -o bin/zip-compressor ./cmd/cli/zip-compressor

# リリース用ビルド（複数プラットフォーム）
GOOS=linux GOARCH=amd64 go build -o bin/zip-compressor-linux ./cmd/cli/zip-compressor
GOOS=windows GOARCH=amd64 go build -o bin/zip-compressor.exe ./cmd/cli/zip-compressor
GOOS=darwin GOARCH=amd64 go build -o bin/zip-compressor-mac ./cmd/cli/zip-compressor
```

### テスト

```bash
# 単体テスト
go test ./internal/zip_compressor/...

# カバレッジ付きテスト
go test -coverprofile=coverage.out ./internal/zip_compressor/...
go tool cover -html=coverage.out -o coverage.html

# 詳細テスト実行
go test -v ./internal/zip_compressor/config
go test -v ./internal/zip_compressor/usecases
```

### テスト結果例

```bash
# 設定層テスト
=== RUN   TestNewConfig_Normal
--- PASS: TestNewConfig_Normal (0.00s)
=== RUN   TestParseFlagsWithParser_Normal
--- PASS: TestParseFlagsWithParser_Normal (0.00s)
PASS
ok      github.com/landmaster135/devbox/internal/zip_compressor/config

# サービス層テスト
=== RUN   TestNewZipCompressorService_Normal
--- PASS: TestNewZipCompressorService_Normal (0.00s)
=== RUN   TestExtractFile_PathTraversalAttack
--- PASS: TestExtractFile_PathTraversalAttack (0.00s)
PASS
ok      github.com/landmaster135/devbox/internal/zip_compressor/usecases
```

### 動作確認

```bash
# テスト用ファイル作成
echo "これはテスト用のファイルです。" > test_file.txt

# ファイル圧縮テスト
go run ./cmd/cli/zip-compressor compress test_file.txt
# 出力: 圧縮が完了しました: test_file.txt.zip

# 展開テスト
go run ./cmd/cli/zip-compressor decompress test_file.txt.zip
# 出力: 展開が完了しました: test_file.txt_decompressed

# 内容確認
cat test_file.txt_decompressed/test_file.txt
# 出力: これはテスト用のファイルです。
```

### 関連ツール

このCLIツールと同じサービスを使用するMCPツールも利用可能です（実装予定）：
- `devbox/cmd/mcp/zip_compressor/mcp.go`

MCPツールを使用することで、Model Context Protocol経由で同じ圧縮・展開機能にアクセスできます。

## トラブルシューティング

### よくある問題

1. **権限エラー**: ファイルやディレクトリへの読み書き権限を確認
2. **ディスク容量不足**: 十分な空き容量があることを確認
3. **パスの問題**: 絶対パスまたは正しい相対パスを使用
4. **Zipファイルの破損**: 有効なZipファイルであることを確認

### デバッグ方法

```bash
# 詳細なエラー情報を確認
go run ./cmd/cli/zip-compressor compress /path/to/file 2>&1

# ファイル権限確認
ls -la /path/to/file

# ディスク容量確認
df -h
```

## ライセンス

このプロジェクトのライセンスについては、プロジェクトルートのLICENSEファイルを参照してください。
