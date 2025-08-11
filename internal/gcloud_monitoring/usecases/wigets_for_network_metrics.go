package usecases

import (
	"time"

	dashboardpb "cloud.google.com/go/monitoring/dashboard/apiv1/dashboardpb"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
)

// createNetworkSentBytesWidget はネットワーク送信バイト数ウィジェットを作成する
func (s *Service) createNetworkSentBytesWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle51,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForNetworkSentBytes(),
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
					Label: yAxisLabelOfNetworkBytesPerSecond,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}

// createNetworkReceivedBytesWidget はネットワーク受信バイト数ウィジェットを作成する
func (s *Service) createNetworkReceivedBytesWidget() *dashboardpb.Widget {
	return &dashboardpb.Widget{
		Title: containerWidgetTitle52,
		Content: &dashboardpb.Widget_XyChart{
			XyChart: &dashboardpb.XyChart{
				DataSets: []*dashboardpb.XyChart_DataSet{
					{
						TimeSeriesQuery: &dashboardpb.TimeSeriesQuery{
							Source: &dashboardpb.TimeSeriesQuery_TimeSeriesFilter{
								TimeSeriesFilter: &dashboardpb.TimeSeriesFilter{
									Filter: s.createPromQLForNetworkReceivedBytes(),
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
					Label: yAxisLabelOfNetworkBytesPerSecond,
					Scale: dashboardpb.XyChart_Axis_LINEAR,
				},
			},
		},
	}
}
