# ツール実装ガイド

## ドキュメント更新ルール

CLIツール実装完了後のドキュメント更新手順は、`docs/tool_implementation/documentation_guide.md` を参照してください。
CLIツール実装完了後は、対応する `cmd/cli/<tool>/README.md` に使い方と実行例を必ず追記してください。

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

### 4. コマンド実行ラッパーによる非シェル化

外部入力を伴う処理で OS コマンドを実行する場合、`exec.Command` を直接 `name` と `args` に分割して呼び出すことで、シェル展開を回避しコマンドインジェクション対策を実装しなくても安全な構造を保てます。

```go
// usecases/command_executor.go
type CommandExecutor struct{}

func (e *CommandExecutor) Execute(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func (e *CommandExecutor) ExecuteInDir(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// 呼び出し側（例）
args := []string{"run", executionFile}
args = append(args, strings.Fields(parameters)...)
output, err := executor.ExecuteInDir(absRootDir, "go", args...)

// ポイント
// - シェル (`sh -c` / `bash -c`) を介さないので、`;` や `&&` を含む入力でも追加コマンド化されない
// - 入力は引数スライスとして渡され、Go ランタイムがエスケープを扱う
// - 出力は呼び出し元で検査・加工し、副作用が必要な場合のみ実施
```

このパターンを守れば、追加のエスケープ処理を実装せずにコマンドインジェクションを防げます。必要に応じて、許可するコマンド名や引数のホワイトリスト検証を組み合わせるとより堅牢になります。

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

## テストコードでのハードコード削減

### 1. テスト関数内での定数定義

テストでしか使わない定数や変数は、個々のテスト関数内で定義します。

```go
func TestCreateDashboardForCloudRun_WithMock(t *testing.T) {
  // テスト固有の定数を関数内で定義
  const (
    testProject     = "test-project"
    testLocation    = "us-central1"
    testServiceName = "test-service"
    expectedResult  = "ダッシュボードが正常に作成されました: projects/test-project/dashboards/test-dashboard-123"
  )
  
  tests := []struct {
    name               string
    setupCloudRunMock  func(*MockCloudRunClient)
    setupDashboardMock func(*MockDashboardClient)
    expectError        bool
    errorMessage       string
    expectedResult     string
  }{
    // テストケース定義
  }
  // ...
}
```

### 2. テーブル駆動テスト

複数のテストケースを構造体の配列で管理し、ループで実行します。

```go
func TestVerifyCloudRunService_WithMock(t *testing.T) {
  tests := []struct {
    name           string
    setupMock      func(*MockCloudRunClient)
    expectedExists bool
    expectError    bool
    errorMessage   string
  }{
    {
      name: "ServiceExists_Normal",
      setupMock: func(mock *MockCloudRunClient) {
        mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
          return &runpb.Service{Name: "projects/test-project/locations/us-central1/services/test-service"}, nil
        }
      },
      expectedExists: true,
      expectError:    false,
    },
    {
      name: "ServiceNotFound_Normal",
      setupMock: func(mock *MockCloudRunClient) {
        mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
          return nil, status.Error(codes.NotFound, "service not found")
        }
      },
      expectedExists: false,
      expectError:    false,
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      mockCloudRunClient := &MockCloudRunClient{}
      tt.setupMock(mockCloudRunClient)
      // テスト実行
    })
  }
}
```

### 3. テストヘルパー関数

共通の初期化処理やモック生成をヘルパー関数として定義します。

```go
// セットアップヘルパー関数
func setupTestService(t *testing.T, cloudRunClient CloudRunClient, dashboardClient DashboardClient) *Service {
    return NewServiceWithClients("test-project", "us-central1", "test-service", "", cloudRunClient, dashboardClient)
}

// モック生成ヘルパー関数
func createSuccessfulCloudRunMock() *MockCloudRunClient {
  mock := &MockCloudRunClient{}
  mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
    return &runpb.Service{Name: req.Name}, nil
  }
  return mock
}
```

### 4. テストデータ構造体

複雑なテストデータは構造体で管理します。

```go
type testCase struct {
  name           string
  input          inputData
  expected       expectedData
  expectError    bool
  errorMessage   string
}

type inputData struct {
  project     string
  location    string
  serviceName string
}

type expectedData struct {
  dashboardName string
  resultMessage string
}
```

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

### 5. フラグパーサーのモック実装の間違い
```go
// ❌ 間違い: フラグ定義後に値を設定しても反映されない
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  *p = value // デフォルト値のみ
}

// ✅ 正しい: 事前設定値をチェックして適用
func (m *MockFlagParser) StringVar(p *string, name string, value string, usage string) {
  if presetValue, exists := m.stringValues[name]; exists {
    *p = presetValue // 事前設定値を優先
  } else {
    *p = value // デフォルト値
  }
  m.stringVars[name] = p
}
```

**参考実装**: `devbox/internal/zip_compressor/config/config_test.go`のMockFlagParserを参照

## MCPツールのテスト
1. `.config/cline/cline_mcp_settings.json`に設定追加
2. Clineから実行してテスト
3. エラーログを確認

## まとめ

- **CLIツール**: 結果を`fmt.Print(result)`で標準出力に表示
- **MCPツール**: 結果を`mcp.NewToolResultText(result)`でクライアントに返却
- **共通**: ビジネスロジックはusecasesパッケージで共有
- **重要**: MCPツールでは標準出力を一切使用しない（タイムアウト対策）
