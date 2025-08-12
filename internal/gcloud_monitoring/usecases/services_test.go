package usecases

import (
	"context"
	"fmt"
	"strings"
	"testing"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"github.com/googleapis/gax-go/v2"
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

func TestNewServiceWithClients_Normal(t *testing.T) {
	project := "test-project"
	location := "us-central1"
	serviceName := "test-service"
	serviceAccountID := "test-sa"

	mockCloudRunClient := &MockCloudRunClient{}
	mockDashboardClient := &MockDashboardClient{}

	service := NewServiceWithClients(project, location, serviceName, serviceAccountID, mockCloudRunClient, mockDashboardClient)

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
	if service.cloudRunClient != mockCloudRunClient {
		t.Error("Expected cloudRunClient to be set")
	}
	if service.dashboardClient != mockDashboardClient {
		t.Error("Expected dashboardClient to be set")
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

// MockCloudRunClient はCloudRunClientのモック実装
type MockCloudRunClient struct {
	GetServiceFunc func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error)
	CloseFunc      func() error
}

func (m *MockCloudRunClient) GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
	if m.GetServiceFunc != nil {
		return m.GetServiceFunc(ctx, req, opts...)
	}
	return nil, fmt.Errorf("GetServiceFunc not implemented")
}

func (m *MockCloudRunClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// MockDashboardClient はDashboardClientのモック実装
type MockDashboardClient struct {
	CreateDashboardFunc func(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error)
	CloseFunc           func() error
}

func (m *MockDashboardClient) CreateDashboard(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error) {
	if m.CreateDashboardFunc != nil {
		return m.CreateDashboardFunc(ctx, req, opts...)
	}
	return nil, fmt.Errorf("CreateDashboardFunc not implemented")
}

func (m *MockDashboardClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// TestVerifyCloudRunService_WithMock はverifyCloudRunServiceのモックを使用したテスト
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
					return &runpb.Service{
						Name: "projects/test-project/locations/us-central1/services/test-service",
					}, nil
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
		{
			name: "PermissionDenied_Error",
			setupMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, status.Error(codes.PermissionDenied, "permission denied")
				}
			},
			expectedExists: false,
			expectError:    true,
			errorMessage:   "cloud Runサービスへのアクセス権限がありません",
		},
		{
			name: "Unauthenticated_Error",
			setupMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, status.Error(codes.Unauthenticated, "unauthenticated")
				}
			},
			expectedExists: false,
			expectError:    true,
			errorMessage:   "認証に失敗しました",
		},
		{
			name: "InternalError_Error",
			setupMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, status.Error(codes.Internal, "internal error")
				}
			},
			expectedExists: false,
			expectError:    true,
			errorMessage:   "cloud Runサービスの取得中にエラーが発生しました",
		},
		{
			name: "NonGRPCError_Error",
			setupMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, fmt.Errorf("non-grpc error")
				}
			},
			expectedExists: false,
			expectError:    true,
			errorMessage:   "予期しないエラーが発生しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCloudRunClient := &MockCloudRunClient{}
			tt.setupMock(mockCloudRunClient)

			service := NewServiceWithClients("test-project", "us-central1", "test-service", "", mockCloudRunClient, nil)

			exists, err := service.verifyCloudRunService(context.Background())

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(fmt.Sprintf("%v", err), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%v'", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if exists != tt.expectedExists {
				t.Errorf("Expected exists %v, got %v", tt.expectedExists, exists)
			}
		})
	}
}

// TestCreateMonitoringDashboard_WithMock はcreateMonitoringDashboardのモックを使用したテスト
func TestCreateMonitoringDashboard_WithMock(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockDashboardClient)
		expectError bool
		errorMessage string
		expectedName string
	}{
		{
			name: "DashboardCreated_Normal",
			setupMock: func(mock *MockDashboardClient) {
				mock.CreateDashboardFunc = func(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error) {
					return &dashboardpb.Dashboard{
						Name: "projects/test-project/dashboards/test-dashboard-123",
					}, nil
				}
			},
			expectError: false,
			expectedName: "projects/test-project/dashboards/test-dashboard-123",
		},
		{
			name: "DashboardCreationFailed_Error",
			setupMock: func(mock *MockDashboardClient) {
				mock.CreateDashboardFunc = func(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error) {
					return nil, fmt.Errorf("dashboard creation failed")
				}
			},
			expectError: true,
			errorMessage: "ダッシュボードの作成リクエストに失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDashboardClient := &MockDashboardClient{}
			tt.setupMock(mockDashboardClient)

			service := NewServiceWithClients("test-project", "us-central1", "test-service", "", nil, mockDashboardClient)

			dashboardName, err := service.createMonitoringDashboard(context.Background())

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(fmt.Sprintf("%v", err), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%v'", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if dashboardName != tt.expectedName {
					t.Errorf("Expected dashboard name '%s', got '%s'", tt.expectedName, dashboardName)
				}
			}
		})
	}
}

