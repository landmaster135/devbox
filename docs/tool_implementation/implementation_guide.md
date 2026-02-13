# ツール実装ガイド

## 実装時の重要な注意点

### ツール種別ごとの注意点
- CLIツール固有の注意点: `docs/tool_implementation/cli/guide.md`
- MCPツール固有の注意点: `docs/tool_implementation/mcp/guide.md`

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

### exec.ExitErrorのExitCode()によるエラーの区別
- `err.(*exec.ExitError)` で型アサーション
- `exitError.ExitCode()` で終了コードを取得し、終了コードに応じて分岐
- **exit status 1**: テスト失敗 → 正常な動作として扱い、出力を返却
- **exit status 2以上**: システムエラー → エラーとして処理
- **その他のエラー**: コマンド未発見等 → エラーとして処理

例えば、Goのテストコマンドでは、テスト失敗とシステムエラーを以下のように区別します。

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

### 1. 共通サービス層の実装

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

### 2. コマンド実行ラッパーによる非シェル化

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

CLIツール固有のチェックリストは `docs/tool_implementation/cli/guide.md` を参照してください。
MCPツール固有のチェックリストは `docs/tool_implementation/mcp/guide.md` を参照してください。

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

CLIツール固有の実装ミスは `docs/tool_implementation/cli/guide.md` を参照してください。
MCPツール固有の実装ミスは `docs/tool_implementation/mcp/guide.md` を参照してください。

## ツール別テスト手順

- CLIツール: `docs/tool_implementation/cli/guide.md`
- MCPツール: `docs/tool_implementation/mcp/guide.md`

## まとめ

- **共通**: ビジネスロジックはusecasesパッケージで共有
- **CLI固有項目**: `docs/tool_implementation/cli/guide.md` を参照
- **MCP固有項目**: `docs/tool_implementation/mcp/guide.md` を参照
