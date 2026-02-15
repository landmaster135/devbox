# cmd配下ツール実装・運用ガイド

`cmd/` 配下の実装作業に関する内容をまとめたものです。

## 1. 実装方式の選定（CLI / MCP）
1. 既存実装を確認する: `docs/project_status/service_implementation_status.md` を参照し、重複を回避
2. 実装ガイドを確認する

## 2. 実装計画の作成
ユーザーから実装計画を作るように指示された場合に、実装計画を下記の手順で作成する。
1. 後述の共通実装フロー（CLI / MCP共通）を参考に実装計画を作成する。
2. 変更予定の対象ファイルも実装計画内で明示的に列挙する。新規作成するファイルは「(新規)」であることを明記する。
3. 実装計画内の各工程で参照するべき資料を明示して確認するように促すように記載する。
4. 実装計画の作成が完了したら、作成した実装計画内で、共通実装フロー（CLI / MCP共通）内の項目の記載漏れが無いかどうかを確かめる。記載漏れがあったら適宜修正する。

## 3. 共通実装フロー（CLI / MCP共通）
1. **要件定義**: 機能仕様、入出力、エラー設計を定義
2. **ディレクトリ作成**: `internal/{tool_name}` を作成して共通ロジックを配置
3. **レイヤー実装**: `docs/tool_implementation/common/guide.md`を参照  
4. **エントリーポイント実装**: 
  - CLI: `docs/tool_implementation/cli/index.md`を参照
  - MCP: `docs/tool_implementation/mcp/index.md`を参照
5. **テスト実装**: `docs/tool_implementation/common/testing.md` を参照
6. **ドキュメント更新**: `docs/tool_implementation/documentation_guide.md` を参照
7. **ビルドスクリプト作成**: `docs/tool_implementation/build_guide.md` を参照

## 4. コード品質管理
- **命名規則**: Go標準に準拠（PascalCase、camelCase、snake_case）
- **コード整形**: `go fmt ./...` または `goimports` による自動整形
- **SOLID原則**: 設計原則の遵守
- **依存関係管理**: go.modによる明示的な依存関係管理
- **静的解析**: `go vet` などの標準ツールを活用

## 5. デプロイメント戦略
- **バイナリ配布**: クロスプラットフォーム対応バイナリの提供
- **バージョン管理**: セマンティックバージョニング
- **リリースプロセス**: 自動化されたビルド・テスト・パッケージング
