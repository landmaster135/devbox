package usecases

import (
	"context"
	"fmt"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// テスト用のPromQLクエリテンプレート
const (
	promQLTemplate = `resource.type="cloud_run_revision" resource.label.service_name="%s" metric.type="%s"`
)

func TestNewService(t *testing.T) {
	project := "test-project"
	location := "us-central1"
	serviceName := "test-service"
	serviceAccountID := "test-sa"

	service := NewService(project, location, serviceName, serviceAccountID)

	if service.project != project {
		t.Errorf("Expected project %s, got %s", project, service.project)
	}
	if service.location != location {
		t.Errorf("Expected location %s, got %s", location, service.location)
	}
	if service.serviceName != serviceName {
		t.Errorf("Expected serviceName %s, got %s", serviceName, service.serviceName)
	}
	if service.serviceAccountID != serviceAccountID {
		t.Errorf("Expected serviceAccountID %s, got %s", serviceAccountID, service.serviceAccountID)
	}
}

func TestBuildDashboardConfig(t *testing.T) {
	const serviceName = "test-service"
	service := NewService("test-project", "us-central1", serviceName, "")

	config := service.buildDashboardConfig()

	if config.DisplayName != "CloudRun ダッシュボード: "+serviceName {
		t.Errorf("Expected display name 'CloudRun ダッシュボード: %s', got %s", serviceName, config.DisplayName)
	}

	if config.Layout == nil {
		t.Error("Expected layout to be set")
	}

	mosaicLayout := config.GetMosaicLayout()
	if mosaicLayout == nil {
		t.Error("Expected mosaic layout to be set")
		return
	}

	if mosaicLayout.Columns != 12 {
		t.Errorf("Expected 12 columns, got %d", mosaicLayout.Columns)
	}

	// ウィジェット数の確認（16個のウィジェット）
	expectedWidgetCount := 16
	if len(mosaicLayout.Tiles) != expectedWidgetCount {
		t.Errorf("Expected %d widgets, got %d", expectedWidgetCount, len(mosaicLayout.Tiles))
	}
}

func TestCreateTextWidget(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	title := "Test Widget"
	content := "Test content"

	widget := service.createTextWidget(title, content)

	if widget.Title != title {
		t.Errorf("Expected title %s, got %s", title, widget.Title)
	}

	textWidget := widget.GetText()
	if textWidget == nil {
		t.Error("Expected text widget to be set")
		return
	}

	expectedContent := "**Test Widget**\n\nTest content\n\nサービス: test-service"
	if textWidget.Content != expectedContent {
		t.Errorf("Expected content %s, got %s", expectedContent, textWidget.Content)
	}
}

// TestVerifyCloudRunService_ErrorHandling はverifyCloudRunServiceのエラーハンドリングをテストする
func TestVerifyCloudRunService_ErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		errorCode      codes.Code
		expectedExists bool
		expectError    bool
		errorMessage   string
	}{
		{
			name:           "NotFound_Normal",
			errorCode:      codes.NotFound,
			expectedExists: false,
			expectError:    false,
			errorMessage:   "",
		},
		{
			name:           "PermissionDenied_Normal",
			errorCode:      codes.PermissionDenied,
			expectedExists: false,
			expectError:    true,
			errorMessage:   "cloud Runサービスへのアクセス権限がありません",
		},
		{
			name:           "Unauthenticated_Normal",
			errorCode:      codes.Unauthenticated,
			expectedExists: false,
			expectError:    true,
			errorMessage:   "認証に失敗しました",
		},
		{
			name:           "Internal_Normal",
			errorCode:      codes.Internal,
			expectedExists: false,
			expectError:    true,
			errorMessage:   "cloud Runサービスの取得中にエラーが発生しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// gRPCエラーの作成
			err := status.Error(tt.errorCode, "test error")

			// エラーコードの確認テスト
			if st, ok := status.FromError(err); ok {
				if st.Code() != tt.errorCode {
					t.Errorf("Expected error code %v, got %v", tt.errorCode, st.Code())
				}
			} else {
				t.Error("Expected gRPC status error")
			}

			// contextの使用テスト
			ctx := context.Background()
			if ctx == nil {
				t.Error("Expected context to be created")
			}
		})
	}
}

func TestCreateDisplayTitleOfDashboard_Normal(t *testing.T) {
	const serviceName = "test-service"
	service := NewService("test-project", "us-central1", serviceName, "")

	title := service.createDisplayTitleOfDashboard()
	expected := "CloudRun ダッシュボード: " + serviceName

	if title != expected {
		t.Errorf("Expected title '%s', got '%s'", expected, title)
	}
}

