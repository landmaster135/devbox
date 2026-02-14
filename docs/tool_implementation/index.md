# cmd配下ツール実装・運用ガイド

`cmd/` 配下の実装作業に関する内容をまとめたものです。

## 1. 実装方式の選定（CLI / MCP）
1. 既存実装を確認する: `docs/project_status/service_implementation_status.md` を参照し、重複を回避
2. 実装ガイドを確認する

## 2. 実装計画の作成
ユーザーから実装計画を作るように指示された場合に、実装計画を下記の手順で作成する。
1. 後述の共通実装フロー（CLI / MCP共通）を参考に実装計画を作成する。
2. 変更予定の対象ファイルも実装計画内で明示的に列挙する。新規作成するファイルは「(新規)」であることを明記する。
3. 実装計画の作成が完了したら、作成した実装計画内で、共通実装フロー（CLI / MCP共通）内の項目の記載漏れが無いかどうかを確かめる。記載漏れがあったら適宜修正する。

## 3. 共通実装フロー（CLI / MCP共通）
1. **要件定義**: 機能仕様、入出力、エラー設計を定義
2. **ディレクトリ作成**: `internal/{tool_name}` を作成して共通ロジックを配置
3. **レイヤー実装**: `domain/`、`usecases/`、`infrastructures/` を実装
4. **エントリーポイント実装**: CLI または MCP から `usecases` を呼び出す構成に統一
5. **テスト実装**: 単体テストと統合テストを追加
6. **ビルドスクリプト作成**: `scripts/build_{tool_name}.sh` を追加
7. **ドキュメント更新**: `docs/tool_implementation/documentation_guide.md` を参照

## 4. 実装ガイド

| 区分 | 参照先 |
|---|---|
| 共通 | `docs/tool_implementation/implementation_guide.md` |
| CLI | `docs/tool_implementation/cli/index.md` |
| MCP | `docs/tool_implementation/mcp/index.md` |

## 5. テスト戦略
- **単体テスト**: 各レイヤーの独立したテスト
- **統合テスト**: レイヤー間の連携テスト
- **CLI観点**: フラグ解析、標準出力、標準エラー出力、終了コード
- **MCP観点**: ツール定義、必須/任意パラメータ取得、`CallToolResult` 返却
- **カバレッジ目標**: 90%以上
- **テストコマンド**: `go test -v ./... -coverpkg=./... -covermode=count -coverprofile=coverage.out`
- **テストツール**: `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock`

## 6. コード品質管理
- **命名規則**: Go標準に準拠（PascalCase、camelCase、snake_case）
- **コード整形**: `go fmt ./...` または `goimports` による自動整形
- **SOLID原則**: 設計原則の遵守
- **依存関係管理**: go.modによる明示的な依存関係管理
- **静的解析**: `go vet` などの標準ツールを活用

## 7. ビルドシステム

### 7.1 ビルドフロー概要
```mermaid
graph LR
    A[scripts/build.sh] --> B[run_all_sh_scripts]
    B --> C[Individual Build Scripts]
    C --> D[Cross-platform Compilation]
    D --> E[pkg/bin/ Output]
    
    F[scripts/build_mcp_tools.sh] --> G[MCP Tools Build]
    G --> H[pkg/bin/mcp/ Output]
```

### 7.2 主要ビルドスクリプト

**build.sh**
- 全ビルドスクリプトの統合実行
- エラーハンドリングと実行結果レポート
- スキップ機能による柔軟な実行制御

**build_mcp_tools.sh**
- MCPツール専用のクロスプラットフォームビルド
- Linux/AMD64、macOS/ARM64、Windows/AMD64対応
- バイナリサイズ最適化（`-ldflags="-s -w" -trimpath`）

### 7.3 クロスプラットフォーム対応
```bash
# 例: MCPツールのマルチプラットフォームビルド
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${LINUX_AMD64_DIR}/${output_name}" "${package}"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -trimpath -o "${WIN_AMD64_DIR}/${WIN_OUTPUT_NAME}" "${package}"
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -trimpath -o "${MAC_ARM64_DIR}/${output_name}" "${package}"
```

## 8. デプロイメント戦略
- **バイナリ配布**: クロスプラットフォーム対応バイナリの提供
- **バージョン管理**: セマンティックバージョニング
- **リリースプロセス**: 自動化されたビルド・テスト・パッケージング
