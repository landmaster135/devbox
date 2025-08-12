# devbox プロジェクト設計仕様書

## 1. プロジェクト概要

### 1.1 プロジェクトの目的
`devbox`は開発者向けの包括的なユーティリティツール集合体です。30以上のCLIツールと20以上のMCP（Model Context Protocol）サーバーを提供し、日常的な開発作業を効率化することを目的としています。

### 1.2 主要機能
- **CLIツール群**: ファイル処理、画像変換、データ変換、コード解析など多様なユーティリティ
- **MCPサーバー群**: AI開発環境との統合を可能にする外部API連携ツール
- **クロスプラットフォーム対応**: Linux、macOS、Windows環境での動作保証
- **自動化ビルドシステム**: 効率的な開発・デプロイメントワークフロー

### 1.3 技術スタック
- **言語**: Go 1.23.5
- **アーキテクチャ**: Clean Architecture
- **ビルドシステム**: Bash スクリプトベース
- **主要依存関係**:
  - `github.com/mark3labs/mcp-go`: MCP プロトコル実装
  - 画像処理: `github.com/anthonynsimon/bild`, `github.com/gen2brain/webp`
  - PDF処理: `github.com/pdfcpu/pdfcpu`
  - データベース: `github.com/lib/pq` (PostgreSQL)
  - 外部API: `golang.org/x/oauth2`, `google.golang.org/api`

## 2. アーキテクチャ設計

### 2.1 設計思想
プロジェクトはClean Architectureの原則に基づいて設計されており、以下の特徴を持ちます：

- **関心の分離**: ビジネスロジック、インフラストラクチャ、プレゼンテーション層の明確な分離
- **依存関係の逆転**: 外部依存への抽象化によるテスタビリティの向上
- **単一責任原則**: 各コンポーネントが明確な責任を持つ設計
- **拡張性**: 新機能追加時の既存コードへの影響最小化

### 2.2 レイヤー構造

```mermaid
graph TD
    A[cmd Layer] --> B[internal/usecases]
    A --> C[internal/interfaces]
    A --> D[internal/domain]
    
    B --> D
    C --> D
    
    E[scripts] --> F[pkg/bin]
    A --> G[util]
    B --> G
    C --> G
    
    F --> H[Cross-platform Binaries]
    G --> I[Standard Library]

    classDef cmd fill:#f96,stroke:#333,stroke-width:2px;
    classDef internal fill:#bbf,stroke:#333,stroke-width:1px;
    classDef scripts fill:#bfb,stroke:#333,stroke-width:1px;
    classDef pkg fill:#fbf,stroke:#333,stroke-width:1px;
    classDef util fill:#ddd,stroke:#333,stroke-width:1px;

    class A cmd;
    class B,C,D internal;
    class E scripts;
    class F,H pkg;
    class G util;
```

## 3. ディレクトリ構造詳細

### 3.1 cmd/ - エントリーポイント層
```
cmd/
├── cli/          # CLIツール群（30以上のツール）
│   ├── arithmetic-calculator/
│   ├── file-processor/
│   ├── image-converter/
│   ├── json-formatter-for-agent-interaction/
│   └── ...
├── mcp/          # MCPサーバー群
│   ├── router.go # MCPサーバールーティング
│   ├── arithmetic_calculator/
│   ├── brave_search/
│   ├── github/
│   └── ...
└── powershell/   # Windows PowerShell スクリプト
```

**責任範囲**:
- コマンドライン引数の解析
- 設定の初期化
- ユースケース層への処理委譲
- エラーハンドリングと出力制御

### 3.2 internal/ - ビジネスロジック層
```
internal/
├── {tool_name}/
│   ├── domain/      # ドメインモデル・エンティティ
│   ├── usecases/    # ビジネスロジック・ユースケース
│   └── interfaces/  # 外部システムインターフェース
```

