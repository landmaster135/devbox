# ツール実装ガイド

## 実装時の重要な注意点

### CLIツール vs MCPツールの出力制御

**重要**: CLIツールとMCPツールでは出力の扱いが根本的に異なります。

CLIツール
```go
// ✅ 正しい実装: 結果を標準出力に表示
func handleOperation(cfg *config.Config) {
    service := usecases.NewService()
    result, err := service.HandleOperation(cfg.Param)
    if err != nil {
        fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
        os.Exit(1)
    }
    
    // 必須: 結果を標準出力に表示
    fmt.Print(result)
}
```

MCPツール
```go
// ✅ 正しい実装: MCPクライアントに結果を返却
func handleOperation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    service := usecases.NewService()
    result, err := service.HandleOperation(param)
    if err != nil {
        return nil, fmt.Errorf("処理に失敗しました: %v", err)
    }
    
    // 必須: MCPクライアントに結果を返却（標準出力は使用しない）
    return mcp.NewToolResultText(result), nil
}
```

### MCPツールでのタイムアウト対策

MCPツールでは長大な標準出力を行うとタイムアウトが発生します。以下を厳守してください：

```go
// ❌ 悪い例: 標準出力への出力
fmt.Printf("処理中...\n")
fmt.Printf("実行コマンド: %s\n", cmd)

// ✅ 良い例: 標準出力を抑制
// 進捗表示や実行コマンドの出力は行わない
```

### ディレクトリ操作の注意点
- `cmd.Dir = dir` の設定は必須
- 相対パスは `filepath.Abs()` で絶対パスに変換
- ディレクトリの存在確認と種別確認を実装

実装例:
```go
// ✅ 正しい実装: ExecuteInDirメソッドでcmd.Dirを設定
func (e *DefaultCommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
    cmd := exec.Command(name, args...)
    cmd.Dir = dir  // 重要: 実行ディレクトリを設定
    return cmd.CombinedOutput()
}

// ✅ 正しい実装: 絶対パス変換とディレクトリ確認
func (s *Service) ExecuteOperation(directory string) (string, error) {
    // ディレクトリの存在確認
    if !s.directoryChecker.Exists(directory) {
        return "", fmt.Errorf("指定されたディレクトリが存在しません: %s", directory)
    }
    if !s.directoryChecker.IsDirectory(directory) {
        return "", fmt.Errorf("指定されたパスはディレクトリではありません: %s", directory)
    }

    // 絶対パスに変換（重要）
    absDir, err := filepath.Abs(directory)
    if err != nil {
        return "", fmt.Errorf("ディレクトリパスの変換に失敗しました: %v", err)
    }

    // 指定されたディレクトリでコマンド実行
    output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "./...")
    // ...
}
```

### exec.ExitErrorのExitCode()による分岐処理
- `err.(*exec.ExitError)` で型アサーション
- `exitError.ExitCode()` で終了コードを取得
- 終了コードに応じた適切な処理分岐

実装例:
```go
// ✅ 正しい実装: ExitErrorの型アサーションと終了コード判定
output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "-cover", "./...")
if err != nil {
    // エラーの種類を判定
    if exitError, ok := err.(*exec.ExitError); ok {
        // exec.ExitErrorの場合は終了コードを確認
        if exitError.ExitCode() == 1 {
            // exit status 1: テスト失敗（正常な動作として扱う）
            // 何もしない - 出力を返却する
        } else {
            // exit status 1以外: 実際のエラー
            return "", fmt.Errorf("コマンド実行でエラーが発生しました: %v\n出力: %s", err, string(output))
        }
    } else {
        // exec.ExitError以外のエラー（コマンドが見つからない等）
        return "", fmt.Errorf("コマンドの実行に失敗しました: %v", err)
    }
}
```

### テスト失敗（exit status 1）と実際のエラーの区別
- **exit status 1**: テスト失敗 → 正常な動作として扱い、出力を返却
- **exit status 2以上**: システムエラー → エラーとして処理
- **その他のエラー**: コマンド未発見等 → エラーとして処理

例えば、Goのテストコマンドでは、テスト失敗とシステムエラーを区別する必要があります：

実装例:
```go
// ✅ 正しい実装: テスト失敗を正常な動作として扱う
func (s *GolangOpsService) ExecuteTestCoverage(directory, grepPattern string) (string, error) {
    // ... ディレクトリ確認等の処理 ...

    // go test -cover ./... を実行
    output, err := s.commandExecutor.ExecuteInDir(absDir, "go", "test", "-cover", "./...")
    if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            if exitError.ExitCode() == 1 {
                // テスト失敗は正常な動作として扱う
                // エラーを返さず、出力をそのまま処理する
            } else {
                // exit code が1以外の場合は実際のエラー
                return "", fmt.Errorf("テストカバレッジの実行でエラーが発生しました: %v\n出力: %s", err, string(output))
            }
        } else {
            // exec.ExitError以外のエラー
            return "", fmt.Errorf("コマンドの実行に失敗しました: %v", err)
        }
    }

    // テスト失敗の場合でも出力を処理して返却
    result.Write(output)
    return result.String(), nil
}
```

テスト実装での注意点:
```go
// テスト用の実際のExitErrorを作成
cmd := exec.Command("sh", "-c", "exit 1")
err := cmd.Run()
exitError, _ := err.(*exec.ExitError)

// モックでExitErrorを返却
mockCommandExecutor.On("ExecuteInDir", mock.AnythingOfType("string"), "go", []string{"test", "-cover", "./..."}).Return(expectedOutput, exitError)
```

## 実装パターン

### 1. CLIツールのmain.go構造

