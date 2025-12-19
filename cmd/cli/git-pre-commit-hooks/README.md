# Git pre-commit hooks

Git pre-commit hook用のシークレット等に対する検知ツールです。JSON設定ファイル内の機密情報と、全ファイル内の禁止されたホームパスを自動検知し、コミット前にブロックします。

## 概要

このツールは、以下の2つの主要機能を提供します：

1. **シークレット検知**: MCP（Model Context Protocol）設定ファイルなどのJSON設定ファイル内に含まれる可能性のある機密情報（APIキー、トークン、パスワードなど）を検知
2. **禁止パス検知**: 全ファイル内のホームパスを検知し、`/home/user`以外の場合はコミットをブロック

## 機能

### シークレット検知対象

- **疑わしいキー名**: `api_key`, `secret_key`, `access_token`, `password`など
- **実際のシークレットパターン**: OpenAI APIキー、GitHub PAT、AWS Access Keyなど
- **高エントロピー値**: ランダムな文字列の検出

### ホームパス検知対象

- **禁止パス**: `/home/[username]`形式のパス（実際のユーザー名を含むパス）
- **許可パス**: `/home/user`（汎用的なユーザーホームパス）

### 許可される値（シークレット検知）

- プレースホルダー値: `YOUR_API_KEY`, `PLACEHOLDER`など
- `YOUR_`で始まる値
- 短い値（8文字未満）
- テスト用の値: `test`, `demo`, `example`など

## インストール

```bash
cd devbox
go build -o bin/git-pre-commit-hooks cmd/cli/git-pre-commit-hooks/main.go
```

## 使用方法

### 基本構文

```bash
./bin/git-pre-commit-hooks [オプション]
```

### オプション

- `--version`: バージョン情報を表示
- `--verbose`: 詳細な出力を表示
- `--config-file <path>`: 特定の設定ファイルのみをチェック
- `--dry-run`: 実際のGit操作なしでテスト実行
- `--help`: ヘルプを表示

## 使用例

### 1. 手動実行

```bash
go run ./cmd/cli/git-pre-commit-hooks
```

### 2. バージョン確認

```bash
go run ./cmd/cli/git-pre-commit-hooks --version
```

出力:
```
Secret Detector v1.0.0
```

### 3. 詳細出力で実行

```bash
go run ./cmd/cli/git-pre-commit-hooks --verbose
```

### 4. 特定ファイルのみチェック

```bash
go run ./cmd/cli/git-pre-commit-hooks --config-file=config.json
```

### 5. ドライランモード

```bash
go run ./cmd/cli/git-pre-commit-hooks --dry-run
```

### 6. Git pre-commitフックとして使用

```bash
# .git/hooks/pre-commit に設定
#!/bin/bash
./bin/git-pre-commit-hooks
```

## 検知例

### 安全な設定例（シークレット検知）

```json
{
  "mcpServers": {
    "weather-server": {
      "command": "weather-mcp",
      "env": {
        "API_KEY": "YOUR_API_KEY",
        "SECRET": "PLACEHOLDER"
      }
    }
  }
}
```

出力:
```
🔍 Checking for secrets and forbidden paths in files...
📋 Checking JSON configuration files for secrets...
📁 Checking configuration file: config.json
  Found 1 server(s) in configuration
  Found 2 environment variable(s) total
    ✅ weather-server.API_KEY: safe placeholder value
    ✅ weather-server.SECRET: safe placeholder value

  Summary: Found 2 placeholder(s), 0 potential secret(s)

📋 Checking all files for /home paths...
  Summary: Found 0 /home path(s) in 1 file(s)

📊 Analysis Results:
  Overall Summary: Found 0 /home path(s), 2 placeholder(s), 0 potential secret(s)

✅ Scan completed. No issues detected in 1 file(s).
```

### 危険な設定例（シークレット検知）

```json
{
  "mcpServers": {
    "openai-server": {
      "command": "openai-mcp",
      "env": {
        "OPENAI_API_KEY": "sk-1234567890abcdef1234567890abcdef1234567890abcdef"
      }
    }
  }
}
```

### 禁止パス検知例

ファイル内容:
```
# 設定ファイル
LOG_PATH=/home/alice/logs/app.log
CONFIG_PATH=/home/user/config/app.conf
```

出力:
```
🔍 Checking for secrets and forbidden paths in files...
📋 Checking JSON configuration files for secrets...
ℹ️  No configuration files found in this commit.

📋 Checking all files for /home paths...
    ❌ config.txt:1: forbidden /home path detected
      Content: LOG_PATH=/home/alice/logs/app.log
    ✅ config.txt:2: allowed /home path
      Content: CONFIG_PATH=/home/user/config/app.conf

  Summary: Found 2 /home path(s) in 1 file(s)

📊 Analysis Results:
  Overall Summary: Found 2 /home path(s), 0 placeholder(s), 0 potential secret(s)

❌ COMMIT BLOCKED: Issues detected!

💡 Suggestions:
📋 For /home paths:
  1. Replace absolute paths with relative paths or environment variables
  2. Use '/home/user' if you need to reference a generic user home
  3. Consider using placeholders like '$HOME' or '~'
  4. If this is intentional, use: git commit --no-verify
```

## 対応ファイル形式

### シークレット検知対象

- `*.json`
- `*.config.js`
- `*.config.ts`
- `mcp_settings.json`
- `claude_desktop_config.json`
- `cline_mcp_settings.json`

### ホームパス検知対象