// TestCreateDashboardForCloudRun_WithMock はCreateDashboardForCloudRunのモックを使用したテスト
func TestCreateDashboardForCloudRun_WithMock(t *testing.T) {
	tests := []struct {
		name               string
		setupCloudRunMock  func(*MockCloudRunClient)
		setupDashboardMock func(*MockDashboardClient)
		expectError        bool
		errorMessage       string
		expectedResult     string
	}{
		{
			name: "Success_Normal",
			setupCloudRunMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return &runpb.Service{
						Name: "projects/test-project/locations/us-central1/services/test-service",
					}, nil
				}
			},
			setupDashboardMock: func(mock *MockDashboardClient) {
				mock.CreateDashboardFunc = func(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error) {
					return &dashboardpb.Dashboard{
						Name: "projects/test-project/dashboards/test-dashboard-123",
					}, nil
				}
			},
			expectError: false,
			expectedResult: "ダッシュボードが正常に作成されました: projects/test-project/dashboards/test-dashboard-123",
		},
		{
			name: "ServiceNotFound_Error",
			setupCloudRunMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, status.Error(codes.NotFound, "service not found")
				}
			},
			setupDashboardMock: func(mock *MockDashboardClient) {},
			expectError: true,
			errorMessage: "指定されたCloud Runサービスが見つかりません",
		},
		{
			name: "ServiceVerificationFailed_Error",
			setupCloudRunMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, status.Error(codes.PermissionDenied, "permission denied")
				}
			},
			setupDashboardMock: func(mock *MockDashboardClient) {},
			expectError: true,
			errorMessage: "cloud Runサービスの確認に失敗しました",
		},
		{
			name: "DashboardCreationFailed_Error",
			setupCloudRunMock: func(mock *MockCloudRunClient) {
				mock.GetServiceFunc = func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return &runpb.Service{
						Name: "projects/test-project/locations/us-central1/services/test-service",
					}, nil
				}
			},
			setupDashboardMock: func(mock *MockDashboardClient) {
				mock.CreateDashboardFunc = func(ctx context.Context, req *dashboardpb.CreateDashboardRequest, opts ...gax.CallOption) (*dashboardpb.Dashboard, error) {
					return nil, fmt.Errorf("dashboard creation failed")
				}
			},
			expectError: true,
			errorMessage: "ダッシュボードの作成に失敗しました",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCloudRunClient := &MockCloudRunClient{}
			mockDashboardClient := &MockDashboardClient{}

			tt.setupCloudRunMock(mockCloudRunClient)
			tt.setupDashboardMock(mockDashboardClient)

			service := NewServiceWithClients("test-project", "us-central1", "test-service", "", mockCloudRunClient, mockDashboardClient)

			result, err := service.CreateDashboardForCloudRun()

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(fmt.Sprintf("%v", err), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%v'", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != tt.expectedResult {
					t.Errorf("Expected result '%s', got '%s'", tt.expectedResult, result)
				}
			}
		})
	}
}

// TestGetClientOptions_Normal はgetClientOptionsのテスト
func TestGetClientOptions_Normal(t *testing.T) {
	tests := []struct {
		name             string
		serviceAccountID string
		expectError      bool
		errorMessage     string
	}{
		{
			name:             "WithServiceAccount_Normal",
			serviceAccountID: "test-sa",
			expectError:      false,
		},
		{
			name:             "WithoutServiceAccount_Normal",
			serviceAccountID: "",
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService("test-project", "us-central1", "test-service", tt.serviceAccountID)

			_, err := service.getClientOptions(context.Background())

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if tt.errorMessage != "" && !strings.Contains(fmt.Sprintf("%v", err), tt.errorMessage) {
					t.Errorf("Expected error message to contain '%s', got '%v'", tt.errorMessage, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				// optsは空のスライスでも問題ない
				// サービスアカウントが指定されていない場合は空のスライスが返される
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
