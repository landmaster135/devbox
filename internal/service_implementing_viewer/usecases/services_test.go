package usecases

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestServiceImplementingViewerService はServiceImplementingViewerServiceのテストクラス
type TestServiceImplementingViewerService struct {
	service *ServiceImplementingViewerService
	tempDir string
}

// setupTestEnvironment はテスト環境をセットアップする
func (t *TestServiceImplementingViewerService) setupTestEnvironment() error {
	// 一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "service_test")
	if err != nil {
		return err
	}
	t.tempDir = tempDir

	// テスト用のディレクトリ構造を作成
	// cli ディレクトリ
	cliDir := filepath.Join(tempDir, "cli")
	if err := os.MkdirAll(cliDir, 0755); err != nil {
		return err
	}

	// mcp ディレクトリ
	mcpDir := filepath.Join(tempDir, "mcp")
	if err := os.MkdirAll(mcpDir, 0755); err != nil {
		return err
	}

	// grpc/handlers ディレクトリ
	grpcHandlersDir := filepath.Join(tempDir, "grpc", "handlers")
	if err := os.MkdirAll(grpcHandlersDir, 0755); err != nil {
		return err
	}

	// http/handlers ディレクトリ
	httpHandlersDir := filepath.Join(tempDir, "http", "handlers")
	if err := os.MkdirAll(httpHandlersDir, 0755); err != nil {
		return err
	}

	// テスト用サービスディレクトリを作成
	testServices := []struct {
		dir  string
		name string
	}{
		// 既存のサービス
		{"cli", "arithmetic_calculator"},
		{"cli", "base64-extractor"},
		{"cli", "git-commit-history-retriever"},
		{"cli", "weather-notificator"},
		{"mcp", "arithmetic-calculator"},
		{"mcp", "git_commit_history_retriever"},
		{"mcp", "github"},
		{"mcp", "weather-notificator"},
		{"grpc/handlers", "weather-notificator"},
		{"grpc/handlers", "user-management"},
		{"grpc/handlers", "notification-service"},
		{"http/handlers", "weather-notificator"},
		{"http/handlers", "api-gateway"},
	}

	for _, service := range testServices {
		serviceDir := filepath.Join(tempDir, service.dir, service.name)
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			return err
		}
	}

	return nil
}

// cleanupTestEnvironment はテスト環境をクリーンアップする
func (t *TestServiceImplementingViewerService) cleanupTestEnvironment() {
	if t.tempDir != "" {
		os.RemoveAll(t.tempDir)
	}
}

// TestNewServiceImplementingViewerService_Normal はNewServiceImplementingViewerServiceの正常系テスト
func TestNewServiceImplementingViewerService_Normal(t *testing.T) {
	// Arrange
	rootDir := "/test/root"
	targetDirs := []string{"cli", "mcp", "grpc/handlers", "http/handlers"}

	// Act
	service := NewServiceImplementingViewerService(rootDir, targetDirs)

	// Assert
	if service == nil {
		t.Fatal("サービスがnilです")
	}
	if service.rootDir != rootDir {
		t.Errorf("rootDirが期待値と異なります。期待値: %s, 実際: %s", rootDir, service.rootDir)
	}
	if len(service.targetDirs) != len(targetDirs) {
		t.Errorf("targetDirsの長さが期待値と異なります。期待値: %d, 実際: %d", len(targetDirs), len(service.targetDirs))
	}
	for i, expected := range targetDirs {
		if service.targetDirs[i] != expected {
			t.Errorf("targetDirs[%d]が期待値と異なります。期待値: %s, 実際: %s", i, expected, service.targetDirs[i])
		}
	}
}

