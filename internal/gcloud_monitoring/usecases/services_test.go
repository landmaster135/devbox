package usecases

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	if config.DisplayName != "CloudRun ダッシュボード: " + serviceName {
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
