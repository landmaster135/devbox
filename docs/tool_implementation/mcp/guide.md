# MCPツール実装ガイド

このドキュメントは、`docs/tool_implementation/implementation_guide.md` から MCPツール専用の内容を切り出したものです。

## 関連ドキュメント

- 共通実装ガイド: `docs/tool_implementation/implementation_guide.md`
- ドキュメント更新手順: `docs/tool_implementation/documentation_guide.md`

## 実装時の重要な注意点

### 出力制御（標準出力を使わない）

CLIツールと異なり、MCPツールは結果を標準出力に表示せず、`mcp.NewToolResultText(...)` で返却します。

```go
func handleOperation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
  service := usecases.NewService()
  result, err := service.HandleOperation("param")
  if err != nil {
    return nil, fmt.Errorf("処理に失敗しました: %v", err)
  }

  // 必須: MCPクライアントに結果を返却（標準出力は使用しない）
  return mcp.NewToolResultText(result), nil
}
```

### タイムアウト対策

MCPツールでは長大な標準出力を行うとタイムアウトが発生します。

```go
// ❌ 悪い例: 標準出力への出力
fmt.Printf("処理中...\n")
fmt.Printf("実行コマンド: %s\n", cmd)

// ✅ 良い例: 標準出力を抑制
// 進捗表示や実行コマンドの出力は行わない
```

### パラメータ取得の使い分け

- 必須パラメータは `request.RequireString()` を使用
- 任意パラメータは `request.GetString("name", "default")` を使用

## 実装パターン

### MCPツールの `mcp.go` 構造

```go
package tool_name

import (
  "context"
  "fmt"
  "os"
  mcp "github.com/mark3labs/mcp-go/mcp"
  server "github.com/mark3labs/mcp-go/server"
  usecases "github.com/landmaster135/devbox/internal/{tool-name}/usecases"
)

// ハンドラー関数
func handleOperation1(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
  param1, err := request.RequireString("param1")
  if err != nil {
    return nil, err
  }

  param2 := request.GetString("param2", "default")

  service := usecases.NewService()
  result, err := service.HandleOperation1(param1, param2)
  if err != nil {
    return nil, fmt.Errorf("操作1の実行に失敗しました: %v", err)
  }

  return mcp.NewToolResultText(result), nil
}

// サーバー設定
func setToolServer(s *server.MCPServer) *server.MCPServer {
  tool := mcp.NewTool(
    "tool_name",
    mcp.WithDescription("ツールの説明"),
    mcp.WithString("param1", mcp.Required(), mcp.Description("必須パラメータ")),
    mcp.WithString("param2", mcp.Description("オプションパラメータ")),
  )
  s.AddTool(tool, handleOperation1)
  return s
}

func createToolServer() *server.MCPServer {
  s := server.NewMCPServer(
    "Tool Name",
    "1.0.0",
    server.WithResourceCapabilities(true, true),
    server.WithPromptCapabilities(true),
    server.WithLogging(),
  )
  return setToolServer(s)
}

func BuildToolServer() {
  s := createToolServer()
  if err := server.ServeStdio(s); err != nil {
    fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
  }
}
```

共通のサービス層実装やコマンド実行ラッパーは `docs/tool_implementation/implementation_guide.md` を参照してください。

## 実装時のチェックリスト

### ハンドラー関数

- [ ] `mcp.NewToolResultText(result)` で結果を返却
- [ ] エラー時は `return nil, fmt.Errorf(...)` でエラーを返却
- [ ] 標準出力への出力を一切行わない（タイムアウト対策）
- [ ] `request.RequireString()` と `request.GetString()` を適切に使い分け
- [ ] 必須パラメータは `request.RequireString()` でエラーハンドリング
- [ ] オプションパラメータは `request.GetString()` でデフォルト値設定

### サーバー設定

- [ ] `mcp.WithDescription()` でツールの説明を設定
- [ ] 必須パラメータに `mcp.Required()` を設定
- [ ] 各パラメータに `mcp.Description()` で説明を設定
- [ ] `s.AddTool(tool, handler)` でツールとハンドラーを関連付け
- [ ] `s.AddPrompt(prompt, handler)` でプロンプトとハンドラーを関連付け
- [ ] `mcp.WithPromptDescription()` でプロンプトの説明を設定
- [ ] `server.WithPromptCapabilities(true)` でプロンプト機能を有効化
- [ ] `server.WithLogging()` でログ機能を有効化
- [ ] `cmd/mcp/router.go` にサーバーを追加

## よくある実装ミス

### 1. MCPツールで標準出力を使用

```go
// ❌ 間違い: MCPツールで標準出力を使用
fmt.Printf("処理中...\n")
fmt.Print(result)

// ✅ 正しい: MCPクライアントに結果を返却
return mcp.NewToolResultText(result), nil
```

### 2. 必須パラメータ取得の間違い

```go
// ❌ 間違い: 必須パラメータでGetStringを使用
param := request.GetString("required_param", "")

// ✅ 正しい: 必須パラメータはRequireStringを使用
param, err := request.RequireString("required_param")
if err != nil {
  return nil, err
}
```

### 3. 任意パラメータ取得の間違い

```go
// ❌ 間違い: GetStringでデフォルト値を設定していない
param := request.GetString("required_param")

// ✅ 正しい: GetStringでデフォルト値を設定する
param := request.GetString("required_param", "")
```

## MCPツールのテスト

1. `.config/cline/cline_mcp_settings.json` に設定追加
2. Cline から実行してテスト
3. エラーログを確認

## まとめ

- 結果は `mcp.NewToolResultText(...)` で返却し、標準出力は使わない
- 必須/任意パラメータで `RequireString` と `GetString` を使い分ける
- チェックリストに沿ってサーバー設定とルーティング登録漏れを防ぐ
