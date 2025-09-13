# Secret Detector

Git pre-commit hook用のGo製シークレット検知ツールです。JSON設定ファイル内の機密情報を自動検知し、コミット前にブロックします。

## 概要

このツールは、MCP（Model Context Protocol）設定ファイルなどのJSON設定ファイル内に含まれる可能性のある機密情報（APIキー、トークン、パスワードなど）を検知します。元々`/home/nov/devbox/scripts/setup-secret-detector.sh`内にあったGoコードを適切なGoプロジェクト構造に移植して実装されています。

## 機能

### 検知対象

- **疑わしいキー名**: `api_key`, `secret_key`, `access_token`, `password`など
- **実際のシークレットパターン**: OpenAI APIキー、GitHub PAT、AWS Access Keyなど
- **高エントロピー値**: ランダムな文字列の検出

### 許可される値

- プレースホルダー値: `YOUR_API_KEY`, `PLACEHOLDER`など
- `YOUR_`で始まる値
- 短い値（8文字未満）
- テスト用の値: `test`, `demo`, `example`など

## インストール

```bash
cd devbox
go build -o bin/secret-detector cmd/cli/secret-detector/main.go
```

## 使用方法

### 基本構文

```bash
./bin/secret-detector [オプション]
```

### オプション

- `--version`: バージョン情報を表示

## 使用例

### 1. 手動実行

```bash
go run ./cmd/cli/secret-detector
```

### 2. バージョン確認

```bash
go run ./cmd/cli/secret-detector --version
```

出力:
```
Secret Detector v1.0.0
```

### 3. Git pre-commitフックとして使用

```bash
# .git/hooks/pre-commit に設定
#!/bin/bash
./bin/secret-detector
```

## 検知例

### 安全な設定例

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
🔍 Checking for secrets in JSON configuration files...
📁 Checking configuration file: config.json
  Found 1 server(s) in configuration
  Found 2 environment variable(s) total
    ✅ weather-server.API_KEY: safe placeholder value
    ✅ weather-server.SECRET: safe placeholder value

  Summary: Found 2 placeholder(s), 0 potential secret(s)

✅ Secret scan completed. No actual secrets detected in 1 file(s).
```

### 危険な設定例

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

出力:
```
🔍 Checking for secrets in JSON configuration files...
📁 Checking configuration file: config.json
  Found 1 server(s) in configuration
  Found 1 environment variable(s) total
    ❌ openai-server.OPENAI_API_KEY: suspicious pattern detected
      Pattern: sk-[a-zA-Z0-9]{48}
      Value: sk-1234567890abcdef1234567890abcdef1234567890abcdef

  Summary: Found 0 placeholder(s), 1 potential secret(s)

❌ COMMIT BLOCKED: Potential secrets detected!

💡 Suggestions:
  1. Replace actual values with placeholders like 'YOUR_API_KEY'
  2. Move secrets to environment variables or secure storage
  3. Use a separate config file that's git-ignored
  4. If this is intentional, use: git commit --no-verify
```

## 対応ファイル形式

- `*.json`
- `*.config.js`
- `*.config.ts`
- `mcp_settings.json`
- `claude_desktop_config.json`
- `cline_mcp_settings.json`

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

## エラーハンドリング

- Gitリポジトリでない場合のエラー
- ステージされたファイルの取得エラー
- JSON解析エラー
- ファイル読み込みエラー

## アーキテクチャ

```
cmd/cli/secret-detector/main.go         # CLIエントリーポイント
internal/secret_detector/
├── config/config.go                    # JSON設定ファイル読み込み
├── domain/                            # ドメインモデル
│   ├── models.go                      # 構造体定義
│   └── patterns.go                    # 検知パターン定義
└── usecases/services.go               # ビジネスロジック
```

## Git Hookとの統合

### 自動セットアップ（元のスクリプト使用）

```bash
./scripts/setup-secret-detector.sh
```

### 手動セットアップ

1. ツールをビルド:
```bash
go build -o .git/hooks/secret-detector cmd/cli/secret-detector/main.go
```

2. pre-commitフックを作成:
```bash
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
.git/hooks/secret-detector
EOF
chmod +x .git/hooks/pre-commit
```

## バイパス方法

緊急時やテスト時にチェックをスキップする場合:

```bash
git commit --no-verify
```

## 終了コード

- `0`: 成功（シークレット未検出）
- `1`: 失敗（シークレット検出またはエラー）

## 設定のカスタマイズ

パターンや許可値をカスタマイズしたい場合は、以下のファイルを編集してください:

- `internal/secret_detector/domain/patterns.go`: 検知パターンの定義
- `internal/secret_detector/usecases/services.go`: 検知ロジックの調整

## トラブルシューティング

### よくある問題

1. **"This directory is not a Git repository"**
   - Gitリポジトリ内で実行してください

2. **"No configuration files found"**
   - ステージされた設定ファイルがない場合の正常な動作です

3. **"Error checking file"**
   - JSONファイルの形式を確認してください

### デバッグ

詳細な情報が必要な場合は、ソースコードの`fmt.Printf`文を参考にしてください。

## ライセンス

このプロジェクトのライセンスに従います。
