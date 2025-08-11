package usecases

import (
	"testing"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	"github.com/stretchr/testify/assert"
)

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

func TestCreateLogByteByHourWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createLogByteByHourWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, "Log Bytes by Hour", widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_STACKED_AREA, xyChart.DataSets[0].PlotType)

	// 1時間のアライメント期間の確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Equal(t, int64(3600), tsFilter.Aggregation.AlignmentPeriod.Seconds) // 1時間 = 3600秒
}