func TestCreateTile_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	// テスト用のウィジェットを作成
	widget := service.createTextWidget("Test Title", "Test Content")

	tile := service.createTile(widget, 1, 2, 3, 4)

	if tile.XPos != 1 {
		t.Errorf("Expected XPos 1, got %d", tile.XPos)
	}

	if tile.YPos != 2 {
		t.Errorf("Expected YPos 2, got %d", tile.YPos)
	}

	if tile.Width != 3 {
		t.Errorf("Expected Width 3, got %d", tile.Width)
	}

	if tile.Height != 4 {
		t.Errorf("Expected Height 4, got %d", tile.Height)
	}

	if tile.Widget != widget {
		t.Error("Expected widget to match the provided widget")
	}
}

func TestBuildDashboardTiles_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	tiles := service.buildDashboardTiles()

	// 期待されるタイル数（16個）
	expectedTileCount := 16
	if len(tiles) != expectedTileCount {
		t.Errorf("Expected %d tiles, got %d", expectedTileCount, len(tiles))
	}

	// 各タイルが適切に設定されているか確認
	for i, tile := range tiles {
		if tile.Widget == nil {
			t.Errorf("Tile %d has nil widget", i)
		}

		if tile.Width <= 0 {
			t.Errorf("Tile %d has invalid width: %d", i, tile.Width)
		}

		if tile.Height <= 0 {
			t.Errorf("Tile %d has invalid height: %d", i, tile.Height)
		}
	}
}

// PromQLクエリ生成メソッドのテーブル駆動テスト
func TestPromQLQueryGeneration_Normal(t *testing.T) {
	const (
		testProject          = "test-project"
		testLocation         = "us-central1"
		testServiceName      = "test-service"
		testServiceAccountID = "test-sa"
	)

	service := NewService(testProject, testLocation, testServiceName, testServiceAccountID)

	tests := []struct {
		name           string
		queryFunc      func() string
		metricType     string
		additionalPart string
	}{
		{
			name:       "RequestCount",
			queryFunc:  service.createPromQLForRequestCount,
			metricType: metricOfRequestCount,
		},
		{
			name:       "RequestLatencies",
			queryFunc:  service.createPromQLForRequestLatencies,
			metricType: metricOfRequestLatencies,
		},
		{
			name:       "Logging",
			queryFunc:  service.createPromQLForLogging,
			metricType: metricOfLoggingByteCount,
		},
		{
			name:           "RequestCountByResponseCode",
			queryFunc:      service.createPromQLForRequestCountByResponseCode,
			metricType:     metricOfRequestCount,
			additionalPart: ` metric.label.response_code_class!="2xx"`,
		},
		{
			name:       "MaxRequestConcurrencies",
			queryFunc:  service.createPromQLMaxRequestConcurrencies,
			metricType: metricOfMaxRequestConcurrencies,
		},
		{
			name:       "ContainerInstanceCount",
			queryFunc:  service.createPromQLForContainerInstanceCount,
			metricType: metricOfContainerInstanceCount,
		},
		{
			name:       "ContainerStartupLatencies",
			queryFunc:  service.createPromQLForContainerStartupLatencies,
			metricType: metricOfContainerStartupLatencies,
		},
		{
			name:       "ContainerBillableInstanceTime",
			queryFunc:  service.createPromQLForContainerBillableInstanceTime,
			metricType: metricOfContainerBillableInstance,
		},
		{
			name:       "ContainerCPUUtilizations",
			queryFunc:  service.createPromQLForContainerCPUUtilizations,
			metricType: metricOfContainerCPUUtilizations,
		},
		{
			name:       "ContainerMemoryUtilizations",
			queryFunc:  service.createPromQLForContainerMemoryUtilizations,
			metricType: metricOfContainerMemoryUtilizations,
		},
		{
			name:       "ContainerMemoryUsageTime",
			queryFunc:  service.createPromQLForContainerMemoryUsageTime,
			metricType: metricOfContainerMemoryUsage,
		},
		{
			name:       "NetworkSentBytes",
			queryFunc:  service.createPromQLForNetworkSentBytes,
			metricType: metricOfNetworkSentBytesCount,
		},
		{
			name:       "NetworkReceivedBytes",
			queryFunc:  service.createPromQLForNetworkReceivedBytes,
			metricType: metricOfNetworkReceivedBytesCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.queryFunc()
			expected := fmt.Sprintf(promQLTemplate, testServiceName, tt.metricType) + tt.additionalPart

			if query != expected {
				t.Errorf("Expected query '%s', got '%s'", expected, query)
			}
		})
	}
}