```go
package main

import (
	"fmt"
	"os"
	config "github.com/landmaster135/devbox/internal/{tool-name}/config"
	usecases "github.com/landmaster135/devbox/internal/{tool-name}/usecases"
)

func main() {
	cfg, err := config.ParseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		config.PrintUsage()
		os.Exit(1)
	}

	if cfg.Help {
		config.PrintUsage()
		return
	}

	switch cfg.Operation {
	case "operation1":
		handleOperation1(cfg)
	case "operation2":
		handleOperation2(cfg)
	default:
		fmt.Fprintf(os.Stderr, "エラー: 未対応の操作タイプです: %s\n", cfg.Operation)
		config.PrintUsage()
		os.Exit(1)
	}
}

func handleOperation1(cfg *config.Config) {
	service := usecases.NewService()
	result, err := service.HandleOperation1(cfg.Param1, cfg.Param2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	fmt.Print(result) // 重要: 結果を標準出力に表示
}
```

### 2. MCPツールのmcp.go構造

```go
package tool_name

import (
	"context"
	"fmt"
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
		fmt.Printf("Server error: %v\n", err)
	}
}
```

### 3. 共通サービス層の実装

```go
package usecases

type Service struct {
	// 依存関係の注入
}

func NewService() *Service {
	return &Service{}
}

// CLIツールとMCPツールで共通のビジネスロジック
func (s *Service) HandleOperation1(param1, param2 string) (string, error) {
	// ビジネスロジックの実装
	// 戻り値は文字列で統一（CLIツールで標準出力、MCPツールでクライアントに返却）
	
	result := fmt.Sprintf("処理結果: %s, %s", param1, param2)
	return result, nil
}
```

## 実装時のチェックリスト

### CLIツール実装チェックリスト

- [ ] `fmt.Print(result)`で結果を標準出力に表示している
- [ ] エラー時は`fmt.Fprintf(os.Stderr, ...)`でエラーメッセージを出力
- [ ] エラー時は`os.Exit(1)`でプロセスを終了
- [ ] ヘルプ機能（`-help`フラグ）を実装
- [ ] 操作タイプのswitch文で未対応操作をハンドリング

### MCPツール実装チェックリスト

**ハンドラー関数**
- [ ] `mcp.NewToolResultText(result)`で結果を返却
- [ ] エラー時は`return nil, fmt.Errorf(...)`でエラーを返却
- [ ] 標準出力への出力を一切行わない（タイムアウト対策）
- [ ] `request.RequireString()`と`request.GetString()`を適切に使い分け
- [ ] 必須パラメータは`request.RequireString()`でエラーハンドリング
- [ ] オプションパラメータは`request.GetString()`でデフォルト値設定

**サーバー設定**
- [ ] `mcp.WithDescription()`でツールの説明を設定
- [ ] 必須パラメータは`mcp.Required()`を設定
- [ ] 各パラメータに`mcp.Description()`で説明を設定
- [ ] `s.AddTool(tool, handler)`でツールとハンドラーを関連付け
- [ ] `s.AddPrompt(prompt, handler)`でプロンプトとハンドラーを関連付け
- [ ] `mcp.WithPromptDescription()`でプロンプトの説明を設定
- [ ] `server.WithPromptCapabilities(true)`でプロンプト機能を有効化
- [ ] `server.WithLogging()`でログ機能を有効化
- [ ] `cmd/mcp/router.go`にサーバーを追加

### 共通チェックリスト

- [ ] ビジネスロジックはusecasesパッケージに実装
- [ ] エラーメッセージは日本語で分かりやすく
- [ ] パラメータの妥当性検証を実装
- [ ] テストコードを作成

## 具体的な実装例の比較

### arithmetic-calculator の実装例
- **CLIツール** (`cmd/cli/arithmetic-calculator/main.go`)
- **MCPツール** (`cmd/mcp/arithmetic_calculator/mcp.go`)

### ops-for-golang の実装例
- **CLIツール** (`cmd/cli/ops-for-golang/main.go`)
- **MCPツール** (`cmd/mcp/ops_for_golang/mcp.go`)

## よくある実装ミス

### 1. MCPツールで標準出力を使用
```go
// ❌ 間違い: MCPツールで標準出力を使用
fmt.Printf("処理中...\n")
fmt.Print(result)

// ✅ 正しい: MCPクライアントに結果を返却
return mcp.NewToolResultText(result), nil
```

### 2. CLIツールで結果を表示しない
```go
// ❌ 間違い: 結果を表示しない
service.HandleOperation(param)

// ✅ 正しい: 結果を標準出力に表示
result, err := service.HandleOperation(param)
if err != nil {
    fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
    os.Exit(1)
}
fmt.Print(result)
```

### 3. パラメータ取得の間違い
```go
// ❌ 間違い: 必須パラメータでGetStringを使用
param := request.GetString("required_param", "")

// ✅ 正しい: 必須パラメータはRequireStringを使用
param, err := request.RequireString("required_param")
if err != nil {
    return nil, err
}
```

### 4. 任意パラメータ取得の間違い
```go
// ❌ 間違い: GetStringでデフォルト値を設定していない
param := request.GetString("required_param")

// ✅ 正しい: GetStringでデフォルト値を設定する
param := request.GetString("required_param", "")
```

## MCPツールのテスト
1. `.config/cline/cline_mcp_settings.json`に設定追加
2. Clineから実行してテスト
3. エラーログを確認

## まとめ

- **CLIツール**: 結果を`fmt.Print(result)`で標準出力に表示
- **MCPツール**: 結果を`mcp.NewToolResultText(result)`でクライアントに返却
- **共通**: ビジネスロジックはusecasesパッケージで共有
- **重要**: MCPツールでは標準出力を一切使用しない（タイムアウト対策）
