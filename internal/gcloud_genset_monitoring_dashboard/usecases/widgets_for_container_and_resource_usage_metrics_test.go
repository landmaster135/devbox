package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
)

func TestCreateContainerInstanceCountWidget_Normal(t *testing.T) {
	const serviceName = "test-service"
	service := NewService("test-project", "us-central1", serviceName, "test-sa")

	widget := service.createContainerInstanceCountWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle11, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
	assert.Equal(t, yAxisLabelOfInstanceCount, xyChart.YAxis.Label)

	// TimeSeriesFilterの確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Contains(t, tsFilter.Filter, serviceName)
	assert.Contains(t, tsFilter.Filter, metricOfContainerInstanceCount)
}

func TestCreateContainerStartupLatenciesWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createContainerStartupLatenciesWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle12, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfStartupLatencyMilliSec, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Contains(t, tsFilter.Filter, metricOfContainerStartupLatencies)
	}
}

func TestCreateContainerBillableInstanceTimeWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createContainerBillableInstanceTimeWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle13, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 1)
	assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, xyChart.DataSets[0].PlotType)
	assert.Equal(t, yAxisLabelOfBillableInstanceTimeSeconds, xyChart.YAxis.Label)

	// TimeSeriesFilterの確認
	tsFilter := xyChart.DataSets[0].TimeSeriesQuery.GetTimeSeriesFilter()
	assert.NotNil(t, tsFilter)
	assert.Contains(t, tsFilter.Filter, metricOfContainerBillableInstance)
}

func TestCreateContainerCPUUtilizationsWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createContainerCPUUtilizationsWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle14, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfCPUUtilizationPercentage, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Contains(t, tsFilter.Filter, metricOfContainerCPUUtilizations)
	}
}

func TestCreateContainerMemoryUtilizationsWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createContainerMemoryUtilizationsWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle15, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfMemoryUtilizationPercentage, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Contains(t, tsFilter.Filter, metricOfContainerMemoryUtilizations)
	}
}

func TestCreateContainerMemoryUsageTimeWidget_Normal(t *testing.T) {
	service := NewService("test-project", "us-central1", "test-service", "test-sa")

	widget := service.createContainerMemoryUsageTimeWidget()

	assert.NotNil(t, widget)
	assert.Equal(t, containerWidgetTitle16, widget.Title)
	assert.NotNil(t, widget.GetXyChart())

	xyChart := widget.GetXyChart()
	assert.Len(t, xyChart.DataSets, 3) // P50, P95, P99
	assert.Equal(t, yAxisLabelOfMemoryUsageBytes, xyChart.YAxis.Label)

	// 各データセットの確認
	for _, dataset := range xyChart.DataSets {
		assert.Equal(t, dashboardpb.XyChart_DataSet_LINE, dataset.PlotType)
		assert.NotEmpty(t, dataset.LegendTemplate)

		tsFilter := dataset.TimeSeriesQuery.GetTimeSeriesFilter()
		assert.NotNil(t, tsFilter)
		assert.Contains(t, tsFilter.Filter, metricOfContainerMemoryUsage)
	}
}