- 全てのテキストファイル（バイナリファイルは除外）
- 除外されるバイナリファイル拡張子:
  - 画像: `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`, `.tiff`, `.webp`
  - 動画: `.mp3`, `.mp4`, `.avi`, `.mov`, `.wmv`, `.flv`, `.mkv`
  - ドキュメント: `.pdf`, `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`
  - アーカイブ: `.zip`, `.tar`, `.gz`, `.bz2`, `.7z`, `.rar`
  - 実行ファイル: `.exe`, `.dll`, `.so`, `.dylib`, `.a`, `.o`, `.obj`
  - データベース: `.bin`, `.dat`, `.db`, `.sqlite`, `.sqlite3`

## 検知パターン

### 疑わしいキー名パターン

- `api[_-]?key`
- `secret[_-]?key`
- `access[_-]?token`
- `private[_-]?key`
- `client[_-]?secret`
- `auth[_-]?token`
- `bearer[_-]?token`
- `webhook[_-]?url`
- `database[_-]?url`
- `password`
- `secret`
- `token`
- `key`

### 実際のシークレットパターン

- OpenAI API key: `sk-[a-zA-Z0-9]{48}`
- Slack token: `xoxp-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-zA-Z0-9]{32}`
- GitHub PAT: `ghp_[a-zA-Z0-9]{36}`
- Google API key: `AIza[0-9A-Za-z_-]{35}`
- AWS access key: `AKIA[0-9A-Z]{16}`
- Google OAuth token: `ya29\.[a-zA-Z0-9_-]+`
- UUID format: `[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`

### ホームパス検知パターン

- **検知対象**: `/home`を含む文字列
- **許可パターン**: `/home/user`を含む文字列
- **禁止パターン**: `/home/[username]`というような実際のユーザー名を含む文字列

## エラーハンドリング

- Gitリポジトリでない場合のエラー
- ステージされたファイルの取得エラー
- JSON解析エラー
- ファイル読み込みエラー
- バイナリファイルの自動スキップ

## アーキテクチャ

```
cmd/cli/git-pre-commit-hooks/main.go         # CLIエントリーポイント
internal/secret_detector/
├── config/
│   ├── config.go                      # JSON設定ファイル読み込み
│   └── file_system.go                 # ファイルシステム操作
├── domain/                            # ドメインモデル
│   ├── models.go                      # 構造体定義
│   └── patterns.go                    # 検知パターン定義
└── usecases/
    ├── services.go                    # ビジネスロジック
    ├── services_test.go               # テストコード
    └── repositories.go                # リポジトリインターフェース
```

## Git Hookとの統合

### 自動セットアップ（元のスクリプト使用）

```bash
./scripts/setup-git-pre-commit-hooks.sh
```

### 手動セットアップ

1. ツールをビルド:
```bash
go build -o .git/hooks/git-pre-commit-hooks cmd/cli/git-pre-commit-hooks/main.go
```

2. pre-commitフックを作成:
```bash
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
.git/hooks/git-pre-commit-hooks
EOF
chmod +x .git/hooks/pre-commit
```

## バイパス方法

緊急時やテスト時にチェックをスキップする場合:

```bash
git commit --no-verify
```

## 終了コード

- `0`: 成功（問題未検出）
- `1`: 失敗（シークレットまたは禁止パス検出、またはエラー）

## 設定のカスタマイズ

パターンや許可値をカスタマイズしたい場合は、以下のファイルを編集してください:

- `internal/secret_detector/domain/patterns.go`: 検知パターンの定義
- `internal/secret_detector/usecases/services.go`: 検知ロジックの調整

### カスタマイズ例

**新しいバイナリファイル拡張子の追加**

`domain/patterns.go`の`GetBinaryFileExtensions()`に追加:

```go
func GetBinaryFileExtensions() []string {
  return []string{
    // 既存の拡張子...
    ".custom", ".binary",  // 新しい拡張子を追加
  }
}
```

## テスト

### 単体テストの実行

```bash
cd devbox
go test ./internal/secret_detector/usecases -v
```

### テストカバレッジの確認

```bash
cd devbox
go test -coverprofile=coverage.out ./internal/secret_detector/usecases
go tool cover -html=coverage.out -o coverage.html
```

## トラブルシューティング

### よくある問題

1. **"This directory is not a Git repository"**
  - Gitリポジトリ内で実行してください

2. **"No configuration files found"**
  - ステージされた設定ファイルがない場合の正常な動作です

3. **"Error checking file"**
  - JSONファイルの形式を確認してください

4. **"Skipping binary file"**
  - バイナリファイルは自動的にスキップされます（正常な動作）

5. **"/home path detected but allowed"**
  - `/home/user`パスは許可されています（正常な動作）

### デバッグ

詳細な情報が必要な場合は、`--verbose`オプションを使用してください:

```bash
go run ./cmd/cli/git-pre-commit-hooks --verbose
```

## 更新履歴

### v1.1.0 (最新)
- 新機能: 全ファイル内の`/home`パス検知機能を追加
- `/home/user`パスは許可、その他の`/home/[username]`パスは禁止
- バイナリファイルの自動除外機能
- 統合された結果表示とエラーメッセージ
- 包括的なテストスイートの追加

### v1.0.0
- 初期リリース
- JSON設定ファイル内のシークレット検知機能
- Git pre-commit hook対応
- 基本的なCLIオプション

## ライセンス

このプロジェクトのライセンスに従います。
