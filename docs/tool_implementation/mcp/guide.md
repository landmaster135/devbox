# MCPツール実装ガイド

このドキュメントは、MCPツール実装に関する内容をまとめたものです。

## 実装時のチェックリスト

- [ ] ハンドラーのエラーは `return nil, fmt.Errorf(...)` で文脈付きに返却する
- [ ] 各パラメータに `mcp.WithDescription()` でツールの説明を設定する
- [ ] `s.AddTool(tool, handler)` でツールとハンドラーを関連付け
- [ ] `s.AddPrompt(prompt, handler)` でプロンプトとハンドラーを関連付け
- [ ] `mcp.WithPromptDescription()` でプロンプトの説明を設定
- [ ] `server.WithPromptCapabilities(true)` でプロンプト機能を有効化
- [ ] `server.WithLogging()` で診断用ログを有効化する

## 実装パターン

### ルーター層

MCPサーバーシステムは `cmd/mcp/router.go` を中心に、起動引数でサーバーを切り替える統一ルーティングを採用しています。

```go
func Router() {
	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run main.go [arguments]")
		os.Exit(1)
	}

	switch args[1] {
	case "arith_calc":
		arithmetic_calculator.BuildArithCalculatorServer()
	case "github":
		github.BuildGitHubServer()
	default:
		fmt.Fprintln(os.Stderr, "argument is invalid")
		os.Exit(1)
	}
}
```

### MCPハンドラー層（`cmd/mcp/arithmetic_calculator/mcp.go` 参考）

ハンドラーでは `request.Require*` で必須入力を取得し、usecase 呼び出し結果を `mcp.CallToolResult` として返却します。

```go
import (
	"context"
	"fmt"

	mcp "github.com/mark3labs/mcp-go/mcp"
	usecases "github.com/landmaster135/devbox/internal/arithmetic_calculator/usecases"
)

func handleToCalculate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	op, err := request.RequireString("operation")
	if err != nil {
		return nil, err
	}

	x, err := request.RequireFloat("x")
	if err != nil {
		return nil, err
	}

	y, err := request.RequireFloat("y")
	if err != nil {
		return nil, err
	}

	service := usecases.NewCalculatorService()
	result, err := service.HandleToCalculate(op, x, y)
	if err != nil {
		return nil, fmt.Errorf("パラメータを用いた算術計算に失敗しました: %v", err)
	}

	return mcp.FormatNumberResult(result), nil
}
```

### ツール定義とハンドラー関連付け

```go
func setTwoNumbersInputtingCalcServer(s *server.MCPServer) *server.MCPServer {
	tool := mcp.NewTool(
		"calculate",
		mcp.WithDescription("Perform basic arithmetic calculations with two numbers"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The arithmetic operation to perform"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x", mcp.Required(), mcp.Description("First number")),
		mcp.WithNumber("y", mcp.Required(), mcp.Description("Second number")),
	)
	s.AddTool(tool, handleToCalculate)
	return s
}
```

## 実装時の注意点

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

共通のサービス層実装やコマンド実行ラッパーは `docs/tool_implementation/implementation_guide.md` を参照してください。

## 実装アンチパターン

### MCPツールで標準出力を使用

```go
// ❌ 間違い: MCPツールで標準出力を使用
fmt.Printf("処理中...\n")
fmt.Print(result)

// ✅ 正しい: MCPクライアントに結果を返却
return mcp.NewToolResultText(result), nil
```

### 必須パラメータ取得の間違い

```go
// ❌ 間違い: 必須パラメータでGetStringを使用
param := request.GetString("required_param", "")

// ✅ 正しい: 必須パラメータはRequireStringを使用
param, err := request.RequireString("required_param")
if err != nil {
  return nil, err
}
```

### 任意パラメータ取得の間違い

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
- 実装手順に沿ってサーバー設定とルーティング登録漏れを防ぐ
