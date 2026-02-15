# MCPツール実装ガイド（入口）

このドキュメントは、MCPツール実装の入口です。

## MCP固有の実装手順

1. `cmd/mcp/{server_name}/` と必要なハンドラーを実装する
2. `mcp.NewTool(...)` でツール定義（説明・必須パラメータ・型）を追加する
3. MCP から `usecases` を呼び出す構成にする
4. 必須入力は `request.Require*`、任意入力は `request.Get*` で取得する
5. 処理結果は `mcp.CallToolResult` として返却し、標準出力は使用しない
6. `cmd/mcp/router.go` にルーティングを追加する

## 実装パターン

実装パターンとアンチパターンの詳細は `docs/tool_implementation/mcp/guide.md` を参照してください。