// TestGetServiceImplementingStatus_Normal はGetServiceImplementingStatusの正常系テスト
func TestGetServiceImplementingStatus_Normal(t *testing.T) {
	// Arrange
	testService := &TestServiceImplementingViewerService{}
	if err := testService.setupTestEnvironment(); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}
	defer testService.cleanupTestEnvironment()

	targetDirs := []string{"cli", "mcp", "grpc/handlers", "http/handlers"}
	service := NewServiceImplementingViewerService(testService.tempDir, targetDirs)

	// Act
	result, statistics, err := service.GetServiceImplementingStatus()

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if statistics == nil {
		t.Fatal("統計情報がnilです")
	}
	if result == "" {
		t.Fatal("結果が空です")
	}

	// 統計情報の検証
	if statistics.TotalServices != 8 {
		t.Errorf("総サービス数が期待値と異なります。期待値: 8, 実際: %d", statistics.TotalServices)
	}
	if statistics.CLICount != 4 {
		t.Errorf("CLICount が期待値と異なります。期待値: 4, 実際: %d", statistics.CLICount)
	}
	if statistics.MCPCount != 4 {
		t.Errorf("MCPCount が期待値と異なります。期待値: 4, 実際: %d", statistics.MCPCount)
	}
	if statistics.GRPCCount != 3 {
		t.Errorf("GRPCCount が期待値と異なります。期待値: 3, 実際: %d", statistics.GRPCCount)
	}
	if statistics.HTTPCount != 2 {
		t.Errorf("HTTPCount が期待値と異なります。期待値: 2, 実際: %d", statistics.HTTPCount)
	}
	if statistics.AllImplementedCount != 1 {
		t.Errorf("AllImplementedCount が期待値と異なります。期待値: 1, 実際: %d", statistics.AllImplementedCount)
	}

	// 結果の内容を検証
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Fatalf("結果の行数が不足しています。期待値: 3行以上, 実際: %d行", len(lines))
	}

	// ヘッダー行の検証
	if !strings.Contains(lines[0], "service") || !strings.Contains(lines[0], "cli") || !strings.Contains(lines[0], "mcp") || !strings.Contains(lines[0], "grpc/handlers") || !strings.Contains(lines[0], "http/handlers") {
		t.Errorf("ヘッダー行が正しくありません: %s", lines[0])
	}

	// セパレーター行の検証
	if !strings.Contains(lines[1], ":") || !strings.Contains(lines[1], "-") {
		t.Errorf("セパレーター行が正しくありません: %s", lines[1])
	}

	// データ行の検証（weather-notificatorが全ディレクトリに存在することをチェック）
	found := false
	for _, line := range lines[2:] {
		if strings.Contains(line, "weather-notificator") {
			found = true
			// weather-notificatorは全ディレクトリに存在するので4つの✅があるはず
			checkCount := strings.Count(line, "✅")
			if checkCount != 4 {
				t.Errorf("weather-notificatorの行の✅の数が期待値と異なります。期待値: 4, 実際: %d, 行: %s", checkCount, line)
			}
			break
		}
	}
	if !found {
		t.Error("weather-notificatorの行が見つかりません")
	}
}

// TestGetServicesInDirectory_Normal はgetServicesInDirectoryの正常系テスト
func TestGetServicesInDirectory_Normal(t *testing.T) {
	// Arrange
	testService := &TestServiceImplementingViewerService{}
	if err := testService.setupTestEnvironment(); err != nil {
		t.Fatalf("テスト環境のセットアップに失敗しました: %v", err)
	}
	defer testService.cleanupTestEnvironment()

	service := NewServiceImplementingViewerService(testService.tempDir, []string{"cli"})
	cliDir := filepath.Join(testService.tempDir, "cli")

	// Act
	services, err := service.getServicesInDirectory(cliDir)

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if len(services) == 0 {
		t.Fatal("サービスが見つかりません")
	}

	expectedServices := []string{"arithmetic_calculator", "base64-extractor", "git-commit-history-retriever", "weather-notificator"}
	if len(services) != len(expectedServices) {
		t.Errorf("サービス数が期待値と異なります。期待値: %d, 実際: %d", len(expectedServices), len(services))
	}

	for _, expected := range expectedServices {
		found := false
		for _, actual := range services {
			if actual == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期待されるサービス %s が見つかりません", expected)
		}
	}
}

// TestGetServicesInDirectory_NonExistentDirectory は存在しないディレクトリのテスト
func TestGetServicesInDirectory_NonExistentDirectory(t *testing.T) {
	// Arrange
	service := NewServiceImplementingViewerService("/non/existent", []string{"cli"})
	nonExistentDir := "/non/existent/directory"

	// Act
	services, err := service.getServicesInDirectory(nonExistentDir)

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if len(services) != 0 {
		t.Errorf("存在しないディレクトリでサービスが見つかりました: %v", services)
	}
}

// TestNormalizeServiceName_Normal はnormalizeServiceNameの正常系テスト
func TestNormalizeServiceName_Normal(t *testing.T) {
	// Arrange
	service := NewServiceImplementingViewerService("/test", []string{"cli"})
	testCases := []struct {
		input    string
		expected string
	}{
		{"arithmetic_calculator", "arithmetic-calculator"},
		{"git-commit-history-retriever", "git-commit-history-retriever"},
		{"base64_extractor", "base64-extractor"},
		{"test_service_name", "test-service-name"},
		{"already-normalized", "already-normalized"},
	}

	for _, testCase := range testCases {
		// Act
		result := service.normalizeServiceName(testCase.input)

		// Assert
		if result != testCase.expected {
			t.Errorf("正規化結果が期待値と異なります。入力: %s, 期待値: %s, 実際: %s", testCase.input, testCase.expected, result)
		}
	}
}

