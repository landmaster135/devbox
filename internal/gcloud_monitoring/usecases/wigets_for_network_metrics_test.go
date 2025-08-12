package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
)

// ネットワークウィジェットのテーブル駆動テスト
func TestNetworkWidgets_Normal(t *testing.T) {
	const (
		// 共通のテストデータ
		testProject     = "test-project"
		testLocation    = "us-central1"
		testServiceName = "test-service"
		testSA          = "test-sa"
	)
	service := NewService(testProject, testLocation, testServiceName, testSA)

	tests := []struct {
		name           string
		widgetFunc     func() *dashboardpb.Widget
		expectedTitle  string
		expectedMetric string
	}{
		{
			name:           containerWidgetTitle51,
			widgetFunc:     service.createNetworkSentBytesWidget,
			expectedTitle:  containerWidgetTitle51,
			expectedMetric: metricOfNetworkSentBytesCount,
		},
		{
			name:           containerWidgetTitle52,
			widgetFunc:     service.createNetworkReceivedBytesWidget,
			expectedTitle:  containerWidgetTitle52,
			expectedMetric: metricOfNetworkReceivedBytesCount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := tt.widgetFunc()

			assert.NotNil(t, widget)
			assert.Equal(t, tt.expectedTitle, widget.Title)
			assert.NotNil(t, widget.GetXyChart())

			xyChart := widget.GetXyChart()
			assert.Len(t, xyChart.DataSets, 1)
			assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
			assert.Equal(t, yAxisLabelOfBytesPerSecond, xyChart.YAxis.Label)

			// TimeSeriesFilterの確認
			tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
			assert.NotNil(t, tsFilter)
			assert.Contains(t, tsFilter.Filter, testServiceName)
			assert.Contains(t, tsFilter.Filter, tt.expectedMetric)
			assert.Equal(t, dashboardpb.Aggregation_ALIGN_RATE, tsFilter.Aggregation.PerSeriesAligner)
		})
	}
}
