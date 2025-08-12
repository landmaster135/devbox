package usecases

import (
	"time"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	"google.golang.org/protobuf/types/known/durationpb"
)

// createRequestCountWidget はリクエスト数（req/sec）ウィジェットを作成する
func (s *Service) createRequestCountWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle01,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestCount(),
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
					Label: yAxisLabelOfRequestsPerSecond,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createRequestsByStatusWidget はステータス別リクエスト数ウィジェットを作成する
func (s *Service) createRequestsByStatusWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle02,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestCount(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										GroupByFields:      []string{"metric.label.response_code_class"},
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_STACKED_AREA,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfRequestsPerSecond,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createTotalRequestsWidget は累積リクエスト数ウィジェットを作成する
func (s *Service) createTotalRequestsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle03,
		Content: &dashboardpb.Widget_Scorecard{
			Scorecard: &dashboardpb.Scorecard{
				TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
					Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
						TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
							Filter: s.createPromQLForRequestCount(),
							Aggregation: &dashboardpb.Aggregation{
								AlignmentPeriod:    durationpb.New(86400 * time.Second), // 24時間
								PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_SUM,
								CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
							},
						},
					},
				},
				DataView: &dashboardpb.Scorecard_SparkChartView_{
					SparkChartView: &dashboardpb.Scorecard_SparkChartView{
						SparkChartType: dashboardpb.SparkChartType_SPARK_LINE,
					},
				},
			},
		},
	}
}

// createLogByteByHourWidget はログバイト数の時間ごとの集計ウィジェットを作成する
func (s *Service) createLogByteByHourWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle04,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForLogging(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(3600 * time.Second), // 1時間
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										GroupByFields:      []string{"metric.label.severity"},
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_STACKED_AREA,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfBytesPerSecond,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
				ChartOptions: &dashboardpb.ChartOptions{
					Mode: dashboardpb.ChartOptions_COLOR,
				},
			},
		},
	}
}

// createRequestLatencyWidget はリクエストレイテンシ (P50,P95,P99) ウィジェットを作成する
func (s *Service) createRequestLatencyWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle05,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfLatencyMilliSec,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createErrorRateWidget はエラー率ウィジェットを作成する
func (s *Service) createErrorRateWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle06,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilterRatio{
								TimeSeriesFilterRatio: &dashboardpb.TimeSeriesFilterRatio{
									Numerator: &dashboardpb.TimeSeriesFilterRatio_RatioPart{
										Filter: s.createPromQLForRequestCountByResponseCode(),
										Aggregation: &dashboardpb.Aggregation{
											AlignmentPeriod:    durationpb.New(60 * time.Second),
											PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
											CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										},
									},
									Denominator: &dashboardpb.TimeSeriesFilterRatio_RatioPart{
										Filter: s.createPromQLForRequestCount(),
										Aggregation: &dashboardpb.Aggregation{
											AlignmentPeriod:    durationpb.New(60 * time.Second),
											PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_RATE,
											CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_SUM,
										},
									},
								},
							},
						},
						PlotType: dashboardpb.XyChart_DataSet_LINE,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfErrorRatePercentage,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createMaxConcurrentRequestsWidget は最大同時リクエストウィジェットを作成する
func (s *Service) createMaxConcurrentRequestsWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle07,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLMaxRequestConcurrencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLMaxRequestConcurrencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLMaxRequestConcurrencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfConcurrentRequests,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createResponseTimeWidget はレスポンス時間 (P50/P95/P99) ウィジェットを作成する
func (s *Service) createResponseTimeWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle08,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_50,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP50,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_95,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP95,
					},
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForRequestLatencies(),
									Aggregation: &dashboardpb.Aggregation{
										AlignmentPeriod:    durationpb.New(60 * time.Second),
										PerSeriesAligner:   dashboardpb.Aggregation_ALIGN_PERCENTILE_99,
										CrossSeriesReducer: dashboardpb.Aggregation_REDUCE_MEAN,
									},
								},
							},
						},
						PlotType:       dashboardpb.XyChart_DataSet_LINE,
						LegendTemplate: datasetLegendTemplateOfP99,
					},
				},
				YAxis: &dashboardpb.XyChart_Axis{
					Label: yAxisLabelOfResponseTimeMilliSec,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}
