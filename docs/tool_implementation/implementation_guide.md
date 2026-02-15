# ツール実装ガイド

## 実装パターン

実装パターンおよび実装時の重要な注意点は `docs/tool_implementation/common/guide.md` を参照してください。

## 実装時のチェックリスト

CLIツール固有のチェックリストは `docs/tool_implementation/cli/guide.md` を参照してください。
MCPツール固有のチェックリストは `docs/tool_implementation/mcp/guide.md` を参照してください。

### 共通チェックリスト

- [ ] **infrastructure層は、domain層だろうが、usecases層だろうが、他のいかなる層にも依存してはならない。**
- [ ] ビジネスロジックはusecasesパッケージに実装
- [ ] エラーメッセージは日本語で分かりやすく
- [ ] パラメータの妥当性検証を実装
- [ ] テストコードを作成

## テストコードでのハードコード削減

テストコードでのハードコード削減は `docs/tool_implementation/common/testing.md` を参照してください。

### テストコード実装チェックリスト

- [ ] テスト固有の定数は各テスト関数内で定義
- [ ] 複数のテストケースはテーブル駆動テストで実装
- [ ] 共通処理はヘルパー関数として抽出
- [ ] 複雑なテストデータは構造体で管理
- [ ] モック生成処理は関数化
- [ ] テスト名は「機能_条件」の形式で命名

**参考実装**: `devbox/internal/gcloud_monitoring/usecases/services_test.go`

## 具体的な実装例の比較

### arithmetic-calculator の実装例
- **CLIツール** (`cmd/cli/arithmetic-calculator/main.go`)
- **MCPツール** (`cmd/mcp/arithmetic_calculator/mcp.go`)

### ops-for-golang の実装例
- **CLIツール** (`cmd/cli/ops-for-golang/main.go`)
- **MCPツール** (`cmd/mcp/ops_for_golang/mcp.go`)

## よくある実装ミス

CLIツール固有の実装ミスは `docs/tool_implementation/cli/guide.md` を参照してください。
MCPツール固有の実装ミスは `docs/tool_implementation/mcp/guide.md` を参照してください。

## ツール別テスト手順

- CLIツール: `docs/tool_implementation/cli/guide.md`
- MCPツール: `docs/tool_implementation/mcp/guide.md`

## まとめ

- **共通**: ビジネスロジックはusecasesパッケージで共有
- **CLI固有項目**: `docs/tool_implementation/cli/guide.md` を参照
- **MCP固有項目**: `docs/tool_implementation/mcp/guide.md` を参照
