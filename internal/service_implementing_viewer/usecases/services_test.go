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

	// テスト用サービスディレクトリを作成
	testServices := []struct {
		dir  string
		name string
	}{
		{"cli", "arithmetic_calculator"},
		{"cli", "base64-extractor"},
		{"cli", "git-commit-history-retriever"},
		{"mcp", "arithmetic-calculator"},
		{"mcp", "git_commit_history_retriever"},
		{"mcp", "github"},
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
	targetDirs := []string{"cli", "mcp"}

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

	targetDirs := []string{"cli", "mcp"}
	service := NewServiceImplementingViewerService(testService.tempDir, targetDirs)

	// Act
	result, err := service.GetServiceImplementingStatus()

	// Assert
	if err != nil {
		t.Fatalf("エラーが発生しました: %v", err)
	}
	if result == "" {
		t.Fatal("結果が空です")
	}

	// 結果の内容を検証
	lines := strings.Split(result, "\n")
	if len(lines) < 3 {
		t.Fatalf("結果の行数が不足しています。期待値: 3行以上, 実際: %d行", len(lines))
	}

	// ヘッダー行の検証
	if !strings.Contains(lines[0], "service") || !strings.Contains(lines[0], "cli") || !strings.Contains(lines[0], "mcp") {
		t.Errorf("ヘッダー行が正しくありません: %s", lines[0])
	}

	// セパレーター行の検証
	if !strings.Contains(lines[1], ":") || !strings.Contains(lines[1], "-") {
		t.Errorf("セパレーター行が正しくありません: %s", lines[1])
	}

	// データ行の検証（arithmetic-calculatorが含まれているかチェック）
	found := false
	for _, line := range lines[2:] {
		if strings.Contains(line, "arithmetic-calculator") {
			found = true
			if !strings.Contains(line, "✅") {
				t.Errorf("arithmetic-calculatorの行に✅が含まれていません: %s", line)
			}
			break
		}
	}
	if !found {
		t.Error("arithmetic-calculatorの行が見つかりません")
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

	expectedServices := []string{"arithmetic_calculator", "base64-extractor", "git-commit-history-retriever"}
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
	service := NewServiceImplementingViewerService("/test", []string{"cli", "mcp"})
	serviceStatuses := []ServiceStatus{
		{
			ServiceName: "arithmetic-calculator",
			Directories: map[string]bool{"cli": false, "mcp": true},
		},
		{
			ServiceName: "base64-extractor",
			Directories: map[string]bool{"cli": true, "mcp": false},
		},
	}

	// Act
	result := service.formatAsTable(serviceStatuses)

	// Assert
	if result == "" {
		t.Fatal("結果が空です")
	}

	lines := strings.Split(result, "\n")
	if len(lines) < 4 {
		t.Fatalf("結果の行数が不足しています。期待値: 4行以上, 実際: %d行", len(lines))
	}

	// ヘッダー行の検証
	if !strings.Contains(lines[0], "service") || !strings.Contains(lines[0], "cli") || !strings.Contains(lines[0], "mcp") {
		t.Errorf("ヘッダー行が正しくありません: %s", lines[0])
	}

	// データ行の検証
	if !strings.Contains(lines[2], "arithmetic-calculator") {
		t.Errorf("arithmetic-calculatorの行が見つかりません: %s", lines[2])
	}
	if !strings.Contains(lines[3], "base64-extractor") {
		t.Errorf("base64-extractorの行が見つかりません: %s", lines[3])
	}

	// 絵文字の検証
	if !strings.Contains(result, "✅") || !strings.Contains(result, "❌️") {
		t.Error("結果に適切な絵文字が含まれていません")
	}
}
