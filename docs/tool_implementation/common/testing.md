# テストコード実装ガイド

## テスト戦略
- **単体テスト**: 各レイヤーの独立したテスト
- **統合テスト**: レイヤー間の連携テスト
- **CLI観点**: フラグ解析、標準出力、標準エラー出力、終了コード
- **MCP観点**: ツール定義、必須/任意パラメータ取得、`CallToolResult` 返却
- **カバレッジ目標**: 90%以上
- **テストコマンド**: `go test -v ./... -coverpkg=./... -covermode=count -coverprofile=coverage.out`
- **テストツール**: `github.com/stretchr/testify`, `github.com/DATA-DOG/go-sqlmock`

## テストコード実装チェックリスト

- [ ] 単体テストと統合テストを追加
- [ ] テスト固有の定数は各テスト関数内で定義
- [ ] 複数のテストケースはテーブル駆動テストで実装
- [ ] 共通処理はヘルパー関数として抽出
- [ ] 複雑なテストデータは構造体で管理
- [ ] モック生成処理は関数化
- [ ] `infrastructures` 層の `Repository` モックは `infrastructures/{resource}/` 直下に配置し、usecases 側で再定義しない
- [ ] テスト名は「機能_条件」の形式で命名
- [ ] operation実装や共通関数をサブディレクトリに分離した場合、対応テストを同ディレクトリの `*_test.go` へ移設した
- [ ] テストカバレッジが実装前より下がっていない

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

## 付録

参考実装: `devbox/internal/gcloud_monitoring/usecases/services_test.go`
