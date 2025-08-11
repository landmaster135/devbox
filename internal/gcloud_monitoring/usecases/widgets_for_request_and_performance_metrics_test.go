package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
)

// テスト用の定数
const (
	oneHourSeconds   = 3600
	twentyFourHours  = 86400
	oneMinuteSeconds = 60
)

func TestCreateRequestCountWidget_Normal(t *testing.T) {
	const serviceName = "test-service"
	service := NewService("test-project", "us-central1", serviceName, "test-sa")

	widget := service.createRequestCountWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle01, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
	assert.Equal(t, yAxisLabelOfRequestsPerSecond, xyChart.YAxis.Label)

	// TimeSeriesFilterの確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Contains(t, tsFilter.Filter, serviceName)
	assert.Contains(t, tsFilter.Filter, metricOfRequestCount)
	assert.Equal(t, dashboardpb.Aggregation_ALIGN_RATE, tsFilter.Aggregation.PerSeriesAligner)
}

func TestCreateRequestsByStatusWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createRequestsByStatusWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle02, widget.Title)
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
	assert.Equal(t, containerWidgetTitle03, widget.Title)
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

func TestCreateLogByteByHourWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createLogByteByHourWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle04, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_STACKED_AREA, xyChart.DataSets[0].PlotType)
	assert.Equal(t, yAxisLabelOfBytesPerSecond, xyChart.YAxis.Label)

	// 1時間のアライメント期間の確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Equal(t, int64(oneHourSeconds), tsFilter.Aggregation.AlignmentPeriod.Seconds)
}

func TestCreateRequestLatencyWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createRequestLatencyWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle05, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfLatencyMilliSec, xyChart.YAxis.Label)

	// 各データセットの確認
	for i, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Equal(t, int64(oneMinuteSeconds), tsFilter.Aggregation.AlignmentPeriod.Seconds)

		// パーセンタイルアライナーの確認
		switch i {
		case 0:
			assert.Equal(t, dashboardpb.Aggregation_ALIGN_PERCENTILE_50, tsFilter.Aggregation.PerSeriesAligner)
		case 1:
			assert.Equal(t, dashboardpb.Aggregation_ALIGN_PERCENTILE_95, tsFilter.Aggregation.PerSeriesAligner)
		case 2:
			assert.Equal(t, dashboardpb.Aggregation_ALIGN_PERCENTILE_99, tsFilter.Aggregation.PerSeriesAligner)
		}
	}
}

func TestCreateErrorRateWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createErrorRateWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle06, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
	assert.Equal(t, yAxisLabelOfErrorRatePercentage, xyChart.YAxis.Label)

	// TimeSeriesFilterRatioの確認
	tsFilterRatio := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilterRatio()
	assert.NotNil(t, tsFilterRatio)
	assert.NotNil(t, tsFilterRatio.Numerator)
	assert.NotNil(t, tsFilterRatio.Denominator)
}

func TestCreateMaxConcurrentRequestsWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createMaxConcurrentRequestsWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle07, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfConcurrentRequests, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Equal(t, int64(oneMinuteSeconds), tsFilter.Aggregation.AlignmentPeriod.Seconds)
	}
}

func TestCreateResponseTimeWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createResponseTimeWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle08, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfResponseTimeMilliSec, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Equal(t, int64(oneMinuteSeconds), tsFilter.Aggregation.AlignmentPeriod.Seconds)
	}
}
