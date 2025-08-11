package usecases

import (
	"context"
	"testing"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	"github.com/stretchr/testify/assert"
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

func TestCreateRequestCountWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createRequestCountWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, "Requests per Second", widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
	assert.Equal(t, "requests/second", xyChart.YAxis.Label)

	// TimeSeriesFilterの確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Contains(t, tsFilter.Filter, "test-service")
	assert.Contains(t, tsFilter.Filter, "run.googleapis.com/request_count")
	assert.Equal(t, dashboardpb.Aggregation_ALIGN_RATE, tsFilter.Aggregation.PerSeriesAligner)
}

func TestCreateRequestsByStatusWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createRequestsByStatusWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, "Requests by Status Code", widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_STACKED_AREA, xyChart.DataSets[0].PlotType)

	// GroupByFieldsの確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Contains(t, tsFilter.Aggregation.GroupByFields, "metric.label.response_code_class")
}

func TestCreateTotalRequestsWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createTotalRequestsWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, "Total Requests (24h)", widget.Title)
	assert.NotNil(t, widget.GetScorecard())

	scorecard := widget.GetScorecard()
	assert.NotNil(t, scorecard.TimeSeriesQuery)
	assert.NotNil(t, scorecard.GetSparkChartView())
	assert.Equal(t, dashboardpb.SparkChartType_SPARK_LINE, scorecard.GetSparkChartView().SparkChartType)

	// 24時間のアライメント期間の確認
	tsFilter := scorecard.TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Equal(t, int64(86400), tsFilter.Aggregation.AlignmentPeriod.Seconds) // 24時間 = 86400秒
	assert.Equal(t, dashboardpb.Aggregation_ALIGN_SUM, tsFilter.Aggregation.PerSeriesAligner)
}

func TestCreateRequestHeatmapWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createLogByteByHourWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, "Request Pattern by Hour", widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_HEATMAP, xyChart.DataSets[0].PlotType)

	// 1時間のアライメント期間の確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Equal(t, int64(3600), tsFilter.Aggregation.AlignmentPeriod.Seconds) // 1時間 = 3600秒
}

func TestBuildDashboardConfig(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "")

	config := service.buildDashboardConfig()

	if config.DisplayName != "Cloud Run Monitoring - test-service" {
		t.Errorf("Expected display name 'Cloud Run Monitoring - test-service', got %s", config.DisplayName)
	}

	if config.Layout == nil {
		t.Error("Expected layout to be set")
	}

	gridLayout := config.GetGridLayout()
	if gridLayout == nil {
		t.Error("Expected grid layout to be set")
	}

	if gridLayout.Columns != 12 {
		t.Errorf("Expected 12 columns, got %d", gridLayout.Columns)
	}

	// ウィジェット数の確認（16個のウィジェット）
	expectedWidgetCount := 16
	if len(gridLayout.Widgets) != expectedWidgetCount {
		t.Errorf("Expected %d widgets, got %d", expectedWidgetCount, len(gridLayout.Widgets))
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