**各サブレイヤーの責任**:
- **domain/**: ビジネスルール、エンティティ、リポジトリインターフェース
- **usecases/**: アプリケーション固有のビジネスロジック
- **interfaces/**: 外部API、データベース、ファイルシステムとの連携

### 3.3 pkg/ - ビルド成果物管理
```
pkg/
├── bin/          # 実行可能バイナリ
│   ├── cli/      # CLIツールバイナリ
│   │   ├── linux_amd64/
│   │   ├── darwin_arm64/
│   │   └── win_amd64/
│   └── mcp/      # MCPサーバーバイナリ
│       ├── linux_amd64/
│       ├── darwin_arm64/
│       └── win_amd64/
└── dos/          # Windows バッチファイル
```

**管理対象**:
- クロスプラットフォーム対応バイナリファイル
- Windows環境向けDOSバッチファイル
- デプロイメント用パッケージ

### 3.4 scripts/ - ビルド自動化
```
scripts/
├── build.sh                    # メインビルドスクリプト
├── build_mcp_tools.sh         # MCPツール専用ビルド
├── create_project_files.sh    # プロジェクト雛形生成
└── build_{tool_name}.sh       # 個別ツールビルド
```

**機能**:
- 統合ビルドプロセスの自動化
- クロスプラットフォームコンパイル
- 依存関係管理
- テスト実行とカバレッジ計測

## 4. MCPサーバーシステム

### 4.1 アーキテクチャ概要
MCPサーバーシステムは`cmd/mcp/router.go`を中心とした統一的なルーティングシステムを採用しています。

```go
// router.goの主要構造
func Router() {
    args := os.Args
    switch args[1] {
    case "arith_calc":
        arithmetic_calculator.BuildArithCalculatorServer()
    case "github":
        github.BuildGitHubServer()
    // ... 20以上のサーバー定義
    }
}
```

### 4.2 提供MCPサーバー一覧

| サーバー名 | 機能概要 | 主要用途 |
|-----------|----------|----------|
| `arithmetic_calculator` | 数値計算処理 | 基本的な算術演算 |
| `datetime_calculator` | 日時計算処理 | 日付・時刻の操作・変換 |
| `brave_search` | Brave検索API | Web検索機能 |
| `duckduckgo_search` | DuckDuckGo検索API | プライバシー重視検索 |
| `github` | GitHub API連携 | リポジトリ操作・Issue管理 |
| `postgresql` | PostgreSQL連携 | データベース操作 |
| `filesystem` | ファイルシステム操作 | ファイル・ディレクトリ管理 |
| `http_request` | HTTP通信 | 外部API呼び出し |
| `figma` | Figma API連携 | デザインファイル操作 |
| `everart` | EverArt画像生成 | AI画像生成 |
| `youtube_transcript` | YouTube字幕取得 | 動画字幕データ抽出 |
| `sequentialthinking` | 段階的思考支援 | 問題解決プロセス支援 |
| `context7` | ライブラリドキュメント | 技術文書検索・取得 |
| `timezone` | タイムゾーン変換 | 時刻変換・管理 |
| `git_diff_recorder` | Git差分記録 | コード変更履歴管理 |
| `git_commit_history_retriever` | Gitコミット履歴 | コミット履歴分析 |
| `service_implementing_viewer` | サービス実装状況 | 開発進捗可視化 |
| `gdrive` | Google Drive連携 | クラウドストレージ操作 |

### 4.3 拡張パターン
新しいMCPサーバーを追加する際の標準的な手順：

1. `cmd/mcp/{server_name}/` ディレクトリ作成
2. `internal/{server_name}/` でビジネスロジック実装
3. `router.go` にルーティング追加
4. `scripts/build_{server_name}.sh` ビルドスクリプト作成

## 5. ビルドシステム

### 5.1 ビルドフロー概要
```mermaid
graph LR
    A[scripts/build.sh] --> B[run_all_sh_scripts]
    B --> C[Individual Build Scripts]
    C --> D[Cross-platform Compilation]
    D --> E[pkg/bin/ Output]
    
    F[scripts/build_mcp_tools.sh] --> G[MCP Tools Build]
    G --> H[pkg/bin/mcp/ Output]
```

### 5.2 主要ビルドスクリプト

**build.sh**
- 全ビルドスクリプトの統合実行
- エラーハンドリングと実行結果レポート
- スキップ機能による柔軟な実行制御

**build_mcp_tools.sh**
- MCPツール専用のクロスプラットフォームビルド
- Linux/AMD64、macOS/ARM64、Windows/AMD64対応
- バイナリサイズ最適化（`-ldflags="-s -w" -trimpath`）

### 5.3 クロスプラットフォーム対応
```bash
# 例: MCPツールのマルチプラットフォームビルド
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${LINUX_AMD64_DIR}/${output_name}" "${package}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${package}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o "${MAC_ARM64_DIR}/${output_name}" "${package}"
```

## 6. 主要CLIツール群

### 6.1 ファイル処理系
- `file-processor`: 汎用ファイル処理
- `file-maneuver`: ファイル操作・移動
- `file-character-replacer`: 文字列置換処理

### 6.2 画像処理系
- `image-converter`: 画像形式変換
- `image-filterer`: 画像フィルタリング
- `image-renamer`: 画像ファイル名変更
- `image-rotator`: 画像回転処理
- `image-trimmer`: 画像トリミング

### 6.3 データ変換系
- `json-formatter-for-agent-interaction`: JSON整形
- `json-modifier`: JSON データ操作
- `yaml-parser`: YAML解析処理
- `base64-extractor`: Base64エンコード・デコード

### 6.4 開発支援系
- `code-analyzer`: コード解析
- `depends-visualizer`: 依存関係可視化
- `git-commit-history-retriever`: Git履歴分析
- `service-implementing-viewer`: サービス実装状況確認

## 7. 開発・運用指針

### 7.1 新機能追加手順
1. **要件定義**: 機能仕様とインターフェース設計
2. **ディレクトリ作成**: `cmd/cli/{tool_name}` および `internal/{tool_name}`
3. **Clean Architecture実装**:
   - `internal/{tool_name}/domain/`: エンティティ・リポジトリインターフェース
   - `internal/{tool_name}/usecases/`: ビジネスロジック
   - `internal/{tool_name}/interfaces/`: 外部システム連携
   - `cmd/cli/{tool_name}/main.go`: エントリーポイント
4. **テスト実装**: TDD原則に基づくテストコード作成
5. **ビルドスクリプト作成**: `scripts/build_{tool_name}.sh`
6. **ドキュメント更新**: README.md および本仕様書の更新

### 7.2 テスト戦略
- **単体テスト**: 各レイヤーの独立したテスト
- **統合テスト**: レイヤー間の連携テスト
- **カバレッジ目標**: 57.0%以上（現在の水準維持・向上）
- **テストツール**: `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock`

### 7.3 コード品質管理
- **命名規則**: Go標準に準拠（PascalCase、camelCase、snake_case）
- **SOLID原則**: 設計原則の遵守
- **依存関係管理**: go.modによる明示的な依存関係管理
- **静的解析**: `go vet`, `golint`の活用

### 7.4 デプロイメント戦略
- **バイナリ配布**: クロスプラットフォーム対応バイナリの提供
- **バージョン管理**: セマンティックバージョニング
- **リリースプロセス**: 自動化されたビルド・テスト・パッケージング

## 8. 今後の拡張計画

### 8.1 機能拡張
- 新しいMCPサーバーの追加（AI/ML関連API連携）
- CLIツールの機能強化（バッチ処理、並列処理対応）
- Web UI提供による操作性向上

### 8.2 技術的改善
- パフォーマンス最適化
- メモリ使用量削減
- エラーハンドリングの強化
- ログ機能の充実

### 8.3 運用改善
- CI/CDパイプラインの構築
- 自動テスト・デプロイメント
- モニタリング・アラート機能
- ドキュメント自動生成

---

## 付録

### A. 依存関係一覧
主要な外部依存関係とその用途：

- **MCP関連**: `github.com/mark3labs/mcp-go`
- **画像処理**: `github.com/anthonynsimon/bild`, `github.com/gen2brain/webp`, `github.com/gen2brain/avif`
- **PDF処理**: `github.com/pdfcpu/pdfcpu`
- **データベース**: `github.com/lib/pq`
- **Web関連**: `github.com/PuerkitoBio/goquery`
- **認証**: `golang.org/x/oauth2`
- **Google API**: `google.golang.org/api`
- **テスト**: `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock`

### B. 設定ファイル
- `go.mod`: Go モジュール定義
- `.gitignore`: Git除外設定
- `LICENSE`: MITライセンス

### C. サンプルデータ
`sample_data/` ディレクトリに各ツールのテスト用サンプルファイルを配置。

---

*本仕様書は devbox プロジェクトの設計思想と実装詳細を記録したものです。プロジェクトの進化に合わせて継続的に更新されます。*