// TestIsServiceImplementedInDirectory_Normal はisServiceImplementedInDirectoryの正常系テスト
func TestIsServiceImplementedInDirectory_Normal(t *testing.T) {
	// Arrange
	service := NewServiceImplementingViewerService("/test", []string{"cli"})
	servicesInDir := []string{"arithmetic_calculator", "base64-extractor", "git-commit-history-retriever"}

	testCases := []struct {
		serviceName string
		expected    bool
	}{
		{"arithmetic-calculator", true},
		{"base64-extractor", true},
		{"git-commit-history-retriever", true},
		{"non-existent-service", false},
		{"github", false},
	}

	for _, testCase := range testCases {
		// Act
		result := service.isServiceImplementedInDirectory(testCase.serviceName, servicesInDir)

		// Assert
		if result != testCase.expected {
			t.Errorf("実装チェック結果が期待値と異なります。サービス名: %s, 期待値: %t, 実際: %t", testCase.serviceName, testCase.expected, result)
		}
	}
}

// TestFormatAsTable_Normal はformatAsTableの正常系テスト
func TestFormatAsTable_Normal(t *testing.T) {
	// Arrange
	service := NewServiceImplementingViewerService("/test", []string{"cli", "mcp", "grpc/handlers", "http/handlers"})
	serviceStatuses := []ServiceStatus{
		{
			ServiceName: "weather-notificator",
			Directories: map[string]bool{"cli": true, "mcp": true, "grpc/handlers": true, "http/handlers": true},
		},
		{
			ServiceName: "arithmetic-calculator",
			Directories: map[string]bool{"cli": true, "mcp": true, "grpc/handlers": false, "http/handlers": false},
		},
		{
			ServiceName: "user-management",
			Directories: map[string]bool{"cli": false, "mcp": false, "grpc/handlers": true, "http/handlers": false},
		},
	}

	// Act
	result := service.formatAsTable(serviceStatuses)

	// Assert
	if result == "" {
		t.Fatal("結果が空です")
	}

	lines := strings.Split(result, "\n")
	if len(lines) < 5 {
		t.Fatalf("結果の行数が不足しています。期待値: 5行以上, 実際: %d行", len(lines))
	}

	// ヘッダー行の検証
	if !strings.Contains(lines[0], "service") || !strings.Contains(lines[0], "cli") || !strings.Contains(lines[0], "mcp") || !strings.Contains(lines[0], "grpc/handlers") || !strings.Contains(lines[0], "http/handlers") {
		t.Errorf("ヘッダー行が正しくありません: %s", lines[0])
	}

	// データ行の検証
	weatherNotificatorFound := false
	arithmeticCalculatorFound := false
	userManagementFound := false

	for _, line := range lines[2:] {
		if strings.Contains(line, "weather-notificator") {
			weatherNotificatorFound = true
			// weather-notificatorは全ディレクトリに存在するので4つの✅があるはず
			checkCount := strings.Count(line, "✅")
			if checkCount != 4 {
				t.Errorf("weather-notificatorの行の✅の数が期待値と異なります。期待値: 4, 実際: %d, 行: %s", checkCount, line)
			}
		} else if strings.Contains(line, "arithmetic-calculator") {
			arithmeticCalculatorFound = true
			// arithmetic-calculatorはcliとmcpのみに存在するので2つの✅があるはず
			checkCount := strings.Count(line, "✅")
			if checkCount != 2 {
				t.Errorf("arithmetic-calculatorの行の✅の数が期待値と異なります。期待値: 2, 実際: %d, 行: %s", checkCount, line)
			}
		} else if strings.Contains(line, "user-management") {
			userManagementFound = true
			// user-managementはgrpc/handlersのみに存在するので1つの✅があるはず
			checkCount := strings.Count(line, "✅")
			if checkCount != 1 {
				t.Errorf("user-managementの行の✅の数が期待値と異なります。期待値: 1, 実際: %d, 行: %s", checkCount, line)
			}
		}
	}

	if !weatherNotificatorFound {
		t.Error("weather-notificatorの行が見つかりません")
	}
	if !arithmeticCalculatorFound {
		t.Error("arithmetic-calculatorの行が見つかりません")
	}
	if !userManagementFound {
		t.Error("user-managementの行が見つかりません")
	}

	// 絵文字の検証
	if !strings.Contains(result, "✅") || !strings.Contains(result, "❌️") {
		t.Error("結果に適切な絵文字が含まれていません")
	}
}
