package usecases

import (
	"time"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createContainerInstanceCountWidget はコンテナインスタンス数ウィジェットを作成する
func (s *Service) createContainerInstanceCountWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Container Instance Count",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerInstanceCount(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_MEAN,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_LINE,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfInstanceCount,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createContainerStartupLatenciesWidget はコンテナ起動レイテンシウィジェットを作成する
func (s *Service) createContainerStartupLatenciesWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Container Startup Latency (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerStartupLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerStartupLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerStartupLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfStartupLatencyMilliSec,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createContainerBillableInstanceTimeWidget は課金対象インスタンス時間ウィジェットを作成する
func (s *Service) createContainerBillableInstanceTimeWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Billable Instance Time",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerBillableInstanceTime(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_LINE,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfBillableInstanceTimeSeconds,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createContainerCPUUtilizationsWidget はCPU使用率ウィジェットを作成する
func (s *Service) createContainerCPUUtilizationsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "CPU Utilization (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerCPUUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerCPUUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerCPUUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfCPUUtilizationPercentage,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createContainerMemoryUtilizationsWidget はメモリ使用率ウィジェットを作成する
func (s *Service) createContainerMemoryUtilizationsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Memory Utilization (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUtilizations(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfMemoryUtilizationPercentage,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createContainerMemoryUsageTimeWidget はメモリ使用量ウィジェットを作成する
func (s *Service) createContainerMemoryUsageTimeWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: "Memory Usage (P50/P95/P99)",
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUsageTime(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUsageTime(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForContainerMemoryUsageTime(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: DatasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: YAxisLabelOfMemoryUsageBytes,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}
